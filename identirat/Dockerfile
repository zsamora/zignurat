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

# Build the application
WORKDIR /build/identirat
RUN go build -o identirat-app

# Ports exposed
EXPOSE 9000

# Start the application
CMD ["/build/identirat/identirat-app"]
