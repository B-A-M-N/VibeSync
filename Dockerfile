# Stage 1: Build the Go binary
# Upgrade to 1.24 to address stdlib vulnerabilities (CVEs)
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY mcp-server/go.mod mcp-server/go.sum $WORKDIR/
RUN go mod download

# Copy the source code
COPY mcp-server/ $WORKDIR/

# Build the binary with advanced security flags:
# -trimpath: Remove file system paths from the binary
# -buildmode=pie: Position Independent Executable (ASLR support)
# -ldflags: Strip symbols (-s -w)
RUN CGO_ENABLED=1 GOOS=linux go build \
    -trimpath \
    -buildmode=pie \
    -ldflags="-s -w -extldflags '-static'" \
    -o vibesync-orchestrator .

# Stage 2: Hardened Runtime
FROM alpine:latest

# Install lightweight security packages
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
