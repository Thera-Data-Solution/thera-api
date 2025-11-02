# Build stage
FROM golang:1.24.9-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Runtime stage
FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8000
CMD ["./server"]
