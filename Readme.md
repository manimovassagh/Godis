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

- In-memory key-value data store
- Append-only file (AOF) persistence
- RESP (REdis Serialization Protocol) implementation
- Custom Godis CLI for server interaction
- Supports basic Redis commands: `SET`, `GET`, `PING`, `ECHO`
- Thread-safe operations using Goroutines and Mutexes
- Modular codebase with clear package separation

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