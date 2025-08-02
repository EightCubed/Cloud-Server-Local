# =======================
# Build Stage
# =======================
FROM golang:alpine AS builder

RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd

# =======================
# Final Stage
# =======================
FROM alpine:latest

RUN apk add --no-cache ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
