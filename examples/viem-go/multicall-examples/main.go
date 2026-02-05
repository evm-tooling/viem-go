package main

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/abi"
	"github.com/ChefBingbong/viem-go/actions/public"
	"github.com/ChefBingbong/viem-go/chain/definitions"
	"github.com/ChefBingbong/viem-go/client"
	"github.com/ChefBingbong/viem-go/client/transport"
)

// Example addresses
var (
	// USDC contract on Polygon
	usdcAddress = common.HexToAddress("0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359")
	// WETH contract on Polygon
	wethAddress = common.HexToAddress("0x7ceB23fD6bC0adD59E62ac25578270cFf1b9f619")
	// WMATIC contract on Polygon
	wmaticAddress = common.HexToAddress("0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270")
	// Vitalik's address for balance checks
	vitalikAddress = common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
	// Random test addresses
	testAddress1 = common.HexToAddress("0x1234567890123456789012345678901234567890")
	testAddress2 = common.HexToAddress("0xdead000000000000000000000000000000000000")
)

// ERC20 ABI for examples
var erc20ABI = []byte(`[
	{"name": "name", "type": "function", "stateMutability": "view", "inputs": [], "outputs": [{"type": "string"}]},
	{"name": "symbol", "type": "function", "stateMutability": "view", "inputs": [], "outputs": [{"type": "string"}]},
	{"name": "decimals", "type": "function", "stateMutability": "view", "inputs": [], "outputs": [{"type": "uint8"}]},
	{"name": "totalSupply", "type": "function", "stateMutability": "view", "inputs": [], "outputs": [{"type": "uint256"}]},
	{"name": "balanceOf", "type": "function", "stateMutability": "view", "inputs": [{"name": "account", "type": "address"}], "outputs": [{"type": "uint256"}]}
]`)

func main() {
	ctx := context.Background()

	printHeader("Multicall Action Examples (viem-go)")

	// Create Public Client with Polygon chain (has multicall3)
	printSection("1. Creating Public Client")
	publicClient, err := client.CreatePublicClient(client.PublicClientConfig{
		Chain:     &definitions.Polygon,
		Transport: transport.HTTP("https://rough-purple-market.matic.quiknode.pro/c1a568726a34041d3c5d58603f5981951e6a8503"),
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}
	defer publicClient.Close()
	fmt.Println("Connected to Polygon Mainnet")

	// Example 1: Basic Multicall - Read multiple token names
	printSection("2. Basic Multicall - Read Multiple Token Names")
	results, err := public.Multicall(ctx, publicClient, public.MulticallParameters{
		Contracts: []public.MulticallContract{
			{
				Address:      usdcAddress,
				ABI:          abi.MustParse(erc20ABI),
				FunctionName: "name",
			},
			{
				Address:      wethAddress,
				ABI:          abi.MustParse(erc20ABI),
				FunctionName: "name",
			},
			{
				Address:      wmaticAddress,
				ABI:          abi.MustParse(erc20ABI),
				FunctionName: "name",
			},
		},
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Token names retrieved in single RPC call:")
		tokens := []string{"USDC", "WETH", "WMATIC"}
		for i, result := range results {
			if result.Status == "success" {
				fmt.Printf("  %s: %v\n", tokens[i], result.Result)
			} else {
				fmt.Printf("  %s: failed - %v\n", tokens[i], result.Error)
			}
		}
	}

	// Example 2: Multicall for Token Metadata
	printSection("3. Multicall - Complete Token Metadata")
	results, err = public.Multicall(ctx, publicClient, public.MulticallParameters{
		Contracts: []public.MulticallContract{
			{
				Address:      usdcAddress,
				ABI:          abi.MustParse(erc20ABI),
				FunctionName: "name",
			},
			{
				Address:      usdcAddress,
				ABI:          abi.MustParse(erc20ABI),
				FunctionName: "symbol",
			},
			{
				Address:      usdcAddress,
				ABI:          abi.MustParse(erc20ABI),
				FunctionName: "decimals",
			},
			{
				Address:      usdcAddress,
				ABI:          abi.MustParse(erc20ABI),
				FunctionName: "totalSupply",
			},
		},
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("USDC Token Metadata:")
		fields := []string{"Name", "Symbol", "Decimals", "Total Supply"}
		for i, result := range results {
			if result.Status == "success" {
				value := result.Result
				// Format total supply with decimals
				if fields[i] == "Total Supply" {
					if supply, ok := value.(*big.Int); ok {
						value = formatTokenAmount(supply, 6) + " USDC"
					}
				}
				fmt.Printf("  %s: %v\n", fields[i], value)
			} else {
				fmt.Printf("  %s: failed - %v\n", fields[i], result.Error)
			}
		}
	}

	// Example 3: Multicall for Multiple Balances
	printSection("4. Multicall - Multiple Balance Queries")
	results, err = public.Multicall(ctx, publicClient, public.MulticallParameters{
		Contracts: []public.MulticallContract{
			{
				Address:      usdcAddress,
				ABI:          abi.MustParse(erc20ABI),
				FunctionName: "balanceOf",
				Args:         []any{vitalikAddress},
			},
			{
				Address:      wethAddress,
				ABI:          abi.MustParse(erc20ABI),
				FunctionName: "balanceOf",
				Args:         []any{vitalikAddress},
			},
			{
				Address:      wmaticAddress,
				ABI:          abi.MustParse(erc20ABI),
				FunctionName: "balanceOf",
				Args:         []any{vitalikAddress},
			},
		},
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Vitalik's Polygon balances (%s):\n", truncateAddress(vitalikAddress))
		tokens := []struct {
			name     string
			decimals int
		}{
			{"USDC", 6},
			{"WETH", 18},
			{"WMATIC", 18},
		}
		for i, result := range results {
			if result.Status == "success" {
				if balance, ok := result.Result.(*big.Int); ok {
					fmt.Printf("  %s: %s\n", tokens[i].name, formatTokenAmount(balance, tokens[i].decimals))
				}
			} else {
				fmt.Printf("  %s: failed - %v\n", tokens[i].name, result.Error)
			}
		}
	}

	// Example 4: Multicall with allowFailure
	printSection("5. Multicall with allowFailure=true (default)")
	// Include an invalid contract call
	invalidAddress := common.HexToAddress("0x0000000000000000000000000000000000000001")
	results, err = public.Multicall(ctx, publicClient, public.MulticallParameters{
		Contracts: []public.MulticallContract{
			{
				Address:      usdcAddress,
				ABI:          abi.MustParse(erc20ABI),
				FunctionName: "name",
			},
			{
				Address:      invalidAddress, // This will fail
				ABI:          abi.MustParse(erc20ABI),
				FunctionName: "name",
			},
			{
				Address:      wethAddress,
				ABI:          abi.MustParse(erc20ABI),
				FunctionName: "name",
			},
		},
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Results (with one failing call):")
		for i, result := range results {
			if result.Status == "success" {
				fmt.Printf("  Call %d: success - %v\n", i+1, result.Result)
			} else {
				fmt.Printf("  Call %d: failure - %v\n", i+1, result.Error)
			}
		}
	}

	// Example 5: Multicall with allowFailure=false
	printSection("6. Multicall with allowFailure=false")
	allowFailure := false
	_, err = public.Multicall(ctx, publicClient, public.MulticallParameters{
		AllowFailure: &allowFailure,
		Contracts: []public.MulticallContract{
			{
				Address:      usdcAddress,
				ABI:          abi.MustParse(erc20ABI),
				FunctionName: "name",
			},
			{
				Address:      invalidAddress, // This will cause the entire multicall to fail
				ABI:          abi.MustParse(erc20ABI),
				FunctionName: "name",
			},
		},
	})
	if err != nil {
		fmt.Printf("Expected error with allowFailure=false: %v\n", truncateError(err, 80))
	} else {
		fmt.Println("Unexpected: no error")
	}

	// Example 6: Large Multicall with Chunking
	printSection("7. Large Multicall with Automatic Chunking")
	// Create many calls to trigger chunking
	var manyContracts []public.MulticallContract
	addresses := []common.Address{usdcAddress, wethAddress, wmaticAddress}
	for i := 0; i < 30; i++ {
		manyContracts = append(manyContracts, public.MulticallContract{
			Address:      addresses[i%3],
			ABI:          abi.MustParse(erc20ABI),
			FunctionName: "balanceOf",
			Args:         []any{vitalikAddress},
		})
	}

	results, err = public.Multicall(ctx, publicClient, public.MulticallParameters{
		Contracts:           manyContracts,
		BatchSize:           512, // Small batch size to force chunking
		MaxConcurrentChunks: 4,   // Execute up to 4 chunks in parallel
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		successCount := 0
		for _, result := range results {
			if result.Status == "success" {
				successCount++
			}
		}
		fmt.Printf("Executed %d calls successfully\n", successCount)
		fmt.Printf("BatchSize: 512 bytes (forces multiple chunks)\n")
		fmt.Printf("MaxConcurrentChunks: 4 (parallel execution)\n")
	}

	// Example 7: Multicall at Specific Block
	printSection("8. Multicall at Specific Block Number")
	blockNum := uint64(52000000) // Historical Polygon block
	results, err = public.Multicall(ctx, publicClient, public.MulticallParameters{
		BlockNumber: &blockNum,
		Contracts: []public.MulticallContract{
			{
				Address:      usdcAddress,
				ABI:          abi.MustParse(erc20ABI),
				FunctionName: "totalSupply",
			},
		},
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		if results[0].Status == "success" {
			if supply, ok := results[0].Result.(*big.Int); ok {
				fmt.Printf("USDC Total Supply at block %d: %s\n", blockNum, formatTokenAmount(supply, 6))
			}
		}
	}

	// Example 8: Deployless Multicall
	printSection("9. Deployless Multicall (no deployed multicall3 needed)")
	results, err = public.Multicall(ctx, publicClient, public.MulticallParameters{
		Deployless: true, // Use bytecode execution instead of deployed contract
		Contracts: []public.MulticallContract{
			{
				Address:      usdcAddress,
				ABI:          abi.MustParse(erc20ABI),
				FunctionName: "symbol",
			},
			{
				Address:      wethAddress,
				ABI:          abi.MustParse(erc20ABI),
				FunctionName: "symbol",
			},
		},
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Symbols via deployless multicall:")
		for i, result := range results {
			if result.Status == "success" {
				fmt.Printf("  Token %d: %v\n", i+1, result.Result)
			}
		}
	}

	// Example 9: Cross-contract Queries in Single Call
	printSection("10. Cross-Contract Queries")
	results, err = public.Multicall(ctx, publicClient, public.MulticallParameters{
		Contracts: []public.MulticallContract{
			// USDC metadata
			{Address: usdcAddress, ABI: abi.MustParse(erc20ABI), FunctionName: "name"},
			{Address: usdcAddress, ABI: abi.MustParse(erc20ABI), FunctionName: "decimals"},
			// WETH metadata
			{Address: wethAddress, ABI: abi.MustParse(erc20ABI), FunctionName: "name"},
			{Address: wethAddress, ABI: abi.MustParse(erc20ABI), FunctionName: "decimals"},
			// WMATIC metadata
			{Address: wmaticAddress, ABI: abi.MustParse(erc20ABI), FunctionName: "name"},
			{Address: wmaticAddress, ABI: abi.MustParse(erc20ABI), FunctionName: "decimals"},
		},
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Multi-contract metadata in single RPC call:")
		for i := 0; i < 6; i += 2 {
			name := "unknown"
			decimals := "?"
			if results[i].Status == "success" {
				name = fmt.Sprintf("%v", results[i].Result)
			}
			if results[i+1].Status == "success" {
				decimals = fmt.Sprintf("%v", results[i+1].Result)
			}
			fmt.Printf("  %s: %s decimals\n", name, decimals)
		}
	}

	// Summary
	printHeader("Examples Complete")
	fmt.Println("Demonstrated Multicall features:")
	fmt.Println("  - Basic multicall batching")
	fmt.Println("  - Token metadata retrieval")
	fmt.Println("  - Multiple balance queries")
	fmt.Println("  - Error handling with allowFailure")
	fmt.Println("  - Large multicalls with chunking")
	fmt.Println("  - Historical block queries")
	fmt.Println("  - Deployless multicall")
	fmt.Println("  - Cross-contract queries")
	fmt.Println()
	fmt.Println("Go concurrency features used:")
	fmt.Println("  - errgroup.WithContext for parallel chunk execution")
	fmt.Println("  - SetLimit() for bounded concurrency")
	fmt.Println("  - Context propagation for cancellation")
	fmt.Println()
}

// Helper functions

func printHeader(title string) {
	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("  %s\n", title)
	fmt.Println(strings.Repeat("=", 60))
}

func printSection(title string) {
	fmt.Printf("\n--- %s ---\n", title)
}

func truncateAddress(addr common.Address) string {
	hex := addr.Hex()
	return hex[:10] + "..." + hex[len(hex)-4:]
}

func formatTokenAmount(amount *big.Int, decimals int) string {
	if amount == nil {
		return "0"
	}
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	whole := new(big.Int).Div(amount, divisor)
	return whole.String()
}

func truncateError(err error, maxLen int) string {
	s := err.Error()
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
