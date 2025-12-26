FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . ./
RUN go build -o api-gateway-go main.go

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/api-gateway-go ./api-gateway-go

EXPOSE 8000

ENTRYPOINT ["/app/api-gateway-go"]
