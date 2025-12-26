from datetime import datetime
from typing import List, Optional

from pydantic import BaseModel, Field


class PaymentRequest(BaseModel):
    user_id: str
    amount: float = Field(gt=0, description="Payment amount must be > 0")
    currency: str = "USD"
    source: Optional[str] = None  # e.g. fake card/token id


class Payment(BaseModel):
    id: str
    user_id: str
    amount: float
    currency: str
    source: Optional[str] = None
    status: str  # e.g. "success"
    created_at: datetime


class PaymentResponse(BaseModel):
    payment_id: str
    status: str


class UserPaymentsResponse(BaseModel):
    user_id: str
    payments: List[Payment]
