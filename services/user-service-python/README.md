# user-service-python

Very simple in-memory user service using FastAPI, intended for Kubernetes demos.

## Users

Users are defined in a static JSON file `users.json` in this directory. Example content:

```json
[
  {
    "id": "11111111-1111-1111-1111-111111111111",
    "username": "alice",
    "password": "password1",
    "full_name": "Alice Example"
  },
  {
    "id": "22222222-2222-2222-2222-222222222222",
    "username": "bob",
    "password": "password2",
    "full_name": "Bob Example"
  }
]
```

The service loads this file at startup (on first request) and keeps the data in memory.

## Endpoints

- `POST /login` – validates `username` and `password` against `users.json` and returns the corresponding `user_id`.
- `GET /user/{user_id}` – returns user details from in-memory storage (backed by `users.json`) or 404.
- `GET /health` – basic health check.

## Running locally

From this directory:

```bash
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
uvicorn app.main:app --host 0.0.0.0 --port 8001 --reload
```

Then call (valid user):

```bash
curl -X POST http://localhost:8001/login \
  -H "Content-Type: application/json" \
  -d '{"username": "alice", "password": "password1"}'
```

You should see a response like:

```json
{"user_id": "11111111-1111-1111-1111-111111111111"}
```

Then:

```bash
curl http://localhost:8001/user/11111111-1111-1111-1111-111111111111
```

For an invalid login (wrong password):

```bash
curl -X POST http://localhost:8001/login \
  -H "Content-Type: application/json" \
  -d '{"username": "alice", "password": "wrong"}'
```

The service responds with HTTP 401 and:

```json
{"detail": "Invalid username or password"}
```
