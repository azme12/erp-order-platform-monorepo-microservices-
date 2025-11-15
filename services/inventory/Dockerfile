FROM golang:1.25.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY services/inventory ./services/inventory
COPY package ./package

WORKDIR /app/services/inventory
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/inventory ./cmd

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/inventory .

EXPOSE 8000
CMD ["./inventory"]

