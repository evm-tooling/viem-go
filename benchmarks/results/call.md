# Benchmark Comparison: viem-go vs viem TypeScript

Generated: 2026-02-07T00:44:51.028Z

## Overall Summary

**🏆 Go is 75.53x faster overall**

| Metric | Go | TypeScript |
|--------|----|-----------|
| Avg ns/op | 201,958 | 15,254,048 |
| Avg ops/s | 4,952 | 66 |
| Wins | 6/6 | 0/6 |

## Detailed Results

| Benchmark | Go (ns/op) | TS (ns/op) | Go (ops/s) | TS (ops/s) | Result |
|-----------|------------|------------|------------|------------|--------|
| Call_Basic | 188,577 | 18,483,844 | 5,303 | 54 | 🟢 Go 98.02x faster |
| Call_WithData | 204,209 | 18,795,827 | 4,897 | 53 | 🟢 Go 92.04x faster |
| Call_WithAccount | 184,555 | 260,339 | 5,418 | 3,841 | 🟢 Go 1.41x faster |
| Call_Decimals | 177,790 | 17,940,985 | 5,625 | 56 | 🟢 Go 100.91x faster |
| Call_Symbol | 180,639 | 17,762,463 | 5,536 | 56 | 🟢 Go 98.33x faster |
| Call_BalanceOfMultiple | 275,976 | 18,280,834 | 3,624 | 55 | 🟢 Go 66.24x faster |

## Win Summary

- 🟢 Go wins: 6 (100%)
- 🔵 TS wins: 0 (0%)
- ⚪ Ties: 0 (0%)

## Notes

- Benchmarks run against the same Anvil instance (mainnet fork) for fair comparison
- ns/op = nanoseconds per operation (lower is better)
- ops/s = operations per second (higher is better)
- 🟢 = Go faster, 🔵 = TS faster, ⚪ = Similar (within 5%)
