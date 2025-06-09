# k8s-controller-tutorial

A starter template for building Kubernetes controllers or CLI tools in Go using [cobra-cli](https://github.com/spf13/cobra-cli).

## Prerequisites

- [Go](https://golang.org/dl/) 1.18 or newer
- [cobra-cli](https://github.com/spf13/cobra-cli) installed:
  ```sh
  go install github.com/spf13/cobra-cli@latest
  ```

## Getting Started

1. **Clone this repository and initialize the Go module:**
   ```sh
   git clone https://github.com/yourusername/k8s-controller-tutorial.git
   cd k8s-controller-tutorial
   go mod init github.com/yourusername/k8s-controller-tutorial
   ```
   
   Make sure the LICENSE file is set to MIT and updated with your name and year (e.g., 2025 Denys Vasyliev).

2. **Add zerolog for structured logging:**
   Integrate [zerolog](https://github.com/rs/zerolog) to enable structured logging with support for log levels: info, debug, trace, warn, and error.
   
   Install zerolog:
   ```sh
   go get github.com/rs/zerolog/log
   ```
   
   Example usage in your code:
   ```go
   import (
       "github.com/rs/zerolog/log"
   )

   func main() {
       log.Info().Msg("info message")
       log.Debug().Msg("debug message")
       log.Trace().Msg("trace message")
       log.Warn().Msg("warn message")
       log.Error().Msg("error message")
   }
   ```
   
   You can configure the log level and output format as needed. See the [zerolog documentation](https://github.com/rs/zerolog) for advanced configuration.

3. **Build your CLI:**
   ```sh
   go build -o controller
   ```

4. **Run your CLI (shows help by default):**
   ```sh
   ./controller --help
   ```

## Project Structure

- `cmd/` — Contains your CLI commands.
- `main.go` — Entry point for your application.

## License

MIT License. See [LICENSE](LICENSE) for details. 