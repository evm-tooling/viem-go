# Error Path Examples (viem-go)

Tests various error failure paths for call/contract operations.

## Run

```bash
cd examples/viem-go/error-examples
go run main.go
```

## Error Scenarios Tested

1. **Invalid Params: code + to** - Mutually exclusive parameters (InvalidCallParamsError)
2. **Invalid Params: code + factory** - Mutually exclusive parameters (InvalidCallParamsError)
3. **Contract Revert: transfer exceeds balance** - ERC20 transfer with insufficient balance (ContractFunctionExecutionError, ContractFunctionRevertedError)
4. **Contract Revert: transfer from zero address** - Simulate transfer as 0x0...0 (ContractFunctionRevertedError)
5. **Call to non-contract address** - EOA with random calldata (CallExecutionError)
6. **Invalid selector** - Non-existent function on contract

## Error Types Demonstrated

- `InvalidCallParamsError` - Invalid parameter combinations
- `ContractFunctionExecutionError` - Contract call failed (wraps revert)
- `ContractFunctionRevertedError` - Decoded revert (reason, signature)
- `CallExecutionError` - Wraps transport errors with call context
