# Cross-Language Benchmarks: viem-go vs viem TypeScript

This directory contains benchmarks to compare the performance of `viem-go` against the original `viem` TypeScript library.

## Prerequisites

- [Go 1.21+](https://golang.org/dl/)
- [Node.js 18+](https://nodejs.org/)
- [Foundry](https://getfoundry.sh/) (for Anvil)
- [Bun](https://bun.sh/) (optional, for comparison script)

Install Foundry:

```bash
curl -L https://foundry.paradigm.xyz | bash
foundryup
```

## Quick Start

Run all benchmarks with a shared Anvil instance:

```bash
make bench
```

This will:
1. Start Anvil with a mainnet fork at the latest block
2. Run Go benchmarks
3. Run TypeScript benchmarks
4. Display results

## Commands

| Command | Description |
|---------|-------------|
| `make bench` | Run all benchmarks with shared Anvil |
| `make bench-go` | Run only Go benchmarks (requires running Anvil) |
| `make bench-ts` | Run only TypeScript benchmarks (requires running Anvil) |
| `make compare` | Generate comparison report from results |
| `make install` | Install TypeScript dependencies |
| `make anvil` | Start Anvil manually (for development) |
| `make clean` | Clean results directory |

## Configuration

Environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `ANVIL_PORT` | `8545` | Port for Anvil RPC |
| `FORK_URL` | `https://rpc.ankr.com/eth` | Mainnet RPC for forking |
| `FORK_BLOCK` | `(latest)` | Block number to fork from (empty = latest) |
| `GO_BENCH_COUNT` | `5` | Number of Go benchmark iterations |

Example with custom configuration:

```bash
FORK_URL=https://eth-mainnet.g.alchemy.com/v2/YOUR_KEY make bench
```

## Directory Structure

```
benchmarks/
├── scripts/
│   └── anvil.sh              # Anvil lifecycle management
├── go/
│   ├── main_test.go          # TestMain with shared client setup
│   └── call_bench_test.go    # Call action benchmarks
├── typescript/
│   ├── package.json          # Dependencies (viem, vitest)
│   ├── vitest.config.ts      # Vitest benchmark config
│   └── call.bench.ts         # Call action benchmarks
├── results/                  # Benchmark output directory
│   └── .gitkeep
├── compare.ts                # Results comparison script
├── Makefile                  # Benchmark orchestration
└── README.md                 # This file
```

## Benchmarks

### Call Action

| Benchmark | Description |
|-----------|-------------|
| `Basic` | Simple `name()` call to USDC contract |
| `WithData` | `balanceOf(address)` call with encoded parameters |
| `WithAccount` | Call with specified sender (msg.sender) |
| `Decimals` | Read token decimals |
| `Symbol` | Read token symbol |
| `BalanceOfMultiple` | Multiple balanceOf calls with different addresses |

## Results

After running benchmarks, results are saved to:

- `results/go-results.txt` - Go benchmark output
- `results/ts-results.txt` - TypeScript benchmark output
- `results/comparison.md` - Comparison report (after `make compare`)

### Sample Output

```
==============================================================================
  BENCHMARK COMPARISON: viem-go vs viem TypeScript
==============================================================================

| Benchmark                    | Go (ns/op)    | TS (ns/op)    | Ratio   | Winner |
|-----------------------------|---------------|---------------|---------|--------|
| Call_Basic                  | 5,234,567     | 6,543,210     | 0.80x   | Go     |
| Call_WithData               | 5,456,789     | 6,789,012     | 0.80x   | Go     |
| Call_WithAccount            | 5,345,678     | 6,654,321     | 0.80x   | Go     |

Summary:
  Go wins:  3
  TS wins:  0
  Ties:     0
```

## Adding New Benchmarks

### Go

1. Create a new `*_bench_test.go` file in `go/`
2. Follow the pattern:

```go
func BenchmarkMyAction_Scenario(b *testing.B) {
    // Setup
    params := public.MyActionParameters{...}
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := public.MyAction(benchCtx, benchClient, params)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

### TypeScript

1. Create a new `*.bench.ts` file in `typescript/`
2. Follow the pattern:

```typescript
import { bench, describe } from 'vitest'

describe('MyAction', () => {
  bench('viem-ts: myAction (scenario)', async () => {
    await client.myAction({...})
  })
})
```

### Update Compare Script

Add the benchmark name mapping in `compare.ts`:

```typescript
const mapping: Record<string, string> = {
  BenchmarkMyAction_Scenario: 'viem-ts: myAction (scenario)',
  // ...
}
```

## Troubleshooting

### Anvil fails to start

Check if another process is using port 8545:

```bash
lsof -i :8545
```

Kill the process or use a different port:

```bash
ANVIL_PORT=8546 make bench
```

### Fork URL rate limited

Use a dedicated RPC endpoint:

```bash
FORK_URL=https://eth-mainnet.g.alchemy.com/v2/YOUR_KEY make bench
```

### TypeScript benchmarks not found

Ensure dependencies are installed:

```bash
make install
```

### No matching benchmarks in comparison

Verify benchmark names match between Go and TypeScript. The `compare.ts` script maps Go benchmark names to TypeScript names.

## Notes

- Both benchmark suites run against the same Anvil instance for fair comparison
- Anvil uses `--no-mining` for consistent block state
- Fork block is pinned for reproducibility
- Go benchmarks use `testing.B` with `-benchmem` for memory allocation info
- TypeScript benchmarks use vitest's `bench` function
