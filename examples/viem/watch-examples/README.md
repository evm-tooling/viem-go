# Watch Actions Examples (viem TypeScript)

This directory contains TypeScript examples demonstrating viem's watch actions for real-time blockchain monitoring.

This is the TypeScript equivalent of the Go `viem-go/watch-examples/main.go`.

## Prerequisites

- Node.js 18+
- npm or bun

## Installation

```bash
npm install
# or
bun install
```

## Running Examples

### Watch Block Number

Monitors incoming block numbers with polling:

```bash
npm run watch-block-number
# or
npx ts-node index.ts watch-block-number
```

### Watch Blocks

Monitors incoming blocks with full block data:

```bash
npm run watch-blocks
# or
npx ts-node index.ts watch-blocks
```

### Watch Pending Transactions

Monitors pending transactions (requires RPC support):

```bash
npm run watch-pending-tx
# or
npx ts-node index.ts watch-pending-tx
```

### Watch Event

Monitors raw event logs (e.g., ERC20 Transfer events):

```bash
npm run watch-event
# or
npx ts-node index.ts watch-event
```

### Watch Contract Event

Monitors contract events with full ABI decoding:

```bash
npm run watch-contract-event
# or
npx ts-node index.ts watch-contract-event
```

### Run All Examples

Runs all examples with timeouts:

```bash
npm run all
# or
npx ts-node index.ts all
```

## API Comparison: TypeScript vs Go

| Feature | TypeScript (viem) | Go (viem-go) |
|---------|-------------------|--------------|
| Result delivery | Callbacks (`onBlockNumber`, `onLogs`) | Channels (`<-chan Event`) |
| Cancellation | `unwatch()` function | `context.Context` cancellation |
| Async model | Event-driven callbacks | Goroutines with channels |
| Error handling | `onError` callback | `Event.Error` field in channel |
| Batching | `batch: true` option | `Batch: true` parameter |

## TypeScript Pattern

```typescript
const unwatch = client.watchBlockNumber({
  emitOnBegin: true,
  pollingInterval: 2_000,
  onBlockNumber: (blockNumber) => {
    console.log(`Block: ${blockNumber}`)
  },
  onError: (error) => {
    console.log(`Error: ${error.message}`)
  },
})

// Later: stop watching
unwatch()
```

## Go Pattern

```go
events := client.WatchBlockNumber(ctx, public.WatchBlockNumberParameters{
    EmitOnBegin: true,
})

for event := range events {
    if event.Error != nil {
        log.Printf("Error: %v", event.Error)
        continue
    }
    fmt.Printf("Block: %d\n", event.BlockNumber)
}
```

## Notes

- Some RPC providers may not support all features (e.g., pending transaction filters)
- WebSocket endpoints provide lower latency than HTTP polling
- Press Ctrl+C for graceful shutdown
