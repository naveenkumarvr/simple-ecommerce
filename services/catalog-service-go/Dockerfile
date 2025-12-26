FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . ./
RUN go build -o catalog-service-go main.go

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/catalog-service-go ./catalog-service-go

EXPOSE 8002

ENTRYPOINT ["/app/catalog-service-go"]
