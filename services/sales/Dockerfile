FROM golang:1.25.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY services/sales ./services/sales
COPY services/contact ./services/contact
COPY services/inventory ./services/inventory
COPY package ./package

WORKDIR /app/services/sales
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/sales ./cmd

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/sales .

EXPOSE 8000
CMD ["./sales"]

