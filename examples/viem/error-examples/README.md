# Error Path Examples (viem TypeScript)

Tests various error failure paths for call/contract operations.

## Run

```bash
cd examples/viem/error-examples
bun run start
```

### Custom Error Extraction

Showcases `err.walk()` for extracting `ContractFunctionRevertedError` and `data.errorName`:

```bash
bun run custom-error
```

Or from the examples root:

```bash
cd examples/viem
bun run error-examples
```

## Error Scenarios Tested

1. **Invalid Params: code + to** - Mutually exclusive parameters
2. **Invalid Params: code + factory** - Mutually exclusive parameters
3. **Contract Revert: transfer exceeds balance** - ERC20 transfer with insufficient balance (ContractFunctionExecutionError, ContractFunctionRevertedError)
4. **Contract Revert: transfer from zero address** - Simulate transfer as 0x0...0 (ContractFunctionRevertedError)
5. **Call to non-contract address** - EOA with random calldata (CallExecutionError)
6. **Invalid selector** - Non-existent function on contract

## Error Types Demonstrated

- `ContractFunctionExecutionError` - Contract call failed (wraps revert)
- `ContractFunctionRevertedError` - Decoded revert (reason, signature, data.errorName)
- `CallExecutionError` - Wraps transport errors with call context

### Custom Error Handling (viem)

Use `BaseError.walk()` to find `ContractFunctionRevertedError` and branch on `data.errorName`:

```ts
if (err instanceof BaseError) {
  const revertError = err.walk((e) => e instanceof ContractFunctionRevertedError)
  if (revertError instanceof ContractFunctionRevertedError) {
    const errorName = revertError.data?.errorName ?? ''
    // branch on errorName: "Error", "Panic", "ERC20InsufficientBalance", etc.
  }
}
```
