# viem-go Examples

This folder contains example scripts demonstrating viem-go usage, alongside equivalent TypeScript viem examples for comparison.

## Prerequisites

- Go 1.21+
- Bun (for TypeScript examples)

## Running Examples

### Go Examples

```bash
# Run all Go examples
bun run go:all

# Run individual examples
bun run go:client        # Basic client usage
bun run go:read-contract # Read from ERC20 contract
bun run go:gas           # Get gas prices
bun run go:block         # Get block information

# Or run directly with Go
go run ./cmd -example=all
go run ./cmd -example=client
```

### TypeScript Examples (viem)

```bash
# Run all TypeScript examples
bun run ts:all

# Run individual examples
bun run ts:client           # Basic client usage
bun run ts:read-contract    # Read from ERC20 contract
bun run ts:send-transaction # Send a transaction (requires PK)
```

### Run Both

```bash
bun run all
```

## Environment Variables

Some examples require environment variables. Create a `.env` file:

```env
TENDERLY_RPC_URL=https://your-rpc-url
PK=your-private-key-for-transactions
```

## Examples

| Example | Go | TypeScript | Description |
|---------|----|-----------:|-------------|
| Client | `client_example.go` | `_viem_ts_client.ts` | Basic client: block number, balance |
| Read Contract | `read_contract_example.go` | `_viem_ts_read_contract.ts` | Read ERC20 token info |
| Gas Prices | `gas_example.go` | - | Get current gas prices |
| Block Info | `block_example.go` | - | Get block information |
| Send Transaction | - | `_viem_ts_send_transaction.ts` | Send ETH transaction |
