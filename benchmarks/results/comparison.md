# Benchmark Comparison: viem-go vs viem TypeScript

Generated: 2026-02-05T01:15:46.095Z

## Results

| Benchmark | Go (ns/op) | TS (ns/op) | Ratio | Winner |
|-----------|------------|------------|-------|--------|
| Call_Basic | 1,425,588,875 | 294,696 | 4837.49x | **TS** |
| Call_Basic | 141,127 | 294,696 | 0.48x | **Go** |
| Call_Basic | 142,534 | 294,696 | 0.48x | **Go** |
| Call_Basic | 141,448 | 294,696 | 0.48x | **Go** |
| Call_Basic | 142,666 | 294,696 | 0.48x | **Go** |
| Call_WithData | 138,816 | 254,411 | 0.55x | **Go** |
| Call_WithData | 137,307 | 254,411 | 0.54x | **Go** |
| Call_WithData | 137,936 | 254,411 | 0.54x | **Go** |
| Call_WithData | 137,105 | 254,411 | 0.54x | **Go** |
| Call_WithData | 137,494 | 254,411 | 0.54x | **Go** |
| Call_WithAccount | 142,945 | 239,621 | 0.60x | **Go** |
| Call_WithAccount | 142,907 | 239,621 | 0.60x | **Go** |
| Call_WithAccount | 142,193 | 239,621 | 0.59x | **Go** |
| Call_WithAccount | 142,806 | 239,621 | 0.60x | **Go** |
| Call_WithAccount | 142,568 | 239,621 | 0.59x | **Go** |
| Call_Decimals | 137,613 | 238,559 | 0.58x | **Go** |
| Call_Decimals | 136,691 | 238,559 | 0.57x | **Go** |
| Call_Decimals | 137,650 | 238,559 | 0.58x | **Go** |
| Call_Decimals | 137,653 | 238,559 | 0.58x | **Go** |
| Call_Decimals | 137,998 | 238,559 | 0.58x | **Go** |
| Call_Symbol | 175,820 | 224,991 | 0.78x | **Go** |
| Call_Symbol | 144,562 | 224,991 | 0.64x | **Go** |
| Call_Symbol | 140,642 | 224,991 | 0.63x | **Go** |
| Call_Symbol | 141,426 | 224,991 | 0.63x | **Go** |
| Call_Symbol | 140,212 | 224,991 | 0.62x | **Go** |
| Call_BalanceOfMultiple | 139,714 | 240,102 | 0.58x | **Go** |
| Call_BalanceOfMultiple | 139,675 | 240,102 | 0.58x | **Go** |
| Call_BalanceOfMultiple | 137,496 | 240,102 | 0.57x | **Go** |
| Call_BalanceOfMultiple | 137,629 | 240,102 | 0.57x | **Go** |
| Call_BalanceOfMultiple | 138,143 | 240,102 | 0.58x | **Go** |

## Summary

- Go wins: 29
- TS wins: 1
- Ties: 0

## Notes

- Ratio > 1.0x means TypeScript is faster
- Ratio < 1.0x means Go is faster
- Benchmarks run against the same Anvil instance for fair comparison
