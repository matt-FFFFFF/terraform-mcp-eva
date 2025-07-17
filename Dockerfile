# Builder stage
FROM golang:latest AS builder

# Set working directory
WORKDIR /app

# Copy source code
COPY . .

# Download dependencies and build the application using TARGETARCH for multi-platform builds
ARG TARGETARCH
RUN go mod download && \
    GOOS=linux GOARCH=${TARGETARCH} CGO_ENABLED=0 go build -o terraform-mcp-eva .

# Runner stage
FROM busybox:latest

# Create a non-root user
RUN adduser -D -s /bin/sh appuser

# Set working directory
WORKDIR /home/appuser

# Copy the binary from builder stage
COPY --from=builder /app/terraform-mcp-eva .

# Change ownership to appuser
RUN chown appuser:appuser terraform-mcp-eva

# Switch to non-root user
USER appuser

# Declare environment variables with default values
ENV TRANSPORT_MODE=stdio
ENV TRANSPORT_HOST=127.0.0.1
ENV TRANSPORT_PORT=8080

# Set the entrypoint
ENTRYPOINT ["./terraform-mcp-eva"]
