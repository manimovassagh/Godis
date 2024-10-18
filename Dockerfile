# Use an official Go runtime as a parent image
FROM golang:1.23 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules manifests and download the dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o godis-server ./cmd/server
RUN go build -o godis-cli ./cmd/client

# Use a lightweight image for the final container
FROM alpine:3.18

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the binaries from the builder stage
COPY --from=builder /app/godis-server /app/godis-cli /app/

# Expose the port the server will run on
EXPOSE 5001

# Command to run the server binary when the container starts
CMD ["./godis-server"]