package brotherql

import (
	"errors"
	"testing"
)

func TestOpenBySerialEmptyError(t *testing.T) {
	_, err := OpenBySerial("")
	if !errors.Is(err, ErrInvalidSerial) {
		t.Errorf("err = %v, want ErrInvalidSerial", err)
	}
}

func TestInfoStringer(t *testing.T) {
	i := Info{Serial: "ABC123", Model: "QL-700", USBPath: "bus 2 addr 5"}
	got := i.String()
	want := "QL-700 serial=ABC123 (bus 2 addr 5)"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
