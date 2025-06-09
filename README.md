# **FastHTTP Server Command**

   - Added a new `server` command using [fasthttp](https://github.com/valyala/fasthttp).
   - The command starts a FastHTTP server with a configurable port (default: 8080).
   - Supports the `--log-level` flag for controlling log verbosity.
   - Uses zerolog for logging.

   **Usage:**
   ```sh
   git switch feature/step4-fasthttp-server

   go run main.go server --port 8080 --log-level debug
   ```

   **What it does:**
   - Starts a FastHTTP server on the specified port.
   - Responds with "Hello from FastHTTP!" to any request.
   - Respects the log level set by the `--log-level` flag.

## Project Structure

- `cmd/` — Contains your CLI commands.
- `main.go` — Entry point for your application.
- `server.go` - fasthttp server

## License

MIT License. See [LICENSE](LICENSE) for details.
