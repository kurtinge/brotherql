package brotherql

import "fmt"

// Status describes the printer's current state.
type Status struct {
	Ready       bool   // true if printer is ready to accept a job
	Error       string // human-readable error message; empty if no error
	MediaWidth  int    // mm
	MediaLength int    // mm; 0 indicates continuous tape
}

// parseStatus parses a 32-byte Brother status response.
// See Brother QL-700 raster command reference for byte layout.
func parseStatus(raw []byte) (Status, error) {
	if len(raw) != 32 {
		return Status{}, fmt.Errorf("brotherql: status response must be 32 bytes, got %d", len(raw))
	}

	s := Status{
		MediaWidth:  int(raw[10]),
		MediaLength: int(raw[17]),
	}

	errorInfo1 := raw[8]
	errorInfo2 := raw[9]

	switch {
	case errorInfo1&0x01 != 0:
		s.Error = "no media"
	case errorInfo1&0x02 != 0:
		s.Error = "end of media"
	case errorInfo1&0x04 != 0:
		s.Error = "tape cutter jam"
	case errorInfo1&0x10 != 0:
		s.Error = "main unit in use"
	case errorInfo1&0x80 != 0:
		s.Error = "fan error"
	case errorInfo2&0x04 != 0:
		s.Error = "transmission error"
	case errorInfo2&0x10 != 0:
		s.Error = "cover open"
	case errorInfo2&0x40 != 0:
		s.Error = "cannot feed"
	case errorInfo2&0x80 != 0:
		s.Error = "system error"
	}

	s.Ready = s.Error == ""
	return s, nil
}
