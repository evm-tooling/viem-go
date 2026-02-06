# Wallet Actions Examples (viem TypeScript)

Comprehensive examples demonstrating all main viem wallet actions on **Polygon mainnet**.
Each action that sends a transaction **simulates before executing**.

## Prerequisites

- [Bun](https://bun.sh) installed

## File Structure

| File | Action |
|------|--------|
| `shared.ts` | Shared setup — clients, account, helpers |
| `sign_message.ts` | `signMessage` — EIP-191 (string, raw hex, raw bytes) |
| `sign_typed_data.ts` | `signTypedData` — EIP-712 structured data |
| `sign_transaction.ts` | `signTransaction` — EIP-1559, legacy, EIP-2930 |
| `sign_authorization.ts` | `signAuthorization` — EIP-7702 |
| `send_transaction.ts` | `sendTransaction` — simulate + send ETH |
| `send_raw_transaction.ts` | `sendRawTransaction` — local sign + broadcast |
| `write_contract.ts` | `writeContract` — simulate + ERC-20 transfer |
| `deploy_contract.ts` | `deployContract` — bytecode + constructor args |
| `send_calls.ts` | `sendCalls` — EIP-5792 batch calls |
| `main.ts` | Entry point with CLI arg routing |

## Running

```bash
# Install dependencies
bun install

# Run all examples
bun run start

# Run a single example
bun run main.ts signMessage

# Run multiple examples
bun run main.ts signMessage signTypedData signTransaction
```

## Available example names

```
signMessage, signTypedData, signTransaction, signAuthorization,
sendTransaction, sendRawTransaction, writeContract, deployContract,
sendCalls
```
