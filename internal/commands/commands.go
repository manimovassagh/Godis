package commands

import (
	"bufio"
	"log"
	"net"
	"strings"

	"github.com/manimovassagh/Godis/internal/aof"
	"github.com/manimovassagh/Godis/internal/datastore"
	"github.com/manimovassagh/Godis/internal/protocol"
)

type Client struct {
	conn      net.Conn
	reader    *bufio.Reader
	datastore *datastore.DataStore
	aof       *aof.AOFHandler
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		conn:      conn,
		reader:    bufio.NewReader(conn),
		datastore: datastore.GetDataStore(),
		aof:       aof.GetAOFHandler(),
	}
}

func (c *Client) Handle() {
	defer c.conn.Close()
	clientAddr := c.conn.RemoteAddr().String()
	log.Printf("Client connected: %s", clientAddr)

	for {
		args, err := protocol.ParseRequest(c.reader)
		if err != nil {
			protocol.WriteError(c.conn, "ERR "+err.Error())
			return
		}
		if len(args) == 0 {
			protocol.WriteError(c.conn, "ERR empty command")
			continue
		}
		cmd := strings.ToUpper(args[0])
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
}

func (c *Client) ping(args []string) {
	var response string
	if len(args) > 1 {
		response = args[1]
	} else {
		response = "PONG"
	}
	protocol.WriteSimpleString(c.conn, response)
}

func (c *Client) echo(args []string) {
	if len(args) != 2 {
		protocol.WriteError(c.conn, "ERR wrong number of arguments for 'ECHO' command")
		return
	}
	protocol.WriteBulkString(c.conn, args[1])
}

func (c *Client) set(args []string) {
	if len(args) != 3 {
		protocol.WriteError(c.conn, "ERR wrong number of arguments for 'SET' command")
		return
	}
	key, value := args[1], args[2]
	c.datastore.Set(key, value)
	c.aof.AppendCommand(args)
	protocol.WriteSimpleString(c.conn, "OK")
}

func (c *Client) get(args []string) {
	if len(args) != 2 {
		protocol.WriteError(c.conn, "ERR wrong number of arguments for 'GET' command")
		return
	}
	key := args[1]
	value, found := c.datastore.Get(key)
	if !found {
		protocol.WriteNullBulkString(c.conn)
	} else {
		protocol.WriteBulkString(c.conn, value)
	}
}
