# DevOps Demo Status Service

A simple Go web service built to demonstrate CI/CD pipelines and GitOps workflows with GitHub Actions and ArgoCD.

## 🎯 Purpose

This application is intentionally simple - it's a **demo vehicle** for showcasing:
- Automated CI/CD pipelines with GitHub Actions
- GitOps deployment with ArgoCD
- Container security scanning with Trivy
- Kubernetes deployment patterns
- Build-time artifact injection (Git SHA, timestamps)

**The app displays:**
- Application version (Git commit SHA)
- Build timestamp
- Environment name
- Kubernetes pod information
- Health status and uptime

## 🚀 Quick Start

### Run Locally

```bash
# Run with Go
go run main.go

# Run tests
go test -v

# Build binary
go build -o devops-demo-app main.go

# Run binary
./devops-demo-app
```

Visit: http://localhost:8080

### Build Docker Image Locally

```bash
# Build with version info
docker build \
  --build-arg VERSION=v1.0.0 \
  --build-arg GIT_COMMIT=$(git rev-parse --short HEAD) \
  --build-arg BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  -t devops-demo-app:latest .

# Run container
docker run -p 8080:8080 \
  -e ENVIRONMENT=local \
  -e POD_NAME=local-container \
  devops-demo-app:latest
```

Visit: http://localhost:8080

## 📡 API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | GET | HTML status dashboard |
| `/health` | GET | Health check (for K8s probes) |
| `/api/status` | GET | JSON status information |
| `/metrics` | GET | Prometheus metrics |

### Example Response (`/api/status`):

```json
{
  "version": "v1.0.0",
  "git_commit": "abc123f",
  "build_time": "2026-10-16T10:23:15Z",
  "environment": "production",
  "pod_name": "demo-app-7d8f9-xz5p",
  "namespace": "devops-demo",
  "status": "healthy",
  "uptime": "2h15m30s",
  "timestamp": "2026-10-16T12:38:45Z"
}
```

## 🏗️ Architecture

### Multi-Stage Dockerfile

1. **Builder stage:** Compiles Go binary with version info injected at build time
2. **Runtime stage:** Minimal distroless image for security (no shell, minimal attack surface)

**Image size:** ~5MB (vs 100MB+ for Python/Node.js apps)

### Build-Time Injection

Version information is injected during compilation using Go's `-ldflags`:

```bash
go build -ldflags="-X main.Version=${VERSION} -X main.GitCommit=${GIT_COMMIT}"
```

This allows each deployment to show exactly which Git commit is running.

## 🔒 Security Features

- ✅ **Non-root user:** Runs as `nonroot:nonroot` (UID 65532)
- ✅ **Distroless base:** No shell, no package managers, minimal attack surface
- ✅ **No CGO:** Static binary, no C dependencies
- ✅ **Small image:** ~5MB, reduced attack surface
- ✅ **Read-only filesystem:** Can run with read-only root filesystem

## ☸️ Kubernetes Deployment

### Environment Variables

The app reads from Kubernetes downward API:

```yaml
env:
  - name: ENVIRONMENT
    value: "production"
  - name: POD_NAME
    valueFrom:
      fieldRef:
        fieldPath: metadata.name
  - name: POD_NAMESPACE
    valueFrom:
      fieldRef:
        fieldPath: metadata.namespace
```

### Health Probes

```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 10

readinessProbe:
  httpGet:
    path: /health
    port: 8080
  initialDelaySeconds: 3
  periodSeconds: 5
```

## 🧪 Testing

```bash
# Run all tests
go test -v

# Run tests with coverage
go test -v -cover

# Generate coverage report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## 🛠️ Development

### Project Structure

```
devops-demo-app/
├── main.go           # Application code
├── main_test.go      # Unit tests
├── go.mod            # Go dependencies
├── Dockerfile        # Multi-stage build
├── .gitignore        # Git ignore rules
└── README.md         # This file
```

### Adding Features

This app is intentionally minimal. To add features:

1. Keep it simple - focus is on the **pipeline**, not the app
2. Add tests for any new handlers
3. Update this README

## 📊 CI/CD Pipeline

*Coming in Milestone 4*

The GitHub Actions pipeline will:
1. Run tests (`go test`)
2. Build Docker image with Git SHA tag
3. Scan image with Trivy
4. Push to Docker Hub
5. Update GitOps repository

## 🔄 GitOps with ArgoCD

*Coming in Milestone 5*

ArgoCD will:
1. Watch the GitOps repository
2. Detect manifest changes
3. Automatically deploy to Kubernetes
4. Provide rollback capabilities

## 📈 Metrics

The `/metrics` endpoint provides Prometheus-compatible metrics:

```
# Application uptime
devops_demo_uptime_seconds 7890

# Application info (labels)
devops_demo_info{version="v1.0.0",git_commit="abc123f",environment="production"} 1
```

## 🎓 Learning Outcomes

Building this app demonstrates:

| Skill | Implementation |
|-------|----------------|
| **Go Development** | HTTP server, JSON APIs, HTML templates |
| **Testing** | Unit tests, table-driven tests, HTTP test recorder |
| **Containerization** | Multi-stage builds, distroless images, non-root |
| **Security** | Static analysis, image scanning, minimal attack surface |
| **12-Factor App** | Environment config, stateless, build/release/run separation |
| **Kubernetes** | Downward API, health probes, environment injection |
| **CI/CD** | Build automation, artifact versioning, GitOps |

## 📝 License

MIT

## 👤 Author

Built as a portfolio project to demonstrate DevOps and Platform Engineering skills.

---

**Note:** This is a demo application. The real value is in the automated pipeline and GitOps workflow, not the application complexity.

## ArgoCD Integration Test - 2026-10-21 07:01:49
