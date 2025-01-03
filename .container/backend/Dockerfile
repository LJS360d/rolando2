# Use Golang official image as builder
FROM golang:1.23-alpine as builder

# Install make utility
RUN apk add --no-cache make

# Set the working directory
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the backend source code
COPY . .

# Build the Go application using 'make build'
RUN make build

# Now, use a minimal image for the final container
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Install necessary libraries (e.g., sqlite3)
RUN apk --no-cache add sqlite-libs

# Copy the binary and database from the build stage
COPY --from=builder /app/main .
COPY --from=builder /app/database.sqlite ./database.sqlite

# Expose the port the backend app listens on
EXPOSE ${BACKEND_PORT}

# Command to run the Go application
CMD ["./main"]
