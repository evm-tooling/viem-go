# viem-go

Go Interface for Ethereum - inspired by [viem](https://viem.sh)

[![CI](https://github.com/ChefBingbong/viem-go/actions/workflows/ci.yml/badge.svg)](https://github.com/ChefBingbong/viem-go/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/ChefBingbong/viem-go.svg)](https://pkg.go.dev/github.com/ChefBingbong/viem-go)

**viem-go** is a Go port of the popular TypeScript [viem](https://github.com/wevm/viem) library. It provides idiomatic Go APIs for interacting with Ethereum, following viem's design philosophy while embracing Go conventions.

> **Note:** This project is under active development. APIs may change.

## Installation

```bash
go get github.com/ChefBingbong/viem-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/ChefBingbong/viem-go/client"
)

func main() {
    // Create a client
    c, err := client.NewClient("https://eth.llamarpc.com")
    if err != nil {
        log.Fatal(err)
    }
    defer c.Close()

    // Get block number
    blockNumber, err := c.GetBlockNumber(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Block number:", blockNumber)
}
```

## Side-by-Side Comparison: viem vs viem-go

### Creating a Client

**viem (TypeScript)**
```typescript
import { createPublicClient, http } from 'viem'
import { mainnet } from 'viem/chains'

const client = createPublicClient({
  chain: mainnet,
  transport: http(),
})

const blockNumber = await client.getBlockNumber()
```

**viem-go**
```go
import (
    "context"
    "github.com/ChefBingbong/viem-go/client"
)

c, err := client.NewClient("https://eth.llamarpc.com")
if err != nil {
    log.Fatal(err)
}
defer c.Close()

blockNumber, err := c.GetBlockNumber(context.Background())
```

### Reading from a Contract (ERC20)

**viem (TypeScript)**
```typescript
import { createPublicClient, http, parseAbi } from 'viem'
import { mainnet } from 'viem/chains'

const client = createPublicClient({
  chain: mainnet,
  transport: http(),
})

const balance = await client.readContract({
  address: '0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48',
  abi: parseAbi(['function balanceOf(address) view returns (uint256)']),
  functionName: 'balanceOf',
  args: ['0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045'],
})
```

**viem-go**
```go
import (
    "context"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ChefBingbong/viem-go/client"
    "github.com/ChefBingbong/viem-go/contracts/erc20"
)

c, _ := client.NewClient("https://eth.llamarpc.com")
defer c.Close()

usdc := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
token, _ := erc20.New(usdc, c)

balance, err := token.BalanceOf(
    context.Background(),
    common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"),
)
```

### Encoding Function Data

**viem (TypeScript)**
```typescript
import { encodeFunctionData, parseAbi } from 'viem'

const data = encodeFunctionData({
  abi: parseAbi(['function transfer(address to, uint256 amount)']),
  functionName: 'transfer',
  args: ['0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045', 1000000n],
})
```

**viem-go**
```go
import (
    "math/big"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ChefBingbong/viem-go/abi"
)

parsed, _ := abi.Parse([]byte(`[{
    "name": "transfer",
    "type": "function",
    "inputs": [
        {"name": "to", "type": "address"},
        {"name": "amount", "type": "uint256"}
    ]
}]`))

data, err := parsed.EncodeFunctionData("transfer",
    common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"),
    big.NewInt(1000000),
)
```

### Parsing Units

**viem (TypeScript)**
```typescript
import { parseEther, parseUnits, formatEther } from 'viem'

const wei = parseEther('1.5')        // 1500000000000000000n
const usdc = parseUnits('100', 6)    // 100000000n
const eth = formatEther(wei)         // "1.5"
```

**viem-go**
```go
import "github.com/ChefBingbong/viem-go/utils/unit"

wei, _ := unit.ParseEther("1.5")       // *big.Int: 1500000000000000000
usdc, _ := unit.ParseUnits("100", 6)   // *big.Int: 100000000
eth := unit.FormatEther(wei)           // "1.5"
```

### Hashing & Signatures

**viem (TypeScript)**
```typescript
import { keccak256, hashMessage, recoverMessageAddress } from 'viem'

const hash = keccak256('0x68656c6c6f')
const messageHash = hashMessage('hello world')
const address = await recoverMessageAddress({
  message: 'hello world',
  signature: '0x...',
})
```

**viem-go**
```go
import (
    "github.com/ChefBingbong/viem-go/utils/hash"
    "github.com/ChefBingbong/viem-go/utils/signature"
)

h := hash.Keccak256([]byte("hello"))
messageHash := signature.HashMessage("hello world")
address, err := signature.RecoverMessageAddress("hello world", sig)
```

## Implementation Status
early 

## Development

```bash
# Install dependencies
go mod download

# Run tests
make test

# Run linter
make lint

# Format code
make fmt
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Acknowledgments

- [viem](https://github.com/wevm/viem) - The original TypeScript library this project is based on
- [go-ethereum](https://github.com/ethereum/go-ethereum) - Ethereum Go implementation used for cryptographic primitives
