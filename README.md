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

## License

MIT License. See [LICENSE](LICENSE) for details. 