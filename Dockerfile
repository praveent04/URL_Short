# Start from golang base image
FROM golang:1.19-alpine as builder

# Add Maintainer info
LABEL maintainer="Praveen"

# Set the current working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source from the current directory to the working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Start a new stage from scratch
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Expose port 3000 to the outside
EXPOSE 3000

# Command to run the executable
CMD ["./main"]