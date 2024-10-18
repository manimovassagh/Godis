package commands_test

import (
	"bufio"
	"net"
	"strings"
	"testing"

	"github.com/manimovassagh/Godis/internal/aof"
	"github.com/manimovassagh/Godis/internal/commands"
	"github.com/manimovassagh/Godis/internal/datastore"
	"github.com/manimovassagh/Godis/internal/protocol"
)

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

// Test SET and GET commands
func TestSetAndGetCommand(t *testing.T) {
	server, client := net.Ppipe()
	defer server.Close()
	defer client.Close()

	// Set up datastore and AOF mocks
	datastore := datastore.NewDataStore()
	aofHandler := aof.NewAOFHandler("test.aof")

	go func() {
		connClient := commands.NewClient(server)
		connClient.Handle()
	}()

	clientWriter := bufio.NewWriter(client)
	clientReader := bufio.NewReader(client)

	// Send SET command
	clientWriter.WriteString("*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n")
	clientWriter.Flush()

	// Read SET response
	response, err := protocol.ReadResponse(clientReader)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}
	if strings.TrimSpace(response) != "OK" {
		t.Errorf("Expected OK, got %s", response)
	}

	// Send GET command
	clientWriter.WriteString("*2\r\n$3\r\nGET\r\n$3\r\nkey\r\n")
	clientWriter.Flush()

	// Read GET response
	response, err = protocol.ReadResponse(clientReader)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}
	if strings.TrimSpace(response) != "value" {
		t.Errorf("Expected value, got %s", response)
	}
}

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