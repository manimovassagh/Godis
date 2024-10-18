package commands_test

import (
	"bufio"
	"net"
	"strings"
	"testing"

	"github.com/manimovassagh/Godis/internal/commands"
	"github.com/manimovassagh/Godis/internal/protocol"
)

// Mock DataStore for testing
type MockDataStore struct {
	data map[string]string
}

func NewMockDataStore() *MockDataStore {
	return &MockDataStore{data: make(map[string]string)}
}

func (m *MockDataStore) Set(key, value string) {
	m.data[key] = value
}

func (m *MockDataStore) Get(key string) (string, bool) {
	value, found := m.data[key]
	return value, found
}

// Mock AOFHandler (not needed for these tests, so it's a no-op)
type MockAOFHandler struct{}

func NewMockAOFHandler() *MockAOFHandler {
	return &MockAOFHandler{}
}

func (m *MockAOFHandler) AppendCommand(args []string) {
	// No-op
}

// Test PING command
func TestPingCommand(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	go func() {
		connClient := commands.NewClient(server)
		connClient.Handle()
	}()

	clientWriter := bufio.NewWriter(client)
	clientReader := bufio.NewReader(client)

	// Send PING command
	clientWriter.WriteString("*1\r\n$4\r\nPING\r\n")
	clientWriter.Flush()

	// Read response
	response, err := protocol.ReadResponse(clientReader)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	if strings.TrimSpace(response) != "PONG" {
		t.Errorf("Expected PONG, got %s", response)
	}
}

// Test ECHO command
func TestEchoCommand(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	go func() {
		connClient := commands.NewClient(server)
		connClient.Handle()
	}()

	clientWriter := bufio.NewWriter(client)
	clientReader := bufio.NewReader(client)

	// Send ECHO command
	clientWriter.WriteString("*2\r\n$4\r\nECHO\r\n$5\r\nHello\r\n")
	clientWriter.Flush()

	// Read response
	response, err := protocol.ReadResponse(clientReader)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	if strings.TrimSpace(response) != "Hello" {
		t.Errorf("Expected Hello, got %s", response)
	}
}

// Test SET and GET commands with mocked DataStore

// Test GET command for non-existent key
func TestGetNonExistentCommand(t *testing.T) {
	server, client := net.Pipe()
	defer server.Close()
	defer client.Close()

	go func() {
		connClient := commands.NewClient(server)
		connClient.Handle()
	}()

	clientWriter := bufio.NewWriter(client)
	clientReader := bufio.NewReader(client)

	// Send GET command for a non-existent key
	clientWriter.WriteString("*2\r\n$3\r\nGET\r\n$4\r\nnone\r\n")
	clientWriter.Flush()

	// Read response
	response, err := protocol.ReadResponse(clientReader)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	if strings.TrimSpace(response) != "(nil)" {
		t.Errorf("Expected (nil), got %s", response)
	}
}
