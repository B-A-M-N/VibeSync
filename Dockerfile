# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY mcp-server/go.mod mcp-server/go.sum ./
RUN go mod download

# Copy the source code
COPY mcp-server/ ./

# Build the binary with security flags
RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o vibesync-orchestrator .

# Stage 2: Hardened Runtime
FROM alpine:latest

# Install lightweight security packages
# tini: Proper signal handling
# libcap: For managing capabilities
RUN apk --no-cache add ca-certificates tini libcap

# Create a non-root user
RUN addgroup -S vibesync && adduser -S vibesync -G vibesync

WORKDIR /home/vibesync/

# Copy the binary
COPY --from=builder /app/vibesync-orchestrator .
RUN chown vibesync:vibesync vibesync-orchestrator && \
    setcap 'cap_net_bind_service=+ep' vibesync-orchestrator

# Create persistence dir with correct permissions
RUN mkdir -p .vibesync/tmp && chown -R vibesync:vibesync .vibesync

USER vibesync

# Use tini to prevent zombie processes
ENTRYPOINT ["/sbin/tini", "--"]

CMD ["./vibesync-orchestrator"]