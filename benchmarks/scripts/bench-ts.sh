#!/bin/bash
#
# TypeScript Benchmark Script with Anvil Management
#
# This script runs TypeScript (viem) benchmarks with automatic Anvil setup and teardown.
# It's a self-contained script that handles everything needed for TypeScript benchmarks.
#
# Usage: ./bench-ts.sh
#
# Configuration (via environment variables):
#   ANVIL_PORT  - Port for Anvil (default: 8545)
#   FORK_URL    - Ethereum RPC URL to fork from (default: https://eth.drpc.org)
#   FORK_BLOCK  - Block number to fork from (default: latest)

set -e

# Change to benchmarks directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.."

# Configuration
ANVIL_PORT=${ANVIL_PORT:-8545}
FORK_URL=${FORK_URL:-"https://eth.drpc.org"}
FORK_BLOCK=${FORK_BLOCK:-}
ANVIL_TIMEOUT=${ANVIL_TIMEOUT:-30}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
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

log_header() {
    echo ""
    echo -e "${BLUE}==================================================${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}==================================================${NC}"
    echo ""
}

# Check if anvil is installed
if ! command -v anvil &> /dev/null; then
    log_error "Anvil not found. Please install Foundry: https://getfoundry.sh"
    exit 1
fi

# Check if bun is installed (preferred) or npm
if command -v bun &> /dev/null; then
    PKG_MANAGER="bun"
    RUN_CMD="bun run"
elif command -v npm &> /dev/null; then
    PKG_MANAGER="npm"
    RUN_CMD="npm run"
else
    log_error "Neither bun nor npm found. Please install bun: https://bun.sh or Node.js: https://nodejs.org"
    exit 1
fi

log_header "TypeScript Benchmark Suite (viem)"

# Ensure results directory exists
mkdir -p results

# Check and install dependencies
if [ ! -d "typescript/node_modules" ]; then
    log_info "Installing TypeScript dependencies..."
    cd typescript
    if [ "$PKG_MANAGER" = "bun" ]; then
        bun install
    else
        npm install
    fi
    cd ..
fi

# Track if we started Anvil (so we know whether to stop it)
ANVIL_PID=""
STARTED_ANVIL=false

# Cleanup function
cleanup() {
    if [ "$STARTED_ANVIL" = true ] && [ -n "$ANVIL_PID" ]; then
        log_info "Stopping Anvil (PID: $ANVIL_PID)..."
        kill $ANVIL_PID 2>/dev/null || true
        wait $ANVIL_PID 2>/dev/null || true
        log_info "Anvil stopped."
    fi
}

# Register cleanup on script exit
trap cleanup EXIT INT TERM

# Check if Anvil is already running
if lsof -Pi :$ANVIL_PORT -sTCP:LISTEN -t >/dev/null 2>&1; then
    log_warn "Port $ANVIL_PORT is already in use. Attempting to use existing Anvil instance."
    export ANVIL_RPC_URL="http://127.0.0.1:$ANVIL_PORT"
    
    # Verify it's responsive
    if curl -s -X POST -H "Content-Type: application/json" \
        --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
        "$ANVIL_RPC_URL" > /dev/null 2>&1; then
        log_info "Using existing Anvil instance at $ANVIL_RPC_URL"
    else
        log_error "Port $ANVIL_PORT is in use but not responding to RPC calls"
        exit 1
    fi
else
    # Start Anvil
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
    STARTED_ANVIL=true

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
    
    export ANVIL_RPC_URL="http://127.0.0.1:$ANVIL_PORT"
fi

# Verify connection and get block number
BLOCK_NUM=$(curl -s -X POST -H "Content-Type: application/json" \
    --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
    "$ANVIL_RPC_URL" | grep -o '"result":"[^"]*"' | cut -d'"' -f4)
log_info "Current block: $BLOCK_NUM"

log_header "Running TypeScript Benchmarks"

log_info "Benchmark settings:"
log_info "  RPC URL: $ANVIL_RPC_URL"
log_info "  Package manager: $PKG_MANAGER"
echo ""

# Run TypeScript benchmarks
cd typescript
$RUN_CMD bench 2>&1 | tee ../results/ts-results.txt
TS_EXIT_CODE=${PIPESTATUS[0]}
cd ..

if [ $TS_EXIT_CODE -eq 0 ]; then
    log_header "TypeScript Benchmarks Complete"
    log_info "Results saved to: results/ts-results.txt"
else
    log_error "TypeScript benchmarks failed with exit code: $TS_EXIT_CODE"
    exit $TS_EXIT_CODE
fi

echo ""
log_info "Done!"
