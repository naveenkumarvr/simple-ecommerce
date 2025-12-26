# catalog-service-go

Very simple Go HTTP service for Kubernetes demos.

## Endpoints

- `GET /products` – returns all products as JSON.
- `GET /products/{id}` – returns a single product by ID (`car`, `bike`, `bus`, `truck`) or 404.

Hardcoded products with prices:

- `car` = 100
- `bike` = 40
- `bus` = 150
- `truck` = 200

## Running locally

From this directory:

```bash
go run ./...
```

Then:

```bash
curl http://localhost:8002/products
curl http://localhost:8002/products/car
```
