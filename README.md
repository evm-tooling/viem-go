# viem-go

Go Interface for Ethereum -- inspired by [viem](https://viem.sh)

[![CI](https://github.com/ChefBingbong/viem-go/actions/workflows/ci.yml/badge.svg)](https://github.com/ChefBingbong/viem-go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/ChefBingbong/viem-go.svg)](https://pkg.go.dev/github.com/ChefBingbong/viem-go)

> **Note:** This project is under active development. APIs may change.

## Why viem-go?

If you've worked with Ethereum in Go, you know the pain. [go-ethereum](https://github.com/ethereum/go-ethereum) (ethclient) is powerful but its contract interaction story is rough -- `abigen` generates thousands of lines of boilerplate, there's no built-in multicall, and even simple reads require wiring together low-level primitives. On the other side, [viem](https://github.com/wevm/viem) in TypeScript nails the developer experience but you're stuck in a single-threaded runtime.

**viem-go sits in the middle.** It brings viem's clean, composable API design into Go -- giving you the ergonomics of a modern Ethereum library backed by Go's concurrency, type safety, and raw speed. You get the simplicity of `readContract` and `multicall` without the ceremony of `abigen`, `bind.NewBoundContract`, and manual ABI packing that go-ethereum demands.

### The Complexity Problem

Here's what reading an ERC20 balance looks like across the three approaches:

| | Lines of setup | ABI handling | Type safety | Multicall |
|---|---|---|---|---|
| **go-ethereum** | ~15-20 | `abigen` codegen or manual ABI packing | Codegen only | Manual |
| **viem (TS)** | ~5 | Inline with `parseAbi` | Runtime | Built-in |
| **viem-go** | ~5 | Inline JSON, typed templates, or codegen | Compile-time generics | Built-in |

viem-go gives you three tiers of contract interaction -- raw `ReadContract` for quick scripts, typed `Fn`/`Call` generics for production code, and full codegen via `viemgen` for large projects -- all without ever touching `abigen`.

## Performance

viem-go outperforms the TypeScript viem library across 59 benchmarks spanning 9 test suites -- winning 54 of them (92%). The benchmarks cover both network-bound RPC operations and pure CPU-bound utilities, giving a complete picture of where Go's advantages lie.

![Summary](benchmarks/results/charts/summary.svg)

### Speedup by Benchmark

![Speedup](benchmarks/results/charts/speedup.svg)

### Latency Comparison

![Latency](benchmarks/results/charts/latency.svg)

**Key takeaways:**

- **ABI encoding/decoding:** Go is **10-26x faster**. `encodeFunctionData` runs in ~300ns in Go vs ~8us in TypeScript. This matters at scale -- indexers and bots encoding thousands of calls per second see a direct throughput improvement.
- **Hashing (keccak256, sha256):** Go is **9-21x faster**. Go's native `crypto/sha3` and `crypto/sha256` implementations outperform the JS WASM/native bindings. The gap widens with larger inputs (16.7x for 1KB keccak256).
- **Signature recovery/verification:** Go is **63-85x faster**. ECDSA operations are computationally heavy and Go's secp256k1 implementation (via go-ethereum) dominates -- `recoverMessageAddress` completes in ~27us vs ~2.3ms in TypeScript.
- **Event log decoding:** Go is **26x faster**. Decoding 100 Transfer events takes ~39us in Go vs ~998us in TypeScript -- critical for indexers processing blocks of logs.
- **ENS (namehash, normalize):** Go is **3-14x faster**. Recursive keccak hashing in `namehash` and Unicode normalization in `normalize` both benefit from Go's efficient string and hash handling.
- **Call actions:** Go is **95-99x faster** for standard contract reads (`eth_call`). Go averages ~0.19ms per call vs ~18ms in TypeScript due to Node.js event loop overhead.
- **Multicall batching:** Go is **1.6-5.6x faster** across batch sizes, widening as batch size grows because Go fans out chunked RPC requests in parallel goroutines.
- **Unit parsing:** Go is **1.5-2.8x faster** for common operations like `parseEther("1.5")` and `parseGwei("20.5")`, using a zero-regex validator and `uint64` fast path that avoids `big.Int` overhead for values that fit in 64 bits. TypeScript still wins on very large decimal strings where V8's native `BigInt` construction is hard to beat.

> 59 benchmarks across 9 suites on Apple M4 Pro against a shared Anvil instance (mainnet fork). Go wins all 9 suites. See [`benchmarks/`](benchmarks/) for full methodology.

### Real-World Comparison: UniswapV2 Pool Extractor

To see how these library-level benchmarks translate into a real application, check out [**viem-go-extractor-demo**](https://github.com/ChefBingbong/viem-go-extractor-demo) -- identical UniswapV2 pool extractors built in both Go (viem-go) and TypeScript (viem + Bun).

Both extractors do the same work: sync 60,000+ pool pairs from on-chain UniswapV2 factories via multicall, decode Sync events from block logs, resolve ERC-20 token metadata, and serve the pool state over an HTTP API. The results show where each language shines:

| Dimension | Winner | Key Finding |
|---|---|---|
| **Multicall (200 contracts)** | Go | **3.5x faster** -- Go's encoding scales linearly while TS shows quadratic-ish overhead |
| **Event decoding** | Go | **9-11x faster** -- compiled byte slicing vs full ABI schema resolution |
| **Factory sync (end-to-end)** | Go | **1.2-2.2x faster** -- compounds across 60K+ pools |
| **JSON serialization** | TypeScript | **1.7-2.3x faster** -- Bun's Zig-optimized `JSON.stringify` |
| **HTTP API under load** | TypeScript | **5-7x lower latency** -- Bun's async I/O model handles 200 concurrent users efficiently |

**The takeaway:** Go (viem-go) is the better choice for the **data pipeline** -- fetching, decoding, and processing on-chain data. TypeScript (viem) is the better choice for the **API layer** -- serving that data to clients. In a production architecture, the optimal design is a Go-based indexer feeding data into a TypeScript API server.

## Installation

```bash
go get github.com/ChefBingbong/viem-go
```

## Quick Start

**Creating a client and reading a contract:**

```go
import (
    "context"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ChefBingbong/viem-go/client"
    "github.com/ChefBingbong/viem-go/contracts/erc20"
)

c, _ := client.NewClient("https://eth.llamarpc.com")
defer c.Close()

// Simple typed binding
usdc := erc20.MustNew(
    common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"), c,
)

name, _ := usdc.Name(context.Background())        // "USD Coin"
balance, _ := usdc.BalanceOf(context.Background(), 
    common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"),
)
```

**Side-by-side with viem (TypeScript):**

```typescript
import { createPublicClient, http } from 'viem'
import { mainnet } from 'viem/chains'

const client = createPublicClient({
  chain: mainnet,
  transport: http(),
})

const balance = await client.readContract({
  address: '0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48',
  abi: parseAbi(['function balanceOf(address) view returns (uint256)']),
  functionName: 'balanceOf',
  args: ['0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045'],
})
```

```go
import (
    "context"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ChefBingbong/viem-go/contracts/erc20"
    "github.com/ChefBingbong/viem-go/client"
)

c, _ := client.NewClient("https://eth.llamarpc.com")
usdc := erc20.MustNew(common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"), c)

balance, _ := usdc.BalanceOf(ctx, common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"))
```

**Parsing units:**

```typescript
import { parseEther, formatEther } from 'viem'

const wei = parseEther('1.5')   // 1500000000000000000n
const eth = formatEther(wei)    // "1.5"
```

```go
import "github.com/ChefBingbong/viem-go/utils/unit"

wei, _ := unit.ParseEther("1.5")    // *big.Int: 1500000000000000000
eth := unit.FormatEther(wei)        // "1.5"
```

## Typed Contract Templates

One of viem-go's unique features is its **typed contract template system** -- something that has no equivalent in either viem (TypeScript) or go-ethereum.

### The Problem

In TypeScript, viem uses runtime ABI parsing and infers types dynamically. That works because TypeScript has powerful type inference from string literals. Go doesn't have that -- so you're normally stuck choosing between:

1. **go-ethereum's `abigen`** -- generates huge files of boilerplate, tightly coupled to geth's internal types, hard to customize.
2. **Raw `[]any` returns** -- call a function, get back `interface{}`, cast everything manually, pray you got the types right.

### The Solution: `Fn` Generics + `viemgen`

viem-go solves this with two complementary approaches:

#### 1. Typed Function Descriptors (Zero Codegen)

Define your contract methods as typed descriptors using Go generics. The compiler enforces both argument types and return types at build time:

```go
import "github.com/ChefBingbong/viem-go/contract"

// Define once -- these are just type descriptors, no codegen needed
var (
    Name      = contract.Fn[string]{Name: "name"}
    BalanceOf = contract.Fn1[common.Address, *big.Int]{Name: "balanceOf"}
    Allowance = contract.Fn2[common.Address, common.Address, *big.Int]{Name: "allowance"}
)

// Bind a contract
token, _ := contract.Bind(tokenAddr, abiJSON, client)

// Fully type-safe calls -- wrong types won't compile
name, err := contract.Call(token, ctx, Name)
balance, err := contract.Call1(token, ctx, BalanceOf, ownerAddr)
allowance, err := contract.Call2(token, ctx, Allowance, owner, spender)
```

The `Fn`, `Fn1`, `Fn2`, `Fn3`, and `Fn4` types encode the argument count and types into the Go type system. `Call1` expects exactly one argument of the type specified in `Fn1` -- pass the wrong type and it's a compile error, not a runtime panic.

#### 2. `viemgen` Code Generator

For larger projects, the `viemgen` CLI generates complete typed bindings from an ABI JSON file:

```bash
# Initialize the directory structure
go run ./cmd/viemgen init

# Place your ABI JSON file
cp MyToken.json _contracts_typed/json/mytoken.json

# Generate typed Go bindings
go run ./cmd/viemgen --pkg mytoken
```

This produces a Go package with:
- Typed method descriptors (`Methods.Name`, `Methods.BalanceOf`, etc.)
- A contract binding struct with methods for every ABI function
- Pre-parsed ABI caching (parsed once, reused across calls)
- Write method helpers with gas estimation
- Event types and parsing

The generated code uses the same `Fn`/`Call` pattern under the hood, so it composes naturally with multicall and other viem-go features.

#### Built-in ERC Standards

Common token standards ship out of the box:

```go
import "github.com/ChefBingbong/viem-go/contracts/erc20"

// Classic binding API
token := erc20.MustNew(usdcAddr, client)
name, _ := token.Name(ctx)
balance, _ := token.BalanceOf(ctx, owner)

// Or use typed descriptors directly
token2 := contract.MustBind(usdcAddr, []byte(erc20.ContractABI), client)
name, _ := contract.Call(token2, ctx, erc20.Methods.Name)
balance, _ := contract.Call1(token2, ctx, erc20.Methods.BalanceOf, owner)
```

## Implementation Status

Early -- core client, public actions (`call`, `multicall`, `getBlockNumber`, `getBalance`, etc.), contract bindings, ABI encoding/decoding, unit parsing, hashing, and signature utilities are implemented.

## Development

```bash
# Install dependencies
go mod download

# Run tests
make test

# Run linter
make lint

# Run benchmarks (requires Foundry for Anvil)
cd benchmarks && make bench
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Acknowledgments

- [viem](https://github.com/wevm/viem) -- The original TypeScript library this project is inspired by
- [go-ethereum](https://github.com/ethereum/go-ethereum) -- Ethereum Go implementation used for cryptographic primitives
