# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
FROM golang:1.21.3 as builder

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed.
RUN go mod download

# Copy the source from the current directory to the working Directory inside the container.
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Use a Docker multi-stage build to create a lean production image.
# Start from the alpine image, which provides a small footprint.
FROM alpine:latest

# Add CA certificates
RUN apk --no-cache add ca-certificates

# Copy the Pre-built binary file from the previous stage. 
COPY --from=builder /app/main .

# Expose port 8080 and 50051 to the outside world
EXPOSE 8080 50051

# Command to run the executable
CMD ["./main"]