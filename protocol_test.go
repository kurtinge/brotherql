package brotherql

import (
	"bytes"
	"image/png"
	"os"
	"testing"
)

func TestBuildPrintJobMatchesGolden(t *testing.T) {
	f, err := os.Open("testdata/text-label.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	raster, err := encodeRaster(img, Label62)
	if err != nil {
		t.Fatal(err)
	}

	got := buildPrintJob(raster, Label62, PrintOptions{
		Label:   Label62,
		Copies:  1,
		AutoCut: true,
	})

	want, err := os.ReadFile("testdata/print-job.golden.bin")
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(got, want) {
		t.Errorf("print job bytes mismatch: got %d bytes, want %d bytes", len(got), len(want))
		minLen := len(got)
		if len(want) < minLen {
			minLen = len(want)
		}
		for i := 0; i < minLen; i++ {
			if got[i] != want[i] {
				t.Errorf("first diff at offset %d: got 0x%02x, want 0x%02x", i, got[i], want[i])
				break
			}
		}
	}
}

func TestBuildPrintJobStartsWithInvalidate(t *testing.T) {
	raster := make([]byte, 90)
	out := buildPrintJob(raster, Label62, PrintOptions{Label: Label62})
	for i := 0; i < 200; i++ {
		if out[i] != 0 {
			t.Errorf("byte %d = 0x%02x, want 0x00 (invalidate)", i, out[i])
			return
		}
	}
}
