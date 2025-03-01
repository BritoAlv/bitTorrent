FROM golang:latest

WORKDIR /app

# Copy only go.mod and go.sum first, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Install required packages
RUN apt-get update && apt-get install -y iproute2 && rm -rf /var/lib/apt/lists/*

# Copy source files separately to avoid invalidating cache
COPY ./ .

# Build the application
RUN go build -o Server ./cmd/serverMain/StartUpServer.go

CMD ["./Server"]