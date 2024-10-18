package aof

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/manimovassagh/Godis/internal/datastore"
	"github.com/manimovassagh/Godis/internal/protocol"
)

type AOFHandler struct {
	file *os.File
	mu   sync.Mutex
}

var (
	instance *AOFHandler
	once     sync.Once
)

func GetAOFHandler() *AOFHandler {
	once.Do(func() {
		file, err := os.OpenFile("appendonly.aof", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		if err != nil {
			panic(fmt.Sprintf("Failed to open AOF file: %v", err))
		}
		instance = &AOFHandler{
			file: file,
		}
	})
	return instance
}

func (a *AOFHandler) AppendCommand(args []string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	cmd := protocol.FormatCommand(args)
	_, err := a.file.WriteString(cmd)
	if err != nil {
		fmt.Printf("Failed to write to AOF: %v\n", err)
	}
}

func (a *AOFHandler) LoadCommands() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	file, err := os.Open("appendonly.aof")
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No AOF file exists yet
		}
		return err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	datastore := datastore.GetDataStore()
	for {
		args, err := protocol.ParseRequest(reader)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if len(args) > 0 {
			cmd := strings.ToUpper(args[0])
			if cmd == "SET" && len(args) == 3 {
				datastore.Set(args[1], args[2])
			}
			// Implement other commands as needed
		}
	}
	return nil
}
