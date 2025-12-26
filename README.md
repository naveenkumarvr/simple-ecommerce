# simple-ecommerce

Simple multi-service e-commerce demo for learning Docker, Kubernetes, and microservices basics.

All services are intentionally very small, use in-memory storage only, and are designed to be easy to understand.

## Services

Each service lives under `services/` and is self-contained with its own `Dockerfile` and README:

1. `frontend-nodejs` (Node.js + Express)
   - Plain HTML pages (no React).
   - Routes: `/login`, `/home`, `/cart`, `/checkout`, `/success`.
   - Uses `fetch` to talk to the API gateway.

2. `api-gateway-go` (Go)
   - Simple reverse proxy.
   - Routes:
     - `/api/login` → user-service
     - `/api/products` → catalog-service
     - `/api/cart/*` → cart-service
     - `/api/order` → order-service

3. `user-service-python` (FastAPI)
   - `POST /login` – validate user against `users.json` and return `user_id`.
   - `GET /user/{id}`.

4. `catalog-service-go` (Go)
   - Hardcoded products: car, bike, bus, truck.
   - `GET /products`, `GET /products/{id}`.

5. `cart-service-python` (FastAPI)
   - `POST /cart/add` – add item to cart.
   - `GET /cart/{user_id}` – view cart.
   - `POST /cart/clear` – clear cart.

6. `payment-service-python` (FastAPI)
   - `POST /pay` – always returns a successful fake payment.
   - `GET /payments/{user_id}` – list payments.

7. `order-service-go` (Go)
   - `POST /order/checkout` – calls cart-service and payment-service and returns `"order successful"`.

> Note: There is also a `notification-service-go` folder reserved for future experiments; it is not wired into the current flow.

## Running everything with Docker Compose

The easiest way to run the full stack locally is with `docker compose` from the project root.

### Prerequisites

- Docker Desktop (or equivalent engine) installed and running.
- Optional: WSL2 on Windows if you prefer using bash.

### Start the stack

From the project root:

```bash
# From bash (WSL)
cd /mnt/d/code/simple-ecommerce
docker compose up --build
```

Or on Windows PowerShell (path adjusted as needed):

```powershell
cd D:\code\simple-ecommerce
docker compose up --build
```

This will build and start these containers:

- `frontend-nodejs` on port `3000`.
- `api-gateway-go` on port `8000`.
- `user-service-python` on port `8001`.
- `catalog-service-go` on port `8002`.
- `cart-service-python` on port `8003`.
- `payment-service-python` on port `8004`.
- `order-service-go` on port `8005`.

All containers share a Docker network where they can reach each other by service name, e.g. `http://user-service-python:8001`.

### Stop the stack

In the same directory:

```bash
docker compose down
```

## End-to-end flow

1. Open the frontend in your browser:
   - <http://localhost:3000>

2. Log in:
   - Go to `/login`.
   - Enter a username and password (for example, `alice` / `password1`).
   - The frontend calls `POST /api/login` on the gateway, which forwards to `user-service-python`.

3. View products:
   - After login, you are redirected to `/home`.
   - The frontend calls `GET /api/products` on the gateway, which forwards to `catalog-service-go`.

4. Add items to cart:
   - On `/home`, click "Add to cart" for a few products.
   - The frontend calls `POST /api/cart/add` on the gateway, which forwards to `cart-service-python`.

5. View cart:
   - Go to `/cart`.
   - The frontend calls `GET /api/cart/{user_id}` on the gateway to show your current cart.

6. Checkout:
   - Go to `/checkout` and click the checkout button.
   - The frontend calls `POST /api/order` on the gateway, which forwards to `order-service-go`.
   - `order-service-go` calls cart-service and payment-service, then returns a simple success response.

7. Success page:
   - You are redirected to `/success`, confirming the order was placed.

This setup is intentionally minimal to make it easy to reason about and to deploy into Kubernetes later (one Deployment/Service per folder under `services/`).
