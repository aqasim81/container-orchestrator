.PHONY: build run test test-coverage lint fmt vet validate clean dev dashboard

# === Build ===

build:
	go build -o bin/orchestrator ./cmd/orchestrator

run: build
	./bin/orchestrator

# === Testing ===

test:
	go test ./... -race -count=1

test-coverage:
	go test ./... -race -count=1 -coverprofile=coverage.out -covermode=atomic
	go tool cover -func=coverage.out

test-integration:
	go test ./tests/integration/... -race -count=1 -tags=integration

test-e2e:
	go test ./tests/e2e/... -race -count=1 -tags=e2e

# === Code Quality ===

fmt:
	gofmt -w .
	goimports -w -local github.com/github-builder/container-orchestrator .

fmt-check:
	@test -z "$$(gofmt -l .)" || (echo "Files need formatting:" && gofmt -l . && exit 1)

lint:
	golangci-lint run ./...

vet:
	go vet ./...

# === Dashboard ===

dashboard-install:
	cd dashboard && npm install

dashboard-dev:
	cd dashboard && npm run dev

dashboard-build:
	cd dashboard && npm run build

dashboard-lint:
	cd dashboard && npx @biomejs/biome check ./src

dashboard-test:
	cd dashboard && npx vitest run

# === Development ===

dev:
	./scripts/dev.sh

seed:
	./scripts/seed.sh

# === Validation (runs ALL checks) ===

validate: fmt-check vet lint test
	@echo "All checks passed."

validate-all: fmt-check vet lint test dashboard-lint dashboard-test
	@echo "All checks (Go + Dashboard) passed."

# === Cleanup ===

clean:
	rm -rf bin/ coverage.out data/
