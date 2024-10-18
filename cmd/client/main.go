package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "os"
    "strings"

    "github.com/manimovassagh/Godis/internal/protocol"
)

func main() {
    // Connect to the Godis server
    conn, err := net.Dial("tcp", "localhost:6379")
    if err != nil {
        log.Fatalf("Failed to connect to Godis server: %v", err)
    }
    defer conn.Close()

    reader := bufio.NewReader(os.Stdin)
    serverReader := bufio.NewReader(conn)

    fmt.Println("Godis CLI connected to localhost:6379")
    fmt.Println("Type 'exit' or 'quit' to close the CLI.")

    for {
        fmt.Print("godis> ")
        input, err := reader.ReadString('\n')
        if err != nil {
            fmt.Printf("Error reading input: %v\n", err)
            continue
        }

        input = strings.TrimSpace(input)
        if input == "" {
            continue
        }

        if input == "exit" || input == "quit" {
            fmt.Println("Bye!")
            break
        }

        // Parse input into arguments
        args := parseInput(input)
        if len(args) == 0 {
            continue
        }

        // Send command to server
        err = protocol.WriteCommand(conn, args)
        if err != nil {
            fmt.Printf("Error sending command: %v\n", err)
            continue
        }

        // Read response from server
        response, err := protocol.ReadResponse(serverReader)
        if err != nil {
            fmt.Printf("Error reading response: %v\n", err)
            continue
        }

        // Print the response
        fmt.Println(response)
    }
}

// parseInput splits the input string into arguments, handling quotes
func parseInput(input string) []string {
    var args []string
    var current strings.Builder
    inQuotes := false

    for i := 0; i < len(input); i++ {
        c := input[i]
        switch c {
        case ' ':
            if inQuotes {
                current.WriteByte(c)
            } else if current.Len() > 0 {
                args = append(args, current.String())
                current.Reset()
            }
        case '"':
            inQuotes = !inQuotes
        default:
            current.WriteByte(c)
        }
    }

    if current.Len() > 0 {
        args = append(args, current.String())
    }

    return args
}