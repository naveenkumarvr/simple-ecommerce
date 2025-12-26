# cart-service-python

Very simple in-memory cart service using FastAPI, intended for Kubernetes demos.

## Endpoints

- `POST /cart/add` – body: `{ "user_id": "...", "product_id": "..." }`; increments quantity of that product in the user's cart.
- `GET /cart/{user_id}` – returns the cart for that user.
- `POST /cart/clear` – body: `{ "user_id": "...", "product_id": "ignored" }`; clears the cart for that user.
- `GET /health` – basic health check.

The cart is stored only in memory as a simple dictionary: `user_id -> [items]`.

## Running locally

From this directory:

```bash
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
uvicorn app.main:app --host 0.0.0.0 --port 8003 --reload
```

Example calls:

```bash
# Add items
curl -X POST http://localhost:8003/cart/add \
  -H "Content-Type: application/json" \
  -d '{"user_id": "11111111-1111-1111-1111-111111111111", "product_id": "car"}'

curl -X POST http://localhost:8003/cart/add \
  -H "Content-Type: application/json" \
  -d '{"user_id": "11111111-1111-1111-1111-111111111111", "product_id": "bike"}'

# Get cart
curl http://localhost:8003/cart/11111111-1111-1111-1111-111111111111

# Clear cart
curl -X POST http://localhost:8003/cart/clear \
  -H "Content-Type: application/json" \
  -d '{"user_id": "11111111-1111-1111-1111-111111111111", "product_id": "ignored"}'
```