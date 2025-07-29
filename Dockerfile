# Builder stage
FROM --platform=${BUILDPLATFORM} golang:1.24.5 AS builder
ARG TARGETARCH
ENV GOARCH=${TARGETARCH}
# Set working directory
WORKDIR /src

# Copy source code
COPY . .

# Download dependencies and build the application using TARGETARCH for multi-platform builds

RUN go mod download && \
  GOOS=linux CGO_ENABLED=0 go build -o terraform-mcp-eva .

# Runner stage
FROM busybox:latest

# Create a non-root user
RUN adduser -D -s /bin/sh appuser

# Set working directory
WORKDIR /home/appuser

# Copy the binary from builder stage
COPY --chown=root:root --from=builder /src/terraform-mcp-eva .

# Set permissions for the binary
RUN chmod 755 terraform-mcp-eva

# Switch to non-root user
USER appuser

# Declare environment variables with default values
ENV TRANSPORT_MODE=stdio
ENV TRANSPORT_HOST=127.0.0.1
ENV TRANSPORT_PORT=8080

# Set the entrypoint
ENTRYPOINT ["./terraform-mcp-eva"]
