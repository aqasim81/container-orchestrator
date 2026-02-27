# Container Orchestrator

A simplified Kubernetes-style container orchestrator that deploys, scales, health-checks, and load-balances Docker containers, with a management API and Next.js dashboard.

## Status

Phase 1 complete — Go backend (config, store, API server, health endpoint), Next.js dashboard skeleton, and dev tooling all implemented and passing `make validate-all`. Next: Phase 2 (Container Lifecycle Management).

## Stack

| Layer | Technology | Version |
|-------|-----------|---------|
| Language (backend) | Go | 1.23 |
| HTTP Router | Chi | 5.1.0 |
| Docker Integration | Docker SDK | 27.4.1 |
| State Store | bbolt | 1.3.11 |
| Logging | zerolog | 1.33.0 |
| Config | caarlos0/env | 11.3.1 |
| Language (dashboard) | TypeScript | 5.7.3 |
| Framework (dashboard) | Next.js | 15.1.6 |
| Styling | Tailwind CSS | 4.0.6 |
| Data Fetching | TanStack Query | 5.66.0 |
| Charts | Recharts | 2.15.1 |
| Linting (Go) | golangci-lint | latest |
| Linting (TS) | Biome | 1.9.4 |
| Testing (Go) | go test + testify | 1.10.0 |
| Testing (TS) | Vitest | 3.0.5 |

## Directory Structure

```
container_orchestrator/
├── cmd/orchestrator/main.go       # Entry point
├── internal/
│   ├── config/                    # Env validation
│   ├── api/                       # Router, handlers, middleware
│   ├── container/                 # Container lifecycle (Docker SDK)
│   ├── node/                      # Node management, heartbeat
│   ├── scheduler/                 # Pluggable scheduling strategies
│   ├── health/                    # Health checks, auto-restart
│   ├── service/                   # Service discovery, load balancing
│   ├── deployment/                # Deployments, reconciliation
│   └── store/                     # State persistence (bbolt)
├── pkg/api/                       # Shared API types
├── dashboard/                     # Next.js management dashboard
│   ├── src/app/                   # Pages (App Router)
│   ├── src/components/            # UI components
│   ├── src/lib/                   # API client, types, env
│   └── src/hooks/                 # Data fetching hooks
├── tests/
│   ├── integration/               # API + scheduler integration tests
│   └── e2e/                       # Full cluster E2E tests
├── scripts/                       # Dev, seed, demo scripts
├── Makefile                       # Build, test, lint, validate
└── .env.example                   # Environment template
```

## Commands

### Go Backend

| Command | Description |
|---------|-------------|
| `make build` | Compile binary to `bin/orchestrator` |
| `make run` | Build and run the API server |
| `make test` | Run all unit tests with race detection |
| `make test-coverage` | Run tests with coverage report |
| `make test-integration` | Run integration tests |
| `make test-e2e` | Run E2E tests (requires Docker) |
| `make fmt` | Format Go code |
| `make fmt-check` | Check formatting without modifying |
| `make lint` | Run golangci-lint |
| `make vet` | Run go vet |
| `make validate` | Run fmt-check + vet + lint + test (all Go checks) |

### Dashboard

| Command | Description |
|---------|-------------|
| `make dashboard-install` | Install dashboard dependencies |
| `make dashboard-dev` | Start Next.js dev server |
| `make dashboard-build` | Production build |
| `make dashboard-lint` | Run Biome linter |
| `make dashboard-test` | Run Vitest |

### Combined

| Command | Description |
|---------|-------------|
| `make dev` | Start both API server and dashboard (concurrent) |
| `make validate-all` | Run all Go + Dashboard checks |
| `make clean` | Remove build artifacts |

## Architecture

### Core Patterns

- **Manager pattern:** Each domain (container, node, service, deployment) has a manager that owns business logic and coordinates between runtime and store.
- **Store interface:** All state persistence goes through a `Store` interface. bbolt for production, in-memory for tests.
- **Strategy pattern:** Scheduler and load balancer use pluggable strategies (bin-packing/spread, round-robin/least-connections).
- **Reconciliation loop:** Deployment controller runs a periodic loop comparing desired state to actual state and taking corrective action.
- **Background goroutines:** Heartbeat monitor, health checker, and reconciler run as managed goroutines with graceful shutdown.

### Data Flow

1. User creates deployment via API/dashboard
2. Deployment controller generates container specs
3. Scheduler selects target node (filter → score → select)
4. Container manager creates container via Docker SDK
5. Health checker monitors container health
6. Service discovery updates endpoints
7. Load balancer distributes traffic to healthy endpoints

## Coding Conventions

### Go

- Follow standard Go project layout (`cmd/`, `internal/`, `pkg/`)
- All exported types and functions have doc comments
- Error handling: return `error`, never panic in business logic
- Use `context.Context` for cancellation propagation
- Interfaces are defined by the consumer, not the producer
- Table-driven tests with `testify/assert` and `testify/require`
- No global state — all dependencies injected via constructors
- Use `zerolog` for structured logging, never `fmt.Println`

### TypeScript (Dashboard)

- Strict mode, zero `any` types
- Functional components only
- Props types defined inline or in same file
- Use `fetch` with typed wrappers, no Axios
- Zod for runtime validation of API responses
- File naming: `kebab-case.tsx` for components, `kebab-case.ts` for utilities

## Error Handling

### Go Backend

- Functions return `(result, error)` — callers must handle errors
- API handlers return structured JSON errors: `{"error": "message", "code": "NOT_FOUND"}`
- Use sentinel errors for known conditions: `ErrNotFound`, `ErrAlreadyExists`, `ErrInvalidState`
- Wrap errors with context: `fmt.Errorf("failed to start container %s: %w", id, err)`
- HTTP status codes: 400 (validation), 404 (not found), 409 (conflict/invalid state), 500 (internal)

### Dashboard

- API client wraps errors with typed error responses
- Components use error boundaries for unexpected failures
- TanStack Query handles retry logic for transient errors

## Testing

### Strategy

- **Unit tests:** All business logic in `internal/` packages. Use in-memory store and mock runtime.
- **Integration tests:** Test API endpoints with real router and in-memory store. Tagged `//go:build integration`.
- **E2E tests:** Full cluster operations with real Docker. Tagged `//go:build e2e`.
- **Dashboard tests:** Vitest + Testing Library for component tests.

### Coverage Targets

- 80%+ on business logic: scheduler, health checker, reconciler, managers
- 70%+ on API handlers
- 70%+ on dashboard components

### Test File Naming

- Go: `*_test.go` in the same package
- Dashboard: `*.test.ts` / `*.test.tsx` next to the source file

## Security

- **Env vars:** All configuration via environment variables. Validated at startup.
- **Never commit:** `.env`, API keys, credentials, `.pem`/`.key` files
- **Secret scanning:** gitleaks runs in pre-commit hook
- **API auth:** API key passed via `X-API-Key` header
- **Input validation:** All API request bodies validated before processing
- **Docker socket:** Controlled access, no privileged containers

## Git Workflow

- **Commits:** `feat:`, `fix:`, `chore:`, `docs:`, `test:`, `refactor:` prefixes
- **Branch naming:** `feat/`, `fix/`, `chore/` + descriptive slug
- **One commit per logical change** — don't mix features and fixes
- **GitHub profile:** `github-builder`

## Anti-Patterns

1. **No `any` types** in TypeScript — use proper types or `unknown`
2. **No `fmt.Println`** — use zerolog for all logging
3. **No global mutable state** — inject dependencies via constructors
4. **No bare goroutines** — always pass context, always handle shutdown
5. **No ignored errors** — handle or explicitly document why it's safe to ignore
6. **No dead code** — delete unused functions, variables, and imports
7. **No `console.log`** — remove debug logs before committing
8. **No hardcoded URLs/ports** — use config/env vars
9. **No raw SQL** — this project uses bbolt (key-value), not relational DB
10. **No panic in business logic** — return errors, let the caller decide

## Session Workflow

1. Read this `CLAUDE.md` file
2. Check `plans/checklist.md` for current phase status
3. Read the relevant phase file in `plans/phases/`
4. Implement the phase deliverables
5. Run `make validate` (Go) or `make validate-all` (Go + Dashboard)
6. Fix any issues
7. Commit with conventional commit message
8. Update `plans/checklist.md` and this file's Status section

## References

- [PRD](plans/prd.md)
- [Implementation Plan](plans/implementation_plan.md)
- [Checklist](plans/checklist.md)
- [Phase Files](plans/phases/)

---

> **No AI attribution** — Never mention Claude, Anthropic, AI-generated, AI-assisted, or any AI tool names in code, comments, commits, README, or documentation. Exception: model ID strings in SDK calls.
