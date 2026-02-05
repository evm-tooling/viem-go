#!/bin/bash
#
# Full Benchmark Suite Script
#
# Runs ALL benchmarks for both Go and TypeScript, then generates
# a comprehensive comparison report with per-benchmark metrics and
# overall global metrics.
#
# Usage: ./bench-full.sh
#
# Configuration (via environment variables):
#   ANVIL_PORT     - Port for Anvil (default: 8545)
#   FORK_URL       - Ethereum RPC URL to fork from (default: https://eth.drpc.org)
#   FORK_BLOCK     - Block number to fork from (default: latest)
#   GO_BENCH_TIME  - Duration for each Go benchmark (default: 10s)
#   GO_BENCH_COUNT - Number of Go benchmark runs (default: 5)

set -e

# Change to benchmarks directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.."

# Configuration
ANVIL_PORT=${ANVIL_PORT:-8545}
FORK_URL=${FORK_URL:-"https://eth.drpc.org"}
FORK_BLOCK=${FORK_BLOCK:-}
ANVIL_TIMEOUT=${ANVIL_TIMEOUT:-30}
GO_BENCH_TIME=${GO_BENCH_TIME:-10s}
GO_BENCH_COUNT=${GO_BENCH_COUNT:-5}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
BOLD='\033[1m'
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

log_subheader() {
    echo ""
    echo -e "${CYAN}--- $1 ---${NC}"
    echo ""
}

log_benchmark() {
    echo -e "${MAGENTA}[BENCH]${NC} $1"
}

# Check if anvil is installed
if ! command -v anvil &> /dev/null; then
    log_error "Anvil not found. Please install Foundry: https://getfoundry.sh"
    exit 1
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    log_error "Go not found. Please install Go: https://go.dev/dl/"
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

log_header "Full Benchmark Suite: viem-go vs viem TypeScript"

# Discover available benchmarks
BENCHMARKS=()
for f in go/*_bench_test.go; do
    if [ -f "$f" ]; then
        name=$(basename "$f" | sed 's/_bench_test\.go$//')
        # Only include if TypeScript counterpart exists
        if [ -f "typescript/${name}.bench.ts" ]; then
            BENCHMARKS+=("$name")
        else
            log_warn "Skipping '$name' - no TypeScript counterpart found"
        fi
    fi
done

if [ ${#BENCHMARKS[@]} -eq 0 ]; then
    log_error "No benchmarks found!"
    exit 1
fi

log_info "Found ${#BENCHMARKS[@]} benchmark suites: ${BENCHMARKS[*]}"
log_info "Package manager: $PKG_MANAGER"
log_info "Go bench settings: time=$GO_BENCH_TIME, count=$GO_BENCH_COUNT"

# Ensure results directory exists
mkdir -p results

# Check and install TypeScript dependencies
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

# Track if we started Anvil
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
    log_warn "Port $ANVIL_PORT is already in use. Using existing Anvil instance."
    export ANVIL_RPC_URL="http://127.0.0.1:$ANVIL_PORT"
    
    if curl -s -X POST -H "Content-Type: application/json" \
        --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
        "$ANVIL_RPC_URL" > /dev/null 2>&1; then
        log_info "Using existing Anvil at $ANVIL_RPC_URL"
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

    ANVIL_CMD="anvil --port $ANVIL_PORT --fork-url $FORK_URL --no-mining --accounts 10 --balance 10000 --silent"
    if [ -n "$FORK_BLOCK" ]; then
        ANVIL_CMD="$ANVIL_CMD --fork-block-number $FORK_BLOCK"
    fi

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
    log_info "Anvil ready in ${READY_TIME}s"
    
    export ANVIL_RPC_URL="http://127.0.0.1:$ANVIL_PORT"
fi

# Get current block number
BLOCK_NUM=$(curl -s -X POST -H "Content-Type: application/json" \
    --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
    "$ANVIL_RPC_URL" | grep -o '"result":"[^"]*"' | cut -d'"' -f4)
log_info "Current block: $BLOCK_NUM"

# Track benchmark results
BENCHMARK_STATUS=()
TOTAL_BENCHMARKS=${#BENCHMARKS[@]}
CURRENT_BENCH=0
START_TIME=$(date +%s)

# Run all Go benchmarks first (more efficient than interleaving)
log_header "Running All Go Benchmarks"

cd ..
go test -bench=. -benchmem -benchtime=$GO_BENCH_TIME -count=$GO_BENCH_COUNT \
    ./benchmarks/go/... 2>&1 | tee benchmarks/results/go-results.txt
GO_EXIT_CODE=${PIPESTATUS[0]}
cd benchmarks

if [ $GO_EXIT_CODE -ne 0 ]; then
    log_warn "Go benchmarks had errors (exit code: $GO_EXIT_CODE)"
fi

log_info "Go results saved to: results/go-results.txt"

# Run all TypeScript benchmarks
log_header "Running All TypeScript Benchmarks"

cd typescript
$RUN_CMD bench 2>&1 | tee ../results/ts-results.txt
TS_EXIT_CODE=${PIPESTATUS[0]}
cd ..

if [ $TS_EXIT_CODE -ne 0 ]; then
    log_warn "TypeScript benchmarks had errors (exit code: $TS_EXIT_CODE)"
fi

log_info "TypeScript results saved to: results/ts-results.txt"

# Calculate total time
END_TIME=$(date +%s)
TOTAL_TIME=$((END_TIME - START_TIME))
TOTAL_MINUTES=$((TOTAL_TIME / 60))
TOTAL_SECONDS=$((TOTAL_TIME % 60))

# Generate comprehensive comparison report
log_header "Generating Comprehensive Comparison Report"

# Run the comparison script in full mode
bun run compare.ts --mode full

log_header "Full Benchmark Suite Complete"

echo ""
echo -e "${BOLD}Summary${NC}"
echo "─────────────────────────────────────────"
echo -e "  Benchmark suites: ${CYAN}${#BENCHMARKS[@]}${NC} (${BENCHMARKS[*]})"
echo -e "  Total time:       ${CYAN}${TOTAL_MINUTES}m ${TOTAL_SECONDS}s${NC}"
echo ""
echo -e "${BOLD}Output Files${NC}"
echo "─────────────────────────────────────────"
echo "  Go results:       results/go-results.txt"
echo "  TS results:       results/ts-results.txt"
echo "  Comparison:       results/comparison.md"
echo "  Full report:      results/full-report.md"
echo ""

# Display quick summary from the report if it exists
if [ -f "results/full-report.md" ]; then
    log_info "Full report generated: results/full-report.md"
fi

echo ""
log_info "Done!"
