# Benchmark Comparison: viem-go vs viem TypeScript

Generated: 2026-02-07T03:09:39.807Z

## Overall Summary

**🏆 Go is 3.19x faster overall**

| Metric | Go | TypeScript |
|--------|----|-----------|
| Avg ns/op | 3,413,495 | 10,903,465 |
| Avg ops/s | 293 | 92 |
| Wins | 55/59 | 4/59 |

## By Suite

| Suite | Benchmarks | Go Wins | TS Wins | Ties | Winner |
|-------|------------|---------|---------|------|--------|
| abi | 6 | 6 | 0 | 0 | 🟢 Go 15.46x faster |
| address | 5 | 2 | 3 | 0 | 🟢 Go 3.24x faster |
| call | 6 | 6 | 0 | 0 | 🟢 Go 84.60x faster |
| ens | 5 | 5 | 0 | 0 | 🟢 Go 11.54x faster |
| event | 3 | 3 | 0 | 0 | 🟢 Go 27.18x faster |
| hash | 7 | 7 | 0 | 0 | 🟢 Go 9.91x faster |
| multicall | 16 | 16 | 0 | 0 | 🟢 Go 2.72x faster |
| signature | 5 | 5 | 0 | 0 | 🟢 Go 76.86x faster |
| unit | 6 | 5 | 1 | 0 | 🟢 Go 1.51x faster |

## Detailed Results

| Benchmark | Go (ns/op) | TS (ns/op) | Go (ops/s) | TS (ops/s) | Result |
|-----------|------------|------------|------------|------------|--------|
| Abi_EncodeSimple | 206 | 6,763 | 4,852,014 | 147,858 | 🟢 Go 32.82x faster |
| Abi_EncodeComplex | 311 | 8,792 | 3,211,304 | 113,738 | 🟢 Go 28.23x faster |
| Abi_EncodeMultiArg | 440 | 8,405 | 2,273,761 | 118,973 | 🟢 Go 19.11x faster |
| Abi_DecodeResult | 92 | 1,018 | 10,909,884 | 982,751 | 🟢 Go 11.10x faster |
| Abi_EncodePacked | 313 | 675 | 3,198,976 | 1,480,831 | 🟢 Go 2.16x faster |
| Abi_EncodePackedMulti | 371 | 1,131 | 2,695,418 | 884,009 | 🟢 Go 3.05x faster |
| Address_IsAddress | 910 | 232 | 1,098,418 | 4,317,001 | 🔵 TS 3.93x faster |
| Address_IsAddressLower | 300 | 232 | 3,335,557 | 4,311,822 | 🔵 TS 1.29x faster |
| Address_Checksum | 801 | 615 | 1,249,063 | 1,626,527 | 🔵 TS 1.30x faster |
| Address_Create | 2,339 | 7,412 | 427,533 | 134,917 | 🟢 Go 3.17x faster |
| Address_Create2 | 2,632 | 14,105 | 379,939 | 70,898 | 🟢 Go 5.36x faster |
| Call_Basic | 184,251 | 18,886,562 | 5,427 | 53 | 🟢 Go 102.50x faster |
| Call_WithData | 179,284 | 18,646,628 | 5,578 | 54 | 🟢 Go 104.01x faster |
| Call_WithAccount | 181,105 | 259,087 | 5,522 | 3,860 | 🟢 Go 1.43x faster |
| Call_Decimals | 183,314 | 18,230,344 | 5,455 | 55 | 🟢 Go 99.45x faster |
| Call_Symbol | 188,984 | 18,337,453 | 5,291 | 55 | 🟢 Go 97.03x faster |
| Call_BalanceOfMultiple | 180,615 | 18,496,083 | 5,537 | 54 | 🟢 Go 102.41x faster |
| Ens_Namehash | 1,556 | 20,659 | 642,674 | 48,404 | 🟢 Go 13.28x faster |
| Ens_NamehashDeep | 3,024 | 40,884 | 330,688 | 24,459 | 🟢 Go 13.52x faster |
| Ens_Labelhash | 429 | 5,358 | 2,332,090 | 186,642 | 🟢 Go 12.49x faster |
| Ens_Normalize | 344 | 1,318 | 2,910,361 | 758,845 | 🟢 Go 3.84x faster |
| Ens_NormalizeLong | 887 | 3,776 | 1,127,015 | 264,859 | 🟢 Go 4.26x faster |
| Event_DecodeTransfer | 374 | 10,329 | 2,672,368 | 96,812 | 🟢 Go 27.60x faster |
| Event_DecodeBatch10 | 3,733 | 103,003 | 267,881 | 9,708 | 🟢 Go 27.59x faster |
| Event_DecodeBatch100 | 37,273 | 1,011,358 | 26,829 | 989 | 🟢 Go 27.13x faster |
| Hash_Keccak256Short | 425 | 5,292 | 2,352,941 | 188,951 | 🟢 Go 12.45x faster |
| Hash_Keccak256Long | 2,667 | 45,453 | 374,953 | 22,001 | 🟢 Go 17.04x faster |
| Hash_Keccak256Hex | 453 | 5,302 | 2,208,968 | 188,614 | 🟢 Go 11.71x faster |
| Hash_Sha256Short | 153 | 1,634 | 6,523,157 | 612,024 | 🟢 Go 10.66x faster |
| Hash_Sha256Long | 700 | 16,668 | 1,428,776 | 59,996 | 🟢 Go 23.81x faster |
| Hash_FunctionSelector | 1,970 | 6,270 | 507,614 | 159,483 | 🟢 Go 3.18x faster |
| Hash_EventSelector | 2,428 | 6,562 | 411,862 | 152,400 | 🟢 Go 2.70x faster |
| Multicall_Basic | 229,551 | 470,754 | 4,356 | 2,124 | 🟢 Go 2.05x faster |
| Multicall_WithArgs | 219,662 | 421,807 | 4,552 | 2,371 | 🟢 Go 1.92x faster |
| Multicall_MultiContract | 265,598 | 392,455 | 3,765 | 2,548 | 🟢 Go 1.48x faster |
| Multicall_10Calls | 315,944 | 512,130 | 3,165 | 1,953 | 🟢 Go 1.62x faster |
| Multicall_30Calls | 523,221 | 880,421 | 1,911 | 1,136 | 🟢 Go 1.68x faster |
| Multicall_Deployless | 381,328 | 561,381 | 2,622 | 1,781 | 🟢 Go 1.47x faster |
| Multicall_TokenMetadata | 238,010 | 364,328 | 4,202 | 2,745 | 🟢 Go 1.53x faster |
| Multicall_50Calls | 1,122,238 | 1,272,864 | 891 | 786 | 🟢 Go 1.13x faster |
| Multicall_100Calls | 1,456,754 | 2,218,131 | 686 | 451 | 🟢 Go 1.52x faster |
| Multicall_200Calls | 2,503,890 | 4,166,840 | 399 | 240 | 🟢 Go 1.66x faster |
| Multicall_500Calls | 3,721,309 | 9,017,946 | 269 | 111 | 🟢 Go 2.42x faster |
| Multicall_MixedContracts_100 | 1,356,606 | 2,375,918 | 737 | 421 | 🟢 Go 1.75x faster |
| Multicall_1000Calls | 5,478,802 | 17,300,649 | 183 | 58 | 🟢 Go 3.16x faster |
| Multicall_10000Calls_SingleRPC | 121,108,049 | 166,933,761 | 8 | 6 | 🟢 Go 1.38x faster |
| Multicall_10000Calls_Chunked | 29,145,681 | 171,865,601 | 34 | 6 | 🟢 Go 5.90x faster |
| Multicall_10000Calls_AggressiveChunking | 32,109,950 | 166,049,516 | 31 | 6 | 🟢 Go 5.17x faster |
| Signature_HashMessage | 750 | 13,327 | 1,333,156 | 75,035 | 🟢 Go 17.77x faster |
| Signature_HashMessageLong | 1,717 | 19,108 | 582,411 | 52,335 | 🟢 Go 11.13x faster |
| Signature_RecoverAddress | 25,956 | 2,568,581 | 38,527 | 389 | 🟢 Go 98.96x faster |
| Signature_VerifyMessage | 27,430 | 1,706,485 | 36,456 | 586 | 🟢 Go 62.21x faster |
| Signature_ParseSignature | 213 | 1,766 | 4,688,233 | 566,376 | 🟢 Go 8.28x faster |
| Unit_ParseEther | 117 | 337 | 8,561,644 | 2,971,394 | 🟢 Go 2.88x faster |
| Unit_ParseEtherLarge | 313 | 234 | 3,198,976 | 4,280,944 | 🔵 TS 1.34x faster |
| Unit_FormatEther | 116 | 153 | 8,591,065 | 6,545,929 | 🟢 Go 1.31x faster |
| Unit_ParseUnits6 | 101 | 209 | 9,871,668 | 4,773,301 | 🟢 Go 2.07x faster |
| Unit_ParseGwei | 103 | 202 | 9,699,321 | 4,958,003 | 🟢 Go 1.96x faster |
| Unit_FormatUnits | 98 | 144 | 10,175,010 | 6,935,395 | 🟢 Go 1.47x faster |

## Win Summary

- 🟢 Go wins: 55 (93%)
- 🔵 TS wins: 4 (7%)
- ⚪ Ties: 0 (0%)

## Notes

- Benchmarks run against the same Anvil instance (mainnet fork) for fair comparison
- ns/op = nanoseconds per operation (lower is better)
- ops/s = operations per second (higher is better)
- 🟢 = Go faster, 🔵 = TS faster, ⚪ = Similar (within 5%)
