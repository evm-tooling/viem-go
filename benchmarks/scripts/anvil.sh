#!/bin/bash
#
# Anvil Management Script for Cross-Language Benchmarks
#
# This script manages the lifecycle of an Anvil instance for benchmarking.
# It starts Anvil with a mainnet fork, waits for readiness, runs the
# provided command(s), and cleans up on exit.
#
# Usage: ./anvil.sh <command>
# Example: ./anvil.sh make _bench-sequential

set -e

# Configuration (can be overridden via environment variables)
ANVIL_PORT=${ANVIL_PORT:-8545}
FORK_URL=${FORK_URL:-"https://eth.drpc.org"}
FORK_BLOCK=${FORK_BLOCK:-}
ANVIL_TIMEOUT=${ANVIL_TIMEOUT:-30}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if anvil is installed
if ! command -v anvil &> /dev/null; then
    log_error "Anvil not found. Please install Foundry: https://getfoundry.sh"
    exit 1
fi

# Check if port is already in use
if lsof -Pi :$ANVIL_PORT -sTCP:LISTEN -t >/dev/null 2>&1; then
    log_warn "Port $ANVIL_PORT is already in use. Attempting to use existing Anvil instance."
    export ANVIL_RPC_URL="http://127.0.0.1:$ANVIL_PORT"
    
    # Verify it's responsive
    if curl -s -X POST -H "Content-Type: application/json" \
        --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
        "$ANVIL_RPC_URL" > /dev/null 2>&1; then
        log_info "Using existing Anvil instance at $ANVIL_RPC_URL"
        exec "$@"
    else
        log_error "Port $ANVIL_PORT is in use but not responding to RPC calls"
        exit 1
    fi
fi

log_info "Starting Anvil..."
log_info "  Port: $ANVIL_PORT"
log_info "  Fork URL: $FORK_URL"
if [ -n "$FORK_BLOCK" ]; then
    log_info "  Fork Block: $FORK_BLOCK"
else
    log_info "  Fork Block: latest"
fi

# Build anvil command
ANVIL_CMD="anvil --port $ANVIL_PORT --fork-url $FORK_URL --no-mining --accounts 10 --balance 10000 --silent"
if [ -n "$FORK_BLOCK" ]; then
    ANVIL_CMD="$ANVIL_CMD --fork-block-number $FORK_BLOCK"
fi

# Start anvil in background
$ANVIL_CMD &

ANVIL_PID=$!

# Cleanup function
cleanup() {
    log_info "Stopping Anvil (PID: $ANVIL_PID)..."
    kill $ANVIL_PID 2>/dev/null || true
    wait $ANVIL_PID 2>/dev/null || true
    log_info "Anvil stopped."
}

# Register cleanup on script exit
trap cleanup EXIT INT TERM

# Wait for Anvil to be ready
log_info "Waiting for Anvil to be ready..."
WAIT_START=$(date +%s)

while true; do
    if curl -s -X POST -H "Content-Type: application/json" \
        --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
        "http://127.0.0.1:$ANVIL_PORT" > /dev/null 2>&1; then
        break
    fi
    
    ELAPSED=$(($(date +%s) - WAIT_START))
    if [ $ELAPSED -ge $ANVIL_TIMEOUT ]; then
        log_error "Timeout waiting for Anvil to start"
        exit 1
    fi
    
    sleep 0.1
done

READY_TIME=$(($(date +%s) - WAIT_START))
log_info "Anvil ready in ${READY_TIME}s on port $ANVIL_PORT"

# Export RPC URL for child processes
export ANVIL_RPC_URL="http://127.0.0.1:$ANVIL_PORT"

# Verify fork block
BLOCK_NUM=$(curl -s -X POST -H "Content-Type: application/json" \
    --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
    "$ANVIL_RPC_URL" | grep -o '"result":"[^"]*"' | cut -d'"' -f4)
log_info "Current block: $BLOCK_NUM"

# Run the provided command
if [ $# -gt 0 ]; then
    log_info "Running: $@"
    echo ""
    "$@"
    EXIT_CODE=$?
    echo ""
    log_info "Command completed with exit code: $EXIT_CODE"
    exit $EXIT_CODE
else
    log_info "No command provided. Anvil running at $ANVIL_RPC_URL"
    log_info "Press Ctrl+C to stop."
    wait $ANVIL_PID
fi
