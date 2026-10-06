# FlowForge — Workflow Automation Engine

> A production-grade workflow automation platform built in Go. Design, execute, and monitor multi-step automated workflows using a visual DAG-based builder with real-time execution tracking.

[![Go](https://img.shields.io/badge/backend-Go-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![React](https://img.shields.io/badge/frontend-React-61DAFB?logo=react&logoColor=black)](https://react.dev/)
[![PostgreSQL](https://img.shields.io/badge/database-PostgreSQL-4169E1?logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/queue-Redis-DC382D?logo=redis&logoColor=white)](https://redis.io/)
[![Docker](https://img.shields.io/badge/deployment-Docker-2496ED?logo=docker&logoColor=white)](https://www.docker.com/)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

---

## What is FlowForge?

FlowForge lets users visually design workflows as **Directed Acyclic Graphs (DAGs)** — where nodes are actions (HTTP requests, conditions, delays, email sends, AI calls) and edges define execution order and branching logic.

Workflows are queued into Redis, picked up by a worker pool, and executed concurrently where possible — with live updates streamed to the frontend via **Server-Sent Events (SSE)**.

---

## Key Technical Features

| Feature | Implementation |
|---|---|
| **DAG Execution Engine** | Kahn's algorithm for topological sort, concurrent branch execution via goroutines |
| **Cycle Detection** | Validated before every execution attempt |
| **Concurrent Branches** | Independent branches run in parallel using `sync.WaitGroup` + channels |
| **Retry + Backoff** | Per-node configurable retry with exponential backoff |
| **Redis Worker Pool** | `BLPOP`-based blocking queue, goroutine worker pool |
| **Workflow Scheduling** | Cron-based scheduler with `next_run_at` tracking in PostgreSQL |
| **SSE Real-time Updates** | Per-workflow pub/sub fan-out using Go channels |
| **JWT Auth** | Access token (15min) + Refresh token (7d) stored in Redis |
| **Webhook Triggers** | Unique URL per workflow, payload passed as initial node input |
| **Execution History** | Full node-level run logs with input/output/status/timestamps |

---

## Architecture

```
React Frontend (Vite + React Flow)
         |
         v
   Go REST API (Gin)
         |
   ┌─────┼──────┐
   |     |      |
 Auth  Workflows Executions
         |
         v
    PostgreSQL
         |
         v
       Redis
         |
         v
    Worker Pool
         |
         v
   DAG Execution Engine
         |
   ┌─────┼─────┐
   v     v     v
 Node A Node B Node C   ← concurrent branches
         |
         v
   Execution Results
         |
         v
    SSE → React UI
```

---

## Node Types

| Node | Purpose |
|---|---|
| `webhook_trigger` | HTTP webhook entry point |
| `http_request` | Outbound HTTP call with headers/body |
| `condition` | Boolean branch (`true`/`false` edges) |
| `delay` | Fixed time pause |
| `email` | SMTP email send |
| `ai` | Gemini AI prompt (optional, per-user API key) |

---

## DAG Engine Design

The core execution engine (`internal/workflow/engine`) solves the classic problem of executing a workflow graph where independent nodes should run **in parallel** and dependent nodes wait for their parents.

### Example: Diamond DAG
```
    A
   / \
  B   C   ← B and C execute concurrently after A
   \ /
    D     ← D waits for BOTH B and C
```

### Algorithm
1. **Kahn's topological sort** validates the graph and detects cycles
2. Root nodes (in-degree = 0) are seeded into a `readyChan`
3. Each ready node is dispatched as a goroutine
4. On completion, the goroutine decrements in-degree of downstream nodes
5. When a node's in-degree reaches 0, it is pushed to `readyChan`
6. Condition nodes evaluate output JSON to pick `true`/`false` branches

### Retry with Exponential Backoff
Each node can have per-config retry settings:
```json
{
  "retry": {
    "enabled": true,
    "max_attempts": 3,
    "delay": "500ms",
    "backoff_exponent": 2.0
  }
}
```

---

## Project Structure

```
.
├── services/
│   └── api-core/
│       ├── cmd/server/         # Entry point
│       └── internal/
│           ├── auth/           # JWT, bcrypt, email verification
│           ├── db/             # SQLC-generated PostgreSQL queries
│           ├── queue/          # Redis queue (BLPOP)
│           ├── scheduler/      # Cron scheduler
│           ├── worker/         # Worker pool
│           ├── workflow/
│           │   ├── engine/     # DAG execution engine
│           │   ├── executor/   # Node type handlers
│           │   ├── handler/    # HTTP handlers + SSE
│           │   ├── repository/ # DB layer
│           │   └── service/    # Business logic
│           └── router/         # Route registration
├── web/                        # React + React Flow frontend
├── sql/                        # SQLC queries + schema
├── migrations/                 # DB migrations
└── docker-compose.yml
```

---

## API Endpoints

### Auth
```
POST /signup
POST /login
POST /logout
POST /refresh
GET  /me
GET  /verify-email
POST /forgot-password
POST /reset-password
```

### Workflows
```
POST   /workflows
GET    /workflows
GET    /workflows/:id/steps
PUT    /workflows/:id/steps
DELETE /workflows/:id
PUT    /workflows/:id
POST   /workflows/:id/run
GET    /workflows/:id/runs
GET    /workflow-runs/:id/steps
GET    /workflows/:id/events     ← SSE stream
PUT    /workflows/:id/schedule
GET    /workflows/:id/schedule
POST   /webhooks/:webhookID
```

---

## Local Development

### Prerequisites
- Docker + Docker Compose
- Go 1.22+
- Node.js 18+

### Setup
```bash
# Clone the repo
git clone https://github.com/harsh2102004/ridegolang.git
cd ridegolang

# Start PostgreSQL + Redis
docker compose up postgres redis -d

# Set up env
cp services/api-core/.env.example services/api-core/.env
# Edit .env with your DB, Redis, SMTP, JWT_SECRET values

# Run backend
cd services/api-core
go run ./cmd/server

# Run frontend
cd web
npm install
npm run dev
```

### Docker (Full Stack)
```bash
docker compose up --build
```
Frontend: http://localhost:5173  
Backend API: http://localhost:8080

---

## Database Schema

Key tables:
- `users` — Auth with bcrypt passwords, email verification
- `workflows` — Workflow metadata + scheduling config
- `workflow_steps` — Nodes with type, config (JSONB), position
- `workflow_edges` — DAG edges with optional condition branches
- `workflow_runs` — Execution records with status + timestamps
- `workflow_step_runs` — Per-node logs with input/output/error

---

## Testing

```bash
cd services/api-core

# DAG engine tests (cycle detection, concurrent branches, retry)
go test ./internal/workflow/engine -v

# All packages
go test ./...
```

---

## Tech Stack

**Backend**: Go · Gin · SQLC · PostgreSQL · Redis  
**Frontend**: React · Vite · React Flow (`@xyflow/react`) · Vanilla CSS  
**Auth**: JWT (golang-jwt) · bcrypt · Redis token store  
**Infra**: Docker · Docker Compose  
**Other**: SSE · Cron · SMTP · Goroutines · sync.WaitGroup

---

## Attribution

This project is substantially modified and extended from [ORG](https://github.com/shinasorg/org) (MIT License). The original project provided inspiration for the workflow automation concept. The DAG execution engine, simplified architecture, authentication system, worker pool, and frontend have been redesigned and rebuilt.

---

## License

MIT — see [LICENSE](LICENSE)
