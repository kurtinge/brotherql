package brotherql

import "testing"

func TestSupportedModels(t *testing.T) {
	cases := []struct {
		name            string
		vendorID        uint16
		productID       uint16
		needsModeSwitch bool
	}{
		{"QL-700", 0x04F9, 0x2042, false},
		{"QL-710W", 0x04F9, 0x2043, true},
	}
	for _, c := range cases {
		m, ok := findModel(c.vendorID, c.productID)
		if !ok {
			t.Errorf("findModel(0x%04X, 0x%04X): not found, want %s", c.vendorID, c.productID, c.name)
			continue
		}
		if m.Name != c.name {
			t.Errorf("findModel(0x%04X, 0x%04X).Name = %q, want %q", c.vendorID, c.productID, m.Name, c.name)
		}
		if m.NeedsModeSwitch != c.needsModeSwitch {
			t.Errorf("%s NeedsModeSwitch = %v, want %v", c.name, m.NeedsModeSwitch, c.needsModeSwitch)
		}
	}
}

func TestFindModelUnknown(t *testing.T) {
	if _, ok := findModel(0x0000, 0x0000); ok {
		t.Errorf("findModel(0,0): ok = true, want false")
	}
	if _, ok := findModel(0x04F9, 0xFFFF); ok {
		t.Errorf("findModel for unsupported Brother PID: ok = true, want false")
	}
}
