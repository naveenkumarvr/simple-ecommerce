from typing import Dict, Optional
import json
from pathlib import Path

from .models import User

# In-memory store: user_id -> User
_USERS: Dict[str, User] = {}
_USERNAME_INDEX: Dict[str, User] = {}
_PASSWORDS: Dict[str, str] = {}


def _load_users_from_file() -> None:
    """Load users from the static JSON file into memory.

    The file is expected to live at ../users.json relative to this module.
    """
    global _USERS, _USERNAME_INDEX, _PASSWORDS

    if _USERS:
        # Already loaded
        return

    users_file = Path(__file__).resolve().parent.parent / "users.json"
    if not users_file.exists():
        # No file found; leave stores empty
        return

    with users_file.open("r", encoding="utf-8") as f:
        data = json.load(f)

    for entry in data:
        user = User(
            id=entry["id"],
            username=entry["username"],
            full_name=entry.get("full_name"),
        )
        _USERS[user.id] = user
        _USERNAME_INDEX[user.username] = user
        _PASSWORDS[user.username] = entry.get("password", "")


def authenticate(username: str, password: str) -> Optional[User]:
    """Return the user if username/password matches one of the dummy users."""
    _load_users_from_file()

    user = _USERNAME_INDEX.get(username)
    if not user:
        return None

    expected_password = _PASSWORDS.get(username)
    if expected_password != password:
        return None

    return user


def get_user(user_id: str) -> Optional[User]:
    """Return a user by ID if found, else None."""
    _load_users_from_file()
    return _USERS.get(user_id)
