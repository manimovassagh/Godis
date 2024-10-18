
# Godis

A simple, lightweight Redis clone written in Go, featuring an in-memory data store with append-only file (AOF) persistence and a custom command-line interface (CLI).

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Architecture](#architecture)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
  - [Starting the Server](#starting-the-server)
  - [Using the CLI](#using-the-cli)
- [Supported Commands](#supported-commands)
- [Project Structure](#project-structure)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)
- [Acknowledgments](#acknowledgments)

## Introduction

**Godis** is a Redis-inspired in-memory data store built from scratch in Go. It supports basic Redis commands and features an append-only file (AOF) persistence mechanism to ensure data durability across restarts. Additionally, Godis includes a custom CLI for interacting with the server without relying on external tools.

## Features

- **In-memory key-value data store**
- **Append-only file (AOF) persistence**
- **RESP (REdis Serialization Protocol) implementation**
- **Custom Godis CLI for server interaction**
- **Supports basic Redis commands**: `SET`, `GET`, `PING`, `ECHO`
- **Thread-safe operations** using Goroutines and Mutexes
- **Modular codebase** with clear package separation

## Architecture

Godis is structured into several packages to promote modularity and maintainability:

- **cmd**: Contains the entry points for the server and CLI.
- **internal/aof**: Manages the append-only file persistence.
- **internal/commands**: Handles client connections and command execution.
- **internal/datastore**: Implements the in-memory data store.
- **internal/protocol**: Parses and constructs RESP messages.
- **internal/server**: Contains the TCP server logic.

## Getting Started

### Prerequisites

- **Go 1.23** or higher installed on your machine
- **Git** for cloning the repository (optional)

### Installation

1. **Clone the repository**

   ```bash
   git clone https://github.com/manimovassagh/Godis.git
   ```

2. **Navigate to the project directory**

   ```bash
   cd Godis
   ```

3. **Build the server**

   ```bash
   go build -o godis-server ./cmd/server
   ```

4. **Build the CLI**

   ```bash
   go build -o godis-cli ./cmd/client
   ```

## Usage

### Starting the Server

Run the Godis server by executing:

```bash
./godis-server
```

The server will start listening on port `6379`.

### Using the CLI

In a new terminal window, start the Godis CLI:

```bash
./godis-cli
```

You should see:

```
Godis CLI connected to localhost:6379
Type 'exit' or 'quit' to close the CLI.
godis>
```

You can now enter commands to interact with the server.

## Supported Commands

- **PING**

  ```bash
  godis> PING
  PONG
  ```

- **ECHO**

  ```bash
  godis> ECHO "Hello, Godis!"
  Hello, Godis!
  ```

- **SET**

  ```bash
  godis> SET mykey "Some value"
  OK
  ```

- **GET**

  ```bash
  godis> GET mykey
  Some value
  ```

- **EXIT / QUIT**

  ```bash
  godis> EXIT
  Bye!
  ```

## Project Structure

```
Godis/
├── cmd/
│   ├── server/
│   │   └── main.go          // Server entry point
│   └── client/
│       └── main.go          // CLI entry point
├── internal/
│   ├── aof/
│   │   └── aof.go           // AOF persistence
│   ├── commands/
│   │   └── commands.go      // Command handling
│   ├── datastore/
│   │   └── datastore.go     // In-memory data store
│   ├── protocol/
│   │   └── protocol.go      // RESP implementation
│   └── server/
│       └── server.go        // TCP server logic
├── go.mod                   // Go module file
└── README.md                // Project documentation
```

## Roadmap

- **Additional Commands**: Implement more Redis commands such as `DEL`, `INCR`, `EXISTS`.
- **Data Structures**: Add support for lists, sets, hashes, and sorted sets.
- **Expiration**: Implement key expiration and TTL functionality.
- **Persistence Enhancements**: Introduce snapshotting (RDB files) and AOF rewriting.
- **Configuration**: Allow server settings via configuration files or command-line flags.
- **Improved CLI**: Enhance the CLI with command history, auto-completion, and syntax highlighting.
- **Testing**: Develop comprehensive unit and integration tests for all components.
- **Logging**: Implement structured logging for better observability.

## Contributing

Contributions are welcome! Please follow these steps:

1. **Fork the repository**

2. **Create a new feature branch**

   ```bash
   git checkout -b feature/my-feature
   ```

3. **Commit your changes**

   ```bash
   git commit -am 'Add new feature'
   ```

4. **Push to the branch**

   ```bash
   git push origin feature/my-feature
   ```

5. **Open a pull request**

Please ensure your code adheres to the project's coding standards and includes appropriate tests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Redis](https://redis.io/) for the inspiration behind this project.
- The Go community for providing excellent tools and documentation.
- Everyone who contributes to open-source projects and promotes knowledge sharing.

---

**Note**: This project is for educational purposes to understand how key-value stores and network servers work. It is not intended for production use.
