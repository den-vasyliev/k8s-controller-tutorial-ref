# Platform API (CRUD + Swagger)

- Added RESTful CRUD API endpoints for the FrontendPage CRD using FastHTTP and fasthttprouter.
- API handlers use the controller-runtime client to create, update, delete, and list FrontendPage resources in Kubernetes, triggering reconciliation.
- Integrated [Swagger](https://swagger.io/) documentation and served Swagger UI for easy API exploration.
- All API endpoints are under `/api/frontendpages`.

**Usage:**
```sh
git switch feature/step12-platform-api 
go run main.go --log-level trace --kubeconfig  ~/.kube/config server --enable-leader-election=0
```
| Method | Endpoint                              | Example Command / Payload                                                                                                                                                                                                                                                                                                                                 | Description                |
|--------|---------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------|
| POST   | /api/frontendpages                    | <pre>curl -X POST http://localhost:8080/api/frontendpages \<br>  -H 'Content-Type: application/json' \<br>  -d '{<br>  "metadata": {<br>    "name": "my-page"<br>  },<br>  "spec": {<br>    "contents": "&lt;h1&gt;Hello&lt;/h1&gt;",<br>    "image": "nginx:latest",<br>    "replicas": 2<br>  }<br>}'</pre>                | Create a new FrontendPage  |
| GET    | /api/frontendpages                    | <pre>curl http://localhost:8080/api/frontendpages</pre>                                                                                                                                                                                                                                                                                                   | List all FrontendPages     |
| GET    | /api/frontendpages/my-page            | <pre>curl http://localhost:8080/api/frontendpages/my-page</pre>                                                                                                                                                                                                                                                                                           | Get a FrontendPage by name |
| PUT    | /api/frontendpages/my-page            | <pre>curl -X PUT http://localhost:8080/api/frontendpages/my-page \<br>  -H 'Content-Type: application/json' \<br>  -d '{<br>  "spec": {<br>    "contents": "&lt;h1&gt;Updated&lt;/h1&gt;",<br>    "image": "nginx:alpine",<br>    "replicas": 1<br>  }<br>}'</pre>                                                    | Update a FrontendPage      |
| DELETE | /api/frontendpages/my-page            | <pre>curl -X DELETE http://localhost:8080/api/frontendpages/my-page</pre>                                                                                                                                                                                                                                                                                | Delete a FrontendPage      |

- Visit `http://localhost:8080/swagger/index.html` for interactive API docs.

**What it does:**
- Exposes CRUD API for FrontendPage resources, backed by Kubernetes CRDs and controller logic.
- Provides OpenAPI/Swagger docs and UI for easy testing and documentation.

## Running Tests

This project uses [envtest](https://book.kubebuilder.io/reference/envtest.html) and controller-runtime for integration and controller tests.

### Prerequisites
- Go (see go.mod for version)
- Make
- The `setup-envtest` binary (automatically handled by the Makefile)
- CRD YAMLs present in `config/crd/`

### Run all tests
```sh
make test
```
This will:
- Download and set up envtest if needed
- Run all Go tests in the project (including controller and utility tests)

### Run only controller tests
```sh
make test-controller
```
This will:
- Run only the tests in `pkg/ctrl/` (controller logic)

### Test output
- Test logs will show simulated etcd state and resource changes for CRDs and controllers.
- JUnit XML and coverage reports are generated as `report.xml` and `coverage.xml`.

### Troubleshooting
- If you see errors about missing CRDs, ensure you have generated CRDs in `config/crd/` (see Usage above for controller-gen command).
- If you see errors about envtest, try running `make envtest` to ensure the binary is present in `bin/`.
- If you add new CRDs or controllers, re-run `controller-gen` and re-run tests.

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

# FrontendPage API - Test Instructions

## Prerequisites
- Go 1.20+
- [kubebuilder test assets](https://book.kubebuilder.io/reference/envtest.html#installing-envtest-binaries) (envtest)
- The following Go dependencies:
  - github.com/valyala/fasthttprouter
  - github.com/google/uuid
  - github.com/stretchr/testify

Install them with:
```sh
go get github.com/valyala/fasthttprouter github.com/google/uuid github.com/stretchr/testify
```

## Running the API E2E/Integration Tests

1. **Install envtest assets** (if not already):
   ```sh
   kubebuilder envtest install
   # or manually download and set KUBEBUILDER_ASSETS
   ```

2. **Run the tests:**
   ```sh
   go test -v -tags=testtools ./pkg/api  -tags=testtools
   ```
   - The tests will spin up a real Kubernetes API server (envtest), start the controller manager, and exercise the full API (create, update, delete FrontendPage resources).
   - Each test uses a unique resource name to avoid collisions.

## GitHub Actions: How to Test Frontend API

To test the frontend API in a GitHub Actions workflow:

1. **Set up Go and envtest in your workflow:**
   ```yaml
   - uses: actions/checkout@v3
   - uses: actions/setup-go@v4
     with:
       go-version: '1.20'
   - name: Install envtest tools
     run: |
       go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest
       setup-envtest use 1.29.0 # or your preferred version
   ```

2. **Run the tests:**
   ```yaml
   - name: Run FrontendPage API tests
     run: go test -v -tags=testtools ./pkg/api
   ```

## Example API Test Payload (for local or CI testing)

To test the API endpoints (e.g., with curl, Postman, or a custom script), use the following JSON payloads:

### Create FrontendPage
```json
{
  "metadata": {
    "name": "test-frontend-page-1234",
    "namespace": "default"
  },
  "spec": {
    "contents": "<h1>Hello</h1>",
    "image": "nginx:latest",
    "replicas": 2
  }
}
```

### Update FrontendPage
- Fetch the existing resource to get its `resourceVersion`.
- Use the same structure as above, but include the `resourceVersion` in `metadata` and update the fields you want.

### Delete FrontendPage
- Send a DELETE request to `/api/frontendpages/{name}`.

## Notes
- The tests do not require a running Kubernetes cluster; everything runs in-process using envtest.
- Pods will not become Ready in envtest; tests only check for resource existence and spec.
- For troubleshooting, check the test logs for API call details and controller reconciliation logs.


## Ngrok
```sh
curl -sSL https://ngrok-agent.s3.amazonaws.com/ngrok.asc \
  | sudo tee /etc/apt/trusted.gpg.d/ngrok.asc >/dev/null \
  && echo "deb https://ngrok-agent.s3.amazonaws.com buster main" \
  | sudo tee /etc/apt/sources.list.d/ngrok.list \
  && sudo apt update \
  && sudo apt install ngrok
  ```
### Configure API key
```sh
ngrok config add-authtoken xxxxxxxxxxxxxxxxxxxxxxxxxx
```

### Run proxy
```sh
ngrok http --url=quietly-just-ferret.ngrok-free.app 8080
```

## License

MIT License. See [LICENSE](LICENSE) for details.
