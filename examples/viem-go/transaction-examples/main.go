// Transaction Examples - viem-go
//
// Comprehensive examples demonstrating transaction-related public actions:
// - GetBlockNumber
// - GetTransaction
// - GetTransactionReceipt
// - GetTransactionConfirmations
// - WaitForTransactionReceipt
// - FillTransaction
package main

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/actions/public"
	"github.com/ChefBingbong/viem-go/chain/definitions"
	"github.com/ChefBingbong/viem-go/client"
	"github.com/ChefBingbong/viem-go/client/transport"
	"github.com/ChefBingbong/viem-go/utils/unit"
)

// Example addresses
var (
	// Vitalik's address
	vitalikAddress = common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
	// Random test address
	testAddress = common.HexToAddress("0x1234567890123456789012345678901234567890")
)

func main() {
	ctx := context.Background()

	printHeader("Transaction Action Examples (viem-go)")

	// Create Public Client
	printSection("1. Creating Public Client")
	publicClient, err := client.CreatePublicClient(client.PublicClientConfig{
		Chain:     &definitions.Mainnet,
		Transport: transport.HTTP("https://rough-purple-market.quiknode.pro/c1a568726a34041d3c5d58603f5981951e6a8503"),
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}
	defer publicClient.Close()
	fmt.Println("Connected to Ethereum Mainnet")

	// =========================================================================
	// GetBlockNumber Examples
	// =========================================================================

	printSection("2. GetBlockNumber - Get Current Block Number")
	blockNumber, err := public.GetBlockNumber(ctx, publicClient, public.GetBlockNumberParameters{})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Current Block Number: %d\n", blockNumber)
	}

	// Example with custom cache time
	printSection("3. GetBlockNumber - With Custom Cache Time")
	cacheTime := 10 * time.Second
	blockNumber, err = public.GetBlockNumber(ctx, publicClient, public.GetBlockNumberParameters{
		CacheTime: &cacheTime,
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Block Number (cached for 10s): %d\n", blockNumber)
	}

	// =========================================================================
	// GetTransaction Examples
	// =========================================================================

	printSection("4. GetTransaction - Get Transaction by Hash")
	// Use a known transaction hash (this is a historical ETH transfer)
	txHash := common.HexToHash("0x88df016429689c079f3b2f6ad39fa052532c56795b733da78a91ebe6a713944b")
	tx, err := public.GetTransaction(ctx, publicClient, public.GetTransactionParameters{
		Hash: &txHash,
	})
	if err != nil {
		fmt.Printf("Error getting transaction: %v\n", err)
	} else {
		fmt.Printf("Transaction Hash: %s\n", tx.Hash.Hex())
		fmt.Printf("  From: %s\n", truncateAddress(tx.From))
		if tx.To != nil {
			fmt.Printf("  To: %s\n", truncateAddress(*tx.To))
		}
		if tx.Value != nil {
			fmt.Printf("  Value: %s ETH\n", unit.FormatEther(tx.Value))
		}
		if tx.BlockNumber != nil {
			fmt.Printf("  Block Number: %d\n", *tx.BlockNumber)
		}
		fmt.Printf("  Nonce: %d\n", tx.Nonce)
		fmt.Printf("  Gas: %d\n", tx.Gas)
	}

	printSection("5. GetTransaction - By Block Number and Index")
	blockNum := uint64(17000000) // A specific block
	txIndex := 0
	tx, err = public.GetTransaction(ctx, publicClient, public.GetTransactionParameters{
		BlockNumber: &blockNum,
		Index:       &txIndex,
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("First transaction in block %d:\n", blockNum)
		fmt.Printf("  Hash: %s\n", tx.Hash.Hex()[:18]+"...")
		fmt.Printf("  From: %s\n", truncateAddress(tx.From))
	}

	// =========================================================================
	// GetTransactionReceipt Examples
	// =========================================================================

	printSection("6. GetTransactionReceipt - Get Receipt by Hash")
	receipt, err := public.GetTransactionReceipt(ctx, publicClient, public.GetTransactionReceiptParameters{
		Hash: txHash,
	})
	if err != nil {
		fmt.Printf("Error getting receipt: %v\n", err)
	} else {
		fmt.Printf("Transaction Receipt:\n")
		fmt.Printf("  Hash: %s\n", receipt.TransactionHash.Hex()[:18]+"...")
		fmt.Printf("  Block Number: %d\n", receipt.BlockNumber)
		fmt.Printf("  Status: %s\n", formatStatus(receipt.Status))
		fmt.Printf("  Gas Used: %d\n", receipt.GasUsed)
		fmt.Printf("  Cumulative Gas Used: %d\n", receipt.CumulativeGasUsed)
		if receipt.ContractAddress != nil {
			fmt.Printf("  Contract Created: %s\n", receipt.ContractAddress.Hex())
		}
		fmt.Printf("  Logs Count: %d\n", len(receipt.Logs))

		// Helper methods
		if receipt.IsSuccess() {
			fmt.Printf("  Transaction SUCCEEDED\n")
		} else {
			fmt.Printf("  Transaction FAILED\n")
		}
	}

	// =========================================================================
	// GetTransactionConfirmations Examples
	// =========================================================================

	printSection("7. GetTransactionConfirmations - By Hash")
	confirmations, err := public.GetTransactionConfirmations(ctx, publicClient, public.GetTransactionConfirmationsParameters{
		Hash: &txHash,
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Transaction %s...\n", txHash.Hex()[:18])
		fmt.Printf("  Confirmations: %d blocks\n", confirmations)
	}

	printSection("8. GetTransactionConfirmations - Using Existing Receipt")
	if receipt != nil {
		confirmations, err = public.GetTransactionConfirmations(ctx, publicClient, public.GetTransactionConfirmationsParameters{
			TransactionReceipt: receipt,
		})
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("Confirmations (from receipt): %d blocks\n", confirmations)
		}
	}

	// =========================================================================
	// WaitForTransactionReceipt Examples
	// =========================================================================

	printSection("9. WaitForTransactionReceipt - Wait for Confirmation")
	fmt.Println("Note: Using an already-mined transaction for demo purposes")

	// For demo, we'll wait for an already-mined transaction
	receipt, err = public.WaitForTransactionReceipt(ctx, publicClient, public.WaitForTransactionReceiptParameters{
		Hash:            txHash,
		Confirmations:   1,
		PollingInterval: 1 * time.Second,
		Timeout:         10 * time.Second,
	})
	if err != nil {
		fmt.Printf("Error waiting for receipt: %v\n", err)
	} else {
		fmt.Printf("Transaction confirmed!\n")
		fmt.Printf("  Block: %d\n", receipt.BlockNumber)
		fmt.Printf("  Status: %s\n", formatStatus(receipt.Status))
	}

	printSection("10. WaitForTransactionReceipt - With Multiple Confirmations")
	receipt, err = public.WaitForTransactionReceipt(ctx, publicClient, public.WaitForTransactionReceiptParameters{
		Hash:            txHash,
		Confirmations:   12, // Wait for 12 confirmations (finality)
		PollingInterval: 2 * time.Second,
		Timeout:         30 * time.Second,
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Transaction has 12+ confirmations!\n")
		fmt.Printf("  Block: %d\n", receipt.BlockNumber)
	}

	printSection("11. WaitForTransactionReceipt - With OnReplaced Callback")
	fmt.Println("Note: Replacement detection works for pending transactions")
	fmt.Println("Demo: Setting up callback for replacement detection...")

	// Demonstrate the callback setup (won't trigger for already-mined tx)
	_, _ = public.WaitForTransactionReceipt(ctx, publicClient, public.WaitForTransactionReceiptParameters{
		Hash:             txHash,
		CheckReplacement: ptrBool(true),
		OnReplaced: func(info public.ReplacementInfo) {
			fmt.Printf("Transaction was replaced!\n")
			fmt.Printf("  Reason: %s\n", info.Reason)
			fmt.Printf("  Original Hash: %s\n", info.ReplacedTransaction.Hash.Hex())
			fmt.Printf("  New Hash: %s\n", info.Transaction.Hash.Hex())
		},
		Timeout: 5 * time.Second,
	})
	fmt.Println("Callback registered (would fire if tx was replaced)")

	// =========================================================================
	// FillTransaction Examples
	// =========================================================================

	printSection("12. FillTransaction - Auto-fill Transaction Fields")
	fmt.Println("Note: eth_fillTransaction may not be supported by all nodes")

	from := vitalikAddress
	to := testAddress
	value := mustParseEther("0.1")

	result, err := public.FillTransaction(ctx, publicClient, public.FillTransactionParameters{
		Account: &from,
		To:      &to,
		Value:   value,
	})
	if err != nil {
		fmt.Printf("FillTransaction error (this is expected on many public nodes): %v\n", err)
		fmt.Println("Note: This method is primarily supported by Geth-based nodes")
	} else {
		fmt.Printf("Filled Transaction:\n")
		fmt.Printf("  From: %s\n", truncateAddress(result.Transaction.From))
		if result.Transaction.To != nil {
			fmt.Printf("  To: %s\n", truncateAddress(*result.Transaction.To))
		}
		fmt.Printf("  Nonce: %d (auto-filled)\n", result.Transaction.Nonce)
		fmt.Printf("  Gas: %d (estimated)\n", result.Transaction.Gas)
		if result.Transaction.MaxFeePerGas != nil {
			fmt.Printf("  Max Fee Per Gas: %s gwei\n", unit.FormatGwei(result.Transaction.MaxFeePerGas))
		}
		if result.Transaction.MaxPriorityFeePerGas != nil {
			fmt.Printf("  Max Priority Fee: %s gwei\n", unit.FormatGwei(result.Transaction.MaxPriorityFeePerGas))
		}
		fmt.Printf("  Raw TX Length: %d bytes\n", len(result.Raw))
	}

	printSection("13. FillTransaction - With Partial Values")
	gas := uint64(50000)
	maxFeePerGas := mustParseGwei("30")

	result, err = public.FillTransaction(ctx, publicClient, public.FillTransactionParameters{
		Account:      &from,
		To:           &to,
		Value:        value,
		Gas:          &gas,         // Provide gas, other fields will be filled
		MaxFeePerGas: maxFeePerGas, // Provide maxFeePerGas, priority fee will be filled
	})
	if err != nil {
		fmt.Printf("FillTransaction error: %v\n", err)
	} else {
		fmt.Printf("Transaction with partial values:\n")
		fmt.Printf("  Gas: %d (supplied)\n", result.Transaction.Gas)
		fmt.Printf("  Max Fee Per Gas: %s gwei (supplied)\n", unit.FormatGwei(result.Transaction.MaxFeePerGas))
		fmt.Printf("  Nonce: %d (auto-filled)\n", result.Transaction.Nonce)
	}

	printSection("14. FillTransaction - Custom Base Fee Multiplier")
	multiplier := 1.5 // 50% buffer instead of default 20%

	result, err = public.FillTransaction(ctx, publicClient, public.FillTransactionParameters{
		Account:           &from,
		To:                &to,
		Value:             mustParseEther("0.01"),
		BaseFeeMultiplier: &multiplier,
	})
	if err != nil {
		fmt.Printf("FillTransaction error: %v\n", err)
	} else {
		fmt.Printf("Transaction with 1.5x base fee multiplier:\n")
		if result.Transaction.MaxFeePerGas != nil {
			fmt.Printf("  Max Fee Per Gas: %s gwei (with 50%% buffer)\n", unit.FormatGwei(result.Transaction.MaxFeePerGas))
		}
	}

	// Summary
	printHeader("Examples Complete")
	fmt.Println("Demonstrated transaction actions:")
	fmt.Println("  - GetBlockNumber: Get current block with caching")
	fmt.Println("  - GetTransaction: By hash, block number, or index")
	fmt.Println("  - GetTransactionReceipt: Full receipt with logs")
	fmt.Println("  - GetTransactionConfirmations: By hash or receipt")
	fmt.Println("  - WaitForTransactionReceipt: Polling with replacement detection")
	fmt.Println("  - FillTransaction: Auto-fill nonce, gas, and fees")
	fmt.Println()
}

// Helper functions

func mustParseEther(s string) *big.Int {
	v, err := unit.ParseEther(s)
	if err != nil {
		panic(fmt.Sprintf("invalid ether value %q: %v", s, err))
	}
	return v
}

func mustParseGwei(s string) *big.Int {
	v, err := unit.ParseGwei(s)
	if err != nil {
		panic(fmt.Sprintf("invalid gwei value %q: %v", s, err))
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

func truncateAddress(addr common.Address) string {
	hex := addr.Hex()
	return hex[:10] + "..." + hex[len(hex)-4:]
}

func formatStatus(status uint64) string {
	if status == 1 {
		return "Success (1)"
	}
	return "Failed (0)"
}

func ptrBool(v bool) *bool {
	return &v
}
