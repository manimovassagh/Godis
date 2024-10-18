package protocol

import (
	"bytes"
	"net"
	"testing"
	"time"
)

// TestFormatCommand tests if commands are correctly formatted into RESP protocol
func TestFormatCommand(t *testing.T) {
	args := []string{"SET", "key", "value"}
	expected := "*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"
	result := FormatCommand(args)

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestWriteSimpleString tests if a simple string is written correctly
func TestWriteSimpleString(t *testing.T) {
	var buf bytes.Buffer
	conn := &mockConn{&buf}

	WriteSimpleString(conn, "PONG")

	expected := "+PONG\r\n"
	result := buf.String()

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestWriteError tests if an error is written correctly
func TestWriteError(t *testing.T) {
	var buf bytes.Buffer
	conn := &mockConn{&buf}

	WriteError(conn, "ERR unknown command")

	expected := "-ERR unknown command\r\n"
	result := buf.String()

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// TestWriteBulkString tests if a bulk string is written correctly
func TestWriteBulkString(t *testing.T) {
	var buf bytes.Buffer
	conn := &mockConn{&buf}

	WriteBulkString(conn, "hello")

	expected := "$5\r\nhello\r\n"
	result := buf.String()

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// mockConn is a mock implementation of net.Conn for testing
type mockConn struct {
	writer *bytes.Buffer
}

func (m *mockConn) Read(b []byte) (n int, err error)   { return 0, nil }
func (m *mockConn) Write(b []byte) (n int, err error)  { return m.writer.Write(b) }
func (m *mockConn) Close() error                       { return nil }
func (m *mockConn) LocalAddr() net.Addr                { return nil }
func (m *mockConn) RemoteAddr() net.Addr               { return nil }
func (m *mockConn) SetDeadline(t time.Time) error      { return nil }
func (m *mockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *mockConn) SetWriteDeadline(t time.Time) error { return nil }
