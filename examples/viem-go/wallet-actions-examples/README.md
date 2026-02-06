# Wallet Actions Examples (viem-go)

Comprehensive examples demonstrating all main viem-go wallet actions on **Polygon mainnet**.
Each action that sends a transaction **simulates before executing**.

## Prerequisites

- Go 1.24+

## File Structure

| File | Action |
|------|--------|
| `shared.go` | Shared setup — clients, account, helpers |
| `sign_message.go` | `signMessage` — EIP-191 (JSON-RPC, local, raw) |
| `sign_typed_data.go` | `signTypedData` — EIP-712 structured data |
| `sign_transaction.go` | `signTransaction` — EIP-1559, legacy, EIP-2930 |
| `sign_authorization.go` | `signAuthorization` — EIP-7702 |
| `send_transaction.go` | `sendTransaction` — simulate + send ETH |
| `send_raw_transaction.go` | `sendRawTransaction` — local sign + broadcast |
| `write_contract.go` | `writeContract` — simulate + ERC-20 transfer |
| `deploy_contract.go` | `deployContract` — bytecode + constructor args |
| `send_calls.go` | `sendCalls` — EIP-5792 batch calls |
| `main.go` | Entry point with CLI arg routing |

## Running

```bash
# Run all examples
go run ./examples/viem-go/wallet-actions-examples/

# Run a single example
go run ./examples/viem-go/wallet-actions-examples/ signMessage

# Run multiple examples
go run ./examples/viem-go/wallet-actions-examples/ signMessage signTypedData signTransaction
```

## Available example names

```
signMessage, signTypedData, signTransaction, signAuthorization,
sendTransaction, sendRawTransaction, writeContract, deployContract,
sendCalls
```
