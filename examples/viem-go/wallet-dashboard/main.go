// Wallet Dashboard - viem-go Example
//
// A simple CLI dashboard demonstrating key viem-go features:
// - Public client creation
// - Fetching network info (block number, chain ID)
// - Gas price estimation
// - Address balance lookup
// - Message signing and verification
package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/accounts"
	"github.com/ChefBingbong/viem-go/chain/definitions"
	"github.com/ChefBingbong/viem-go/client"
	"github.com/ChefBingbong/viem-go/client/transport"
	"github.com/ChefBingbong/viem-go/utils/signature"
	"github.com/ChefBingbong/viem-go/utils/unit"
)

// Example address to check balance (Vitalik's address)
var exampleAddress = common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")

// Example private key for signing demo (DO NOT use real keys!)
const demoPrivateKey = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

func formatDuration(d time.Duration) string {
	if d < time.Microsecond {
		return fmt.Sprintf("%dns", d.Nanoseconds())
	}
	if d < time.Millisecond {
		return fmt.Sprintf("%dµs", d.Microseconds())
	}
	if d < time.Second {
		return fmt.Sprintf("%.2fms", float64(d.Microseconds())/1000)
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

func main() {
	totalStart := time.Now()
	ctx := context.Background()

	fmt.Println("╔══════════════════════════════════════════════════╗")
	fmt.Println("║          Wallet Dashboard (viem-go)              ║")
	fmt.Println("║        Ethereum Network Information              ║")
	fmt.Println("╚══════════════════════════════════════════════════╝")

	// 1. Create Public Client
	printSection("Creating Public Client")
	publicClient, err := client.CreatePublicClient(client.PublicClientConfig{
		Chain:     &definitions.Polygon,
		Transport: transport.HTTP("https://polygon-rpc.com"),
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}
	defer publicClient.Close()

	fmt.Printf("Connected to: %s (Chain ID: %d)\n", definitions.Mainnet.Name, definitions.Mainnet.ID)

	// 2. Fetch Network Information
	printSection("Network Information")

	blockNumber, err := publicClient.GetBlockNumber(ctx)
	if err != nil {
		fmt.Printf("Error getting block number: %v\n", err)
		return
	}
	fmt.Printf("Current Block Number: %d\n", blockNumber)

	chainID, err := publicClient.GetChainID(ctx)
	if err != nil {
		fmt.Printf("Error getting chain ID: %v\n", err)
		return
	}
	fmt.Printf("Chain ID: %d\n", chainID)

	// 3. Get Gas Prices
	printSection("Gas Prices")

	gasPrice, err := publicClient.GetGasPrice(ctx)
	if err != nil {
		fmt.Printf("Error getting gas price: %v\n", err)
		return
	}
	fmt.Printf("Current Gas Price: %s Gwei\n", unit.FormatGwei(gasPrice))

	maxPriorityFee, err := publicClient.GetMaxPriorityFeePerGas(ctx)
	if err != nil {
		fmt.Println("Max Priority Fee: Not available on this network")
	} else {
		fmt.Printf("Max Priority Fee: %s Gwei\n", unit.FormatGwei(maxPriorityFee))
	}

	// 4. Check Address Balance
	printSection("Address Balance")

	balance, err := publicClient.GetBalance(ctx, exampleAddress)
	if err != nil {
		fmt.Printf("Error getting balance: %v\n", err)
		return
	}

	fmt.Printf("Address: %s\n", exampleAddress.Hex())
	fmt.Printf("Balance: %s ETH\n", unit.FormatEther(balance))

	// Get transaction count (nonce)
	txCount, err := publicClient.GetTransactionCount(ctx, exampleAddress)
	if err != nil {
		fmt.Printf("Error getting transaction count: %v\n", err)
		return
	}
	fmt.Printf("Transaction Count: %d\n", txCount)

	// 5. Message Signing Demo
	printSection("Message Signing Demo")

	account, err := accounts.PrivateKeyToAccount(demoPrivateKey)
	if err != nil {
		fmt.Printf("Error creating account: %v\n", err)
		return
	}
	fmt.Printf("Demo Account Address: %s\n", account.GetAddress())

	message := "Hello from Wallet Dashboard!"
	fmt.Printf("Message: \"%s\"\n", message)

	sig, err := account.SignMessage(signature.NewSignableMessage(message))
	if err != nil {
		fmt.Printf("Error signing message: %v\n", err)
		return
	}

	// Truncate signature for display
	sigDisplay := sig
	if len(sig) > 42 {
		sigDisplay = sig[:42] + "..."
	}
	fmt.Printf("Signature: %s\n", sigDisplay)

	// Verify the signature
	isValid, err := signature.VerifyMessage(account.GetAddress(), signature.NewSignableMessage(message), sig)
	if err != nil {
		fmt.Printf("Error verifying signature: %v\n", err)
	} else {
		validStr := "No"
		if isValid {
			validStr = "Yes"
		}
		fmt.Printf("Signature Valid: %s\n", validStr)
	}

	// 6. Get Latest Block Info
	printSection("Latest Block Info")

	block, err := publicClient.GetBlock(ctx, client.BlockTagLatest, false)
	if err != nil {
		fmt.Printf("Error getting block: %v\n", err)
		return
	}

	if block != nil {
		hashDisplay := block.Hash.Hex()
		if len(hashDisplay) > 42 {
			hashDisplay = hashDisplay[:42] + "..."
		}
		fmt.Printf("Block Hash: %s\n", hashDisplay)

		timestamp := time.Unix(int64(block.Timestamp), 0).UTC()
		fmt.Printf("Timestamp: %s\n", timestamp.Format(time.RFC3339))
		fmt.Printf("Transactions: %d\n", len(block.Transactions))
		fmt.Printf("Gas Used: %d\n", block.GasUsed)
		fmt.Printf("Gas Limit: %d\n", block.GasLimit)
	}

	// Summary
	totalElapsed := time.Since(totalStart)
	printHeader("Dashboard Summary")
	fmt.Printf("  Network: %s\n", definitions.Mainnet.Name)
	fmt.Printf("  Block: #%d\n", blockNumber)
	fmt.Printf("  Gas Price: %s Gwei\n", unit.FormatGwei(gasPrice))

	accountDisplay := account.GetAddress()
	if len(accountDisplay) > 10 {
		accountDisplay = accountDisplay[:10] + "..."
	}
	fmt.Printf("  Demo Account: %s\n", accountDisplay)
	fmt.Printf("  Total Runtime: %s\n", formatDuration(totalElapsed))
	fmt.Println()
}

func printHeader(title string) {
	fmt.Println()
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("  %s\n", title)
	fmt.Println(strings.Repeat("=", 50))
}

func printSection(title string) {
	fmt.Printf("\n--- %s ---\n", title)
}
