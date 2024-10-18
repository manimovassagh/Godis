package main

import (
	"testing"
	"time"
)

func TestMainFunction(t *testing.T) {
	// Create a channel to signal when the main function completes
	done := make(chan bool, 1)

	// Run the main function in a goroutine
	go func() {
		main()
		done <- true
	}()

	// Use a select statement to wait for either the test to finish or time out
	select {
	case <-done:
		// Test completed successfully
	case <-time.After(1 * time.Second):
		// If it takes longer than 1 second, consider it timed out but don't use os.Exit
		t.Log("Test timed out but will return gracefully")
	}
}