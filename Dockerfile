# Multi-stage build for zenjob-api-linter
FROM --platform=${BUILDPLATFORM} golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH} \
    go build -ldflags '-s -w -extldflags "-static"' -o api-linter ./cmd/api-linter

# Final stage
FROM alpine:3

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Set labels
LABEL org.opencontainers.image.source=https://github.com/zenjob/api-linter
LABEL org.opencontainers.image.description="Zenjob API Linter with custom rules"
LABEL org.opencontainers.image.licenses=Apache-2.0

# Copy binary from builder stage
COPY --from=builder /app/api-linter /usr/local/bin/

# Create non-root user
RUN addgroup -g 10001 -S api-linter && \
    adduser -u 10001 -S api-linter -G api-linter

# Switch to non-root user
USER 10001:10001

# Set entrypoint
ENTRYPOINT ["api-linter"]
