from pydantic import BaseModel
from typing import Optional


class LoginRequest(BaseModel):
    username: str
    password: str  # Not actually validated in this demo


class LoginResponse(BaseModel):
    user_id: str


class User(BaseModel):
    id: str
    username: str
    full_name: Optional[str] = None
