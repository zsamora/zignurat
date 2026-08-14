# Golang Docker image
FROM golang:1.26 AS base

# Move to working directory /build
WORKDIR /build

# Copy the go.mod and go.sum files first (better layer caching)
COPY notariat/go.mod notariat/go.sum ./notariat/
COPY utils/go.mod utils/go.sum ./utils/

# Install dependencies
WORKDIR /build/notariat
RUN go mod download

# Copy the entire source code into the container (notariat + the shared utils module)
WORKDIR /build
COPY notariat/ ./notariat/
COPY utils/ ./utils/

# Build the application
WORKDIR /build/notariat
RUN go build -o notariat-app

# Ports exposed
EXPOSE 9002

# Start the application
CMD ["/build/notariat/notariat-app"]
