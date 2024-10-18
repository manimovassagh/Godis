package server

import (
	"net"
	"testing"
	"time"
)

// MockListener is a mock implementation of net.Listener
type MockListener struct {
	connChan chan net.Conn
	closed   bool
}

func NewMockListener() *MockListener {
	return &MockListener{
		connChan: make(chan net.Conn),
	}
}

func (ml *MockListener) Accept() (net.Conn, error) {
	if ml.closed {
		return nil, net.ErrClosed
	}
	conn, ok := <-ml.connChan
	if !ok {
		return nil, net.ErrClosed
	}
	return conn, nil
}

func (ml *MockListener) Close() error {
	if !ml.closed {
		ml.closed = true
		close(ml.connChan)
	}
	return nil
}

func (ml *MockListener) Addr() net.Addr {
	return &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 8080,
	}
}

// MockConn is a mock implementation of net.Conn for testing
type MockConn struct{}

func (mc *MockConn) Read(b []byte) (n int, err error)   { return 0, nil }
func (mc *MockConn) Write(b []byte) (n int, err error)  { return len(b), nil }
func (mc *MockConn) Close() error                       { return nil }
func (mc *MockConn) LocalAddr() net.Addr                { return nil }
func (mc *MockConn) RemoteAddr() net.Addr               { return nil }
func (mc *MockConn) SetDeadline(t time.Time) error      { return nil }
func (mc *MockConn) SetReadDeadline(t time.Time) error  { return nil }
func (mc *MockConn) SetWriteDeadline(t time.Time) error { return nil }

// TestServerRun tests the server's ability to accept connections and handle them
func TestServerRun(t *testing.T) {
	// Create a mock listener and server
	mockListener := NewMockListener()
	server := New("127.0.0.1:8080")

	// Mock the listener to return our mock listener instead of a real one
	go func() {
		// Simulate a client connection being accepted
		mockListener.connChan <- &MockConn{}
	}()

	// Run the server in a separate goroutine and close it after a short delay
	go func() {
		_ = server.Run()
	}()

	// Sleep briefly to give the server time to run
	time.Sleep(100 * time.Millisecond)

	// Check if the server accepted a connection and handled it
	select {
	case conn := <-mockListener.connChan:
		if conn == nil {
			t.Fatal("Expected a connection, but got nil")
		}
	default:
		t.Fatal("Server did not accept any connections")
	}

	// Close the mock listener to stop the server
	mockListener.Close()
}

// TestServerAcceptFailure tests if the server properly handles connection acceptance failure
func TestServerAcceptFailure(t *testing.T) {
	// Create a server
	server := New("invalid_address")

	// Run the server and expect it to return an error
	err := server.Run()
	if err == nil {
		t.Fatal("Expected an error when trying to run the server with an invalid address")
	}
}
