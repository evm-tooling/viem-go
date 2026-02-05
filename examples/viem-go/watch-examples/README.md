# Watch Actions Examples

This directory contains examples demonstrating the viem-go watch actions for real-time blockchain monitoring.

## Prerequisites

- Go 1.22 or later
- An Ethereum RPC endpoint (HTTP or WebSocket)

## Running Examples

You can set a custom RPC URL via environment variable:

```bash
export RPC_URL="https://your-rpc-endpoint.com"
```

### Watch Block Number

Monitors incoming block numbers:

```bash
go run main.go watch-block-number
```

### Watch Blocks

Monitors incoming blocks with full block data:

```bash
go run main.go watch-blocks
```

### Watch Pending Transactions

Monitors pending transactions (requires RPC support):

```bash
go run main.go watch-pending-tx
```

### Watch Event

Monitors raw event logs (e.g., ERC20 Transfer events):

```bash
go run main.go watch-event
```

### Watch Contract Event

Monitors contract events with ABI decoding:

```bash
go run main.go watch-contract-event
```

### Run All Examples

Runs all examples with timeouts:

```bash
go run main.go all
```

## Features Demonstrated

1. **Polling Mode**: Works with HTTP transports by periodically querying the RPC
2. **Subscription Mode**: Uses WebSocket subscriptions when available for real-time updates
3. **Channel-Based API**: All watch functions return channels for Go-idiomatic event consumption
4. **Context Cancellation**: Clean shutdown via context cancellation
5. **Error Handling**: Non-blocking error delivery through event structs
6. **Batching**: Optional batching for high-throughput scenarios
7. **Emit on Begin**: Option to emit the current state immediately
8. **Emit Missed**: Option to emit any blocks missed between polling intervals

## Notes

- Some RPC providers may not support all features (e.g., pending transaction filters)
- WebSocket endpoints provide lower latency than HTTP polling
- Press Ctrl+C for graceful shutdown
