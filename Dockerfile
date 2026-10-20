# Multi-stage build for Go application
# Stage 1: Build the Go binary
FROM golang:1.26-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build arguments for version injection
ARG VERSION=dev
ARG GIT_COMMIT=unknown
ARG BUILD_TIME=unknown

# Build the binary with ldflags to inject version info
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s -X main.Version=${VERSION} -X main.GitCommit=${GIT_COMMIT} -X main.BuildTime=${BUILD_TIME}" \
    -o /app/devops-demo-app \
    main.go

# Stage 2: Create minimal runtime image
FROM gcr.io/distroless/static-debian11:nonroot

# Copy binary from builder
COPY --from=builder /app/devops-demo-app /devops-demo-app

# Use non-root user (distroless provides this by default)
USER nonroot:nonroot

# Expose port
EXPOSE 8080

# Run the application
ENTRYPOINT ["/devops-demo-app"]
