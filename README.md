# Golang Kubernetes Controller Tutorial

This project is a step-by-step tutorial for DevOps and SRE engineers to learn about building Golang CLI applications and Kubernetes controllers. Each step is implemented as a feature branch and includes a README section with explanations and command history.

---

## Step 1: Golang CLI Application using Cobra

- Initialized a new CLI application using [cobra-cli](https://github.com/spf13/cobra).
- Provides a basic command-line interface.

**Command history:**
```sh
git checkout -b step1-cobra-cli
cobra-cli init --pkg-name github.com/yourusername/k8s-controller-tutorial
# edited main.go, cmd/root.go
```

---

## Step 2: Zerolog for Log Levels

- Integrated [zerolog](https://github.com/rs/zerolog) for structured logging.
- Supports log levels: info, debug, trace, warn, error.

**Command history:**
```sh
git checkout -b step2-zerolog
go get github.com/rs/zerolog
# edited cmd/root.go to add zerolog logging
```

---

## Step 3: pflag for Log Level Flags

- Added [pflag](https://github.com/spf13/pflag) to support a `--log-level` flag.
- Users can set log level via CLI flag.

**Usage:**
```sh
go run main.go --log-level debug
```

**Command history:**
```sh
git checkout -b step3-pflag-loglevel
# edited cmd/root.go to add log-level flag
```

---

## Step 4: FastHTTP Server Command

- Added a new `server` command using [fasthttp](https://github.com/valyala/fasthttp).
- The command starts a FastHTTP server with a configurable port (default: 8080).
- Uses zerolog for logging.

**Usage:**
```sh
go run main.go server --port 8080
```

**What it does:**
- Starts a FastHTTP server on the specified port.
- Responds with "Hello from FastHTTP!" to any request.

**Command history:**
```sh
git checkout -b step4-fasthttp-server
go get github.com/valyala/fasthttp
# created cmd/server.go, added server command
# added cmd/server_test.go for basic tests
go mod tidy
git add .
git commit -m "step4: add fasthttp server command with port flag"
```

---

## Step 6: List Kubernetes Deployments with client-go

- Added a new `list` command using [k8s.io/client-go](https://github.com/kubernetes/client-go).
- Lists deployments in the default namespace.
- Supports a `--kubeconfig` flag to specify the kubeconfig file for authentication.
- Uses zerolog for error logging.

**Usage:**
```sh
go run main.go list --kubeconfig ~/.kube/config
```

**What it does:**
- Connects to the Kubernetes cluster using the provided kubeconfig file.
- Lists all deployments in the `default` namespace and prints their names.

**Command history:**
```sh
go get k8s.io/client-go@v0.29.0
# created cmd/list.go, added list command
# edited go.mod, go.sum
# run go mod tidy
git add .
git commit -m "step6: add list command for Kubernetes deployments using client-go"
```

---

## Step 7: Deployment Informer with client-go

- Added a new `informer` command using [k8s.io/client-go](https://github.com/kubernetes/client-go).
- Runs a shared informer for Deployments in the default namespace.
- Supports both kubeconfig and in-cluster authentication (flags: `--kubeconfig`, `--in-cluster`).
- Logs add, update, and delete events for Deployments using zerolog.

**Usage:**
```sh
go run main.go informer --kubeconfig ~/.kube/config
go run main.go informer --in-cluster
```

**What it does:**
- Connects to the Kubernetes cluster using the provided kubeconfig file or in-cluster config.
- Watches for Deployment events (add, update, delete) in the `default` namespace and logs them.

**Command history:**
```sh
# created cmd/informer.go, added informer command
git add .
git commit -m "step7: add informer command for Kubernetes deployments using client-go"
```

---

## Step 8: /deployments JSON API Endpoint

- Added a `/deployments` endpoint to the FastHTTP server.
- Returns a JSON array of deployment names from the informer's cache (default namespace).
- Uses the informer's local cache, not a live API call.

**Usage:**
```sh
curl http://localhost:8080/deployments
# Output: ["deployment1","deployment2",...]
```

**What it does:**
- Serves a JSON array of deployment names currently in the informer cache.
- Does not query the Kubernetes API directly for each request (fast, efficient).

**Command history:**
```sh
# updated pkg/informer/informer.go, cmd/server.go
git add .
git commit -m "step8: add /deployments JSON API endpoint to server using informer cache"
```

---

## Step 10: Leader Election and Metrics for Controller Manager

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
go run main.go server --enable-leader-election=false --metrics-port=9090
```

**Command history:**
```sh
# updated cmd/server.go to add leader election and metrics flags
git add .
git commit -m "step10: add leader election and metrics flags to controller-runtime manager"
```

---

## Step 11: FrontendPage CRD and Advanced Controller Implementation

- Added the Go type for the FrontendPage custom resource in `pkg/apis/frontend/v1alpha1/frontendpage_types.go`.
- Created `groupversion_info.go` to define the group, version, and scheme for the CRD.
- Used [controller-gen](https://github.com/kubernetes-sigs/controller-tools) to generate CRD manifests and deepcopy code.
- Implemented a controller for the FrontendPage CRD using controller-runtime in `pkg/ctrl/frontendpage_controller.go`.
- The controller watches FrontendPage resources and manages both a Deployment and a ConfigMap:
  - Creates/updates a ConfigMap containing the `spec.contents` from the FrontendPage CR.
  - Creates/updates a Deployment that mounts the ConfigMap as a volume and uses the image/replicas from the CR spec.
  - Cleans up both the Deployment and ConfigMap when the FrontendPage is deleted.
- Registered the controller with the manager in `cmd/server.go`.

**What it does:**
- Defines the FrontendPage CRD structure and registers it with the Kubernetes API machinery.
- Generates the CRD YAML and deepcopy methods required for Kubernetes controllers.
- Reconciles FrontendPage resources to ensure a matching Deployment and ConfigMap exist in the cluster.
- Updates the Deployment and ConfigMap if the FrontendPage spec changes.
- Handles creation, update, and cleanup logic for Deployments and ConfigMaps owned by FrontendPage resources.

**Command history:**
```sh
# Add Go types and group version info for FrontendPage
# (edit pkg/apis/frontend/v1alpha1/frontendpage_types.go and groupversion_info.go)

# Run controller-gen to generate CRD and deepcopy code
controller-gen crd:crdVersions=v1 paths=./pkg/apis/... output:crd:dir=./config/crd object paths=./pkg/apis/...

# Scaffold and implement the advanced FrontendPage controller
mkdir -p pkg/ctrl
# created pkg/ctrl/frontendpage_controller.go and implemented controller logic for Deployment and ConfigMap management
# registered the controller in cmd/server.go

# Run the server to start the controller
make run
```

---

## Step 12: Platform API (CRUD + Swagger)

- Added RESTful CRUD API endpoints for the FrontendPage CRD using FastHTTP and fasthttprouter.
- API handlers use the controller-runtime client to create, update, delete, and list FrontendPage resources in Kubernetes, triggering reconciliation.
- Integrated [Swagger](https://swagger.io/) documentation and served Swagger UI for easy API exploration.
- All API endpoints are under `/api/frontendpages`.

**Usage:**
```sh
curl -X POST http://localhost:8080/api/frontendpages -H 'Content-Type: application/json' -d '{"metadata":{"name":"my-page"},"spec":{"contents":"<h1>Hello</h1>","image":"nginx:latest","replicas":2}}'
curl http://localhost:8080/api/frontendpages
curl http://localhost:8080/api/frontendpages/my-page
curl -X PUT http://localhost:8080/api/frontendpages/my-page -H 'Content-Type: application/json' -d '{"spec":{"contents":"<h1>Updated</h1>","image":"nginx:alpine","replicas":1}}'
curl -X DELETE http://localhost:8080/api/frontendpages/my-page
```
- Visit `http://localhost:8080/swagger/index.html` for interactive API docs.

**What it does:**
- Exposes CRUD API for FrontendPage resources, backed by Kubernetes CRDs and controller logic.
- Provides OpenAPI/Swagger docs and UI for easy testing and documentation.

**Command history:**
```sh
# Add CRUD API handlers and FastHTTP router integration
# Add Swagger docs and serve Swagger UI
# Update README with API usage examples
# Commit: "step12: add platform API CRUD endpoints and Swagger docs"
```

---

## Step 13: MCP Integration (Machine Control Protocol)

- Integrated [MCP server](https://github.com/mark3labs/mcp-go) for programmatic control and automation.
- Added `--enable-mcp` and `--mcp-port` flags to the server command.
- MCP server runs in SSE (Server-Sent Events) mode for real-time tool execution and feedback.
- Registered MCP tools for listing and creating FrontendPage resources (extensible for more tools).

**Usage:**
```sh
go run main.go server --enable-mcp --mcp-port 9090
# MCP server will be available on http://localhost:9090
```
- Use an MCP client or compatible tool to connect and invoke registered tools.

**What it does:**
- Enables external systems to interact with the controller via the MCP protocol (list/create FrontendPages, etc.).
- SSE mode provides real-time updates for tool execution.

**Command history:**
```sh
# Add MCP server integration and flags
# Register MCP tools for FrontendPage
# Start MCP server in SSE mode if enabled
# Commit: "step13: add MCP integration and SSE server mode"
```

---

## Step 14: JWT Authentication

- Added JWT authentication middleware to protect all `/api/frontendpages` endpoints.
- Only requests with a valid JWT in the `Authorization: Bearer <token>` header are allowed.
- Added a `/api/token` endpoint for local testing, which issues a JWT for a test user (valid for 1 hour).
- Uses a hardcoded secret for development (update for production use).

**Usage:**
```sh
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

**Command history:**
```sh
# Add JWT middleware and protect /api/frontendpages endpoints
# Add /api/token endpoint for local testing
# Commit: "step14: add JWT authentication middleware and token endpoint"
```

---

Continue to the next steps for more advanced Kubernetes and controller features! 