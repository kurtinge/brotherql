// Package brotherql is a Go library and CLI for printing to Brother QL
// label printers over USB.
package brotherql

// LabelType describes a Brother label specification.
type LabelType struct {
	Name     string // identifier like "62" or "12"
	WidthMM  int    // physical width of the label media in mm
	HeightMM int    // physical height in mm; 0 indicates continuous tape
	WidthPx  int    // printable dots (brother_ql dots_printable)
	OffsetPx int    // brother_ql offset_r; left pad (dots) of the emitted raster row
}

// Predefined label types supported by QL printers.
var (
	Label12    = LabelType{Name: "12", WidthMM: 12, HeightMM: 0, WidthPx: 106, OffsetPx: 29}
	Label62    = LabelType{Name: "62", WidthMM: 62, HeightMM: 0, WidthPx: 696, OffsetPx: 12}
	Label62x29 = LabelType{Name: "62x29", WidthMM: 62, HeightMM: 29, WidthPx: 696, OffsetPx: 12}
)
