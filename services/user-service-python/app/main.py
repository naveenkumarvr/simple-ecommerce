from fastapi import FastAPI, HTTPException

from .models import LoginRequest, LoginResponse, User
from . import storage

app = FastAPI(
    title="User Service (Python)",
    description="Very simple in-memory user service for Kubernetes demo.",
    version="1.0.0",
)


@app.post("/login", response_model=LoginResponse)
async def login(payload: LoginRequest) -> LoginResponse:
    """Validate login credentials against users defined in users.json.

    If username/password match one of the dummy users, return its user_id.
    Otherwise return 401.
    """
    user = storage.authenticate(payload.username, payload.password)
    if not user:
        raise HTTPException(status_code=401, detail="Invalid username or password")
    return LoginResponse(user_id=user.id)


@app.get("/user/{user_id}", response_model=User)
async def get_user(user_id: str) -> User:
    user = storage.get_user(user_id)
    if not user:
        raise HTTPException(status_code=404, detail="User not found")
    return user


@app.get("/health")
async def health() -> dict:
    return {"status": "ok"}
