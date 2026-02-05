# Multicall Examples (viem-go)

Comprehensive examples demonstrating the Multicall action in viem-go.

## Features Demonstrated

- **Basic Multicall**: Batch multiple contract calls into a single RPC request
- **Token Metadata**: Retrieve name, symbol, decimals, and supply in one call
- **Balance Queries**: Check multiple token balances efficiently
- **Error Handling**: `allowFailure` mode for graceful error handling
- **Chunking**: Automatic batching of large multicalls
- **Parallel Execution**: Concurrent chunk execution with bounded concurrency
- **Historical Queries**: Query at specific block numbers
- **Deployless Multicall**: Execute without deployed multicall3 contract

## Go Concurrency Features

The multicall implementation leverages Go's concurrency primitives:

- `errgroup.WithContext` for parallel chunk execution
- `SetLimit()` for bounded concurrency (prevents RPC overload)
- Context propagation for cancellation support
- Index-based result assignment for ordered results

## Running

```bash
go run main.go
```

## Example Output

```
=== Basic Multicall - Read Multiple Token Names ===
Token names retrieved in single RPC call:
  USDC: USD Coin
  WETH: Wrapped Ether
  WMATIC: Wrapped Matic

=== Large Multicall with Automatic Chunking ===
Executed 30 calls successfully
BatchSize: 512 bytes (forces multiple chunks)
MaxConcurrentChunks: 4 (parallel execution)
```

## Configuration Options

```go
public.MulticallParameters{
    Contracts:           contracts,      // List of contract calls
    AllowFailure:        &allowFailure,  // Continue on individual failures (default: true)
    BatchSize:           1024,           // Max bytes per chunk (default: 1024)
    MaxConcurrentChunks: 4,              // Parallel chunk limit (default: 4)
    Deployless:          true,           // Use bytecode execution
    BlockNumber:         &blockNum,      // Historical block query
    MulticallAddress:    &addr,          // Custom multicall3 address
}
```
