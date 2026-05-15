package brotherql

import (
	"fmt"
	"image"
)

// Printer represents an open connection to a Brother QL-700.
// Always Close() when done; defer is recommended.
type Printer struct {
	tr     transport
	serial string
}

// Close releases the underlying transport.
func (p *Printer) Close() error {
	if p.tr == nil {
		return nil
	}
	return p.tr.Close()
}

// Status sends a status request to the printer and parses the response.
func (p *Printer) Status() (Status, error) {
	if _, err := p.tr.Write([]byte{0x1B, 0x69, 0x53}); err != nil {
		return Status{}, fmt.Errorf("brotherql: write status request: %w", err)
	}
	buf := make([]byte, 32)
	n, err := p.tr.Read(buf)
	if err != nil {
		return Status{}, fmt.Errorf("brotherql: read status: %w", err)
	}
	if n != 32 {
		return Status{}, fmt.Errorf("brotherql: short status read: %d bytes", n)
	}
	return parseStatus(buf)
}

// Print sends an image to the printer and waits for completion.
// The image is resized and thresholded to fit opts.Label.
func (p *Printer) Print(img image.Image, opts PrintOptions) error {
	if opts.Label.WidthPx == 0 {
		return fmt.Errorf("brotherql: PrintOptions.Label is required")
	}

	raster, err := encodeRaster(img, opts.Label)
	if err != nil {
		return fmt.Errorf("brotherql: raster encode: %w", err)
	}

	cmd := buildPrintJob(raster, opts.Label, opts)
	if _, err := p.tr.Write(cmd); err != nil {
		return fmt.Errorf("brotherql: write print job: %w", err)
	}

	buf := make([]byte, 32)
	if _, err := p.tr.Read(buf); err != nil {
		return fmt.Errorf("brotherql: read post-print status: %w", err)
	}
	s, err := parseStatus(buf)
	if err != nil {
		return fmt.Errorf("brotherql: parse post-print status: %w", err)
	}
	if s.Error != "" {
		return fmt.Errorf("brotherql: printer error after print: %s", s.Error)
	}
	return nil
}

// PrintOptions configures a single print job.
type PrintOptions struct {
	Label   LabelType // required: which label media is loaded
	Copies  int       // number of copies; 0 or 1 prints once
	AutoCut bool      // cut paper after print
	HighDPI bool      // 600 DPI mode (default 300)
}
