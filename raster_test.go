package brotherql

import (
	"bytes"
	"image/png"
	"os"
	"testing"
)

func TestEncodeRasterMatchesGolden(t *testing.T) {
	f, err := os.Open("testdata/text-label.png")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f.Close() }()
	img, err := png.Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	got, err := encodeRaster(img, Label62)
	if err != nil {
		t.Fatalf("encodeRaster: %v", err)
	}

	want, err := os.ReadFile("testdata/text-label.golden.bin")
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(got, want) {
		t.Errorf("raster bytes mismatch: got %d bytes, want %d bytes", len(got), len(want))
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

func TestEncodeRasterRowCount(t *testing.T) {
	f, err := os.Open("testdata/text-label.png")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f.Close() }()
	img, err := png.Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	got, err := encodeRaster(img, Label62)
	if err != nil {
		t.Fatalf("encodeRaster: %v", err)
	}
	// 100 rows × 90 bytes per row = 9000
	if len(got) != 9000 {
		t.Errorf("len(raster) = %d, want 9000", len(got))
	}
}

func TestEncodeRaster12mmMatchesGolden(t *testing.T) {
	f, err := os.Open("testdata/text-label.12mm.png")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f.Close() }()
	img, err := png.Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	got, err := encodeRaster(img, Label12)
	if err != nil {
		t.Fatalf("encodeRaster: %v", err)
	}

	want, err := os.ReadFile("testdata/text-label.12mm.golden.bin")
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(got, want) {
		t.Errorf("raster bytes mismatch: got %d bytes, want %d bytes", len(got), len(want))
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
