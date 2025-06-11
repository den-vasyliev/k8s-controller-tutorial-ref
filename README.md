# MCP Integration

- Integrated [MCP (Multi-Cluster Platform)](https://github.com/mark3labs/mcp-go) server into the project.
- MCP server can be enabled with the `--enable-mcp` flag and runs on a configurable port (default: 9090).
- MCP server runs alongside the FastHTTP API server and controller-runtime manager.
- Provides a real-time event stream and management interface for Kubernetes resources via the MCP protocol.

**Usage:**
```sh
git switch feature/step13-mcp-integration 

go run main.go server --log-level trace --enable-mcp --mcp-port 9090
# MCP server will be available on http://localhost:9090
```
- Use an MCP client or compatible tool to connect and invoke registered tools.

**What it does:**
- Enables external systems to interact with the controller via the MCP protocol (list/create FrontendPages, etc.).
- SSE mode provides real-time updates for tool execution.

---

## Project Structure

- `cmd/` — Contains your CLI commands.
- `main.go` — Entry point for your application.
- `server.go` - fasthttp server
- `Makefile` — Build automation tasks.
- `Dockerfile` — Distroless Dockerfile for secure containerization.
- `.github/workflows/` — GitHub Actions workflows for CI/CD.
- `list.go` - list cli command
- `charts/app` - helm chart
- `pkg/informer` - informer implementation
- `pkg/testutil` - envtest kit
- `pkg/ctrl` - controller implementation
- `config/crd` - CRD definition
- `pkg/apis` - CRD types and deepcopy
- `pkg/api` - API for PE integration
- `cmd/mcp.go` - MCP implementation

## License

MIT License. See [LICENSE](LICENSE) for details.
