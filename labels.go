// Package brotherql is a Go library and CLI for printing to Brother QL-700
// label printers over USB.
package brotherql

// LabelType describes a Brother label specification.
type LabelType struct {
	Name     string // identifier like "62" or "62x29"
	WidthMM  int    // physical width of the label media in mm
	HeightMM int    // physical height in mm; 0 indicates continuous tape
	WidthPx  int    // pixel width the printer expects for raster data
}

// Predefined label types supported by QL-700.
var (
	Label62    = LabelType{Name: "62", WidthMM: 62, HeightMM: 0, WidthPx: 696}
	Label62x29 = LabelType{Name: "62x29", WidthMM: 62, HeightMM: 29, WidthPx: 696}
)
