from pydantic import BaseModel
from typing import List


class CartAddRequest(BaseModel):
    user_id: str
    product_id: str


class CartItem(BaseModel):
    product_id: str
    quantity: int


class CartResponse(BaseModel):
    user_id: str
    items: List[CartItem]
