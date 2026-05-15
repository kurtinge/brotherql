package brotherql

import (
	"bytes"
	"image/png"
	"os"
	"testing"
)

func TestBuildPrintJobMatchesGolden(t *testing.T) {
	m, ok := findModel(0x04F9, 0x2042)
	if !ok {
		t.Fatal("QL-700 not in supportedModels")
	}
	assertPrintJobMatchesGolden(t, m, "testdata/print-job.QL-700.golden.bin")
}

func TestBuildPrintJobMatchesGoldenQL710W(t *testing.T) {
	m, ok := findModel(0x04F9, 0x2043)
	if !ok {
		t.Fatal("QL-710W not in supportedModels")
	}
	assertPrintJobMatchesGolden(t, m, "testdata/print-job.QL-710W.golden.bin")
}

func assertPrintJobMatchesGolden(t *testing.T, model modelInfo, goldenPath string) {
	t.Helper()
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
	}, model)

	want, err := os.ReadFile(goldenPath)
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
	out := buildPrintJob(raster, Label62, PrintOptions{Label: Label62}, modelInfo{})
	for i := 0; i < 200; i++ {
		if out[i] != 0 {
			t.Errorf("byte %d = 0x%02x, want 0x00 (invalidate)", i, out[i])
			return
		}
	}
}

func TestBuildPrintJobModeSwitchPrefix(t *testing.T) {
	raster := make([]byte, 90)
	out := buildPrintJob(raster, Label62, PrintOptions{Label: Label62}, modelInfo{NeedsModeSwitch: true})
	want := []byte{0x1B, 0x69, 0x61, 0x01}
	if !bytes.Equal(out[:4], want) {
		t.Errorf("first 4 bytes = % x, want % x (mode switch prefix)", out[:4], want)
	}
}
