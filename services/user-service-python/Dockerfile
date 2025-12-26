FROM python:3.11-slim

WORKDIR /app

# Now copy the rest of the app (including app/ and users.json)
COPY . /app/

RUN pip install --no-cache-dir -r requirements.txt

EXPOSE 8001

ENTRYPOINT ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8001"]