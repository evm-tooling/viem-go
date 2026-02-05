# Transaction Examples

This example demonstrates transaction-related public actions in viem-go.

## Actions Covered

### GetBlockNumber
Returns the current block number with optional caching.

```go
blockNumber, err := public.GetBlockNumber(ctx, client, public.GetBlockNumberParameters{})
```

### GetTransaction
Retrieves transaction details by hash, block number + index, or block hash + index.

```go
// By hash
tx, err := public.GetTransaction(ctx, client, public.GetTransactionParameters{
    Hash: &txHash,
})

// By block number and index
tx, err := public.GetTransaction(ctx, client, public.GetTransactionParameters{
    BlockNumber: &blockNum,
    Index:       &index,
})
```

### GetTransactionReceipt
Returns the transaction receipt including logs, status, and gas usage.

```go
receipt, err := public.GetTransactionReceipt(ctx, client, public.GetTransactionReceiptParameters{
    Hash: txHash,
})

if receipt.IsSuccess() {
    fmt.Println("Transaction succeeded!")
}
```

### GetTransactionConfirmations
Returns the number of confirmations (blocks passed) since the transaction was mined.

```go
confirmations, err := public.GetTransactionConfirmations(ctx, client, public.GetTransactionConfirmationsParameters{
    Hash: &txHash,
})
```

### WaitForTransactionReceipt
Polls until the transaction is mined and returns the receipt. Supports:
- Configurable confirmation count
- Transaction replacement detection
- Timeout and polling interval

```go
receipt, err := public.WaitForTransactionReceipt(ctx, client, public.WaitForTransactionReceiptParameters{
    Hash:            txHash,
    Confirmations:   12,              // Wait for finality
    PollingInterval: 2 * time.Second,
    Timeout:         5 * time.Minute,
    OnReplaced: func(info public.ReplacementInfo) {
        fmt.Printf("Transaction replaced: %s\n", info.Reason)
    },
})
```

### FillTransaction
Fills missing transaction fields (nonce, gas, gas prices) via `eth_fillTransaction`.

```go
result, err := public.FillTransaction(ctx, client, public.FillTransactionParameters{
    Account: &from,
    To:      &to,
    Value:   big.NewInt(1e18), // 1 ETH
})

fmt.Printf("Nonce: %d\n", result.Transaction.Nonce)
fmt.Printf("Gas: %d\n", result.Transaction.Gas)
```

## Running the Example

```bash
cd examples/viem-go/transaction-examples
go run main.go
```

## Notes

- `eth_fillTransaction` is not supported by all RPC providers. It's primarily available on Geth-based nodes.
- `WaitForTransactionReceipt` supports replacement detection for transactions that are sped up or cancelled.
- Transaction confirmations return 0 for pending transactions.
