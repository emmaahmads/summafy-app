# Use the official Go image as the base
FROM golang:alpine

# Set the working directory to /app
WORKDIR /app

# Copy the Go source code into the container
COPY main.go .

# Build the Go binary
RUN go build -o server main.go

# Expose port 8080
EXPOSE 8080

# Run the Go binary when the container starts
CMD ["./server"]