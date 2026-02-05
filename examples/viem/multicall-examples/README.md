# Multicall Examples (viem TypeScript)

Comprehensive examples demonstrating the multicall action in viem TypeScript.

## Features Demonstrated

- **Basic Multicall**: Batch multiple contract calls into a single RPC request
- **Token Metadata**: Retrieve name, symbol, decimals, and supply in one call
- **Balance Queries**: Check multiple token balances efficiently
- **Error Handling**: `allowFailure` mode for graceful error handling
- **Batching**: Automatic batching of large multicalls
- **Historical Queries**: Query at specific block numbers
- **Deployless Multicall**: Execute without deployed multicall3 contract

## Running

```bash
# Install dependencies
bun install

# Run examples
bun run start

# Or with watch mode
bun run dev
```

## Example Output

```
=== Basic Multicall - Read Multiple Token Names ===
Token names retrieved in single RPC call:
  USDC: USD Coin
  WETH: Wrapped Ether
  WMATIC: Wrapped Matic

=== Large Multicall with Automatic Batching ===
Executed 30 calls successfully
BatchSize: 512 bytes (forces multiple batches)
```

## Configuration Options

```typescript
await client.multicall({
  contracts: [...],           // Array of contract calls
  allowFailure: true,         // Continue on individual failures (default: true)
  batchSize: 1024,            // Max bytes per batch (default: 1024)
  deployless: true,           // Use bytecode execution
  blockNumber: 12345678n,     // Historical block query
  multicallAddress: '0x...',  // Custom multicall3 address
})
```
