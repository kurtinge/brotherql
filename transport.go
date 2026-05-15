package brotherql

// transport is the minimal IO surface for talking to a printer.
// Implementations: usbTransport (real hardware), mockTransport (testing).
type transport interface {
	Write(p []byte) (int, error)
	Read(p []byte) (int, error)
	Close() error
}

// mockTransport is an in-memory transport for tests. It records all bytes
// written and serves successive reads from a queue of canned responses.
type mockTransport struct {
	written   []byte
	responses [][]byte
	readIdx   int
	closed    bool
}

func (m *mockTransport) Write(p []byte) (int, error) {
	m.written = append(m.written, p...)
	return len(p), nil
}

func (m *mockTransport) Read(p []byte) (int, error) {
	if m.readIdx >= len(m.responses) {
		return 0, nil
	}
	resp := m.responses[m.readIdx]
	m.readIdx++
	return copy(p, resp), nil
}

func (m *mockTransport) Close() error {
	m.closed = true
	return nil
}
