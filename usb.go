package brotherql

import (
	"fmt"
	"runtime"

	"github.com/google/gousb"
)

// usbTransport implements transport over a real USB connection to a Brother QL printer.
type usbTransport struct {
	ctx      *gousb.Context
	dev      *gousb.Device
	intf     *gousb.Interface
	intfDone func()
	out      *gousb.OutEndpoint
	in       *gousb.InEndpoint
}

// listUSB enumerates connected Brother QL printers without claiming them.
func listUSB() ([]Info, error) {
	ctx := gousb.NewContext()
	defer func() { _ = ctx.Close() }()

	var infos []Info
	devs, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		_, ok := findModel(uint16(desc.Vendor), uint16(desc.Product))
		return ok
	})
	if err != nil {
		return nil, fmt.Errorf("brotherql: enumerate USB: %w", err)
	}
	defer func() {
		for _, d := range devs {
			_ = d.Close()
		}
	}()

	for _, d := range devs {
		serial, err := d.SerialNumber()
		if err != nil {
			serial = "unknown"
		}
		m, _ := findModel(uint16(d.Desc.Vendor), uint16(d.Desc.Product))
		infos = append(infos, Info{
			Serial:  serial,
			Model:   m.Name,
			USBPath: fmt.Sprintf("bus %d addr %d", d.Desc.Bus, d.Desc.Address),
		})
	}
	return infos, nil
}

// openUSB opens a Brother QL printer by serial. If serial is empty, opens the first found.
func openUSB(serial string) (*Printer, error) {
	ctx := gousb.NewContext()

	devs, err := ctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		_, ok := findModel(uint16(desc.Vendor), uint16(desc.Product))
		return ok
	})
	if err != nil {
		_ = ctx.Close()
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
			_ = d.Close()
		}
	}
	if dev == nil {
		_ = ctx.Close()
		return nil, ErrPrinterNotFound
	}

	// SetAutoDetach asks libusb to detach any kernel driver claiming the
	// interface before we claim it ourselves. On Linux the usblp driver
	// commonly holds the interface and we genuinely need this. On macOS
	// libusb implements detach via USBDeviceReEnumerate, which requires
	// root or a special entitlement — calling it from an unprivileged
	// process fails with EACCES and prevents the subsequent claim. So
	// skip the call on macOS and let the claim either succeed (driver
	// only matched, not actively holding the interface) or fail with a
	// clearer BUSY error.
	if runtime.GOOS != "darwin" {
		_ = dev.SetAutoDetach(true)
	}

	intf, intfDone, err := dev.DefaultInterface()
	if err != nil {
		_ = dev.Close()
		_ = ctx.Close()
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
		_ = dev.Close()
		_ = ctx.Close()
		return nil, fmt.Errorf("brotherql: missing bulk endpoints")
	}

	serial2, _ := dev.SerialNumber()
	m, _ := findModel(uint16(dev.Desc.Vendor), uint16(dev.Desc.Product))
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
		model:  m,
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
		_ = u.dev.Close()
	}
	if u.ctx != nil {
		_ = u.ctx.Close()
	}
	return nil
}
