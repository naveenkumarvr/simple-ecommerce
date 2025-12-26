from typing import Dict, List

from .models import CartItem, CartResponse

# In-memory cart: user_id -> list of CartItem
_CARTS: Dict[str, List[CartItem]] = {}


def add_to_cart(user_id: str, product_id: str) -> CartResponse:
    """Add a product to the user's cart (quantity +1 each time)."""
    items = _CARTS.setdefault(user_id, [])

    for item in items:
        if item.product_id == product_id:
            item.quantity += 1
            break
    else:
        items.append(CartItem(product_id=product_id, quantity=1))

    return CartResponse(user_id=user_id, items=items)


def get_cart(user_id: str) -> CartResponse:
    items = _CARTS.get(user_id, [])
    return CartResponse(user_id=user_id, items=items)


def clear_cart(user_id: str) -> None:
    _CARTS.pop(user_id, None)
