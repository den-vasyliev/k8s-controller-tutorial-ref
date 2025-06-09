# Leader Election and Metrics for Controller Manager

- Added leader election support using a Lease resource (enabled by default, can be disabled with a flag).
- Added a flag to set the metrics port for the controller manager.
- Both features are configurable via CLI flags.

**New flags:**
- `--enable-leader-election` (default: true) — Enable/disable leader election for the controller manager.
- `--metrics-port` (default: 8081) — Port for controller manager metrics endpoint.

**What it does:**
- Ensures only one instance of the controller manager is active at a time (HA support).
- Exposes controller metrics on the specified port.

**Usage:**
```sh
git switch feature/step10-leader-election 

go run main.go server --enable-leader-election=false --metrics-port=9090
```
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

## License

MIT License. See [LICENSE](LICENSE) for details.