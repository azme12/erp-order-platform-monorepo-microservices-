FROM golang:1.25.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY services/purchase ./services/purchase
COPY package ./package

WORKDIR /app/services/purchase
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/purchase ./cmd

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/purchase .

EXPOSE 8000
CMD ["./purchase"]

