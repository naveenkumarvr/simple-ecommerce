# api-gateway-go

Very simple Go HTTP API gateway for the simple-ecommerce demo.

Acts only as a reverse proxy; no authentication, no middleware, no business logic.

## Routes

- `POST /api/login` → `user-service-python` `/login`
- `GET /api/products` → `catalog-service-go` `/products`
- `/api/cart/*` → `cart-service-python` `/cart/*`
- `POST /api/order` → `order-service-go` `/order/checkout`
- `GET /health` – health check for the gateway itself.

The service assumes the backend services are reachable at:

- `http://user-service-python:8001`
- `http://catalog-service-go:8002`
- `http://cart-service-python:8003`
- `http://order-service-go:8005`

You can adjust these URLs in `main.go` if needed.

## Running locally

From this directory:

```bash
go run ./...
```

Make sure the backend services are already running and reachable at the configured URLs.

Example calls:

```bash
# Login (proxied to user-service)
curl -X POST http://localhost:8000/api/login \
  -H "Content-Type: application/json" \
  -d '{"username": "alice", "password": "password1"}'

# Get products (proxied to catalog-service)
curl http://localhost:8000/api/products

# Add to cart (proxied to cart-service)
curl -X POST http://localhost:8000/api/cart/add \
  -H "Content-Type: application/json" \
  -d '{"user_id": "11111111-1111-1111-1111-111111111111", "product_id": "car"}'

# Checkout order (proxied to order-service)
curl -X POST http://localhost:8000/api/order \
  -H "Content-Type: application/json" \
  -d '{"user_id": "11111111-1111-1111-1111-111111111111"}'
```

## Docker

To build and run the Docker image:

```bash
docker build -t api-gateway-go .
docker run --rm -p 8000:8000 api-gateway-go
```
