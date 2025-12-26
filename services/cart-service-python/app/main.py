from fastapi import FastAPI

from .models import CartAddRequest, CartResponse
from . import storage

app = FastAPI(
    title="Cart Service (Python)",
    description="Very simple in-memory cart service for Kubernetes demo.",
    version="1.0.0",
)


@app.post("/cart/add", response_model=CartResponse)
async def add_to_cart(payload: CartAddRequest) -> CartResponse:
    """Add a product to the user's cart (no validation of user/product here)."""
    return storage.add_to_cart(payload.user_id, payload.product_id)


@app.get("/cart/{user_id}", response_model=CartResponse)
async def get_cart(user_id: str) -> CartResponse:
    return storage.get_cart(user_id)


@app.post("/cart/clear")
async def clear_cart(payload: CartAddRequest) -> dict:
    # We only use user_id from the payload for clearing.
    storage.clear_cart(payload.user_id)
    return {"status": "cleared", "user_id": payload.user_id}


@app.get("/health")
async def health() -> dict:
    return {"status": "ok"}
