package brotherql

import (
	"fmt"

	"github.com/google/gousb"
)

// usbTransport implements transport over a real USB connection to a QL-700.
type usbTransport struct {
	ctx      *gousb.Context
	dev      *gousb.Device
	intf     *gousb.Interface
	intfDone func()
	out      *gousb.OutEndpoint
	in       *gousb.InEndpoint
}

// listUSB enumerates connected QL-700 printers without claiming them.
func listUSB() ([]Info, error) {
	ctx := gousb.NewContext()
	defer ctx.Close()

	var infos []Info
	devs, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		return desc.Vendor == qlVendorID && desc.Product == qlProductID
	})
	if err != nil {
		return nil, fmt.Errorf("brotherql: enumerate USB: %w", err)
	}
	defer func() {
		for _, d := range devs {
			d.Close()
		}
	}()

	for _, d := range devs {
		serial, err := d.SerialNumber()
		if err != nil {
			serial = "unknown"
		}
		infos = append(infos, Info{
			Serial:  serial,
			Model:   "QL-700",
			USBPath: fmt.Sprintf("bus %d addr %d", d.Desc.Bus, d.Desc.Address),
		})
	}
	return infos, nil
}

// openUSB opens a QL-700 by serial. If serial is empty, opens the first found.
func openUSB(serial string) (*Printer, error) {
	ctx := gousb.NewContext()

	devs, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		return desc.Vendor == qlVendorID && desc.Product == qlProductID
	})
	if err != nil {
		ctx.Close()
		return nil, fmt.Errorf("brotherql: enumerate USB: %w", err)
	}

	var dev *gousb.Device
	for _, d := range devs {
		if serial == "" {
			dev = d
			break
		}
		s, err := d.SerialNumber()
		if err == nil && s == serial {
			dev = d
			break
		}
	}
	for _, d := range devs {
		if d != dev {
			d.Close()
		}
	}
	if dev == nil {
		ctx.Close()
		return nil, ErrPrinterNotFound
	}

	// SetAutoDetach is a no-op on macOS where there's no kernel driver
	// to detach; safe to ignore the error.
	_ = dev.SetAutoDetach(true)

	intf, intfDone, err := dev.DefaultInterface()
	if err != nil {
		dev.Close()
		ctx.Close()
		return nil, fmt.Errorf("brotherql: claim interface: %w", err)
	}

	var outEP *gousb.OutEndpoint
	var inEP *gousb.InEndpoint
	for _, ep := range intf.Setting.Endpoints {
		switch ep.Direction {
		case gousb.EndpointDirectionOut:
			if e, err := intf.OutEndpoint(ep.Number); err == nil {
				outEP = e
			}
		case gousb.EndpointDirectionIn:
			if e, err := intf.InEndpoint(ep.Number); err == nil {
				inEP = e
			}
		}
	}
	if outEP == nil || inEP == nil {
		intfDone()
		dev.Close()
		ctx.Close()
		return nil, fmt.Errorf("brotherql: missing bulk endpoints")
	}

	serial2, _ := dev.SerialNumber()
	return &Printer{
		tr: &usbTransport{
			ctx:      ctx,
			dev:      dev,
			intf:     intf,
			intfDone: intfDone,
			out:      outEP,
			in:       inEP,
		},
		serial: serial2,
	}, nil
}

func (u *usbTransport) Write(p []byte) (int, error) {
	return u.out.Write(p)
}

func (u *usbTransport) Read(p []byte) (int, error) {
	return u.in.Read(p)
}

func (u *usbTransport) Close() error {
	if u.intfDone != nil {
		u.intfDone()
	}
	if u.dev != nil {
		u.dev.Close()
	}
	if u.ctx != nil {
		u.ctx.Close()
	}
	return nil
}
