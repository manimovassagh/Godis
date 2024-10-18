package commands

import (
	"bufio"
	"bytes"
	"net"
	"testing"

	"github.com/manimovassagh/Godis/internal/datastore"
	"github.com/manimovassagh/Godis/internal/protocol"
)

// MockConn simulates an in-memory net.Conn
type MockConn struct {
	net.Conn
	readBuffer  *bytes.Buffer
	writeBuffer *bytes.Buffer
}

// NewMockConn creates a new mock connection
func NewMockConn() *MockConn {
	return &MockConn{
		readBuffer:  new(bytes.Buffer),
		writeBuffer: new(bytes.Buffer),
	}
}

// Write simulates writing to the connection
func (m *MockConn) Write(b []byte) (int, error) {
	return m.writeBuffer.Write(b)
}

// Read simulates reading from the connection
func (m *MockConn) Read(b []byte) (int, error) {
	return m.readBuffer.Read(b)
}

// SimulateInput simulates input from the client
func (m *MockConn) SimulateInput(input string) {
	m.readBuffer.WriteString(input)
}

// GetOutput gets the output written to the connection
func (m *MockConn) GetOutput() string {
	return m.writeBuffer.String()
}

// TestPingCommand tests the PING command
func TestPingCommand(t *testing.T) {
	client, mockConn := createMockClient()

	// Simulate client sending "PING" command in Redis protocol format
	mockConn.SimulateInput("*1\r\n$4\r\nPING\r\n")
	client.HandleOnce()

	expected := "+PONG\r\n"
	if mockConn.GetOutput() != expected {
		t.Errorf("Expected %q, got %q", expected, mockConn.GetOutput())
	}
}

// TestEchoCommand tests the ECHO command
func TestEchoCommand(t *testing.T) {
	client, mockConn := createMockClient()

	// Simulate client sending "ECHO hello" command in Redis protocol format
	mockConn.SimulateInput("*2\r\n$4\r\nECHO\r\n$5\r\nhello\r\n")
	client.HandleOnce()

	expected := "$5\r\nhello\r\n"
	if mockConn.GetOutput() != expected {
		t.Errorf("Expected %q, got %q", expected, mockConn.GetOutput())
	}
}

// TestSetGetCommand tests the SET and GET commands
// TestSetGetCommand tests the SET and GET commands
func TestSetGetCommand(t *testing.T) {
	client, mockConn := createMockClient()

	// Simulate client sending "SET key value" command in Redis protocol format
	mockConn.SimulateInput("*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n")
	client.HandleOnce()

	// Check the output for SET command
	expectedSet := "+OK\r\n"
	if mockConn.GetOutput() != expectedSet {
		t.Errorf("Expected %q for SET, got %q", expectedSet, mockConn.GetOutput())
	}

	// Clear the mock connection output buffer for the next command
	mockConn.writeBuffer.Reset()

	// Simulate client sending "GET key" command in Redis protocol format
	mockConn.SimulateInput("*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n")
	client.HandleOnce()

	// Check the output for GET command
	expectedGet := "$5\r\nvalue\r\n"
	if mockConn.GetOutput() != expectedGet {
		t.Errorf("Expected %q for GET, got %q", expectedGet, mockConn.GetOutput())
	}
}
// Helper function to create a mock client with in-memory connection
func createMockClient() (*Client, *MockConn) {
	mockConn := NewMockConn()
	client := NewClient(mockConn)
	client.reader = bufio.NewReader(mockConn)
	client.datastore = datastore.GetDataStore() // Ensure the real datastore is used
	return client, mockConn
}

// Add this method to handle just one command and return, avoiding infinite loop
func (c *Client) HandleOnce() {
	args, err := protocol.ParseRequest(c.reader)
	if err != nil {
		protocol.WriteError(c.conn, "ERR "+err.Error())
		return
	}
	if len(args) == 0 {
		protocol.WriteError(c.conn, "ERR empty command")
		return
	}
	cmd := args[0]
	switch cmd {
	case "PING":
		c.ping(args)
	case "ECHO":
		c.echo(args)
	case "SET":
		c.set(args)
	case "GET":
		c.get(args)
	default:
		protocol.WriteError(c.conn, "ERR unknown command '"+cmd+"'")
	}
}