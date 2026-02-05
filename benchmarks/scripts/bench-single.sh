#!/bin/bash
#
# Single Benchmark Comparison Script
#
# Runs a specific benchmark by name for both Go and TypeScript,
# then generates a comparison report for just that benchmark.
#
# Usage: ./bench-single.sh --bench <name>
#        ./bench-single.sh -b call
#        ./bench-single.sh -b multicall
#
# Benchmark naming conventions:
#   - Go:         {name}_bench_test.go
#   - TypeScript: {name}.bench.ts
#
# Configuration (via environment variables):
#   ANVIL_PORT     - Port for Anvil (default: 8545)
#   FORK_URL       - Ethereum RPC URL to fork from (default: https://eth.drpc.org)
#   FORK_BLOCK     - Block number to fork from (default: latest)
#   GO_BENCH_TIME  - Duration for each Go benchmark (default: 10s)
#   GO_BENCH_COUNT - Number of Go benchmark runs (default: 3)

set -e

# Change to benchmarks directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.."

# Configuration
ANVIL_PORT=${ANVIL_PORT:-8545}
FORK_URL=${FORK_URL:-"https://rough-purple-market.quiknode.pro/c1a568726a34041d3c5d58603f5981951e6a8503"}
FORK_BLOCK=${FORK_BLOCK:-}
ANVIL_TIMEOUT=${ANVIL_TIMEOUT:-30}
GO_BENCH_TIME=${GO_BENCH_TIME:-2s}
GO_BENCH_COUNT=${GO_BENCH_COUNT:-1}
BENCH_NAME=""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
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

usage() {
    echo "Usage: $0 --bench <name> | -b <name>"
    echo ""
    echo "Run a specific benchmark by name and generate a comparison report."
    echo ""
    echo "Options:"
    echo "  -b, --bench <name>  Benchmark name (e.g., 'call', 'multicall')"
    echo "  -h, --help          Show this help message"
    echo "  -l, --list          List available benchmarks"
    echo ""
    echo "Examples:"
    echo "  $0 -b call          Run 'call' benchmarks"
    echo "  $0 --bench multicall Run 'multicall' benchmarks"
    echo ""
    echo "Environment Variables:"
    echo "  ANVIL_PORT      Port for Anvil (default: 8545)"
    echo "  FORK_URL        Ethereum RPC URL to fork (default: https://eth.drpc.org)"
    echo "  GO_BENCH_TIME   Duration for Go benchmarks (default: 10s)"
    echo "  GO_BENCH_COUNT  Number of benchmark runs (default: 3)"
}

list_benchmarks() {
    echo "Available benchmarks:"
    echo ""
    echo "Go benchmarks (in go/):"
    for f in go/*_bench_test.go; do
        if [ -f "$f" ]; then
            name=$(basename "$f" | sed 's/_bench_test\.go$//')
            echo "  - $name"
        fi
    done
    echo ""
    echo "TypeScript benchmarks (in typescript/):"
    for f in typescript/*.bench.ts; do
        if [ -f "$f" ]; then
            name=$(basename "$f" | sed 's/\.bench\.ts$//')
            echo "  - $name"
        fi
    done
}

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -b|--bench)
            BENCH_NAME="$2"
            shift 2
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        -l|--list)
            list_benchmarks
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            usage
            exit 1
            ;;
    esac
done

# Validate benchmark name
if [ -z "$BENCH_NAME" ]; then
    log_error "Benchmark name is required"
    usage
    exit 1
fi

# Check if benchmark files exist
GO_BENCH_FILE="go/${BENCH_NAME}_bench_test.go"
TS_BENCH_FILE="typescript/${BENCH_NAME}.bench.ts"

if [ ! -f "$GO_BENCH_FILE" ]; then
    log_error "Go benchmark not found: $GO_BENCH_FILE"
    log_info "Use --list to see available benchmarks"
    exit 1
fi

if [ ! -f "$TS_BENCH_FILE" ]; then
    log_error "TypeScript benchmark not found: $TS_BENCH_FILE"
    log_info "Use --list to see available benchmarks"
    exit 1
fi

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

log_header "Single Benchmark: ${BENCH_NAME}"

log_info "Benchmark files:"
log_info "  Go: $GO_BENCH_FILE"
log_info "  TS: $TS_BENCH_FILE"

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

# Result files for this specific benchmark
GO_RESULT_FILE="results/${BENCH_NAME}-go-results.txt"
TS_RESULT_FILE="results/${BENCH_NAME}-ts-results.txt"

# Run Go benchmark
log_subheader "Running Go Benchmark: ${BENCH_NAME}"
log_info "Settings: time=$GO_BENCH_TIME, count=$GO_BENCH_COUNT"

# Convert bench name to Go benchmark pattern (e.g., "call" -> "Benchmark.*" in call_bench_test.go)
# We need to run only benchmarks from the specific file
GO_BENCH_PATTERN="."

cd ..
go test -bench=$GO_BENCH_PATTERN -benchmem -benchtime=$GO_BENCH_TIME -count=$GO_BENCH_COUNT \
    ./benchmarks/go/${BENCH_NAME}_bench_test.go ./benchmarks/go/main_test.go 2>&1 | tee benchmarks/$GO_RESULT_FILE
GO_EXIT_CODE=${PIPESTATUS[0]}
cd benchmarks

if [ $GO_EXIT_CODE -ne 0 ]; then
    log_error "Go benchmark failed with exit code: $GO_EXIT_CODE"
    exit $GO_EXIT_CODE
fi

log_info "Go results saved to: $GO_RESULT_FILE"

# Run TypeScript benchmark
log_subheader "Running TypeScript Benchmark: ${BENCH_NAME}"
log_info "Package manager: $PKG_MANAGER"

cd typescript
# Run only the specific benchmark file
$RUN_CMD bench -- "${BENCH_NAME}.bench.ts" 2>&1 | tee ../$TS_RESULT_FILE
TS_EXIT_CODE=${PIPESTATUS[0]}
cd ..

if [ $TS_EXIT_CODE -ne 0 ]; then
    log_error "TypeScript benchmark failed with exit code: $TS_EXIT_CODE"
    exit $TS_EXIT_CODE
fi

log_info "TypeScript results saved to: $TS_RESULT_FILE"

# Generate comparison report
log_subheader "Generating Comparison Report"

# Run the comparison script in single-bench mode
bun run compare.ts --bench "$BENCH_NAME" --go-results "$GO_RESULT_FILE" --ts-results "$TS_RESULT_FILE"

log_header "Benchmark Complete: ${BENCH_NAME}"

echo ""
log_info "Results files:"
log_info "  Go:         $GO_RESULT_FILE"
log_info "  TypeScript: $TS_RESULT_FILE"
log_info "  Comparison: results/${BENCH_NAME}-comparison.md"
echo ""
