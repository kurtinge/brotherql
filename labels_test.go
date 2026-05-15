package brotherql

import "testing"

func TestLabel62Constants(t *testing.T) {
	if Label62.Name != "62" {
		t.Errorf("Label62.Name = %q, want %q", Label62.Name, "62")
	}
	if Label62.WidthMM != 62 {
		t.Errorf("Label62.WidthMM = %d, want 62", Label62.WidthMM)
	}
	if Label62.HeightMM != 0 {
		t.Errorf("Label62.HeightMM = %d, want 0 (continuous)", Label62.HeightMM)
	}
	if Label62.WidthPx != 696 {
		t.Errorf("Label62.WidthPx = %d, want 696", Label62.WidthPx)
	}
}

func TestLabel62x29Constants(t *testing.T) {
	if Label62x29.WidthMM != 62 {
		t.Errorf("Label62x29.WidthMM = %d, want 62", Label62x29.WidthMM)
	}
	if Label62x29.HeightMM != 29 {
		t.Errorf("Label62x29.HeightMM = %d, want 29", Label62x29.HeightMM)
	}
}
