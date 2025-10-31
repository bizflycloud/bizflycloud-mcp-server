# Build stage
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /build

# Install git and ca-certificates (needed for go modules)
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o bizfly-mcp-server .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 mcp && \
    adduser -D -u 1000 -G mcp mcp

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/bizfly-mcp-server .

# Change ownership to non-root user
RUN chown -R mcp:mcp /app

# Switch to non-root user
USER mcp

# Expose port (if needed in future)
# EXPOSE 8080

# Set environment variables with defaults
ENV BIZFLY_REGION=HaNoi
ENV BIZFLY_API_URL=https://manage.bizflycloud.vn

# Run the application
ENTRYPOINT ["./bizfly-mcp-server"]

