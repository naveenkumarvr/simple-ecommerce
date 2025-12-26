# frontend-nodejs

Very simple Node.js + Express frontend for the `simple-ecommerce` demo.

## Features

- Plain HTML pages (no React)
- Pages:
  - `/login` – username + password form, calls backend login
  - `/home` – shows 4 products (car, bike, bus, truck) from catalog-service via API gateway
  - `/cart` – shows current cart items
  - `/checkout` – triggers order checkout
  - `/success` – simple success page
- Uses `fetch` to call backend APIs via the API gateway.
- Login is intentionally simple; any credentials are accepted as long as the backend responds successfully.

By default, it expects the API gateway to be accessible at `http://localhost:8000/api`.

## Prerequisites

- Node.js 18+ recommended
- The backend services and `api-gateway-go` running locally or via Docker/Kubernetes.

## Installation (local dev)

```bash
cd services/frontend-nodejs
npm install
```

## Run locally

```bash
cd services/frontend-nodejs
npm start
```

This starts the frontend on <http://localhost:3000>.

If your API gateway is not running on `http://localhost:8000/api`, set the `API_BASE` environment variable before starting, for example in PowerShell:

```powershell
cd services/frontend-nodejs
$env:API_BASE = "http://my-gateway-host:8000/api"
npm start
```

Or in bash:

```bash
cd services/frontend-nodejs
API_BASE="http://my-gateway-host:8000/api" npm start
```

## Docker

A simple Dockerfile is provided in this directory.

### Build image

```bash
cd services/frontend-nodejs
docker build -t frontend-nodejs .
```

### Run container

Assuming your API gateway is reachable on `http://localhost:8000/api` from the host, you can run:

```bash
docker run --rm -p 3000:3000 -e API_BASE="http://host.docker.internal:8000/api" frontend-nodejs
```

Then open <http://localhost:3000> in your browser.

If everything is running in a Docker or Kubernetes network where the gateway is known as `api-gateway-go` on port `8000`, you can rely on the default `API_BASE` baked into the image (`http://api-gateway-go:8000/api`):

```bash
docker run --rm -p 3000:3000 frontend-nodejs
```

## How it talks to the backend

- `POST /api/login` → `user-service-python`
- `GET /api/products` → `catalog-service-go`
- `POST /api/cart/add` and `GET /api/cart/{user_id}` → `cart-service-python`
- `POST /api/order` → `order-service-go`

These are all routed through `api-gateway-go`.

This is intentionally very simple and uses in-browser `localStorage` to remember the logged-in user.
