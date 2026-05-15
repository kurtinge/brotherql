package brotherql

import (
	"bytes"
	"image/png"
	"os"
	"testing"
)

func TestPrinterPrintSendsCommandSequence(t *testing.T) {
	mock := &mockTransport{}
	statusReady := make([]byte, 32)
	statusReady[10] = 62
	mock.responses = [][]byte{statusReady, statusReady}

	p := &Printer{tr: mock}

	f, err := os.Open("testdata/text-label.png")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f.Close() }()
	img, err := png.Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	err = p.Print(img, PrintOptions{Label: Label62, AutoCut: true})
	if err != nil {
		t.Fatalf("Print: %v", err)
	}

	if len(mock.written) < 200 {
		t.Fatalf("written too short: %d", len(mock.written))
	}
	for i := 0; i < 200; i++ {
		if mock.written[i] != 0 {
			t.Errorf("byte %d = 0x%02x, want 0x00", i, mock.written[i])
			return
		}
	}
	last := mock.written[len(mock.written)-1]
	if last != 0x1A {
		t.Errorf("last byte = 0x%02x, want 0x1A (print+cut)", last)
	}
}

func TestPrinterStatus(t *testing.T) {
	mock := &mockTransport{}
	raw := make([]byte, 32)
	raw[10] = 62
	mock.responses = [][]byte{raw}

	p := &Printer{tr: mock}
	s, err := p.Status()
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if s.MediaWidth != 62 {
		t.Errorf("MediaWidth = %d, want 62", s.MediaWidth)
	}
	if !s.Ready {
		t.Errorf("Ready = false, want true")
	}
	want := []byte{0x1B, 0x69, 0x53}
	if !bytes.Equal(mock.written, want) {
		t.Errorf("written = %v, want %v", mock.written, want)
	}
}

func TestPrinterClose(t *testing.T) {
	mock := &mockTransport{}
	p := &Printer{tr: mock}
	if err := p.Close(); err != nil {
		t.Errorf("Close: %v", err)
	}
	if !mock.closed {
		t.Errorf("transport not closed")
	}
}
