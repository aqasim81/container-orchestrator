# Container Orchestrator

A simplified Kubernetes-style container orchestrator that deploys, scales, health-checks, and load-balances Docker containers, with a management API and Next.js dashboard.

## Status

Phase 1 complete — Go backend (config, store, API server, health endpoint), Next.js dashboard skeleton, and dev tooling all implemented and passing `make validate-all`. Next: Phase 2 (Container Lifecycle Management).

## Stack

- **Backend:** Go 1.23, Chi 5.1, Docker SDK 27.4, bbolt 1.3, zerolog 1.33, caarlos0/env 11.3
- **Dashboard:** TypeScript 5.7, Next.js 15.1, Tailwind CSS 4.0, TanStack Query 5.66, Recharts 2.15
- **Linting:** golangci-lint (Go), Biome 1.9 (TS)
- **Testing:** go test + testify 1.10 (Go), Vitest 3.0 (TS)

## Key Commands

- `make validate-all` — Run all Go + Dashboard checks (fmt, vet, lint, test)
- `make validate` — Go checks only
- `make test` / `make test-coverage` — Unit tests (with optional coverage)
- `make test-integration` / `make test-e2e` — Integration / E2E tests (E2E requires Docker)
- `make build` / `make run` — Compile to `bin/orchestrator` / build and run
- `make dashboard-dev` / `make dashboard-build` / `make dashboard-test` — Dashboard dev/build/test
- `make dev` — Start both API server and dashboard concurrently

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

- Standard Go project layout (`cmd/`, `internal/`, `pkg/`)
- All exported types and functions have doc comments
- Error handling: return `(result, error)`, never panic in business logic
- Use `context.Context` for cancellation propagation
- Interfaces defined by consumer, not producer
- Table-driven tests with `testify/assert` and `testify/require`
- No global state — all dependencies injected via constructors
- Use `zerolog` for structured logging, never `fmt.Println`
- No bare goroutines — always pass context, always handle shutdown
- No dead code — delete unused functions, variables, and imports
- API errors as structured JSON: `{"error": "message", "code": "NOT_FOUND"}`
- Sentinel errors: `ErrNotFound`, `ErrAlreadyExists`, `ErrInvalidState`
- Wrap errors: `fmt.Errorf("failed to start container %s: %w", id, err)`
- HTTP status codes: 400 (validation), 404 (not found), 409 (conflict), 500 (internal)

### TypeScript (Dashboard)

- Strict mode, zero `any` types
- Functional components only
- Props types defined inline or in same file
- Use `fetch` with typed wrappers, no Axios
- Zod for runtime validation of API responses
- File naming: `kebab-case.tsx` for components, `kebab-case.ts` for utilities
- No `console.log` — remove debug logs before committing

## Testing

- **Unit tests:** All business logic in `internal/` packages. Use in-memory store and mock runtime.
- **Integration tests:** API endpoints with real router and in-memory store. Tagged `//go:build integration`.
- **E2E tests:** Full cluster operations with real Docker. Tagged `//go:build e2e`.
- **Coverage targets:** 80%+ business logic (scheduler, health checker, reconciler, managers), 70%+ API handlers, 70%+ dashboard components.

## Security

- All configuration via environment variables, validated at startup
- Never commit `.env`, API keys, credentials, `.pem`/`.key` files
- gitleaks runs in pre-commit hook
- API auth via `X-API-Key` header
- All API request bodies validated before processing
- Docker socket: controlled access, no privileged containers

## Git Workflow

- Conventional commits: `feat:`, `fix:`, `chore:`, `docs:`, `test:`, `refactor:`
- Branch naming: `feat/`, `fix/`, `chore/` + descriptive slug
- One commit per logical change — don't mix features and fixes
- GitHub profile: `github-builder`

## Session Workflow

1. Read this `CLAUDE.md`
2. Check `plans/checklist.md` for current phase status
3. Read the relevant phase file in `plans/phases/`
4. Implement phase deliverables
5. Run `make validate-all`, fix issues, commit
6. Update `plans/checklist.md` and this file's Status section

## References

[PRD](plans/prd.md) | [Implementation Plan](plans/implementation_plan.md) | [Checklist](plans/checklist.md) | [Phase Files](plans/phases/)

> **No AI attribution** — Never mention Claude, Anthropic, AI-generated, AI-assisted, or any AI tool names in code, comments, commits, README, or documentation.
