## Benchmark Makefile (single entrypoint)
##
## All benchmarking is driven via `benchmarks/bench.sh`.
##
## Usage:
##   make bench                         # full run (defaults: iter=100, count=4)
##   make bench ITER=100 COUNT=2         # full run (100, 1000)
##   make bench DRY_RUN=1                # dry run (no files)
##   make bench-single BENCH=abi         # single suite compare (Go+TS)
##   make bench-single BENCH=call RUNTIME=ts ITER=10 COUNT=1
##   make chart SUITE=abi ITER=5000       # generate comparison chart for one suite
##   make charts ITER=5000                # generate comparison charts for all suites

.PHONY: help bench bench-single install clean clean-all chart charts

# Benchmark params
ITER ?= 100
COUNT ?= 4
DRY_RUN ?= 0

# Optional runtime selection in single mode (go|ts). Empty => run both.
RUNTIME ?=

# Anvil config (only used by RPC-heavy suites)
ANVIL_PORT ?= 8545
FORK_URL ?= https://eth.drpc.org
# FORK_BLOCK ?=
export ANVIL_PORT
export FORK_URL
export FORK_BLOCK

help:
	@echo "Benchmark Suite (single entrypoint)"
	@echo ""
	@echo "Targets:"
	@echo "  make bench                         Full run (CPU suites ladder by default)"
	@echo "  make bench-single BENCH=<suite>     Single suite (Go+TS by default)"
	@echo "  make chart SUITE=<suite>            Generate comparison chart for one suite"
	@echo "  make charts                         Generate comparison charts for all suites"
	@echo "  make install                       Install TS deps"
	@echo "  make clean                         Remove results artifacts"
	@echo ""
	@echo "Vars:"
	@echo "  ITER=$(ITER)   (base iterations)"
	@echo "  COUNT=$(COUNT) (levels 1..4: iter, iter*10, iter*100, iter*1000)"
	@echo "  DRY_RUN=$(DRY_RUN) (0|1)"
	@echo "  RUNTIME=$(if $(RUNTIME),$(RUNTIME),<both>) (only for bench-single)"
	@echo ""
	@echo "Anvil:"
	@echo "  ANVIL_PORT=$(ANVIL_PORT)"
	@echo "  FORK_URL=$(FORK_URL)"
	@echo "  FORK_BLOCK=$(if $(FORK_BLOCK),$(FORK_BLOCK),latest)"

install:
	cd typescript && npm install

bench:
	bash ./bench.sh --type full --iter $(ITER) --count $(COUNT) --isdryrun $(DRY_RUN)

bench-single:
ifndef BENCH
	@echo "Error: BENCH is required (e.g. BENCH=abi, BENCH=call)"; exit 1
endif
	@if [ -n "$(RUNTIME)" ]; then \
		./bench.sh --type single --benchname $(BENCH) --runtime $(RUNTIME) --iter $(ITER) --count $(COUNT) --isdryrun $(DRY_RUN); \
	else \
		./bench.sh --type single --benchname $(BENCH) --iter $(ITER) --count $(COUNT) --isdryrun $(DRY_RUN); \
	fi

chart:
ifndef SUITE
	@echo "Error: SUITE is required (e.g. SUITE=abi, SUITE=address)"; exit 1
endif
	bun run generate-chart.ts --suite $(SUITE) --iter $(ITER)

charts:
	bun run generate-chart.ts --all --iter $(ITER)

clean:
	rm -rf results/*.txt results/*.json results/*.md results/charts/ results/*charts-iter*/ results/*go-results.iter*.txt results/*ts-results.iter*.txt results/*comparison.iter*.md results/*full-report.iter*.md results/single-run results/bench-comparison-charts

clean-all: clean
	rm -rf typescript/node_modules
