package server

import (
	"log"
	"net"

	"github.com/manimovassagh/Godis/internal/commands"
)

type Server struct {
	address string
}

// New returns a new Server instance that will listen on the given address.
func New(address string) *Server {
	return &Server{
		address: address,
	}
}

// Run starts the TCP server and begins listening for incoming connections.
// When a connection is established, it creates a new Client and runs it in a
// goroutine.
func (s *Server) Run() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}
	defer listener.Close()
	log.Printf("Listening on %s...", s.address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		client := commands.NewClient(conn)
		go client.Handle()
	}
}
