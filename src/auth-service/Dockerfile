FROM python:3.11-slim

WORKDIR /app

# Install dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends gcc libpq-dev sqlite3 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY src/auth-service/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt && \
    pip install --no-cache-dir email-validator

COPY src/auth-service/app ./app
COPY src/auth-service/config ./config

# Create a non-root user
RUN adduser --disabled-password --gecos "" appuser
USER appuser

ENV PYTHONUNBUFFERED=1

EXPOSE 8080

ENTRYPOINT [ "opentelemetry-instrument", "python", "-m", "app.main" ] 