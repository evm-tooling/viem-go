# Benchmark Makefile
#
# Orchestrates cross-language benchmarks for viem-go and viem TypeScript.
# Uses a shared Anvil instance for fair comparison.
#
# Usage:
#   make bench          - Run all benchmarks with shared Anvil
#   make bench-go       - Run only Go benchmarks (requires running Anvil)
#   make bench-ts       - Run only TypeScript benchmarks (requires running Anvil)
#   make compare        - Generate comparison report from results
#   make install        - Install TypeScript dependencies
#   make clean          - Clean results directory
#   make anvil          - Start Anvil manually (for development)

.PHONY: bench bench-go bench-ts bench-go-standalone bench-ts-standalone bench-single bench-full bench-all compare install clean anvil anvil-stop setup help

# Configuration
ANVIL_PORT ?= 8545
FORK_URL ?= https://eth.drpc.org
# FORK_BLOCK ?=
RPC_URL := http://127.0.0.1:$(ANVIL_PORT)
GO_BENCH_COUNT ?= 5
GO_BENCH_TIME ?= 10s

# Export for child processes
export ANVIL_PORT
export FORK_URL
export FORK_BLOCK
export ANVIL_RPC_URL := $(RPC_URL)

# Default target
help:
	@echo "Benchmark Suite for viem-go vs viem TypeScript"
	@echo ""
	@echo "Usage:"
	@echo "  make bench                    Run all benchmarks with shared Anvil"
	@echo "  make bench-single BENCH=name  Run single benchmark comparison (e.g., BENCH=call)"
	@echo "  make bench-full               Run full suite with comprehensive report"
	@echo "  make bench-go                 Run only Go benchmarks (requires running Anvil)"
	@echo "  make bench-ts                 Run only TypeScript benchmarks (requires running Anvil)"
	@echo "  make bench-go-standalone      Run Go benchmarks with own Anvil lifecycle"
	@echo "  make bench-ts-standalone      Run TypeScript benchmarks with own Anvil lifecycle"
	@echo "  make compare                  Generate comparison report"
	@echo "  make install                  Install TypeScript dependencies"
	@echo "  make clean                    Clean results directory"
	@echo "  make anvil                    Start Anvil manually"
	@echo "  make anvil-stop               Stop running Anvil instance"
	@echo ""
	@echo "Single Benchmark Examples:"
	@echo "  make bench-single BENCH=call      # Compare call benchmarks"
	@echo "  make bench-single BENCH=multicall # Compare multicall benchmarks"
	@echo ""
	@echo "Configuration (via environment):"
	@echo "  ANVIL_PORT=$(ANVIL_PORT)"
	@echo "  FORK_URL=$(FORK_URL)"
	@echo "  FORK_BLOCK=$(if $(FORK_BLOCK),$(FORK_BLOCK),latest)"
	@echo "  GO_BENCH_COUNT=$(GO_BENCH_COUNT)"
	@echo "  GO_BENCH_TIME=$(GO_BENCH_TIME)"

# Setup: ensure results directory exists and dependencies are installed
setup: results/.gitkeep
	@if [ ! -d "typescript/node_modules" ]; then \
		echo "Installing TypeScript dependencies..."; \
		$(MAKE) install; \
	fi

results/.gitkeep:
	@mkdir -p results
	@touch results/.gitkeep

# Install TypeScript dependencies
install:
	@echo "Installing TypeScript dependencies..."
	cd typescript && npm install

# Run all benchmarks with shared Anvil instance
bench: setup
	@echo "=============================================="
	@echo "Running Cross-Language Benchmarks"
	@echo "=============================================="
	@echo ""
	./scripts/anvil.sh $(MAKE) _bench-sequential

# Internal target: run benchmarks sequentially (called by anvil.sh)
_bench-sequential:
	@echo ""
	@echo "=============================================="
	@echo "Go Benchmarks (viem-go)"
	@echo "=============================================="
	$(MAKE) bench-go
	@echo ""
	@echo "=============================================="
	@echo "TypeScript Benchmarks (viem)"
	@echo "=============================================="
	$(MAKE) bench-ts
	@echo ""
	@echo "=============================================="
	@echo "Benchmarks Complete"
	@echo "=============================================="
	@echo ""
	@echo "Results saved to:"
	@echo "  - results/go-results.txt"
	@echo "  - results/ts-results.json"
	@echo ""
	$(MAKE) compare
	@echo ""
	@echo "=============================================="
	@echo "Done"
	@echo "=============================================="
	@echo ""

# Run Go benchmarks only
bench-go:
	@echo "Running Go benchmarks against $(RPC_URL)..."
	cd .. && go test -bench=. -benchmem -benchtime=$(GO_BENCH_TIME) -count=$(GO_BENCH_COUNT) ./benchmarks/go/... 2>&1 | tee benchmarks/results/go-results.txt

# Run TypeScript benchmarks only
bench-ts:
	@echo "Running TypeScript benchmarks against $(RPC_URL)..."
	cd typescript && npm run bench 2>&1 | tee ../results/ts-results.txt

# Run Go benchmarks with standalone Anvil lifecycle
bench-go-standalone:
	@echo "Running standalone Go benchmarks..."
	./scripts/bench-go.sh

# Run TypeScript benchmarks with standalone Anvil lifecycle
bench-ts-standalone:
	@echo "Running standalone TypeScript benchmarks..."
	./scripts/bench-ts.sh

# Run a single benchmark comparison (Go vs TypeScript)
# Usage: make bench-single BENCH=call
bench-single:
ifndef BENCH
	@echo "Error: BENCH variable is required"
	@echo "Usage: make bench-single BENCH=<name>"
	@echo "Examples:"
	@echo "  make bench-single BENCH=call"
	@echo "  make bench-single BENCH=multicall"
	@echo ""
	@echo "Available benchmarks:"
	@./scripts/bench-single.sh --list | grep -E "^\s+-"
	@exit 1
endif
	@echo "Running single benchmark comparison: $(BENCH)"
	./scripts/bench-single.sh --bench $(BENCH)

# Run full benchmark suite with comprehensive report
bench-full:
	@echo "Running full benchmark suite..."
	./scripts/bench-full.sh

# Generate comparison report
compare:
	@echo "Generating comparison report..."
	@if [ -f "results/go-results.txt" ] && [ -f "results/ts-results.txt" ]; then \
		bun run compare.ts; \
	else \
		echo "Error: Results files not found. Run 'make bench' first."; \
		exit 1; \
	fi

# Start Anvil manually (for development/debugging)
anvil:
	@echo "Starting Anvil manually..."
	@echo "Press Ctrl+C to stop."
	./scripts/anvil.sh

# Stop any running Anvil instance on the configured port
anvil-stop:
	@echo "Stopping Anvil on port $(ANVIL_PORT)..."
	@PID=$$(lsof -ti :$(ANVIL_PORT) 2>/dev/null); \
	if [ -n "$$PID" ]; then \
		kill $$PID 2>/dev/null && echo "Killed Anvil (PID: $$PID)" || echo "Failed to kill process"; \
	else \
		echo "No process found on port $(ANVIL_PORT)"; \
	fi

# Clean results
clean:
	@echo "Cleaning results..."
	rm -rf results/*.txt results/*.json results/*.md
	@echo "Done."

# Clean everything including node_modules
clean-all: clean
	rm -rf typescript/node_modules
