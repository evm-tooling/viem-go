# Benchmark Comparison: viem-go vs viem TypeScript

Generated: 2026-02-07T04:32:16.864Z

## Overall Summary

**🏆 Go is 3.26x faster overall**

| Metric | Go | TypeScript |
|--------|----|-----------|
| Avg ns/op | 13,172,358 | 43,007,142 |
| Avg ops/s | 76 | 23 |
| Wins | 16/16 | 0/16 |

## Detailed Results

| Benchmark | Go (ns/op) | TS (ns/op) | Go (ops/s) | TS (ops/s) | Result |
|-----------|------------|------------|------------|------------|--------|
| Multicall_Basic | 207,937 | 513,271 | 4,809 | 1,948 | 🟢 Go 2.47x faster |
| Multicall_WithArgs | 212,253 | 446,949 | 4,711 | 2,237 | 🟢 Go 2.11x faster |
| Multicall_MultiContract | 293,359 | 486,976 | 3,409 | 2,053 | 🟢 Go 1.66x faster |
| Multicall_10Calls | 279,306 | 598,365 | 3,580 | 1,671 | 🟢 Go 2.14x faster |
| Multicall_30Calls | 570,403 | 1,032,098 | 1,753 | 969 | 🟢 Go 1.81x faster |
| Multicall_Deployless | 392,310 | 634,007 | 2,549 | 1,577 | 🟢 Go 1.62x faster |
| Multicall_TokenMetadata | 232,237 | 447,399 | 4,306 | 2,235 | 🟢 Go 1.93x faster |
| Multicall_50Calls | 785,560 | 1,504,212 | 1,273 | 665 | 🟢 Go 1.91x faster |
| Multicall_100Calls | 1,356,091 | 2,637,966 | 737 | 379 | 🟢 Go 1.95x faster |
| Multicall_200Calls | 2,483,017 | 5,048,720 | 403 | 198 | 🟢 Go 2.03x faster |
| Multicall_500Calls | 3,859,493 | 10,886,048 | 259 | 92 | 🟢 Go 2.82x faster |
| Multicall_MixedContracts_100 | 1,390,738 | 2,597,875 | 719 | 385 | 🟢 Go 1.87x faster |
| Multicall_1000Calls | 5,939,041 | 20,884,282 | 168 | 48 | 🟢 Go 3.52x faster |
| Multicall_10000Calls_SingleRPC | 126,200,298 | 213,552,011 | 8 | 5 | 🟢 Go 1.69x faster |
| Multicall_10000Calls_Chunked | 29,002,749 | 215,731,113 | 34 | 5 | 🟢 Go 7.44x faster |
| Multicall_10000Calls_AggressiveChunking | 37,552,933 | 211,112,988 | 27 | 5 | 🟢 Go 5.62x faster |

## Win Summary

- 🟢 Go wins: 16 (100%)
- 🔵 TS wins: 0 (0%)
- ⚪ Ties: 0 (0%)

## Notes

- Benchmarks run against the same Anvil instance (mainnet fork) for fair comparison
- ns/op = nanoseconds per operation (lower is better)
- ops/s = operations per second (higher is better)
- 🟢 = Go faster, 🔵 = TS faster, ⚪ = Similar (within 5%)
