package main

import (
	"testing"
)

// TestParseInput tests the parseInput function for correctness
func TestParseInput(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"SET key value", []string{"SET", "key", "value"}},
		{"GET \"key with spaces\"", []string{"GET", "key with spaces"}},
		{"ECHO \"Hello, World!\"", []string{"ECHO", "Hello, World!"}},
		{"PING", []string{"PING"}},
		{"\"quoted\" \"multiple words\"", []string{"quoted", "multiple words"}},
	}

	for _, test := range tests {
		args := parseInput(test.input)
		if len(args) != len(test.expected) {
			t.Errorf("Expected %d args, but got %d", len(test.expected), len(args))
		}
		for i, arg := range args {
			if arg != test.expected[i] {
				t.Errorf("Expected arg %q, but got %q", test.expected[i], arg)
			}
		}
	}
}
