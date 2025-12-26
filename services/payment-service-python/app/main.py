from fastapi import FastAPI

from .models import PaymentRequest, PaymentResponse, UserPaymentsResponse
from . import storage

app = FastAPI(
    title="Payment Service (Python)",
    description="Very simple in-memory payment service for Kubernetes demo.",
    version="1.0.0",
)


@app.post("/pay", response_model=PaymentResponse)
async def pay(req: PaymentRequest) -> PaymentResponse:
    """Create a fake payment and return its id and status. Always succeeds."""
    payment = storage.create_payment(req)
    return PaymentResponse(payment_id=payment.id, status=payment.status)


@app.get("/payments/{user_id}", response_model=UserPaymentsResponse)
async def get_payments(user_id: str) -> UserPaymentsResponse:
    """Return all payments for a given user."""
    return storage.get_payments_for_user(user_id)


@app.get("/health")
async def health() -> dict:
    return {"status": "ok"}
