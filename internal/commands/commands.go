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

// NewClient returns a new Client instance that will handle the given connection.
//
// It initializes the Client with the given connection, a new bufio.Reader,
// the global DataStore, and the global AOFHandler.
func NewClient(conn net.Conn) *Client {
	return &Client{
		conn:      conn,
		reader:    bufio.NewReader(conn),
		datastore: datastore.GetDataStore(),
		aof:       aof.GetAOFHandler(),
	}
}

// Handle starts a loop that reads commands from the client and executes them.
// It uses the Client's reader to read requests from the client and the
// protocol package to parse the requests. It then executes the commands
// and writes the response back to the client. If the client disconnects,
// the loop exits and the connection is closed.
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

// ping handles the PING command for the client. 
// It takes an array of arguments and responds with the appropriate message ("PONG" if no argument provided).
func (c *Client) ping(args []string) {
	var response string
	if len(args) > 1 {
		response = args[1]
	} else {
		response = "PONG"
	}
	protocol.WriteSimpleString(c.conn, response)
}

// echo handles the ECHO command for the client.
// It takes an array of arguments and responds with the same string sent by the client.
func (c *Client) echo(args []string) {
	if len(args) != 2 {
		protocol.WriteError(c.conn, "ERR wrong number of arguments for 'ECHO' command")
		return
	}
	protocol.WriteBulkString(c.conn, args[1])
}

// set handles the SET command for the client.
// It takes an array of arguments with the following format: ["SET", key, value].
// It sets the given key-value pair in the in-memory data store and appends the command to the AOF file, then responds with "OK".
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

// get handles the GET command for the client.
// It takes an array of arguments with the following format: ["GET", key].
// It looks up the given key in the in-memory data store and responds with the
// associated value. If the key is not found, it responds with a null bulk string.
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
