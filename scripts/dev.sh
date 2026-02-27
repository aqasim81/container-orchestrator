#!/usr/bin/env bash
set -euo pipefail

# Start both the Go API server and the Next.js dashboard concurrently.
# Kills both processes on exit.

cleanup() {
    echo ""
    echo "Shutting down..."
    kill 0 2>/dev/null || true
    wait 2>/dev/null || true
    echo "Done."
}

trap cleanup EXIT INT TERM

echo "Starting API server and Dashboard..."
echo ""

# Start Go API server.
make run &

# Start Next.js dashboard.
make dashboard-dev &

echo ""
echo "API server:  http://localhost:${ORCHESTRATOR_PORT:-8080}"
echo "Dashboard:   http://localhost:3000"
echo ""
echo "Press Ctrl+C to stop both."
echo ""

# Wait for either process to exit.
wait
