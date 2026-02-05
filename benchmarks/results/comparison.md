# Benchmark Comparison: viem-go vs viem TypeScript

Generated: 2026-02-05T17:27:24.998Z

## Overall Summary

**🏆 Go is 1.89x faster overall**

| Metric | Go | TypeScript |
|--------|----|-----------|
| Avg ns/op | 239,789 | 452,385 |
| Avg ops/s | 4,170 | 2,211 |
| Wins | 68/70 | 1/70 |

## By Suite

| Suite | Benchmarks | Go Wins | TS Wins | Ties | Winner |
|-------|------------|---------|---------|------|--------|
| call | 30 | 29 | 0 | 1 | 🟢 Go 1.43x faster |
| multicall | 40 | 39 | 1 | 0 | 🟢 Go 2.07x faster |

## Detailed Results

| Benchmark | Go (ns/op) | TS (ns/op) | Go (ops/s) | TS (ops/s) | Result |
|-----------|------------|------------|------------|------------|--------|
| Call_Basic | 155,413 | 227,160 | 6,434 | 4,402 | 🟢 Go 1.46x faster |
| Call_Basic | 165,052 | 227,160 | 6,059 | 4,402 | 🟢 Go 1.38x faster |
| Call_Basic | 167,995 | 227,160 | 5,953 | 4,402 | 🟢 Go 1.35x faster |
| Call_Basic | 155,179 | 227,160 | 6,444 | 4,402 | 🟢 Go 1.46x faster |
| Call_Basic | 164,573 | 227,160 | 6,076 | 4,402 | 🟢 Go 1.38x faster |
| Call_WithData | 155,350 | 218,990 | 6,437 | 4,566 | 🟢 Go 1.41x faster |
| Call_WithData | 175,161 | 218,990 | 5,709 | 4,566 | 🟢 Go 1.25x faster |
| Call_WithData | 152,108 | 218,990 | 6,574 | 4,566 | 🟢 Go 1.44x faster |
| Call_WithData | 171,290 | 218,990 | 5,838 | 4,566 | 🟢 Go 1.28x faster |
| Call_WithData | 148,793 | 218,990 | 6,721 | 4,566 | 🟢 Go 1.47x faster |
| Call_WithAccount | 153,954 | 219,699 | 6,495 | 4,552 | 🟢 Go 1.43x faster |
| Call_WithAccount | 155,034 | 219,699 | 6,450 | 4,552 | 🟢 Go 1.42x faster |
| Call_WithAccount | 173,141 | 219,699 | 5,776 | 4,552 | 🟢 Go 1.27x faster |
| Call_WithAccount | 155,549 | 219,699 | 6,429 | 4,552 | 🟢 Go 1.41x faster |
| Call_WithAccount | 158,752 | 219,699 | 6,299 | 4,552 | 🟢 Go 1.38x faster |
| Call_Decimals | 177,428 | 216,204 | 5,636 | 4,625 | 🟢 Go 1.22x faster |
| Call_Decimals | 183,648 | 216,204 | 5,445 | 4,625 | 🟢 Go 1.18x faster |
| Call_Decimals | 156,643 | 216,204 | 6,384 | 4,625 | 🟢 Go 1.38x faster |
| Call_Decimals | 181,777 | 216,204 | 5,501 | 4,625 | 🟢 Go 1.19x faster |
| Call_Decimals | 141,277 | 216,204 | 7,078 | 4,625 | 🟢 Go 1.53x faster |
| Call_Symbol | 144,042 | 216,639 | 6,942 | 4,616 | 🟢 Go 1.50x faster |
| Call_Symbol | 158,880 | 216,639 | 6,294 | 4,616 | 🟢 Go 1.36x faster |
| Call_Symbol | 177,846 | 216,639 | 5,623 | 4,616 | 🟢 Go 1.22x faster |
| Call_Symbol | 147,443 | 216,639 | 6,782 | 4,616 | 🟢 Go 1.47x faster |
| Call_Symbol | 226,148 | 216,639 | 4,422 | 4,616 | ⚪ Similar |
| Call_BalanceOfMultiple | 145,611 | 276,218 | 6,868 | 3,620 | 🟢 Go 1.90x faster |
| Call_BalanceOfMultiple | 139,128 | 276,218 | 7,188 | 3,620 | 🟢 Go 1.99x faster |
| Call_BalanceOfMultiple | 138,889 | 276,218 | 7,200 | 3,620 | 🟢 Go 1.99x faster |
| Call_BalanceOfMultiple | 143,384 | 276,218 | 6,974 | 3,620 | 🟢 Go 1.93x faster |
| Call_BalanceOfMultiple | 138,884 | 276,218 | 7,200 | 3,620 | 🟢 Go 1.99x faster |
| Multicall_Basic | 79,525 | 465,095 | 12,575 | 2,150 | 🟢 Go 5.85x faster |
| Multicall_Basic | 79,463 | 465,095 | 12,584 | 2,150 | 🟢 Go 5.85x faster |
| Multicall_Basic | 80,539 | 465,095 | 12,416 | 2,150 | 🟢 Go 5.77x faster |
| Multicall_Basic | 90,260 | 465,095 | 11,079 | 2,150 | 🟢 Go 5.15x faster |
| Multicall_Basic | 82,162 | 465,095 | 12,171 | 2,150 | 🟢 Go 5.66x faster |
| Multicall_WithArgs | 81,471 | 382,162 | 12,274 | 2,617 | 🟢 Go 4.69x faster |
| Multicall_WithArgs | 81,197 | 382,162 | 12,316 | 2,617 | 🟢 Go 4.71x faster |
| Multicall_WithArgs | 82,106 | 382,162 | 12,179 | 2,617 | 🟢 Go 4.65x faster |
| Multicall_WithArgs | 93,663 | 382,162 | 10,677 | 2,617 | 🟢 Go 4.08x faster |
| Multicall_WithArgs | 82,049 | 382,162 | 12,188 | 2,617 | 🟢 Go 4.66x faster |
| Multicall_MultiContract | 105,056 | 536,414 | 9,519 | 1,864 | 🟢 Go 5.11x faster |
| Multicall_MultiContract | 113,872 | 536,414 | 8,782 | 1,864 | 🟢 Go 4.71x faster |
| Multicall_MultiContract | 123,861 | 536,414 | 8,074 | 1,864 | 🟢 Go 4.33x faster |
| Multicall_MultiContract | 104,902 | 536,414 | 9,533 | 1,864 | 🟢 Go 5.11x faster |
| Multicall_MultiContract | 103,751 | 536,414 | 9,638 | 1,864 | 🟢 Go 5.17x faster |
| Multicall_10Calls | 243,439 | 529,568 | 4,108 | 1,888 | 🟢 Go 2.18x faster |
| Multicall_10Calls | 245,016 | 529,568 | 4,081 | 1,888 | 🟢 Go 2.16x faster |
| Multicall_10Calls | 243,733 | 529,568 | 4,103 | 1,888 | 🟢 Go 2.17x faster |
| Multicall_10Calls | 243,627 | 529,568 | 4,105 | 1,888 | 🟢 Go 2.17x faster |
| Multicall_10Calls | 272,166 | 529,568 | 3,674 | 1,888 | 🟢 Go 1.95x faster |
| Multicall_30Calls | 769,597 | 1,068,273 | 1,299 | 936 | 🟢 Go 1.39x faster |
| Multicall_30Calls | 769,688 | 1,068,273 | 1,299 | 936 | 🟢 Go 1.39x faster |
| Multicall_30Calls | 833,296 | 1,068,273 | 1,200 | 936 | 🟢 Go 1.28x faster |
| Multicall_30Calls | 828,576 | 1,068,273 | 1,207 | 936 | 🟢 Go 1.29x faster |
| Multicall_30Calls | 2,495,194 | 1,068,273 | 401 | 936 | 🔵 TS 2.34x faster |
| Multicall_ChunkedParallel | 563,068 | 979,058 | 1,776 | 1,021 | 🟢 Go 1.74x faster |
| Multicall_ChunkedParallel | 520,211 | 979,058 | 1,922 | 1,021 | 🟢 Go 1.88x faster |
| Multicall_ChunkedParallel | 552,542 | 979,058 | 1,810 | 1,021 | 🟢 Go 1.77x faster |
| Multicall_ChunkedParallel | 526,846 | 979,058 | 1,898 | 1,021 | 🟢 Go 1.86x faster |
| Multicall_ChunkedParallel | 529,303 | 979,058 | 1,889 | 1,021 | 🟢 Go 1.85x faster |
| Multicall_Deployless | 84,016 | 572,433 | 11,902 | 1,747 | 🟢 Go 6.81x faster |
| Multicall_Deployless | 83,462 | 572,433 | 11,982 | 1,747 | 🟢 Go 6.86x faster |
| Multicall_Deployless | 83,626 | 572,433 | 11,958 | 1,747 | 🟢 Go 6.85x faster |
| Multicall_Deployless | 87,646 | 572,433 | 11,410 | 1,747 | 🟢 Go 6.53x faster |
| Multicall_Deployless | 82,910 | 572,433 | 12,061 | 1,747 | 🟢 Go 6.90x faster |
| Multicall_TokenMetadata | 106,116 | 425,474 | 9,424 | 2,350 | 🟢 Go 4.01x faster |
| Multicall_TokenMetadata | 106,266 | 425,474 | 9,410 | 2,350 | 🟢 Go 4.00x faster |
| Multicall_TokenMetadata | 107,969 | 425,474 | 9,262 | 2,350 | 🟢 Go 3.94x faster |
| Multicall_TokenMetadata | 106,300 | 425,474 | 9,407 | 2,350 | 🟢 Go 4.00x faster |
| Multicall_TokenMetadata | 108,342 | 425,474 | 9,230 | 2,350 | 🟢 Go 3.93x faster |

## Win Summary

- 🟢 Go wins: 68 (97%)
- 🔵 TS wins: 1 (1%)
- ⚪ Ties: 1 (1%)

## Notes

- Benchmarks run against the same Anvil instance (mainnet fork) for fair comparison
- ns/op = nanoseconds per operation (lower is better)
- ops/s = operations per second (higher is better)
- 🟢 = Go faster, 🔵 = TS faster, ⚪ = Similar (within 5%)
