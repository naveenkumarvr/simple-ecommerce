from datetime import datetime
from typing import Dict, List
import uuid

from .models import Payment, PaymentRequest, UserPaymentsResponse

# In-memory store: user_id -> list of Payment
_PAYMENTS: Dict[str, List[Payment]] = {}


def create_payment(req: PaymentRequest) -> Payment:
    payment_id = str(uuid.uuid4())
    payment = Payment(
        id=payment_id,
        user_id=req.user_id,
        amount=req.amount,
        currency=req.currency,
        source=req.source,
        status="success",  # always successful in this demo
        created_at=datetime.utcnow(),
    )
    user_payments = _PAYMENTS.setdefault(req.user_id, [])
    user_payments.append(payment)
    return payment


def get_payments_for_user(user_id: str) -> UserPaymentsResponse:
    payments = _PAYMENTS.get(user_id, [])
    return UserPaymentsResponse(user_id=user_id, payments=payments)
