// Call Examples - viem-go
//
// Comprehensive examples demonstrating the Call action features:
// - Basic contract calls
// - Call with various parameters (gas, value, fees)
// - State and block overrides
// - Deployless calls
// - Access lists
// - Error handling
package main

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/actions/public"
	"github.com/ChefBingbong/viem-go/chain/definitions"
	"github.com/ChefBingbong/viem-go/client"
	"github.com/ChefBingbong/viem-go/client/transport"
	"github.com/ChefBingbong/viem-go/contracts/erc20"
	"github.com/ChefBingbong/viem-go/types"
	"github.com/ChefBingbong/viem-go/utils/unit"
)

// Example addresses
var (
	// USDC contract on Mainnet
	usdcAddress = common.HexToAddress("0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359")
	// Vitalik's address for balance checks
	vitalikAddress = common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
	// Random test address
	testAddress = common.HexToAddress("0x1234567890123456789012345678901234567890")
)

// ERC20 function selectors
var (
	// balanceOf(address) selector
	balanceOfSelector = common.Hex2Bytes("70a08231")
	// name() selector
	nameSelector = common.Hex2Bytes("06fdde03")
	// decimals() selector
	decimalsSelector = common.Hex2Bytes("313ce567")
	// totalSupply() selector
	totalSupplySelector = common.Hex2Bytes("18160ddd")
)

func main() {
	ctx := context.Background()

	printHeader("Call Action Examples (viem-go)")

	// Create Public Client
	printSection("1. Creating Public Client")
	publicClient, err := client.CreatePublicClient(client.PublicClientConfig{
		Chain:     &definitions.Polygon,
		Transport: transport.HTTP("https://polygon-rpc.com"),
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}
	defer publicClient.Close()
	fmt.Println("Connected to Ethereum Mainnet")

	// Example 1: Basic Call - Read contract name
	printSection("2. Basic Call - Read ERC20 Name")
	result, err := public.Call(ctx, publicClient, public.CallParameters{
		To:   &usdcAddress,
		Data: nameSelector,
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		// Decode string result (skip first 64 bytes of offset and length)
		if len(result.Data) > 64 {
			name := strings.TrimRight(string(result.Data[64:]), "\x00")
			fmt.Printf("Contract Name: %s\n", name)
		}
		fmt.Printf("Raw Result: 0x%x\n", truncateBytes(result.Data, 32))
	}

	// Example 2: Call with Address Parameter - balanceOf
	printSection("3. Call with Parameter - balanceOf(address)")
	// Encode balanceOf(vitalikAddress) - selector + padded address
	balanceOfData := append(balanceOfSelector, common.LeftPadBytes(vitalikAddress.Bytes(), 32)...)
	result, err = public.Call(ctx, publicClient, public.CallParameters{
		To:   &usdcAddress,
		Data: balanceOfData,
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		balance := new(big.Int).SetBytes(result.Data)
		// USDC has 6 decimals
		fmt.Printf("Vitalik's USDC Balance: %s (raw: %s)\n",
			formatTokenAmount(balance, 6), balance.String())
	}

	// Example 3: Call with From Address (Account)
	printSection("4. Call with Account (from address)")
	result, err = public.Call(ctx, publicClient, public.CallParameters{
		Account: &vitalikAddress,
		To:      &usdcAddress,
		Data:    decimalsSelector,
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		decimals := new(big.Int).SetBytes(result.Data)
		fmt.Printf("USDC Decimals: %d (called from %s)\n", decimals.Uint64(), truncateAddress(vitalikAddress))
	}

	// Example 4: Call with Block Number
	printSection("5. Call at Specific Block Number")
	blockNum := uint64(82563588) // Historical block
	result, err = public.Call(ctx, publicClient, public.CallParameters{
		To:          &usdcAddress,
		Data:        totalSupplySelector,
		BlockNumber: &blockNum,
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		supply := new(big.Int).SetBytes(result.Data)
		fmt.Printf("USDC Total Supply at block %d: %s\n", blockNum, formatTokenAmount(supply, 6))
	}

	// Example 5: Call with Block Tag
	printSection("6. Call with Block Tag (pending)")
	result, err = public.Call(ctx, publicClient, public.CallParameters{
		To:       &usdcAddress,
		Data:     totalSupplySelector,
		BlockTag: public.BlockTagPending,
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		supply := new(big.Int).SetBytes(result.Data)
		fmt.Printf("USDC Total Supply (pending): %s\n", formatTokenAmount(supply, 6))
	}

	// Example 6: Call with Gas Parameters
	printSection("7. Call with Gas Parameters")
	gas := uint64(100000)
	gasPrice := mustParseGwei("700")
	result, err = public.Call(ctx, publicClient, public.CallParameters{
		To:       &usdcAddress,
		Data:     nameSelector,
		Gas:      &gas,
		GasPrice: gasPrice,
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Call succeeded with gas=%d, gasPrice=%s gwei\n", gas, unit.FormatGwei(gasPrice))
		fmt.Printf("Result length: %d bytes\n", len(result.Data))
	}

	// Example 7: Call with EIP-1559 Fees
	printSection("8. Call with EIP-1559 Fees")
	maxFeePerGas := mustParseGwei("700")
	maxPriorityFeePerGas := mustParseGwei("2")
	result, err = public.Call(ctx, publicClient, public.CallParameters{
		To:                   &usdcAddress,
		Data:                 decimalsSelector,
		MaxFeePerGas:         maxFeePerGas,
		MaxPriorityFeePerGas: maxPriorityFeePerGas,
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Call succeeded with maxFeePerGas=%s gwei, maxPriorityFeePerGas=%s gwei\n",
			unit.FormatGwei(maxFeePerGas), unit.FormatGwei(maxPriorityFeePerGas))
	}

	// Example 8: Call with Value (simulating ETH transfer check)
	printSection("9. Call with Value")
	value := mustParseEther("0.1")
	result, err = public.Call(ctx, publicClient, public.CallParameters{
		Account: &vitalikAddress,
		To:      &testAddress,
		Value:   value,
		Data:    []byte{}, // Empty call with value
	})
	if err != nil {
		fmt.Printf("Simulated transfer of %s ETH: Error - %v\n", unit.FormatEther(value), err)
	} else {
		fmt.Printf("Simulated transfer of %s ETH would succeed\n", unit.FormatEther(value))
	}

	// Example 9: Call with State Override
	printSection("10. Call with State Override")
	// Override an address to have a large ETH balance
	overrideBalance, _ := new(big.Int).SetString("1000000000000000000000", 10) // 1000 ETH
	stateOverride := types.StateOverride{
		testAddress: types.StateOverrideAccount{
			Balance: overrideBalance,
		},
	}
	result, err = public.Call(ctx, publicClient, public.CallParameters{
		Account:       &testAddress,
		To:            &vitalikAddress,
		Value:         mustParseEther("100"), // Transfer 100 ETH
		StateOverride: stateOverride,
	})
	if err != nil {
		fmt.Printf("Error with state override: %v\n", err)
	} else {
		fmt.Printf("State override successful! Test address had balance overridden to 1000 ETH\n")
		fmt.Printf("Simulated 100 ETH transfer succeeded\n")
	}

	// Example 10: Call with Block Override
	printSection("11. Call with Block Override")
	overrideGasLimit := uint64(50000000)    // 50M gas limit
	overrideBaseFee := mustParseGwei("1")   // 1 gwei base fee
	overrideTimestamp := uint64(1700000000) // Fixed timestamp
	blockOverrides := &types.BlockOverrides{
		GasLimit:      &overrideGasLimit,
		BaseFeePerGas: overrideBaseFee,
		Time:          &overrideTimestamp,
	}
	result, err = public.Call(ctx, publicClient, public.CallParameters{
		To:             &usdcAddress,
		Data:           totalSupplySelector,
		BlockOverrides: blockOverrides,
	})
	if err != nil {
		fmt.Printf("Error with block override: %v\n", err)
	} else {
		fmt.Printf("Block override successful!\n")
		fmt.Printf("  Simulated gasLimit: %d\n", overrideGasLimit)
		fmt.Printf("  Simulated baseFee: %s gwei\n", unit.FormatGwei(overrideBaseFee))
	}

	// Example 11: Call with Access List
	printSection("12. Call with Access List (EIP-2930)")
	storageSlot := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
	accessList := types.AccessList{
		{
			Address:     usdcAddress,
			StorageKeys: []common.Hash{storageSlot},
		},
	}
	result, err = public.Call(ctx, publicClient, public.CallParameters{
		To:         &usdcAddress,
		Data:       nameSelector,
		AccessList: accessList,
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Call with access list succeeded\n")
		fmt.Printf("  Pre-warmed contract: %s\n", truncateAddress(usdcAddress))
		fmt.Printf("  Pre-warmed storage slot: %s\n", storageSlot.Hex()[:18]+"...")
	}

	// Example 12: Deployless Call (Code parameter)
	printSection("13. Deployless Call (execute bytecode without deployment)")
	// Simple bytecode that returns 42: PUSH1 0x2a PUSH1 0x00 MSTORE PUSH1 0x20 PUSH1 0x00 RETURN
	// This bytecode stores 42 at memory position 0 and returns 32 bytes
	simpleBytecode := common.Hex2Bytes("602a60005260206000f3")
	result, err = public.Call(ctx, publicClient, public.CallParameters{
		Code: simpleBytecode,
		Data: []byte{}, // No calldata needed for this simple contract
	})
	if err != nil {
		fmt.Printf("Deployless call error: %v\n", err)
	} else {
		if len(result.Data) > 0 {
			returnValue := new(big.Int).SetBytes(result.Data)
			fmt.Printf("Deployless call returned: %d\n", returnValue.Uint64())
		}
	}

	// Example 13: Combined State + Block Override
	printSection("14. Combined State and Block Overrides")
	combinedStateOverride := types.StateOverride{
		testAddress: types.StateOverrideAccount{
			Balance: mustParseEther("10000"), // 10000 ETH
			Nonce:   ptrUint64(100),
		},
	}
	combinedBlockOverride := &types.BlockOverrides{
		GasLimit: ptrUint64(100000000),
	}

	// client.EncodeFunctionData()
	// Encode balanceOf(testAddress) - selector + padded address
	combinedBalanceOfData, err := client.EncodeFunctionData(client.EncodeFunctionDataOptions{
		ABI:          erc20.ContractABI,
		FunctionName: "balanceOf",
		Args:         []any{testAddress},
	})
	result, err = public.Call(ctx, publicClient, public.CallParameters{
		Account:        &testAddress,
		To:             &usdcAddress,
		Data:           combinedBalanceOfData,
		StateOverride:  combinedStateOverride,
		BlockOverrides: combinedBlockOverride,
	})
	if err != nil {
		fmt.Printf("Combined override error: %v\n", err)
	} else {
		fmt.Printf("Combined state + block override successful!\n")
	}

	// Example 14: Error Handling - Invalid Parameters
	printSection("15. Error Handling Examples")

	// Test: Code + To (mutually exclusive)
	fmt.Println("\nTest: Code + To (should fail)...")
	_, err = public.Call(ctx, publicClient, public.CallParameters{
		Code: simpleBytecode,
		To:   &usdcAddress,
		Data: nameSelector,
	})
	if err != nil {
		if _, ok := err.(*public.InvalidCallParamsError); ok {
			fmt.Printf("  Correctly rejected: %v\n", err)
		} else {
			fmt.Printf("  Error (unexpected type): %v\n", err)
		}
	}

	// Test: Code + Factory (mutually exclusive)
	fmt.Println("\nTest: Code + Factory (should fail)...")
	factoryAddr := common.HexToAddress("0xfactory0000000000000000000000000000000")
	_, err = public.Call(ctx, publicClient, public.CallParameters{
		Code:        simpleBytecode,
		Factory:     &factoryAddr,
		FactoryData: []byte{0x12, 0x34},
		Data:        nameSelector,
	})
	if err != nil {
		if _, ok := err.(*public.InvalidCallParamsError); ok {
			fmt.Printf("  Correctly rejected: %v\n", err)
		} else {
			fmt.Printf("  Error (unexpected type): %v\n", err)
		}
	}

	// Summary
	printHeader("Examples Complete")
	fmt.Println("Demonstrated Call features:")
	fmt.Println("  - Basic contract calls")
	fmt.Println("  - Calls with parameters (account, gas, value)")
	fmt.Println("  - Block number and block tag queries")
	fmt.Println("  - EIP-1559 fee parameters")
	fmt.Println("  - State overrides (modify account state)")
	fmt.Println("  - Block overrides (modify block context)")
	fmt.Println("  - Access lists (EIP-2930)")
	fmt.Println("  - Deployless calls (execute bytecode)")
	fmt.Println("  - Error handling for invalid parameters")
	fmt.Println()
}

// Helper functions

func mustParseGwei(s string) *big.Int {
	v, err := unit.ParseGwei(s)
	if err != nil {
		panic(fmt.Sprintf("invalid gwei value %q: %v", s, err))
	}
	return v
}

func mustParseEther(s string) *big.Int {
	v, err := unit.ParseEther(s)
	if err != nil {
		panic(fmt.Sprintf("invalid ether value %q: %v", s, err))
	}
	return v
}

func printHeader(title string) {
	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("  %s\n", title)
	fmt.Println(strings.Repeat("=", 60))
}

func printSection(title string) {
	fmt.Printf("\n--- %s ---\n", title)
}

func truncateBytes(data []byte, maxLen int) []byte {
	if len(data) <= maxLen {
		return data
	}
	return data[:maxLen]
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

func ptrUint64(v uint64) *uint64 {
	return &v
}
