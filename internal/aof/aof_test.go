package aof

import (
	"os"
	"testing"

	"github.com/manimovassagh/Godis/internal/protocol"
)

// TestGetAOFHandler ensures that GetAOFHandler returns a singleton instance
func TestGetAOFHandler(t *testing.T) {
	handler1 := GetAOFHandler()
	handler2 := GetAOFHandler()

	if handler1 != handler2 {
		t.Errorf("Expected singleton instance, but got different instances")
	}
}

// TestAppendCommand tests if commands are correctly appended to the AOF file
func TestAppendCommand(t *testing.T) {
	// Create a temporary AOF file
	tmpFile, err := os.CreateTemp("", "appendonly.aof")
	if err != nil {
		t.Fatalf("Failed to create temporary AOF file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up the file after the test

	// Replace the real AOF file with the temporary file
	handler := &AOFHandler{file: tmpFile}

	// Simulate a command and append it
	args := []string{"SET", "key", "value"}
	handler.AppendCommand(args)

	// Close the file to ensure all data is written
	tmpFile.Close()

	// Reopen the file for reading
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read AOF file: %v", err)
	}

	// Check if the command was correctly written to the file
	expected := protocol.FormatCommand(args)
	if string(content) != expected {
		t.Errorf("Expected %q, but got %q", expected, string(content))
	}
}
