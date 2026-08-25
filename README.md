# Neighbor Relay

Neighbor Relay helps community field workers spend more time visiting households by coordinating case follow-up, visit execution and review evidence. It is a production-shaped Go backend with SQLite persistence, migrations, sessions, role checks, audit trails and a resumable outbox worker.

## Run

```bash
go run ./cmd/server
curl http://localhost:8080/healthz
```

Default demo users are created by tests only; production deployments create users through the repository/service API. Configuration uses `ADDR`, `DB_PATH`, `SESSION_TTL`, and `WORKER_INTERVAL`.

## API

`POST /v1/auth/login`, `POST /v1/auth/logout`, `GET /v1/households`, `POST /v1/households`, `POST /v1/cases`, `POST /v1/visits`, `POST /v1/visits/{id}/start`, `POST /v1/visits/{id}/complete`, `GET /healthz`, and `GET /readyz`.

## Design

The service stores users, revocable sessions, households, cases, visits, notes, audit events, idempotency keys and outbox jobs. A visit transaction validates the household/case relationship, writes the state transition, audit record and notification job atomically. A worker retries jobs with bounded exponential backoff and recovers in-flight work after restart.
