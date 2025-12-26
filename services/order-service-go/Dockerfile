FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . ./
RUN go build -o order-service-go main.go

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/order-service-go ./order-service-go

EXPOSE 8005

ENTRYPOINT ["/app/order-service-go"]
