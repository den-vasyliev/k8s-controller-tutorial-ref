# **Support log-level flags**

   Example usage in your `main.go`:
   ```go

       rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Set log level: trace, debug, info, warn, error")
       rootCmd.Execute()
   }
   ```

   Build your CLI:
   ```sh
   go build -o controller
   ```

   You can now run your CLI with different log levels:
   ```sh
   ./controller --log-level debug
   ./controller --log-level trace
   ```

## Project Structure

- `cmd/` — Contains your CLI commands.
- `main.go` — Entry point for your application.

## License

MIT License. See [LICENSE](LICENSE) for details. 