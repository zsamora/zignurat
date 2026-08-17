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

# Static Tailwind CSS v3.4.17
WORKDIR /build/notariat
RUN ARCH=$(uname -m) && \
    case "$ARCH" in \
        x86_64) TW_ARCH=x64 ;; \
        aarch64) TW_ARCH=arm64 ;; \
        *) echo "Unsupported architecture: $ARCH" && exit 1 ;; \
    esac && \
    curl -sLo /usr/local/bin/tailwindcss "https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.17/tailwindcss-linux-${TW_ARCH}" && \
    chmod +x /usr/local/bin/tailwindcss
RUN mkdir -p static && tailwindcss -i ./tailwind/input.css -o ./static/tailwind.css --minify

# Build the application
RUN go build -o notariat-app

# Ports exposed
EXPOSE 9002

# Start the application
CMD ["/build/notariat/notariat-app"]
