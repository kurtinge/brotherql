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
	if Label62.OffsetPx != 12 {
		t.Errorf("Label62.OffsetPx = %d, want 12", Label62.OffsetPx)
	}
}

func TestLabel62x29Constants(t *testing.T) {
	if Label62x29.WidthMM != 62 {
		t.Errorf("Label62x29.WidthMM = %d, want 62", Label62x29.WidthMM)
	}
	if Label62x29.HeightMM != 29 {
		t.Errorf("Label62x29.HeightMM = %d, want 29", Label62x29.HeightMM)
	}
	if Label62x29.OffsetPx != 12 {
		t.Errorf("Label62x29.OffsetPx = %d, want 12", Label62x29.OffsetPx)
	}
}

func TestLabel12Constants(t *testing.T) {
	if Label12.Name != "12" {
		t.Errorf("Label12.Name = %q, want %q", Label12.Name, "12")
	}
	if Label12.WidthMM != 12 {
		t.Errorf("Label12.WidthMM = %d, want 12", Label12.WidthMM)
	}
	if Label12.HeightMM != 0 {
		t.Errorf("Label12.HeightMM = %d, want 0 (continuous)", Label12.HeightMM)
	}
	if Label12.WidthPx != 106 {
		t.Errorf("Label12.WidthPx = %d, want 106", Label12.WidthPx)
	}
	if Label12.OffsetPx != 29 {
		t.Errorf("Label12.OffsetPx = %d, want 29", Label12.OffsetPx)
	}
}
