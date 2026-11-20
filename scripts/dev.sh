#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"

export USE_MEMORY=true
export HTTP_ADDR="${HTTP_ADDR:-:8080}"
export TICK_INTERVAL_MS="${TICK_INTERVAL_MS:-2000}"
export CORS_ORIGIN="${CORS_ORIGIN:-*}"

echo "→ starting game-api on ${HTTP_ADDR}"
(cd "$ROOT/backend" && go run ./cmd/server) &
API_PID=$!

cleanup() {
  kill "$API_PID" 2>/dev/null || true
}
trap cleanup EXIT

sleep 1
echo "→ starting web on :5173"
cd "$ROOT/web"
npm run dev
