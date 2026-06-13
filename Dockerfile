# ==============================================================================
# Builder Stage
# ==============================================================================
FROM golang:alpine AS builder

# Install system packages required for compiling Go binaries
RUN apk add --no-cache git ca-certificates tzdata

RUN go install github.com/swaggo/swag/cmd/swag@v1.16.6

# Set the working directory inside the container
WORKDIR /src

# Copy go.mod and go.sum files to cache dependencies
COPY go.mod go.sum ./

# Download all dependencies. They will be cached if the go.mod/go.sum files don't change.
RUN go mod download

# Copy the entire source code
COPY . .

# Generate Swagger documentation before compiling
RUN swag init --parseDependency --parseInternal

# Build the application binary.
# CGO_ENABLED=0 builds a statically linked binary (ideal for scratch/alpine runners).
# -ldflags="-s -w" strips debugging symbols to reduce binary size.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/lepas main.go

# ==============================================================================
# Runner Stage
# ==============================================================================
FROM alpine:3.19 AS runner

# Create a non-root system group and user for running the application securely
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set the working directory
WORKDIR /app

# Copy CA certificates and timezone information from the builder stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the compiled Go binary
COPY --from=builder /bin/lepas ./lepas

# Copy the web templates directory containing layout/views (loaded at runtime)
COPY --from=builder /src/web ./web

# Change ownership of the application folder to the non-root user
RUN chown -R appuser:appgroup /app

# Switch to the non-root user
USER appuser

# Expose the application port (defaults to 7002/7001, configurable at runtime)
EXPOSE 7002

# Set default production environment variables
ENV PORT=7002 \
    SERVER_ENV=prod

# Define the entrypoint and default command
ENTRYPOINT ["./lepas"]
CMD ["serve"]
