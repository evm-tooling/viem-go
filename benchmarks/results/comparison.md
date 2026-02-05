# Benchmark Comparison: viem-go vs viem TypeScript

Generated: 2026-02-05T01:59:05.015Z

## Overall Summary

**🏆 Go is 1.60x faster overall**

| Metric | Go | TypeScript |
|--------|----|-----------|
| Avg ns/op | 141,771 | 226,825 |
| Avg ops/s | 7,054 | 4,409 |
| Wins | 30/30 | 0/30 |

## Detailed Results

| Benchmark | Go (ns/op) | TS (ns/op) | Go (ops/s) | TS (ops/s) | Result |
|-----------|------------|------------|------------|------------|--------|
| Call_Basic | 154,375 | 222,725 | 6,478 | 4,490 | 🟢 Go 1.44x faster |
| Call_Basic | 142,532 | 222,725 | 7,016 | 4,490 | 🟢 Go 1.56x faster |
| Call_Basic | 142,017 | 222,725 | 7,041 | 4,490 | 🟢 Go 1.57x faster |
| Call_Basic | 142,123 | 222,725 | 7,036 | 4,490 | 🟢 Go 1.57x faster |
| Call_Basic | 141,330 | 222,725 | 7,076 | 4,490 | 🟢 Go 1.58x faster |
| Call_WithData | 139,134 | 218,133 | 7,187 | 4,584 | 🟢 Go 1.57x faster |
| Call_WithData | 138,392 | 218,133 | 7,226 | 4,584 | 🟢 Go 1.58x faster |
| Call_WithData | 138,059 | 218,133 | 7,243 | 4,584 | 🟢 Go 1.58x faster |
| Call_WithData | 142,945 | 218,133 | 6,996 | 4,584 | 🟢 Go 1.53x faster |
| Call_WithData | 139,468 | 218,133 | 7,170 | 4,584 | 🟢 Go 1.56x faster |
| Call_WithAccount | 145,126 | 221,723 | 6,891 | 4,510 | 🟢 Go 1.53x faster |
| Call_WithAccount | 144,207 | 221,723 | 6,934 | 4,510 | 🟢 Go 1.54x faster |
| Call_WithAccount | 143,599 | 221,723 | 6,964 | 4,510 | 🟢 Go 1.54x faster |
| Call_WithAccount | 153,601 | 221,723 | 6,510 | 4,510 | 🟢 Go 1.44x faster |
| Call_WithAccount | 144,631 | 221,723 | 6,914 | 4,510 | 🟢 Go 1.53x faster |
| Call_Decimals | 138,271 | 217,033 | 7,232 | 4,608 | 🟢 Go 1.57x faster |
| Call_Decimals | 138,338 | 217,033 | 7,229 | 4,608 | 🟢 Go 1.57x faster |
| Call_Decimals | 138,790 | 217,033 | 7,205 | 4,608 | 🟢 Go 1.56x faster |
| Call_Decimals | 136,389 | 217,033 | 7,332 | 4,608 | 🟢 Go 1.59x faster |
| Call_Decimals | 139,792 | 217,033 | 7,153 | 4,608 | 🟢 Go 1.55x faster |
| Call_Symbol | 141,822 | 244,578 | 7,051 | 4,089 | 🟢 Go 1.72x faster |
| Call_Symbol | 142,393 | 244,578 | 7,023 | 4,089 | 🟢 Go 1.72x faster |
| Call_Symbol | 142,481 | 244,578 | 7,018 | 4,089 | 🟢 Go 1.72x faster |
| Call_Symbol | 146,772 | 244,578 | 6,813 | 4,089 | 🟢 Go 1.67x faster |
| Call_Symbol | 143,502 | 244,578 | 6,969 | 4,089 | 🟢 Go 1.70x faster |
| Call_BalanceOfMultiple | 138,520 | 236,758 | 7,219 | 4,224 | 🟢 Go 1.71x faster |
| Call_BalanceOfMultiple | 138,658 | 236,758 | 7,212 | 4,224 | 🟢 Go 1.71x faster |
| Call_BalanceOfMultiple | 138,659 | 236,758 | 7,212 | 4,224 | 🟢 Go 1.71x faster |
| Call_BalanceOfMultiple | 138,504 | 236,758 | 7,220 | 4,224 | 🟢 Go 1.71x faster |
| Call_BalanceOfMultiple | 138,697 | 236,758 | 7,210 | 4,224 | 🟢 Go 1.71x faster |

## By Category

- 🟢 **Basic Calls**: Go 1.54x faster
- 🟢 **With Parameters**: Go 1.56x faster
- 🟢 **With Account**: Go 1.52x faster
- 🟢 **Other**: Go 1.64x faster
- 🟢 **Batch Operations**: Go 1.71x faster

## Win Summary

- 🟢 Go wins: 30 (100%)
- 🔵 TS wins: 0 (0%)
- ⚪ Ties: 0 (0%)

## Notes

- Benchmarks run against the same Anvil instance (mainnet fork) for fair comparison
- ns/op = nanoseconds per operation (lower is better)
- ops/s = operations per second (higher is better)
- 🟢 = Go faster, 🔵 = TS faster, ⚪ = Similar (within 5%)
