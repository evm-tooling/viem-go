# Benchmark Comparison: viem-go vs viem TypeScript

Generated: 2026-02-05T17:26:39.109Z

## Overall Summary

**🏆 Go is 3.25x faster overall**

| Metric | Go | TypeScript |
|--------|----|-----------|
| Avg ns/op | 17,874,996 | 58,180,284 |
| Avg ops/s | 56 | 17 |
| Wins | 16/16 | 0/16 |

## Detailed Results

| Benchmark | Go (ns/op) | TS (ns/op) | Go (ops/s) | TS (ops/s) | Result |
|-----------|------------|------------|------------|------------|--------|
| Multicall_Basic | 241,418 | 455,660 | 4,142 | 2,195 | 🟢 Go 1.89x faster |
| Multicall_WithArgs | 221,758 | 387,133 | 4,509 | 2,583 | 🟢 Go 1.75x faster |
| Multicall_MultiContract | 253,772 | 462,263 | 3,941 | 2,163 | 🟢 Go 1.82x faster |
| Multicall_10Calls | 284,814 | 555,713 | 3,511 | 1,799 | 🟢 Go 1.95x faster |
| Multicall_30Calls | 547,743 | 1,012,361 | 1,826 | 988 | 🟢 Go 1.85x faster |
| Multicall_Deployless | 394,342 | 640,623 | 2,536 | 1,561 | 🟢 Go 1.62x faster |
| Multicall_TokenMetadata | 249,351 | 429,251 | 4,010 | 2,330 | 🟢 Go 1.72x faster |
| Multicall_50Calls | 823,730 | 1,448,478 | 1,214 | 690 | 🟢 Go 1.76x faster |
| Multicall_100Calls | 1,708,452 | 2,693,966 | 585 | 371 | 🟢 Go 1.58x faster |
| Multicall_200Calls | 3,403,431 | 4,908,216 | 294 | 204 | 🟢 Go 1.44x faster |
| Multicall_500Calls | 8,354,617 | 12,045,682 | 120 | 83 | 🟢 Go 1.44x faster |
| Multicall_MixedContracts_100 | 1,732,611 | 2,568,779 | 577 | 389 | 🟢 Go 1.48x faster |
| Multicall_1000Calls | 16,887,841 | 24,288,646 | 59 | 41 | 🟢 Go 1.44x faster |
| Multicall_10000Calls_SingleRPC | 165,489,012 | 322,997,416 | 6 | 3 | 🟢 Go 1.95x faster |
| Multicall_10000Calls_Chunked | 36,930,845 | 214,518,620 | 27 | 5 | 🟢 Go 5.81x faster |
| Multicall_10000Calls_AggressiveChunking | 48,476,204 | 341,471,743 | 21 | 3 | 🟢 Go 7.04x faster |

## Win Summary

- 🟢 Go wins: 16 (100%)
- 🔵 TS wins: 0 (0%)
- ⚪ Ties: 0 (0%)

## Notes

- Benchmarks run against the same Anvil instance (mainnet fork) for fair comparison
- ns/op = nanoseconds per operation (lower is better)
- ops/s = operations per second (higher is better)
- 🟢 = Go faster, 🔵 = TS faster, ⚪ = Similar (within 5%)
