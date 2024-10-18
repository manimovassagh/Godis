# Use an official Go runtime as a parent image
FROM golang:1.23 AS BUILDER

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules manifest (go.mod) and download dependencies
COPY go.mod ./
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app and place the binaries in the build directory
RUN mkdir -p build
RUN go build -o build/godis-server ./cmd/server
RUN go build -o build/godis-cli ./cmd/client

# Use a lightweight image for the final container
FROM alpine:3.18

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the binaries and Makefile from the builder stage
COPY --from=BUILDER /app/build /app/build/
COPY --from=BUILDER /app/Makefile /app/

# Expose the port the server will run on
EXPOSE 5001

# Command to run the server binary
CMD ["./build/godis-server"]