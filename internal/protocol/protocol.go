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

func WriteSimpleString(conn net.Conn, message string) {
    fmt.Fprintf(conn, "+%s\r\n", message)
}

func WriteError(conn net.Conn, message string) {
    fmt.Fprintf(conn, "-%s\r\n", message)
}

func WriteBulkString(conn net.Conn, message string) {
    fmt.Fprintf(conn, "$%d\r\n%s\r\n", len(message), message)
}

func WriteNullBulkString(conn net.Conn) {
    fmt.Fprint(conn, "$-1\r\n")
}

func FormatCommand(args []string) string {
    var sb strings.Builder
    sb.WriteString(fmt.Sprintf("*%d\r\n", len(args)))
    for _, arg := range args {
        sb.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg))
    }
    return sb.String()
}