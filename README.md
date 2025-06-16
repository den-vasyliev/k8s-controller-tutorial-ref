# k8s-controller-tutorial

A starter template for building Kubernetes controllers or CLI tools in Go using [cobra-cli](https://github.com/spf13/cobra-cli).

## Prerequisites

- [Go](https://golang.org/dl/) 1.24 or newer
- [cobra-cli](https://github.com/spf13/cobra-cli) installed:
  ```sh
  go install github.com/spf13/cobra-cli@latest
  ```

## Getting Started

1. **Clone this repository:**
   ```sh
   git clone https://github.com/yourusername/k8s-controller-tutorial.git
   cd k8s-controller-tutorial
   ```

2. **Initialize Go module (if not already):**
   ```sh
   go mod init github.com/yourusername/k8s-controller-tutorial
   ```

3. **Initialize Cobra:**
   ```sh
   cobra-cli init
   ```

4. **Build your CLI:**
   ```sh
   go build -o controller
   ```

5. **Run your CLI (shows help by default):**
   ```sh
   ./controller --help
   ```

## Project Structure

- `cmd/` — Contains your CLI commands.
- `main.go` — Entry point for your application.
- `cmd/go_basic.go`: Implements the command and struct logic
- `cmd/go_basic_test.go`: Unit tests for the struct methods 

This directory contains the `go_basic.go` file, which demonstrates basic usage of Go structs and methods within a Cobra CLI command.

## go_basic.go Overview
- Defines a `Kubernetes` struct with fields for name, version, users, and node number.
- Implements methods to print users and add a new user.
- Registers a `go-basic` Cobra command that
  - Initializes a sample `Kubernetes` struct
  - Prints the list of users
  - Adds a new user
  - Prints the updated list of users

## Usage

To run the `go-basic` command:

```sh
# From the project root
go run main.go go-basic
```

You should see output listing the initial users, then the updated list after adding a new user.

## Testing

Unit tests for the `Kubernetes` struct are provided in `go_basic_test.go`.
To run the tests:

```sh
go test ./cmd
```

## License

MIT License. See [LICENSE](LICENSE) for details. 