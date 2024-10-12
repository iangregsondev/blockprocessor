# Stage 1: Build stage
FROM golang:1.23.2-alpine AS build

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Set the working directory
WORKDIR /app

# Copy and download dependencies separately for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=1 go build -o bitcoin-block-processor ./cmd/bitcoin/blockprocessor/main.go

# Stage 2: Final stage
FROM alpine:3.20

# Set the working directory
WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/bitcoin-block-processor .

# Set up to run the app, with arguments provided at runtime
ENTRYPOINT ["/app/bitcoin-block-processor"]