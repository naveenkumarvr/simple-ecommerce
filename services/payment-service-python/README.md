# payment-service-python

Very simple in-memory payment service using FastAPI, intended for Kubernetes demos.

## Endpoints

- `POST /pay` – body: `{ "user_id": "...", "amount": 123.45, "currency": "USD", "source": "fake-card-1" }`; returns a `payment_id` and `status`.
- `GET /payments/{user_id}` – returns all payments for that user.
- `GET /health` – basic health check.

The payments are stored only in memory as a simple dictionary: `user_id -> [payments]`.

## Running locally

From this directory:

```bash
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
uvicorn app.main:app --host 0.0.0.0 --port 8004 --reload
```

Example calls:

```bash
# Create a payment
curl -X POST http://localhost:8004/pay \
  -H "Content-Type: application/json" \
  -d '{"user_id": "11111111-1111-1111-1111-111111111111", "amount": 99.99, "currency": "USD", "source": "card-123"}'

# Get payments for a user
curl http://localhost:8004/payments/11111111-1111-1111-1111-111111111111
```

## Docker

To build and run the Docker image:

```bash
docker build -t payment-service-python .
docker run --rm -p 8004:8004 payment-service-python
```
