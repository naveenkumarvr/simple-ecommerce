# order-service-go

Very simple Go HTTP order service for Kubernetes demos.

This service exposes a single endpoint that coordinates between the cart and payment services to simulate a checkout.

## Endpoints

- `POST /order/checkout` – perform a very simple checkout flow.
  - Calls `cart-service-python` to fetch the user's cart.
  - Computes a naive total (sum of item quantities).
  - Calls `payment-service-python` to create a payment.
  - Returns a JSON response indicating success.

### Request body

```json
{
  "user_id": "11111111-1111-1111-1111-111111111111"
}
```

### Example success response

```json
{
  "message": "order successful",
  "payment_id": "<payment-id-from-payment-service>"
}
```

## Running locally

From this directory:

```bash
go run ./...
```

Make sure the dependent services are running and reachable at:

- `http://cart-service-python:8003` (cart service)
- `http://payment-service-python:8004` (payment service)

For local testing without Kubernetes/docker-compose, you can adapt `cartServiceURL` and `paymentServiceURL` in `main.go` to point to `localhost` instead.

### Example call

```bash
curl -X POST http://localhost:8005/order/checkout \
  -H "Content-Type: application/json" \
  -d '{"user_id": "11111111-1111-1111-1111-111111111111"}'
```

## Docker

To build and run the Docker image:

```bash
docker build -t order-service-go .
docker run --rm -p 8005:8005 order-service-go
```
