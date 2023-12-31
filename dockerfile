# Use official Golang image
FROM golang:1.21.3 as builder

# Set working directory
WORKDIR /app

# Copy the source code
COPY . .

# Download and install the dependencies
RUN go mod download

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Use Alpine for the final image
FROM alpine:latest

# Add CA certificates
RUN apk --no-cache add ca-certificates

# Copy the built binary from the builder stage
COPY --from=builder /app/server /server

# Expose the port the app runs on
EXPOSE 8080 50051

# Run the server
CMD ["/server"]