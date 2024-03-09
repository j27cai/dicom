# Start from a Go image with version 1.16
FROM golang:latest

# Set the current working directory inside the container
WORKDIR /api

# Copy the go.mod and go.sum files into the container
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the binary executable
RUN go test ./...
RUN go build -o dicom api/cmd/main.go

# Expose the port that the container will listen on
EXPOSE 8000

# Set the entry point for the container
CMD ["./dicom"]
