# Benchmark Comparison: viem-go vs viem TypeScript

Generated: 2026-02-07T06:17:34.965Z

## Overall Summary

**🏆 Go is 6.81x faster overall**

| Metric | Go | TypeScript |
|--------|----|-----------|
| Avg ns/op | 2,036,775 | 13,861,286 |
| Avg ops/s | 491 | 72 |
| Wins | 53/59 | 4/59 |

## By Suite

| Suite | Benchmarks | Go Wins | TS Wins | Ties | Winner |
|-------|------------|---------|---------|------|--------|
| abi | 6 | 6 | 0 | 0 | 🟢 Go 17.29x faster |
| address | 5 | 2 | 2 | 1 | 🟢 Go 4.07x faster |
| call | 6 | 6 | 0 | 0 | 🟢 Go 71.61x faster |
| ens | 5 | 5 | 0 | 0 | 🟢 Go 14.88x faster |
| event | 3 | 3 | 0 | 0 | 🟢 Go 29.82x faster |
| hash | 7 | 7 | 0 | 0 | 🟢 Go 12.47x faster |
| multicall | 16 | 15 | 0 | 1 | 🟢 Go 6.07x faster |
| signature | 5 | 5 | 0 | 0 | 🟢 Go 59.63x faster |
| unit | 6 | 4 | 2 | 0 | 🟢 Go 1.25x faster |

## Detailed Results

| Benchmark | Go (ns/op) | TS (ns/op) | Go (ops/s) | TS (ops/s) | Result |
|-----------|------------|------------|------------|------------|--------|
| Abi_EncodeSimple | 215 | 8,644 | 4,642,526 | 115,685 | 🟢 Go 40.13x faster |
| Abi_EncodeComplex | 328 | 9,533 | 3,050,641 | 104,895 | 🟢 Go 29.08x faster |
| Abi_EncodeMultiArg | 458 | 10,030 | 2,184,360 | 99,700 | 🟢 Go 21.91x faster |
| Abi_DecodeResult | 94 | 1,073 | 10,624,734 | 931,629 | 🟢 Go 11.40x faster |
| Abi_EncodePacked | 322 | 703 | 3,106,555 | 1,422,858 | 🟢 Go 2.18x faster |
| Abi_EncodePackedMulti | 391 | 1,281 | 2,556,891 | 780,533 | 🟢 Go 3.28x faster |
| Address_IsAddress | 913 | 293 | 1,095,290 | 3,418,491 | 🔵 TS 3.12x faster |
| Address_IsAddressLower | 299 | 295 | 3,348,962 | 3,395,274 | ⚪ Similar |
| Address_Checksum | 812 | 722 | 1,231,527 | 1,384,784 | 🔵 TS 1.12x faster |
| Address_Create | 2,358 | 9,761 | 424,088 | 102,451 | 🟢 Go 4.14x faster |
| Address_Create2 | 2,623 | 17,410 | 381,243 | 57,437 | 🟢 Go 6.64x faster |
| Call_Basic | 299,953 | 19,322,365 | 3,334 | 52 | 🟢 Go 64.42x faster |
| Call_WithData | 206,460 | 18,632,279 | 4,844 | 54 | 🟢 Go 90.25x faster |
| Call_WithAccount | 186,727 | 289,804 | 5,355 | 3,451 | 🟢 Go 1.55x faster |
| Call_Decimals | 190,463 | 17,701,183 | 5,250 | 56 | 🟢 Go 92.94x faster |
| Call_Symbol | 202,769 | 17,923,138 | 4,932 | 56 | 🟢 Go 88.39x faster |
| Call_BalanceOfMultiple | 201,730 | 18,375,462 | 4,957 | 54 | 🟢 Go 91.09x faster |
| Ens_Namehash | 1,607 | 28,027 | 622,278 | 35,680 | 🟢 Go 17.44x faster |
| Ens_NamehashDeep | 3,058 | 55,594 | 327,011 | 17,987 | 🟢 Go 18.18x faster |
| Ens_Labelhash | 434 | 7,205 | 2,303,617 | 138,784 | 🟢 Go 16.60x faster |
| Ens_Normalize | 355 | 1,000 | 2,816,108 | 1,000,321 | 🟢 Go 2.82x faster |
| Ens_NormalizeLong | 896 | 2,663 | 1,116,196 | 375,460 | 🟢 Go 2.97x faster |
| Event_DecodeTransfer | 398 | 11,837 | 2,511,301 | 84,484 | 🟢 Go 29.73x faster |
| Event_DecodeBatch10 | 4,098 | 123,131 | 244,021 | 8,121 | 🟢 Go 30.05x faster |
| Event_DecodeBatch100 | 39,847 | 1,187,409 | 25,096 | 842 | 🟢 Go 29.80x faster |
| Hash_Keccak256Short | 444 | 7,543 | 2,253,775 | 132,569 | 🟢 Go 17.00x faster |
| Hash_Keccak256Long | 2,679 | 60,971 | 373,274 | 16,401 | 🟢 Go 22.76x faster |
| Hash_Keccak256Hex | 452 | 7,119 | 2,213,369 | 140,464 | 🟢 Go 15.76x faster |
| Hash_Sha256Short | 159 | 1,487 | 6,289,308 | 672,569 | 🟢 Go 9.35x faster |
| Hash_Sha256Long | 623 | 14,108 | 1,605,136 | 70,882 | 🟢 Go 22.65x faster |
| Hash_FunctionSelector | 1,933 | 8,651 | 517,331 | 115,598 | 🟢 Go 4.48x faster |
| Hash_EventSelector | 2,388 | 8,365 | 418,760 | 119,545 | 🟢 Go 3.50x faster |
| Multicall_Basic | 197,920 | 469,942 | 5,053 | 2,128 | 🟢 Go 2.37x faster |
| Multicall_WithArgs | 233,293 | 394,218 | 4,286 | 2,537 | 🟢 Go 1.69x faster |
| Multicall_MultiContract | 481,325 | 466,651 | 2,078 | 2,143 | ⚪ Similar |
| Multicall_10Calls | 322,710 | 522,065 | 3,099 | 1,915 | 🟢 Go 1.62x faster |
| Multicall_30Calls | 508,492 | 1,020,637 | 1,967 | 980 | 🟢 Go 2.01x faster |
| Multicall_Deployless | 370,241 | 687,205 | 2,701 | 1,455 | 🟢 Go 1.86x faster |
| Multicall_TokenMetadata | 230,417 | 430,263 | 4,340 | 2,324 | 🟢 Go 1.87x faster |
| Multicall_50Calls | 512,859 | 1,499,093 | 1,950 | 667 | 🟢 Go 2.92x faster |
| Multicall_100Calls | 926,012 | 2,662,194 | 1,080 | 376 | 🟢 Go 2.87x faster |
| Multicall_200Calls | 1,448,256 | 5,728,689 | 690 | 175 | 🟢 Go 3.96x faster |
| Multicall_500Calls | 2,170,623 | 10,795,604 | 461 | 93 | 🟢 Go 4.97x faster |
| Multicall_MixedContracts_100 | 865,080 | 2,553,952 | 1,156 | 392 | 🟢 Go 2.95x faster |
| Multicall_1000Calls | 3,002,460 | 20,927,595 | 333 | 48 | 🟢 Go 6.97x faster |
| Multicall_10000Calls_SingleRPC | 67,045,603 | 209,200,644 | 15 | 5 | 🟢 Go 3.12x faster |
| Multicall_10000Calls_Chunked | 20,333,353 | 218,966,914 | 49 | 5 | 🟢 Go 10.77x faster |
| Multicall_10000Calls_AggressiveChunking | 20,108,967 | 244,385,249 | 50 | 4 | 🟢 Go 12.15x faster |
| Signature_HashMessage | 768 | 8,647 | 1,302,762 | 115,648 | 🟢 Go 11.26x faster |
| Signature_HashMessageLong | 1,790 | 18,234 | 558,659 | 54,843 | 🟢 Go 10.19x faster |
| Signature_RecoverAddress | 26,092 | 1,673,780 | 38,326 | 597 | 🟢 Go 64.15x faster |
| Signature_VerifyMessage | 26,082 | 1,572,154 | 38,341 | 636 | 🟢 Go 60.28x faster |
| Signature_ParseSignature | 185 | 1,908 | 5,408,329 | 524,037 | 🟢 Go 10.32x faster |
| Unit_ParseEther | 116 | 246 | 8,635,579 | 4,068,068 | 🟢 Go 2.12x faster |
| Unit_ParseEtherLarge | 318 | 233 | 3,145,643 | 4,293,864 | 🔵 TS 1.37x faster |
| Unit_FormatEther | 118 | 143 | 8,503,401 | 6,975,713 | 🟢 Go 1.22x faster |
| Unit_ParseUnits6 | 140 | 218 | 7,122,507 | 4,586,885 | 🟢 Go 1.55x faster |
| Unit_ParseGwei | 105 | 203 | 9,569,378 | 4,927,703 | 🟢 Go 1.94x faster |
| Unit_FormatUnits | 145 | 133 | 6,901,311 | 7,507,014 | 🔵 TS 1.09x faster |

## Win Summary

- 🟢 Go wins: 53 (90%)
- 🔵 TS wins: 4 (7%)
- ⚪ Ties: 2 (3%)

## Notes

- Benchmarks run against the same Anvil instance (mainnet fork) for fair comparison
- ns/op = nanoseconds per operation (lower is better)
- ops/s = operations per second (higher is better)
- 🟢 = Go faster, 🔵 = TS faster, ⚪ = Similar (within 5%)
