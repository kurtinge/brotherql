//go:build hardware

package brotherql

import (
	"errors"
	"image"
	"image/color"
	"testing"
)

func TestHardwareList(t *testing.T) {
	infos, err := List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(infos) == 0 {
		t.Skip("no printer connected")
	}
	for _, i := range infos {
		t.Logf("found: %s", i)
	}
}

func TestHardwareStatus(t *testing.T) {
	p, err := Open()
	if errors.Is(err, ErrPrinterNotFound) {
		t.Skip("no printer connected")
	}
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer p.Close()

	s, err := p.Status()
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	t.Logf("status: %+v", s)
}

func TestHardwarePrint(t *testing.T) {
	p, err := Open()
	if errors.Is(err, ErrPrinterNotFound) {
		t.Skip("no printer connected")
	}
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer p.Close()

	img := image.NewRGBA(image.Rect(0, 0, 696, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 696; x++ {
			img.Set(x, y, color.White)
		}
	}
	for x := 0; x < 696; x++ {
		img.Set(x, 0, color.Black)
		img.Set(x, 99, color.Black)
	}
	for y := 0; y < 100; y++ {
		img.Set(0, y, color.Black)
		img.Set(695, y, color.Black)
	}

	if err := p.Print(img, PrintOptions{Label: Label62, AutoCut: true}); err != nil {
		t.Fatalf("Print: %v", err)
	}
	t.Log("printed test label — verify a bordered rectangle came out")
}
