# **Add zerolog for structured logging:**

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