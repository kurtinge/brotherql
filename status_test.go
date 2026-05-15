package brotherql

import "testing"

func TestParseStatusReady(t *testing.T) {
	// 32-byte status response from a ready QL-700 with 62mm continuous tape.
	raw := []byte{
		0x80, 0x20, 0x42, 0x34, 0x39, 0x30, 0x00, 0x00,
		0x00, 0x00, 62, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
	s, err := parseStatus(raw)
	if err != nil {
		t.Fatalf("parseStatus error: %v", err)
	}
	if !s.Ready {
		t.Errorf("Ready = false, want true")
	}
	if s.Error != "" {
		t.Errorf("Error = %q, want empty", s.Error)
	}
	if s.MediaWidth != 62 {
		t.Errorf("MediaWidth = %d, want 62", s.MediaWidth)
	}
	if s.MediaLength != 0 {
		t.Errorf("MediaLength = %d, want 0", s.MediaLength)
	}
}

func TestParseStatusError(t *testing.T) {
	raw := make([]byte, 32)
	raw[0] = 0x80
	raw[1] = 0x20
	raw[8] = 0x01 // error info 1: no media
	s, err := parseStatus(raw)
	if err != nil {
		t.Fatalf("parseStatus error: %v", err)
	}
	if s.Ready {
		t.Errorf("Ready = true, want false")
	}
	if s.Error == "" {
		t.Errorf("Error empty, want non-empty")
	}
}

func TestParseStatusWrongLength(t *testing.T) {
	_, err := parseStatus(make([]byte, 16))
	if err == nil {
		t.Errorf("parseStatus expected error for short buffer")
	}
}
