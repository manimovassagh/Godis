package server

import (
	"log"
	"net"

	"github.com/manimovassagh/Godis/internal/commands"
)

type Server struct {
	address string
}

func New(address string) *Server {
	return &Server{
		address: address,
	}
}

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
