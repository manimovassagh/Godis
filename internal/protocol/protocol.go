package protocol

import (
    "bufio"
    "errors"
    "fmt"
    "io"
    "net"
    "strconv"
    "strings"
)

// ParseRequest parses a client request from the connection
func ParseRequest(reader *bufio.Reader) ([]string, error) {
    line, err := reader.ReadString('\n')
    if err != nil {
        return nil, err
    }
    line = strings.TrimSpace(line)
    if len(line) == 0 || line[0] != '*' {
        return nil, errors.New("protocol error: expected '*'")
    }
    numArgs, err := strconv.Atoi(line[1:])
    if err != nil {
        return nil, errors.New("protocol error: invalid array length")
    }
    args := make([]string, numArgs)
    for i := 0; i < numArgs; i++ {
        line, err = reader.ReadString('\n')
        if err != nil {
            return nil, err
        }
        if len(line) == 0 || line[0] != '$' {
            return nil, errors.New("protocol error: expected '$'")
        }
        argLen, err := strconv.Atoi(strings.TrimSpace(line[1:]))
        if err != nil {
            return nil, errors.New("protocol error: invalid bulk string length")
        }
        arg := make([]byte, argLen+2) // +2 for \r\n
        _, err = io.ReadFull(reader, arg)
        if err != nil {
            return nil, err
        }
        args[i] = string(arg[:argLen])
    }
    return args, nil
}

// WriteSimpleString writes a simple string response to the client
func WriteSimpleString(conn net.Conn, message string) {
    fmt.Fprintf(conn, "+%s\r\n", message)
}

// WriteError writes an error response to the client
func WriteError(conn net.Conn, message string) {
    fmt.Fprintf(conn, "-%s\r\n", message)
}

// WriteBulkString writes a bulk string response to the client
func WriteBulkString(conn net.Conn, message string) {
    fmt.Fprintf(conn, "$%d\r\n%s\r\n", len(message), message)
}

// WriteNullBulkString writes a null bulk string response to the client
func WriteNullBulkString(conn net.Conn) {
    fmt.Fprint(conn, "$-1\r\n")
}

// FormatCommand formats a command for sending to the server
func FormatCommand(args []string) string {
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("*%d\r\n", len(args)))
    for _, arg := range args {
        sb.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg))
    }
    return sb.String()
}

// WriteCommand sends a command to the server using the RESP protocol
func WriteCommand(conn net.Conn, args []string) error {
    cmd := FormatCommand(args)
    _, err := conn.Write([]byte(cmd))
    return err
}

// ReadResponse reads and parses the server's response
func ReadResponse(reader *bufio.Reader) (string, error) {
    line, err := reader.ReadString('\n')
    if err != nil {
        return "", err
    }

    line = strings.TrimSpace(line)
    if len(line) == 0 {
        return "", errors.New("empty response")
    }

    switch line[0] {
    case '+':
        // Simple string
        return line[1:], nil
    case '-':
        // Error
        return "(error) " + line[1:], nil
    case ':':
        // Integer
        return line[1:], nil
    case '$':
        // Bulk string
        length, err := strconv.Atoi(line[1:])
        if err != nil {
            return "", errors.New("invalid bulk string length")
        }
        if length == -1 {
            return "(nil)", nil
        }
        buf := make([]byte, length+2) // +2 for \r\n
        _, err = io.ReadFull(reader, buf)
        if err != nil {
            return "", err
        }
        return string(buf[:length]), nil
    default:
        return "", errors.New("unknown response type")
    }
}