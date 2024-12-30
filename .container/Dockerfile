FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY ../go.mod ../go.sum ./
RUN go mod tidy
COPY ../ ./
RUN go build -o bot .
EXPOSE 8080

ENTRYPOINT ["/root/entrypoint.sh"]
