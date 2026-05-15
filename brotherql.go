package brotherql

import (
	"errors"
	"fmt"
)

// USB identifiers for the QL-700.
const (
	qlVendorID  = 0x04F9
	qlProductID = 0x2042
)

// Error sentinels for predictable error matching.
var (
	ErrPrinterNotFound = errors.New("brotherql: no QL-700 printer found")
	ErrInvalidSerial   = errors.New("brotherql: invalid serial number")
)

// Info describes a discovered printer (returned by List, not opened).
type Info struct {
	Serial  string // printer serial number
	Model   string // e.g. "QL-700"
	USBPath string // human-readable USB location
}

// String formats Info as a one-line human-readable description.
func (i Info) String() string {
	return fmt.Sprintf("%s serial=%s (%s)", i.Model, i.Serial, i.USBPath)
}

// Open finds and opens the first connected QL-700 over USB.
// Returns ErrPrinterNotFound if none are connected.
func Open() (*Printer, error) {
	return openUSB("")
}

// OpenBySerial opens a specific QL-700 by serial number.
// Useful when multiple printers are connected.
func OpenBySerial(serial string) (*Printer, error) {
	if serial == "" {
		return nil, ErrInvalidSerial
	}
	return openUSB(serial)
}

// List returns all currently connected QL-700 printers without opening them.
// The returned Infos can be used to call OpenBySerial.
func List() ([]Info, error) {
	return listUSB()
}
