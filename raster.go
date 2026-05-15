package brotherql

import (
	"fmt"
	"image"
	"image/color"

	"golang.org/x/image/draw"
)

// bytesPerRow is the QL-700 raster row width in bytes.
// 90 bytes = 720 bits = 12 left padding + 696 printable + 12 right padding.
const bytesPerRow = 90

// leftPaddingBits is the QL-700's mechanical left margin: pixels are shifted
// right by this many bits within each row.
const leftPaddingBits = 12

// encodeRaster converts an image to Brother raster bytes for the given label.
// The image is resized to label.WidthPx wide while preserving aspect ratio,
// then thresholded to 1-bit black/white.
//
// Each output row is 90 bytes (720 bits). The image is mirrored horizontally
// to match the QL-700's print orientation, then placed at dots
// [leftPaddingBits, leftPaddingBits+WidthPx). Within each byte, bits are
// packed MSB-first (bit 7 = leftmost dot).
func encodeRaster(img image.Image, label LabelType) ([]byte, error) {
	if label.WidthPx == 0 {
		return nil, fmt.Errorf("brotherql: label.WidthPx must be set")
	}

	resized := resizeToWidth(img, label.WidthPx)
	bw := toBlackWhite(resized)

	h := bw.Bounds().Dy()
	out := make([]byte, h*bytesPerRow)

	// Image x=0 maps to dot (leftPaddingBits + WidthPx - 1), so the image
	// is mirrored across the printable area as the QL-700 expects.
	mirroredBase := leftPaddingBits + label.WidthPx - 1

	for y := 0; y < h; y++ {
		rowOffset := y * bytesPerRow
		for x := 0; x < label.WidthPx; x++ {
			if bw.GrayAt(x, y).Y == 0 { // black
				dotIdx := mirroredBase - x
				byteIdx := dotIdx / 8
				bitIdx := uint(7 - (dotIdx % 8)) // MSB-first
				out[rowOffset+byteIdx] |= 1 << bitIdx
			}
		}
	}
	return out, nil
}

func resizeToWidth(src image.Image, targetWidth int) image.Image {
	b := src.Bounds()
	if b.Dx() == targetWidth {
		return src
	}
	ratio := float64(targetWidth) / float64(b.Dx())
	targetHeight := int(float64(b.Dy()) * ratio)
	dst := image.NewRGBA(image.Rect(0, 0, targetWidth, targetHeight))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, b, draw.Over, nil)
	return dst
}

// toBlackWhite thresholds an image to pure black/white at 50% gray.
// Returns a Gray image where 0 = black, 255 = white.
func toBlackWhite(src image.Image) *image.Gray {
	b := src.Bounds()
	g := image.NewGray(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, gg, bb, _ := src.At(x, y).RGBA()
			lum := (299*r + 587*gg + 114*bb) / 1000 / 257
			if lum < 128 {
				g.SetGray(x, y, color.Gray{Y: 0})
			} else {
				g.SetGray(x, y, color.Gray{Y: 255})
			}
		}
	}
	return g
}
