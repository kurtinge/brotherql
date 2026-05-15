package brotherql

import (
	"bytes"
	"testing"
)

func TestMockTransportWriteRead(t *testing.T) {
	mock := &mockTransport{}
	mock.responses = [][]byte{{0x01, 0x02, 0x03}}

	n, err := mock.Write([]byte{0xAA, 0xBB})
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Errorf("Write returned %d, want 2", n)
	}
	if !bytes.Equal(mock.written, []byte{0xAA, 0xBB}) {
		t.Errorf("written = %v, want [AA BB]", mock.written)
	}

	buf := make([]byte, 32)
	n, err = mock.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	if n != 3 {
		t.Errorf("Read returned %d, want 3", n)
	}
	if !bytes.Equal(buf[:n], []byte{0x01, 0x02, 0x03}) {
		t.Errorf("read = %v, want [01 02 03]", buf[:n])
	}
}

func TestMockTransportClose(t *testing.T) {
	mock := &mockTransport{}
	if err := mock.Close(); err != nil {
		t.Errorf("Close returned %v, want nil", err)
	}
	if !mock.closed {
		t.Errorf("closed = false, want true after Close()")
	}
}
