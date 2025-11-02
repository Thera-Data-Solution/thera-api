FROM golang:1.24.9-alpine AS builder

WORKDIR /app

# install dependency untuk libwebp & build tools
RUN apk add --no-cache gcc g++ make libc-dev libwebp-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o server .

FROM alpine:3.19
RUN apk add --no-cache libwebp
WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]
