FROM golang:1.25.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy service code and shared packages
COPY services/auth ./services/auth
COPY package ./package

# Build
WORKDIR /app/services/auth
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/auth ./cmd

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/auth .

EXPOSE 8000
CMD ["./auth"]

