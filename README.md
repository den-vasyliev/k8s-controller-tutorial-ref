# JWT Authentication

- Added JWT authentication middleware to protect all `/api/frontendpages` endpoints.
- Only requests with a valid JWT in the `Authorization: Bearer <token>` header are allowed.
- Added a `/api/token` endpoint for local testing, which issues a JWT for a test user (valid for 1 hour).
- Uses a hardcoded secret for development (update for production use).

## Running with JWT Secret

To start the server with a custom JWT secret, use the `--jwt-secret` flag:

```sh
go run main.go --log-level trace --kubeconfig  ~/.kube/config server --enable-leader-election=0 --jwt-secret "your-strong-secret"
```
This secret will be used for signing and validating JWT tokens. Make sure to use a strong, unique value in production environments.

**Usage:**
```sh
git switch feature/step14-jwt-auth
go run main.go --log-level trace --kubeconfig  ~/.kube/config server --enable-leader-election=0
# Obtain a JWT token (for dev/testing)
curl -X POST http://localhost:8080/api/token
# Response: {"token":"<JWT>"}

# Use the token to access protected endpoints
TOKEN=$(curl -s -X POST http://localhost:8080/api/token | jq -r .token)
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/frontendpages
```

**What it does:**
- Secures the platform API with JWT authentication.
- Provides a simple way to test authentication locally.

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
- `pkg/api/jwt*` - API jwt implementation

## License

MIT License. See [LICENSE](LICENSE) for details.

