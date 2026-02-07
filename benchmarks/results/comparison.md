# Benchmark Comparison: viem-go vs viem TypeScript

Generated: 2026-02-07T00:36:17.123Z

## Overall Summary

**🏆 Go is 3.60x faster overall**

| Metric | Go | TypeScript |
|--------|----|-----------|
| Avg ns/op | 9,833,455 | 35,400,609 |
| Avg ops/s | 102 | 28 |
| Wins | 22/22 | 0/22 |

## By Suite

| Suite | Benchmarks | Go Wins | TS Wins | Ties | Winner |
|-------|------------|---------|---------|------|--------|
| call | 6 | 6 | 0 | 0 | 🟢 Go 80.89x faster |
| multicall | 16 | 16 | 0 | 0 | 🟢 Go 3.19x faster |

## Detailed Results

| Benchmark | Go (ns/op) | TS (ns/op) | Go (ops/s) | TS (ops/s) | Result |
|-----------|------------|------------|------------|------------|--------|
| Call_Basic | 184,930 | 18,309,921 | 5,407 | 55 | 🟢 Go 99.01x faster |
| Call_WithData | 179,602 | 18,611,646 | 5,568 | 54 | 🟢 Go 103.63x faster |
| Call_WithAccount | 191,890 | 259,973 | 5,211 | 3,847 | 🟢 Go 1.35x faster |
| Call_Decimals | 194,174 | 18,327,371 | 5,150 | 55 | 🟢 Go 94.39x faster |
| Call_Symbol | 202,911 | 18,409,527 | 4,928 | 54 | 🟢 Go 90.73x faster |
| Call_BalanceOfMultiple | 184,987 | 18,169,001 | 5,406 | 55 | 🟢 Go 98.22x faster |
| Multicall_Basic | 218,827 | 495,285 | 4,570 | 2,019 | 🟢 Go 2.26x faster |
| Multicall_WithArgs | 222,541 | 1,243,797 | 4,494 | 804 | 🟢 Go 5.59x faster |
| Multicall_MultiContract | 264,798 | 553,018 | 3,776 | 1,808 | 🟢 Go 2.09x faster |
| Multicall_10Calls | 303,373 | 541,688 | 3,296 | 1,846 | 🟢 Go 1.79x faster |
| Multicall_30Calls | 549,894 | 1,038,508 | 1,819 | 963 | 🟢 Go 1.89x faster |
| Multicall_Deployless | 391,033 | 617,936 | 2,557 | 1,618 | 🟢 Go 1.58x faster |
| Multicall_TokenMetadata | 238,988 | 425,579 | 4,184 | 2,350 | 🟢 Go 1.78x faster |
| Multicall_50Calls | 780,430 | 1,428,796 | 1,281 | 700 | 🟢 Go 1.83x faster |
| Multicall_100Calls | 1,373,991 | 2,528,573 | 728 | 395 | 🟢 Go 1.84x faster |
| Multicall_200Calls | 2,585,980 | 5,202,914 | 387 | 192 | 🟢 Go 2.01x faster |
| Multicall_500Calls | 6,271,189 | 10,647,768 | 159 | 94 | 🟢 Go 1.70x faster |
| Multicall_MixedContracts_100 | 1,383,611 | 2,557,610 | 723 | 391 | 🟢 Go 1.85x faster |
| Multicall_1000Calls | 6,041,294 | 20,660,047 | 166 | 48 | 🟢 Go 3.42x faster |
| Multicall_10000Calls_SingleRPC | 125,677,167 | 216,651,862 | 8 | 5 | 🟢 Go 1.72x faster |
| Multicall_10000Calls_Chunked | 32,944,897 | 210,930,414 | 30 | 5 | 🟢 Go 6.40x faster |
| Multicall_10000Calls_AggressiveChunking | 35,949,493 | 211,202,163 | 28 | 5 | 🟢 Go 5.87x faster |

## Win Summary

- 🟢 Go wins: 22 (100%)
- 🔵 TS wins: 0 (0%)
- ⚪ Ties: 0 (0%)

## Notes

- Benchmarks run against the same Anvil instance (mainnet fork) for fair comparison
- ns/op = nanoseconds per operation (lower is better)
- ops/s = operations per second (higher is better)
- 🟢 = Go faster, 🔵 = TS faster, ⚪ = Similar (within 5%)
