#!/bin/bash
#
# Unified benchmark entrypoint for viem-go vs viem (TypeScript).
#
# Supports:
#   - type=full: run all suites (CPU suites default to iteration ladder)
#   - type=single: run one suite (Go+TS by default, optional runtime=go|ts)
#
# Iteration model:
#   - Go: go test -benchtime=<N>x -count=1
#   - TS: vitest bench with BENCH_ITERATIONS=<N>
#
# Notes:
#   - CPU-only suites: abi, address, ens, event, hash, signature, unit
#   - RPC-heavy suites: call, multicall (requires Anvil)
#
# Usage examples:
#   ./bench.sh --type full
#   ./bench.sh --type full --iter 100 --count 4
#   ./bench.sh --type single --benchname abi --iter 100 --count 2
#   ./bench.sh --type single --benchname call --runtime ts --iter 10 --count 1
#
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BENCH_DIR="$SCRIPT_DIR"
REPO_ROOT="$(cd "$BENCH_DIR/.." && pwd)"

RESULTS_DIR="$BENCH_DIR/results"
SINGLE_RUN_ROOT="$RESULTS_DIR/single-run"
RUN_ID=""
RUN_DIR=""

# ---------------------------------------------------------------------------
# Defaults
# ---------------------------------------------------------------------------

TYPE=""
RUNTIME=""          # only meaningful for type=single
BENCHNAME=""        # only meaningful for type=single
BASE_ITER=100
LEVELS=4            # 1..4
IS_DRY_RUN=0

USER_SET_ITER=0
USER_SET_COUNT=0

# Anvil config
ANVIL_PORT="${ANVIL_PORT:-8545}"
FORK_URL="${FORK_URL:-https://eth.drpc.org}"
FORK_BLOCK="${FORK_BLOCK:-}"
ANVIL_TIMEOUT="${ANVIL_TIMEOUT:-30}"
ANVIL_WARMUP_SECONDS="${ANVIL_WARMUP_SECONDS:-10}"

# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------

log() { echo "[bench] $*"; }
warn() { echo "[bench][warn] $*" >&2; }
die() { echo "[bench][error] $*" >&2; exit 1; }

usage() {
  cat <<'EOF'
Unified benchmark entrypoint.

Args:
  --type full|single
  --benchname <suite>        (required if --type single)
  --runtime go|ts            (optional if --type single; default runs both)
  --iter <baseIterations>    (default: 100)
  --count <levels>           (1..4, default: 4) => iter, iter*10, iter*100, iter*1000
  --isdryrun 0|1             (default: 0)

Env (optional):
  ANVIL_PORT, FORK_URL, FORK_BLOCK, ANVIL_TIMEOUT

Examples:
  ./bench.sh --type full
  ./bench.sh --type full --iter 100 --count 4
  ./bench.sh --type single --benchname abi --iter 100 --count 2
  ./bench.sh --type single --benchname call --runtime ts --iter 10 --count 1
EOF
}

timestamp_id() {
  # e.g. 20260207-153045
  date +"%Y%m%d-%H%M%S"
}

mkdirp() {
  mkdir -p "$1"
}

write_codeblock_md() {
  local out="$1"
  local title="$2"
  local raw="$3"

  mkdirp "$(dirname "$out")"
  {
    echo "# $title"
    echo ""
    echo '```'
    cat "$raw"
    echo '```'
    echo ""
  } >"$out"
}

capitalize() {
  # Portable (macOS BSD userland) capitalization of first letter.
  awk '{ printf toupper(substr($0, 1, 1)) substr($0, 2) }' <<<"$1"
}

is_cpu_suite() {
  case "$1" in
    abi|address|ens|event|hash|signature|unit) return 0 ;;
    *) return 1 ;;
  esac
}

is_rpc_suite() {
  case "$1" in
    call|multicall) return 0 ;;
    *) return 1 ;;
  esac
}

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || die "Missing required command: $1"
}

ensure_ts_deps() {
  if [ -d "$BENCH_DIR/typescript/node_modules" ]; then
    return 0
  fi
  log "Installing TypeScript benchmark dependencies..."
  (cd "$BENCH_DIR/typescript" && npm install)
}

ensure_anvil() {
  require_cmd anvil

  if lsof -Pi :"$ANVIL_PORT" -sTCP:LISTEN -t >/dev/null 2>&1; then
    export ANVIL_RPC_URL="http://127.0.0.1:$ANVIL_PORT"
    if curl -s -X POST -H "Content-Type: application/json" \
      --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
      "$ANVIL_RPC_URL" >/dev/null 2>&1; then
      log "Using existing Anvil at $ANVIL_RPC_URL"
      return 0
    fi
    die "Port $ANVIL_PORT is in use but not responding to JSON-RPC"
  fi

  log "Starting Anvil on port $ANVIL_PORT (fork: $FORK_URL${FORK_BLOCK:+ @ $FORK_BLOCK})..."
  local cmd=(anvil --port "$ANVIL_PORT" --fork-url "$FORK_URL" --no-mining --accounts 10 --balance 10000 --silent)
  if [ -n "$FORK_BLOCK" ]; then
    cmd+=(--fork-block-number "$FORK_BLOCK")
  fi

  "${cmd[@]}" &
  ANVIL_PID=$!

  cleanup_anvil() {
    if [ -n "${ANVIL_PID:-}" ]; then
      log "Stopping Anvil (PID: $ANVIL_PID)..."
      kill "$ANVIL_PID" 2>/dev/null || true
      wait "$ANVIL_PID" 2>/dev/null || true
      log "Anvil stopped."
      ANVIL_PID=""
    fi
  }
  trap cleanup_anvil EXIT INT TERM

  log "Waiting for Anvil to be ready..."
  local start_ts
  start_ts="$(date +%s)"
  while true; do
    if curl -s -X POST -H "Content-Type: application/json" \
      --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
      "http://127.0.0.1:$ANVIL_PORT" >/dev/null 2>&1; then
      break
    fi
    local elapsed
    elapsed="$(($(date +%s) - start_ts))"
    if [ "$elapsed" -ge "$ANVIL_TIMEOUT" ]; then
      die "Timeout waiting for Anvil to start"
    fi
    sleep 0.1
  done

  export ANVIL_RPC_URL="http://127.0.0.1:$ANVIL_PORT"
  log "Anvil ready at $ANVIL_RPC_URL"

  warmup_anvil "$ANVIL_RPC_URL" "$ANVIL_WARMUP_SECONDS"
}

ensure_anvil_fresh() {
  # Like ensure_anvil, but refuses to reuse an existing process on the port.
  require_cmd anvil
  if lsof -Pi :"$ANVIL_PORT" -sTCP:LISTEN -t >/dev/null 2>&1; then
    die "Port $ANVIL_PORT is already in use. Stop the existing Anvil/process to run with a fresh Anvil per runtime."
  fi
  ensure_anvil
}

stop_anvil_if_owned() {
  # Only stops Anvil if this script started it.
  if [ -n "${ANVIL_PID:-}" ] && declare -F cleanup_anvil >/dev/null 2>&1; then
    cleanup_anvil || true
  fi
}

warmup_anvil() {
  local rpc_url="$1"
  local seconds="$2"

  # Allow opting out by setting ANVIL_WARMUP_SECONDS=0
  if [ "$seconds" -le 0 ]; then
    return 0
  fi

  log "Warming up Anvil for ${seconds}s..."
  local end_ts
  end_ts="$(($(date +%s) + seconds))"

  # Exercise common RPC paths the benchmarks will hit (mainnet fork).
  # This warms fork/cache without doing any per-runtime warmup inside tests.
  while [ "$(date +%s)" -lt "$end_ts" ]; do
    curl -s -X POST -H "Content-Type: application/json" \
      --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
      "$rpc_url" >/dev/null 2>&1 || true

    curl -s -X POST -H "Content-Type: application/json" \
      --data '{"jsonrpc":"2.0","method":"eth_call","params":[{"to":"0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48","data":"0x06fdde03"},"latest"],"id":1}' \
      "$rpc_url" >/dev/null 2>&1 || true

    sleep 0.2
  done

  log "Anvil warmup complete."
}
truncate_results_files() {
  if [ "$IS_DRY_RUN" = "1" ]; then
    return 0
  fi
  mkdirp "$RESULTS_DIR"
  : >"$RESULTS_DIR/go-results.txt"
  : >"$RESULTS_DIR/ts-results.txt"
}

run_go_suite() {
  local suite="$1"
  local iterations="$2"
  local out_raw="${3:-}" # optional: where to write this suite's raw output

  require_cmd go
  # All Go benchmarks share a TestMain that connects to Anvil.
  ensure_anvil

  local cap
  cap="$(capitalize "$suite")"
  local bench_regex="^Benchmark${cap}_"

  local cmd=(go test -run '^$' -bench "$bench_regex" -benchmem -benchtime "${iterations}x" -count 1 ./benchmarks/go)

  log "Go: suite=$suite iterations=$iterations"
  if [ "$IS_DRY_RUN" = "1" ]; then
    (cd "$REPO_ROOT" && "${cmd[@]}")
  else
    if [ -n "$out_raw" ]; then
      mkdirp "$(dirname "$out_raw")"
      : >"$out_raw"
      (cd "$REPO_ROOT" && "${cmd[@]}" 2>&1 | tee -a "$RESULTS_DIR/go-results.txt" >"$out_raw")
    else
      (cd "$REPO_ROOT" && "${cmd[@]}" 2>&1 | tee -a "$RESULTS_DIR/go-results.txt")
    fi
  fi
}

run_ts_suite() {
  local suite="$1"
  local iterations="$2"
  local out_raw="${3:-}" # optional: where to write this suite's raw output

  ensure_ts_deps

  local file_path="$BENCH_DIR/typescript/${suite}.bench.ts"
  if [ ! -f "$file_path" ]; then
    die "TypeScript benchmark file not found: $file_path"
  fi

  # RPC-heavy TS suites require Anvil.
  if is_rpc_suite "$suite"; then
    ensure_anvil
  fi

  log "TS: suite=$suite iterations=$iterations"
  if [ "$IS_DRY_RUN" = "1" ]; then
    (cd "$BENCH_DIR/typescript" && BENCH_ITERATIONS="$iterations" BENCH_SUITE="$suite" ./node_modules/.bin/vitest bench --run)
  else
    if [ -n "$out_raw" ]; then
      mkdirp "$(dirname "$out_raw")"
      : >"$out_raw"
      (cd "$BENCH_DIR/typescript" && BENCH_ITERATIONS="$iterations" BENCH_SUITE="$suite" ./node_modules/.bin/vitest bench --run 2>&1 | tee -a "$RESULTS_DIR/ts-results.txt" >"$out_raw")
    else
      (cd "$BENCH_DIR/typescript" && BENCH_ITERATIONS="$iterations" BENCH_SUITE="$suite" ./node_modules/.bin/vitest bench --run 2>&1 | tee -a "$RESULTS_DIR/ts-results.txt")
    fi
  fi
}

snapshot_iter_overall() {
  local iter_dir="$1"   # .../iter-<N>

  if [ "$IS_DRY_RUN" = "1" ]; then
    log "Dry run: skipping report generation"
    return 0
  fi

  require_cmd bun

  # Generate reports/charts from canonical file locations.
  (cd "$BENCH_DIR" && bun run compare.ts --mode full)

  local go_src="$RESULTS_DIR/go-results.txt"
  local ts_src="$RESULTS_DIR/ts-results.txt"
  local cmp_src="$RESULTS_DIR/comparison.md"
  local full_src="$RESULTS_DIR/full-report.md"
  local charts_src="$RESULTS_DIR/charts"

  local overall_dir="$iter_dir/_overall"
  mkdirp "$overall_dir"
  local go_dst="$overall_dir/go-results.txt"
  local ts_dst="$overall_dir/ts-results.txt"
  local cmp_dst="$overall_dir/comparison.md"
  local full_dst="$overall_dir/full-report.md"
  local charts_dst="$overall_dir/charts"

  for f in "$go_dst" "$ts_dst" "$cmp_dst" "$full_dst"; do
    if [ -e "$f" ]; then
      die "Refusing to overwrite existing file: $f (run make clean or delete results first)"
    fi
  done
  if [ -e "$charts_dst" ]; then
    die "Refusing to overwrite existing dir: $charts_dst (run make clean or delete results first)"
  fi

  mv "$go_src" "$go_dst"
  mv "$ts_src" "$ts_dst"
  mv "$cmp_src" "$cmp_dst"
  mv "$full_src" "$full_dst"
  if [ -d "$charts_src" ]; then
    mv "$charts_src" "$charts_dst"
  fi

  log "Saved overall report to: $overall_dir"
}

snapshot_iter_runtime_results() {
  local iter_dir="$1" # .../iter-<N>
  local runtime="$2"  # go|ts

  if [ "$IS_DRY_RUN" = "1" ]; then
    return 0
  fi

  local overall_dir="$iter_dir/_overall"
  mkdirp "$overall_dir"

  if [ "$runtime" = "go" ]; then
    local src="$RESULTS_DIR/go-results.txt"
    local dst="$overall_dir/go-results.txt"
    [ -e "$dst" ] && die "Refusing to overwrite existing file: $dst"
    mv "$src" "$dst"
  elif [ "$runtime" = "ts" ]; then
    local src="$RESULTS_DIR/ts-results.txt"
    local dst="$overall_dir/ts-results.txt"
    [ -e "$dst" ] && die "Refusing to overwrite existing file: $dst"
    mv "$src" "$dst"
  else
    die "snapshot_iter_runtime_results: invalid runtime '$runtime'"
  fi
}

generate_iter_overall_reports() {
  local iter_dir="$1" # .../iter-<N>

  if [ "$IS_DRY_RUN" = "1" ]; then
    return 0
  fi

  require_cmd bun

  local overall_dir="$iter_dir/_overall"
  local go_saved="$overall_dir/go-results.txt"
  local ts_saved="$overall_dir/ts-results.txt"
  if [ ! -f "$go_saved" ] || [ ! -f "$ts_saved" ]; then
    die "Missing saved results for $iter_dir (expected $go_saved and $ts_saved)"
  fi

  # Hydrate canonical locations for compare.ts
  mkdirp "$RESULTS_DIR"
  cp "$go_saved" "$RESULTS_DIR/go-results.txt"
  cp "$ts_saved" "$RESULTS_DIR/ts-results.txt"

  # Generate reports/charts from canonical file locations.
  (cd "$BENCH_DIR" && bun run compare.ts --mode full)

  local cmp_src="$RESULTS_DIR/comparison.md"
  local full_src="$RESULTS_DIR/full-report.md"
  local charts_src="$RESULTS_DIR/charts"

  local cmp_dst="$overall_dir/comparison.md"
  local full_dst="$overall_dir/full-report.md"
  local charts_dst="$overall_dir/charts"

  for f in "$cmp_dst" "$full_dst"; do
    if [ -e "$f" ]; then
      die "Refusing to overwrite existing file: $f (run make clean or delete results first)"
    fi
  done
  if [ -e "$charts_dst" ]; then
    die "Refusing to overwrite existing dir: $charts_dst (run make clean or delete results first)"
  fi

  mv "$cmp_src" "$cmp_dst"
  mv "$full_src" "$full_dst"
  if [ -d "$charts_src" ]; then
    mv "$charts_src" "$charts_dst"
  fi
}

discover_suites() {
  # Intersection of go/*_bench_test.go and typescript/*.bench.ts by basename.
  python3 - <<'PY'
import os
bench_dir = os.environ["BENCH_DIR"]
go_dir = os.path.join(bench_dir, "go")
ts_dir = os.path.join(bench_dir, "typescript")

go = set()
for f in os.listdir(go_dir):
  if f.endswith("_bench_test.go"):
    go.add(f.replace("_bench_test.go", ""))

ts = set()
for f in os.listdir(ts_dir):
  if f.endswith(".bench.ts"):
    ts.add(f.replace(".bench.ts", ""))

both = sorted(go.intersection(ts))
print("\n".join(both))
PY
}

# ---------------------------------------------------------------------------
# Arg parsing
# ---------------------------------------------------------------------------

while [ $# -gt 0 ]; do
  case "$1" in
    --type)
      TYPE="${2:-}"; shift 2 ;;
    --runtime)
      RUNTIME="${2:-}"; shift 2 ;;
    --benchname)
      BENCHNAME="${2:-}"; shift 2 ;;
    --iter)
      BASE_ITER="${2:-}"; USER_SET_ITER=1; shift 2 ;;
    --count)
      LEVELS="${2:-}"; USER_SET_COUNT=1; shift 2 ;;
    --isdryrun)
      IS_DRY_RUN="${2:-}"; shift 2 ;;
    -h|--help)
      usage; exit 0 ;;
    *)
      die "Unknown arg: $1 (use --help)" ;;
  esac
done

if [ -z "$TYPE" ]; then
  die "--type is required (full|single)"
fi
if [ "$TYPE" != "full" ] && [ "$TYPE" != "single" ]; then
  die "--type must be 'full' or 'single'"
fi
if [ "$TYPE" = "single" ] && [ -z "$BENCHNAME" ]; then
  die "--benchname is required when --type single"
fi
if [ -n "$RUNTIME" ] && [ "$RUNTIME" != "go" ] && [ "$RUNTIME" != "ts" ]; then
  die "--runtime must be 'go' or 'ts'"
fi

if ! [[ "$BASE_ITER" =~ ^[0-9]+$ ]] || [ "$BASE_ITER" -le 0 ]; then
  die "--iter must be a positive integer"
fi
if ! [[ "$LEVELS" =~ ^[0-9]+$ ]] || [ "$LEVELS" -lt 1 ] || [ "$LEVELS" -gt 4 ]; then
  die "--count must be an integer in [1..4]"
fi
if [ "$IS_DRY_RUN" != "0" ] && [ "$IS_DRY_RUN" != "1" ]; then
  die "--isdryrun must be 0 or 1"
fi

# In full mode, runtime/benchname are ignored even if provided.
if [ "$TYPE" = "full" ]; then
  RUNTIME=""
  BENCHNAME=""
fi

export BENCH_DIR

log "type=$TYPE runtime=${RUNTIME:-<both>} benchname=${BENCHNAME:-<all>} iter=$BASE_ITER count=$LEVELS dry=$IS_DRY_RUN"

# ---------------------------------------------------------------------------
# Run
# ---------------------------------------------------------------------------

SUITES=()
while IFS= read -r s; do
  [ -n "$s" ] && SUITES+=("$s")
done < <(discover_suites)

if [ "${#SUITES[@]}" -eq 0 ]; then
  die "No suites found (expected benchmarks/go/*_bench_test.go and benchmarks/typescript/*.bench.ts)"
fi

if [ "$TYPE" = "single" ]; then
  found=0
  for s in "${SUITES[@]}"; do
    if [ "$s" = "$BENCHNAME" ]; then found=1; break; fi
  done
  if [ "$found" -ne 1 ]; then
    die "Unknown benchname '$BENCHNAME' (available: ${SUITES[*]})"
  fi
fi

RUN_ID="run-$(timestamp_id)"
RUN_DIR="$SINGLE_RUN_ROOT/$RUN_ID"
if [ "$IS_DRY_RUN" != "1" ]; then
  mkdirp "$RUN_DIR"
  log "Run directory: $RUN_DIR"
fi

run_suite_once() {
  local suite="$1"
  local iter="$2"
  local iter_dir="$3" # .../iter-<N>

  local suite_dir="$iter_dir/$suite"
  local go_raw="$suite_dir/go.txt"
  local ts_raw="$suite_dir/ts.txt"

  if [ -z "$RUNTIME" ] || [ "$RUNTIME" = "go" ]; then
    run_go_suite "$suite" "$iter" "$go_raw"
    if [ "$IS_DRY_RUN" != "1" ]; then
      write_codeblock_md "$suite_dir/go.md" "Go (${suite}, iter=${iter})" "$go_raw"
      rm -f "$go_raw"
    fi
  fi

  if [ -z "$RUNTIME" ] || [ "$RUNTIME" = "ts" ]; then
    run_ts_suite "$suite" "$iter" "$ts_raw"
    if [ "$IS_DRY_RUN" != "1" ]; then
      write_codeblock_md "$suite_dir/ts.md" "TypeScript (${suite}, iter=${iter})" "$ts_raw"
      rm -f "$ts_raw"
    fi
  fi
}

RPC_FIXED_ITERS=(1 5 10 15)

contains_iter() {
  local needle="$1"
  shift
  for x in "$@"; do
    if [ "$x" = "$needle" ]; then
      return 0
    fi
  done
  return 1
}

if [ "$TYPE" = "single" ]; then
  if is_rpc_suite "$BENCHNAME"; then
    # Hard override: call + multicall always run at fixed iteration levels.
    for iter in "${RPC_FIXED_ITERS[@]}"; do
      iter_dir="$RUN_DIR/iter-$iter"
      if [ "$IS_DRY_RUN" != "1" ]; then
        mkdirp "$iter_dir"
      fi
      truncate_results_files
      run_suite_once "$BENCHNAME" "$iter" "$iter_dir"
      if [ -z "$RUNTIME" ]; then
        snapshot_iter_overall "$iter_dir"
      fi
    done
  else
    # CPU suite: use iteration ladder from --iter/--count.
    for ((lvl=0; lvl<LEVELS; lvl++)); do
      iter=$((BASE_ITER * (10 ** lvl)))
      iter_dir="$RUN_DIR/iter-$iter"
      if [ "$IS_DRY_RUN" != "1" ]; then
        mkdirp "$iter_dir"
      fi
      truncate_results_files
      run_suite_once "$BENCHNAME" "$iter" "$iter_dir"
      if [ -z "$RUNTIME" ]; then
        snapshot_iter_overall "$iter_dir"
      fi
    done
  fi
else
  # full
  # CPU iteration ladder (driven by --iter/--count)
  CPU_ITERS=()
  for ((lvl=0; lvl<LEVELS; lvl++)); do
    CPU_ITERS+=($((BASE_ITER * (10 ** lvl))))
  done

  # Run on the union of CPU ladder iters + fixed RPC iters.
  ALL_ITERS="$(printf "%s\n" "${CPU_ITERS[@]}" "${RPC_FIXED_ITERS[@]}" | sort -n -u)"

  run_full_pass_for_runtime() {
    local rt="$1" # go|ts

    # Fresh Anvil per runtime pass for most similar fork/cache conditions.
    stop_anvil_if_owned
    ensure_anvil_fresh

    RUNTIME="$rt"
    while IFS= read -r iter; do
      [ -z "$iter" ] && continue
      iter_dir="$RUN_DIR/iter-$iter"
      if [ "$IS_DRY_RUN" != "1" ]; then
        mkdirp "$iter_dir"
      fi

      truncate_results_files

      if contains_iter "$iter" "${CPU_ITERS[@]}"; then
        for suite in "${SUITES[@]}"; do
          if is_cpu_suite "$suite"; then
            run_suite_once "$suite" "$iter" "$iter_dir"
          fi
        done
      fi

      if contains_iter "$iter" "${RPC_FIXED_ITERS[@]}"; then
        for suite in "${SUITES[@]}"; do
          if is_rpc_suite "$suite"; then
            run_suite_once "$suite" "$iter" "$iter_dir"
          fi
        done
      fi

      snapshot_iter_runtime_results "$iter_dir" "$rt"
    done <<<"$ALL_ITERS"

    RUNTIME=""
    stop_anvil_if_owned
  }

  run_full_pass_for_runtime go
  run_full_pass_for_runtime ts

  # Now that both runtimes have produced results for each iter, generate comparisons.
  if [ "$IS_DRY_RUN" != "1" ]; then
    while IFS= read -r iter; do
      [ -z "$iter" ] && continue
      iter_dir="$RUN_DIR/iter-$iter"
      generate_iter_overall_reports "$iter_dir"
      log "Saved overall report to: $iter_dir/_overall"
    done <<<"$ALL_ITERS"
  fi
fi

# Final aggregate report across iteration levels (only meaningful if we generated overall comparisons)
if [ "$IS_DRY_RUN" != "1" ] && [ -z "$RUNTIME" ]; then
  if command -v bun >/dev/null 2>&1; then
    (cd "$BENCH_DIR" && bun run summarize-iter-runs.ts --run-dir "$RUN_DIR") || warn "Failed to generate summary.md"
  else
    warn "bun not found; skipping summary.md generation"
  fi
fi

log "Done."

