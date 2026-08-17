# Golang Docker image
FROM golang:1.26 AS base

# Move to working directory /build
WORKDIR /build

# Copy the go.mod and go.sum files first (better layer caching)
COPY identirat/go.mod identirat/go.sum ./identirat/
COPY utils/go.mod utils/go.sum ./utils/

# Install dependencies
WORKDIR /build/identirat
RUN go mod download

# Copy the entire source code into the container (identirat + the shared utils module)
WORKDIR /build
COPY identirat/ ./identirat/
COPY utils/ ./utils/

# Static Tailwind CSS
WORKDIR /build/identirat
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
RUN go build -o identirat-app

# Ports exposed
EXPOSE 9000

# Start the application
CMD ["/build/identirat/identirat-app"]
