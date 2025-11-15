FROM golang:1.25.0-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy gateway code and shared packages
COPY gateway ./gateway
COPY package ./package

# Build
WORKDIR /app/gateway
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/gateway ./cmd

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/gateway .

EXPOSE 8080
CMD ["./gateway"]

