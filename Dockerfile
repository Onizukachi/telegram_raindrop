# Build Stage
FROM golang:1.25-bookworm AS base

# Development stage
# =============================================================================
# Create a development stage based on the "base" image
FROM base AS development

# Change the working directory to /app
WORKDIR /app

# Install the air CLI for auto-reloading
RUN go install github.com/air-verse/air@latest

# Copy the go.mod and go.sum files to the /app directory
COPY go.mod go.sum ./

# Install dependencies
RUN go mod download

# Start air for live reloading
CMD ["air", "--build.cmd", "go build -o ./tmp/main -buildvcs=false ./cmd/bot", "--build.bin", "./tmp/main"]

# Builder stage
# =============================================================================
# Create a builder stage based on the "base" image
FROM base AS builder

# Move to working directory /build
WORKDIR /build

# Copy the go.mod and go.sum files to the /build directory
COPY go.mod go.sum ./

# Install dependencies
RUN go mod download

# Copy the entire source code into the container
COPY . .

# Build the application
# Turn off CGO to ensure static binaries
RUN CGO_ENABLED=0  go build -o go-bot ./cmd/bot

# Production stage
# =============================================================================
# Create a production stage to run the application binary
FROM alpine:3.18 AS production

# Move to working directory /prod
WORKDIR /prod

# Copy binary from builder stage
COPY --from=builder /build/go-bot ./

# Copy necessary files
COPY messages.yml ./
COPY ./.env ./
COPY ./internal/storage/migrations ./internal/storage/migrations

# Document the port that may need to be published
EXPOSE 8000

CMD ["/prod/go-bot"]
