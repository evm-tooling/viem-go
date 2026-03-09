# Error Path Examples (viem TypeScript)

Tests various error failure paths for call/contract operations.

## Run

```bash
cd examples/viem/error-examples
bun run start
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
- `ContractFunctionRevertedError` - Decoded revert (reason, signature)
- `CallExecutionError` - Wraps transport errors with call context
