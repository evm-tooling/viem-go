# Full Benchmark Report: viem-go vs viem TypeScript

Generated: 2026-02-05T14:43:18.555Z

---

## Executive Summary

This report compares **70** benchmarks across **2** test suites.

### 🏆 Winner: Go (viem-go)

Go is **1.89x faster** on average across all benchmarks.

### Quick Stats

| Metric | Value |
|--------|-------|
| Total Benchmarks | 70 |
| Test Suites | 2 |
| Go Wins | 68 (97.1%) |
| TypeScript Wins | 1 (1.4%) |
| Ties | 1 (1.4%) |
| Avg Go Latency | 239.79 µs |
| Avg TS Latency | 452.38 µs |
| Go Throughput | 4,170 ops/s |
| TS Throughput | 2,211 ops/s |

---

## Suite-by-Suite Analysis

### Call Suite

**Result:** 🟢 Go 1.43x faster

| Benchmark | Go | TS | Diff | Winner |
|-----------|----|----|------|--------|
| Basic | 155.41 µs | 227.16 µs | 1.46x | 🟢 |
| Basic | 165.05 µs | 227.16 µs | 1.38x | 🟢 |
| Basic | 168.00 µs | 227.16 µs | 1.35x | 🟢 |
| Basic | 155.18 µs | 227.16 µs | 1.46x | 🟢 |
| Basic | 164.57 µs | 227.16 µs | 1.38x | 🟢 |
| WithData | 155.35 µs | 218.99 µs | 1.41x | 🟢 |
| WithData | 175.16 µs | 218.99 µs | 1.25x | 🟢 |
| WithData | 152.11 µs | 218.99 µs | 1.44x | 🟢 |
| WithData | 171.29 µs | 218.99 µs | 1.28x | 🟢 |
| WithData | 148.79 µs | 218.99 µs | 1.47x | 🟢 |
| WithAccount | 153.95 µs | 219.70 µs | 1.43x | 🟢 |
| WithAccount | 155.03 µs | 219.70 µs | 1.42x | 🟢 |
| WithAccount | 173.14 µs | 219.70 µs | 1.27x | 🟢 |
| WithAccount | 155.55 µs | 219.70 µs | 1.41x | 🟢 |
| WithAccount | 158.75 µs | 219.70 µs | 1.38x | 🟢 |
| Decimals | 177.43 µs | 216.20 µs | 1.22x | 🟢 |
| Decimals | 183.65 µs | 216.20 µs | 1.18x | 🟢 |
| Decimals | 156.64 µs | 216.20 µs | 1.38x | 🟢 |
| Decimals | 181.78 µs | 216.20 µs | 1.19x | 🟢 |
| Decimals | 141.28 µs | 216.20 µs | 1.53x | 🟢 |
| Symbol | 144.04 µs | 216.64 µs | 1.50x | 🟢 |
| Symbol | 158.88 µs | 216.64 µs | 1.36x | 🟢 |
| Symbol | 177.85 µs | 216.64 µs | 1.22x | 🟢 |
| Symbol | 147.44 µs | 216.64 µs | 1.47x | 🟢 |
| Symbol | 226.15 µs | 216.64 µs | 1.04x | ⚪ |
| BalanceOfMultiple | 145.61 µs | 276.22 µs | 1.90x | 🟢 |
| BalanceOfMultiple | 139.13 µs | 276.22 µs | 1.99x | 🟢 |
| BalanceOfMultiple | 138.89 µs | 276.22 µs | 1.99x | 🟢 |
| BalanceOfMultiple | 143.38 µs | 276.22 µs | 1.93x | 🟢 |
| BalanceOfMultiple | 138.88 µs | 276.22 µs | 1.99x | 🟢 |

**Suite Statistics:**
- Benchmarks: 30
- Go wins: 29, TS wins: 0, Ties: 1
- Avg Go: 160.28 µs | Avg TS: 229.15 µs

### Multicall Suite

**Result:** 🟢 Go 2.07x faster

| Benchmark | Go | TS | Diff | Winner |
|-----------|----|----|------|--------|
| Basic | 79.53 µs | 465.09 µs | 5.85x | 🟢 |
| Basic | 79.46 µs | 465.09 µs | 5.85x | 🟢 |
| Basic | 80.54 µs | 465.09 µs | 5.77x | 🟢 |
| Basic | 90.26 µs | 465.09 µs | 5.15x | 🟢 |
| Basic | 82.16 µs | 465.09 µs | 5.66x | 🟢 |
| WithArgs | 81.47 µs | 382.16 µs | 4.69x | 🟢 |
| WithArgs | 81.20 µs | 382.16 µs | 4.71x | 🟢 |
| WithArgs | 82.11 µs | 382.16 µs | 4.65x | 🟢 |
| WithArgs | 93.66 µs | 382.16 µs | 4.08x | 🟢 |
| WithArgs | 82.05 µs | 382.16 µs | 4.66x | 🟢 |
| MultiContract | 105.06 µs | 536.41 µs | 5.11x | 🟢 |
| MultiContract | 113.87 µs | 536.41 µs | 4.71x | 🟢 |
| MultiContract | 123.86 µs | 536.41 µs | 4.33x | 🟢 |
| MultiContract | 104.90 µs | 536.41 µs | 5.11x | 🟢 |
| MultiContract | 103.75 µs | 536.41 µs | 5.17x | 🟢 |
| 10Calls | 243.44 µs | 529.57 µs | 2.18x | 🟢 |
| 10Calls | 245.02 µs | 529.57 µs | 2.16x | 🟢 |
| 10Calls | 243.73 µs | 529.57 µs | 2.17x | 🟢 |
| 10Calls | 243.63 µs | 529.57 µs | 2.17x | 🟢 |
| 10Calls | 272.17 µs | 529.57 µs | 1.95x | 🟢 |
| 30Calls | 769.60 µs | 1.07 ms | 1.39x | 🟢 |
| 30Calls | 769.69 µs | 1.07 ms | 1.39x | 🟢 |
| 30Calls | 833.30 µs | 1.07 ms | 1.28x | 🟢 |
| 30Calls | 828.58 µs | 1.07 ms | 1.29x | 🟢 |
| 30Calls | 2.50 ms | 1.07 ms | 2.34x | 🔵 |
| ChunkedParallel | 563.07 µs | 979.06 µs | 1.74x | 🟢 |
| ChunkedParallel | 520.21 µs | 979.06 µs | 1.88x | 🟢 |
| ChunkedParallel | 552.54 µs | 979.06 µs | 1.77x | 🟢 |
| ChunkedParallel | 526.85 µs | 979.06 µs | 1.86x | 🟢 |
| ChunkedParallel | 529.30 µs | 979.06 µs | 1.85x | 🟢 |
| Deployless | 84.02 µs | 572.43 µs | 6.81x | 🟢 |
| Deployless | 83.46 µs | 572.43 µs | 6.86x | 🟢 |
| Deployless | 83.63 µs | 572.43 µs | 6.85x | 🟢 |
| Deployless | 87.65 µs | 572.43 µs | 6.53x | 🟢 |
| Deployless | 82.91 µs | 572.43 µs | 6.90x | 🟢 |
| TokenMetadata | 106.12 µs | 425.47 µs | 4.01x | 🟢 |
| TokenMetadata | 106.27 µs | 425.47 µs | 4.00x | 🟢 |
| TokenMetadata | 107.97 µs | 425.47 µs | 3.94x | 🟢 |
| TokenMetadata | 106.30 µs | 425.47 µs | 4.00x | 🟢 |
| TokenMetadata | 108.34 µs | 425.47 µs | 3.93x | 🟢 |

**Suite Statistics:**
- Benchmarks: 40
- Go wins: 39, TS wins: 1, Ties: 0
- Avg Go: 299.42 µs | Avg TS: 619.81 µs

---

## Category Analysis

### Basic Operations

🟢 **Go 2.84x faster**

Benchmarks: 10 | Go wins: 10 | TS wins: 0 | Ties: 0

### With Parameters

🟢 **Go 2.92x faster**

Benchmarks: 15 | Go wins: 15 | TS wins: 0 | Ties: 0

### With Account

🟢 **Go 1.38x faster**

Benchmarks: 5 | Go wins: 5 | TS wins: 0 | Ties: 0

### Simple Reads

🟢 **Go 1.28x faster**

Benchmarks: 10 | Go wins: 9 | TS wins: 0 | Ties: 1

### Batch Operations

🟢 **Go 1.22x faster**

Benchmarks: 15 | Go wins: 14 | TS wins: 1 | Ties: 0

### Multi-Contract

🟢 **Go 4.86x faster**

Benchmarks: 5 | Go wins: 5 | TS wins: 0 | Ties: 0

### Parallel Execution

🟢 **Go 1.82x faster**

Benchmarks: 5 | Go wins: 5 | TS wins: 0 | Ties: 0

### Deployless

🟢 **Go 6.79x faster**

Benchmarks: 5 | Go wins: 5 | TS wins: 0 | Ties: 0

---

## Memory Analysis (Go)

| Benchmark | Bytes/op | Allocs/op |
|-----------|----------|----------|
| Call_Basic | 8,795 | 106 |
| Call_Basic | 8,808 | 106 |
| Call_Basic | 8,805 | 106 |
| Call_Basic | 8,810 | 106 |
| Call_Basic | 8,815 | 106 |
| Call_WithData | 8,697 | 107 |
| Call_WithData | 8,675 | 107 |
| Call_WithData | 8,679 | 107 |
| Call_WithData | 8,676 | 107 |
| Call_WithData | 8,681 | 107 |
| Call_WithAccount | 9,137 | 112 |
| Call_WithAccount | 9,142 | 112 |
| Call_WithAccount | 9,110 | 112 |
| Call_WithAccount | 9,128 | 112 |
| Call_WithAccount | 9,139 | 112 |
| Call_Decimals | 8,453 | 106 |
| Call_Decimals | 8,479 | 106 |
| Call_Decimals | 8,481 | 106 |
| Call_Decimals | 8,477 | 106 |
| Call_Decimals | 8,498 | 106 |
| Call_Symbol | 8,802 | 106 |
| Call_Symbol | 8,830 | 106 |
| Call_Symbol | 8,819 | 106 |
| Call_Symbol | 8,807 | 106 |
| Call_Symbol | 8,767 | 106 |
| Call_BalanceOfMultiple | 8,688 | 107 |
| Call_BalanceOfMultiple | 8,713 | 107 |
| Call_BalanceOfMultiple | 8,694 | 107 |
| Call_BalanceOfMultiple | 8,705 | 107 |
| Call_BalanceOfMultiple | 8,704 | 107 |
| Multicall_Basic | 55,351 | 641 |
| Multicall_Basic | 55,347 | 641 |
| Multicall_Basic | 55,435 | 641 |
| Multicall_Basic | 55,464 | 641 |
| Multicall_Basic | 55,514 | 641 |
| Multicall_WithArgs | 56,092 | 656 |
| Multicall_WithArgs | 56,059 | 656 |
| Multicall_WithArgs | 56,049 | 656 |
| Multicall_WithArgs | 56,033 | 656 |
| Multicall_WithArgs | 56,060 | 656 |
| Multicall_MultiContract | 72,303 | 837 |
| Multicall_MultiContract | 72,372 | 837 |
| Multicall_MultiContract | 72,340 | 837 |
| Multicall_MultiContract | 72,354 | 837 |
| Multicall_MultiContract | 72,262 | 837 |
| Multicall_10Calls | 173,581 | 1,997 |
| Multicall_10Calls | 173,583 | 1,997 |
| Multicall_10Calls | 173,522 | 1,997 |
| Multicall_10Calls | 173,518 | 1,997 |
| Multicall_10Calls | 173,705 | 1,997 |
| Multicall_30Calls | 536,312 | 5,902 |
| Multicall_30Calls | 537,491 | 5,903 |
| Multicall_30Calls | 538,192 | 5,903 |
| Multicall_30Calls | 537,003 | 5,902 |
| Multicall_30Calls | 534,680 | 5,902 |
| Multicall_ChunkedParallel | 370,739 | 4,066 |
| Multicall_ChunkedParallel | 366,875 | 4,064 |
| Multicall_ChunkedParallel | 367,687 | 4,065 |
| Multicall_ChunkedParallel | 367,225 | 4,064 |
| Multicall_ChunkedParallel | 368,717 | 4,065 |
| Multicall_Deployless | 56,217 | 641 |
| Multicall_Deployless | 56,101 | 641 |
| Multicall_Deployless | 56,011 | 641 |
| Multicall_Deployless | 55,988 | 641 |
| Multicall_Deployless | 55,975 | 641 |
| Multicall_TokenMetadata | 72,494 | 827 |
| Multicall_TokenMetadata | 72,453 | 827 |
| Multicall_TokenMetadata | 72,423 | 827 |
| Multicall_TokenMetadata | 72,329 | 827 |
| Multicall_TokenMetadata | 72,327 | 827 |

---

## Detailed Raw Data

| Benchmark | Suite | Go ns/op | TS ns/op | Go ops/s | TS ops/s | Ratio | Winner |
|-----------|-------|----------|----------|----------|----------|-------|--------|
| Call_Basic | call | 155,413 | 227,160 | 6,434 | 4,402 | 0.684 | 🟢 |
| Call_Basic | call | 165,052 | 227,160 | 6,059 | 4,402 | 0.727 | 🟢 |
| Call_Basic | call | 167,995 | 227,160 | 5,953 | 4,402 | 0.740 | 🟢 |
| Call_Basic | call | 155,179 | 227,160 | 6,444 | 4,402 | 0.683 | 🟢 |
| Call_Basic | call | 164,573 | 227,160 | 6,076 | 4,402 | 0.724 | 🟢 |
| Call_WithData | call | 155,350 | 218,990 | 6,437 | 4,566 | 0.709 | 🟢 |
| Call_WithData | call | 175,161 | 218,990 | 5,709 | 4,566 | 0.800 | 🟢 |
| Call_WithData | call | 152,108 | 218,990 | 6,574 | 4,566 | 0.695 | 🟢 |
| Call_WithData | call | 171,290 | 218,990 | 5,838 | 4,566 | 0.782 | 🟢 |
| Call_WithData | call | 148,793 | 218,990 | 6,721 | 4,566 | 0.679 | 🟢 |
| Call_WithAccount | call | 153,954 | 219,699 | 6,495 | 4,552 | 0.701 | 🟢 |
| Call_WithAccount | call | 155,034 | 219,699 | 6,450 | 4,552 | 0.706 | 🟢 |
| Call_WithAccount | call | 173,141 | 219,699 | 5,776 | 4,552 | 0.788 | 🟢 |
| Call_WithAccount | call | 155,549 | 219,699 | 6,429 | 4,552 | 0.708 | 🟢 |
| Call_WithAccount | call | 158,752 | 219,699 | 6,299 | 4,552 | 0.723 | 🟢 |
| Call_Decimals | call | 177,428 | 216,204 | 5,636 | 4,625 | 0.821 | 🟢 |
| Call_Decimals | call | 183,648 | 216,204 | 5,445 | 4,625 | 0.849 | 🟢 |
| Call_Decimals | call | 156,643 | 216,204 | 6,384 | 4,625 | 0.725 | 🟢 |
| Call_Decimals | call | 181,777 | 216,204 | 5,501 | 4,625 | 0.841 | 🟢 |
| Call_Decimals | call | 141,277 | 216,204 | 7,078 | 4,625 | 0.653 | 🟢 |
| Call_Symbol | call | 144,042 | 216,639 | 6,942 | 4,616 | 0.665 | 🟢 |
| Call_Symbol | call | 158,880 | 216,639 | 6,294 | 4,616 | 0.733 | 🟢 |
| Call_Symbol | call | 177,846 | 216,639 | 5,623 | 4,616 | 0.821 | 🟢 |
| Call_Symbol | call | 147,443 | 216,639 | 6,782 | 4,616 | 0.681 | 🟢 |
| Call_Symbol | call | 226,148 | 216,639 | 4,422 | 4,616 | 1.044 | ⚪ |
| Call_BalanceOfMultiple | call | 145,611 | 276,218 | 6,868 | 3,620 | 0.527 | 🟢 |
| Call_BalanceOfMultiple | call | 139,128 | 276,218 | 7,188 | 3,620 | 0.504 | 🟢 |
| Call_BalanceOfMultiple | call | 138,889 | 276,218 | 7,200 | 3,620 | 0.503 | 🟢 |
| Call_BalanceOfMultiple | call | 143,384 | 276,218 | 6,974 | 3,620 | 0.519 | 🟢 |
| Call_BalanceOfMultiple | call | 138,884 | 276,218 | 7,200 | 3,620 | 0.503 | 🟢 |
| Multicall_Basic | multicall | 79,525 | 465,095 | 12,575 | 2,150 | 0.171 | 🟢 |
| Multicall_Basic | multicall | 79,463 | 465,095 | 12,584 | 2,150 | 0.171 | 🟢 |
| Multicall_Basic | multicall | 80,539 | 465,095 | 12,416 | 2,150 | 0.173 | 🟢 |
| Multicall_Basic | multicall | 90,260 | 465,095 | 11,079 | 2,150 | 0.194 | 🟢 |
| Multicall_Basic | multicall | 82,162 | 465,095 | 12,171 | 2,150 | 0.177 | 🟢 |
| Multicall_WithArgs | multicall | 81,471 | 382,162 | 12,274 | 2,617 | 0.213 | 🟢 |
| Multicall_WithArgs | multicall | 81,197 | 382,162 | 12,316 | 2,617 | 0.212 | 🟢 |
| Multicall_WithArgs | multicall | 82,106 | 382,162 | 12,179 | 2,617 | 0.215 | 🟢 |
| Multicall_WithArgs | multicall | 93,663 | 382,162 | 10,677 | 2,617 | 0.245 | 🟢 |
| Multicall_WithArgs | multicall | 82,049 | 382,162 | 12,188 | 2,617 | 0.215 | 🟢 |
| Multicall_MultiContract | multicall | 105,056 | 536,414 | 9,519 | 1,864 | 0.196 | 🟢 |
| Multicall_MultiContract | multicall | 113,872 | 536,414 | 8,782 | 1,864 | 0.212 | 🟢 |
| Multicall_MultiContract | multicall | 123,861 | 536,414 | 8,074 | 1,864 | 0.231 | 🟢 |
| Multicall_MultiContract | multicall | 104,902 | 536,414 | 9,533 | 1,864 | 0.196 | 🟢 |
| Multicall_MultiContract | multicall | 103,751 | 536,414 | 9,638 | 1,864 | 0.193 | 🟢 |
| Multicall_10Calls | multicall | 243,439 | 529,568 | 4,108 | 1,888 | 0.460 | 🟢 |
| Multicall_10Calls | multicall | 245,016 | 529,568 | 4,081 | 1,888 | 0.463 | 🟢 |
| Multicall_10Calls | multicall | 243,733 | 529,568 | 4,103 | 1,888 | 0.460 | 🟢 |
| Multicall_10Calls | multicall | 243,627 | 529,568 | 4,105 | 1,888 | 0.460 | 🟢 |
| Multicall_10Calls | multicall | 272,166 | 529,568 | 3,674 | 1,888 | 0.514 | 🟢 |
| Multicall_30Calls | multicall | 769,597 | 1,068,273 | 1,299 | 936 | 0.720 | 🟢 |
| Multicall_30Calls | multicall | 769,688 | 1,068,273 | 1,299 | 936 | 0.720 | 🟢 |
| Multicall_30Calls | multicall | 833,296 | 1,068,273 | 1,200 | 936 | 0.780 | 🟢 |
| Multicall_30Calls | multicall | 828,576 | 1,068,273 | 1,207 | 936 | 0.776 | 🟢 |
| Multicall_30Calls | multicall | 2,495,194 | 1,068,273 | 401 | 936 | 2.336 | 🔵 |
| Multicall_ChunkedParallel | multicall | 563,068 | 979,058 | 1,776 | 1,021 | 0.575 | 🟢 |
| Multicall_ChunkedParallel | multicall | 520,211 | 979,058 | 1,922 | 1,021 | 0.531 | 🟢 |
| Multicall_ChunkedParallel | multicall | 552,542 | 979,058 | 1,810 | 1,021 | 0.564 | 🟢 |
| Multicall_ChunkedParallel | multicall | 526,846 | 979,058 | 1,898 | 1,021 | 0.538 | 🟢 |
| Multicall_ChunkedParallel | multicall | 529,303 | 979,058 | 1,889 | 1,021 | 0.541 | 🟢 |
| Multicall_Deployless | multicall | 84,016 | 572,433 | 11,902 | 1,747 | 0.147 | 🟢 |
| Multicall_Deployless | multicall | 83,462 | 572,433 | 11,982 | 1,747 | 0.146 | 🟢 |
| Multicall_Deployless | multicall | 83,626 | 572,433 | 11,958 | 1,747 | 0.146 | 🟢 |
| Multicall_Deployless | multicall | 87,646 | 572,433 | 11,410 | 1,747 | 0.153 | 🟢 |
| Multicall_Deployless | multicall | 82,910 | 572,433 | 12,061 | 1,747 | 0.145 | 🟢 |
| Multicall_TokenMetadata | multicall | 106,116 | 425,474 | 9,424 | 2,350 | 0.249 | 🟢 |
| Multicall_TokenMetadata | multicall | 106,266 | 425,474 | 9,410 | 2,350 | 0.250 | 🟢 |
| Multicall_TokenMetadata | multicall | 107,969 | 425,474 | 9,262 | 2,350 | 0.254 | 🟢 |
| Multicall_TokenMetadata | multicall | 106,300 | 425,474 | 9,407 | 2,350 | 0.250 | 🟢 |
| Multicall_TokenMetadata | multicall | 108,342 | 425,474 | 9,230 | 2,350 | 0.255 | 🟢 |

---

## Methodology

### Test Environment

- **Network:** Anvil (Mainnet fork)
- **Go Benchmark:** `go test -bench=. -benchmem -benchtime=10s -count=5`
- **TS Benchmark:** `vitest bench` with 10s per benchmark

### Measurement Notes

- **ns/op:** Nanoseconds per operation (lower is better)
- **ops/s:** Operations per second (higher is better)
- **Ratio:** Go time / TS time (>1 means TS is faster)
- **Tie:** Within 5% of each other

### Caveats

- Network latency dominates most benchmarks (RPC calls)
- Results may vary based on network conditions
- CPU-bound operations may show different characteristics
