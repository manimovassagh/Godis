package main

import (
	"log"

	"github.com/manimovassagh/Godis/internal/aof"
	"github.com/manimovassagh/Godis/internal/server"
)

func main() {
	// Initialize AOF handler and load existing data
	aofHandler := aof.GetAOFHandler()
	if err := aofHandler.LoadCommands(); err != nil {
		log.Fatalf("Failed to load AOF file: %v", err)
	}

	// Start the server
	srv := server.New(":6379")
	log.Println("Server is starting on port 6379...")
	if err := srv.Run(); err != nil {
		log.Fatalf("Server stopped with error: %v", err)
	}
}
