// Package main demonstrates the viem-go watch actions.
// These examples show how to use WatchBlockNumber, WatchBlocks, WatchPendingTransactions,
// WatchEvent, and WatchContractEvent with both polling and subscription modes.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"

	viemabi "github.com/ChefBingbong/viem-go/abi"
	"github.com/ChefBingbong/viem-go/actions/public"
	"github.com/ChefBingbong/viem-go/chain/definitions"
	"github.com/ChefBingbong/viem-go/client"
	"github.com/ChefBingbong/viem-go/client/transport"
)

// Example ERC20 Transfer event ABI
var erc20TransferEventABI = `[{
	"anonymous": false,
	"inputs": [
		{"indexed": true, "name": "from", "type": "address"},
		{"indexed": true, "name": "to", "type": "address"},
		{"indexed": false, "name": "value", "type": "uint256"}
	],
	"name": "Transfer",
	"type": "event"
}]`

func main() {
	// Parse command line arguments for the example to run
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <example>")
		fmt.Println("Examples:")
		fmt.Println("  watch-block-number    - Watch for new block numbers")
		fmt.Println("  watch-blocks          - Watch for new blocks with full data")
		fmt.Println("  watch-pending-tx      - Watch for pending transactions")
		fmt.Println("  watch-event           - Watch for generic events")
		fmt.Println("  watch-contract-event  - Watch for ERC20 Transfer events")
		fmt.Println("  all                   - Run all examples")
		os.Exit(1)
	}

	example := os.Args[1]

	// Create a public client
	// Using Ankr public RPC for demonstration
	rpcURL := "https://rough-purple-market.matic.quiknode.pro/c1a568726a34041d3c5d58603f5981951e6a8503"
	if envURL := os.Getenv("RPC_URL"); envURL != "" {
		rpcURL = envURL
	}

	publicClient, err := client.CreatePublicClient(client.PublicClientConfig{
		Chain:           &definitions.Polygon,
		Transport:       transport.HTTP(rpcURL),
		PollingInterval: 1 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer publicClient.Close()

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())

	// Handle graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\nShutting down...")
		cancel()
	}()

	switch example {
	case "watch-block-number":
		watchBlockNumberExample(ctx, publicClient)
	case "watch-blocks":
		watchBlocksExample(ctx, publicClient)
	case "watch-pending-tx":
		watchPendingTransactionsExample(ctx, publicClient)
	case "watch-event":
		watchEventExample(ctx, publicClient)
	case "watch-contract-event":
		watchContractEventExample(ctx, publicClient)
	case "all":
		fmt.Println("Running all examples in sequence...")
		fmt.Println("\n=== Watch Block Number ===")
		runWithTimeout(ctx, 15*time.Second, func(ctx context.Context) {
			watchBlockNumberExample(ctx, publicClient)
		})
		fmt.Println("\n=== Watch Blocks ===")
		runWithTimeout(ctx, 15*time.Second, func(ctx context.Context) {
			watchBlocksExample(ctx, publicClient)
		})
		fmt.Println("\n=== Watch Pending Transactions ===")
		runWithTimeout(ctx, 10*time.Second, func(ctx context.Context) {
			watchPendingTransactionsExample(ctx, publicClient)
		})
		fmt.Println("\n=== Watch Event ===")
		runWithTimeout(ctx, 15*time.Second, func(ctx context.Context) {
			watchEventExample(ctx, publicClient)
		})
		fmt.Println("\n=== Watch Contract Event ===")
		runWithTimeout(ctx, 15*time.Second, func(ctx context.Context) {
			watchContractEventExample(ctx, publicClient)
		})
		fmt.Println("\nAll examples completed!")
	default:
		fmt.Printf("Unknown example: %s\n", example)
		os.Exit(1)
	}
}

// runWithTimeout runs a function with a timeout
func runWithTimeout(parent context.Context, timeout time.Duration, fn func(ctx context.Context)) {
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()
	fn(ctx)
}

// watchBlockNumberExample demonstrates WatchBlockNumber
func watchBlockNumberExample(ctx context.Context, c *client.PublicClient) {
	fmt.Println("Watching for new block numbers...")
	fmt.Println("Press Ctrl+C to stop")

	// Watch block numbers with emit on begin and emit missed
	events := c.WatchBlockNumber(ctx, public.WatchBlockNumberParameters{
		EmitOnBegin: true,
		EmitMissed:  true,
	})

	count := 0
	for event := range events {
		if event.Error != nil {
			fmt.Printf("Error: %v\n", event.Error)
			continue
		}

		prev := "nil"
		if event.PrevBlockNumber != nil {
			prev = fmt.Sprintf("%d", *event.PrevBlockNumber)
		}

		fmt.Printf("Block: %d (prev: %s)\n", event.BlockNumber, prev)

		count++
		if count >= 5 {
			fmt.Println("Received 5 block numbers, stopping...")
			return
		}
	}
}

// watchBlocksExample demonstrates WatchBlocks
func watchBlocksExample(ctx context.Context, c *client.PublicClient) {
	fmt.Println("Watching for new blocks with full data...")
	fmt.Println("Press Ctrl+C to stop")

	// Watch blocks with transactions included
	events := c.WatchBlocks(ctx, public.WatchBlocksParameters{
		EmitOnBegin:         true,
		EmitMissed:          true,
		IncludeTransactions: false, // Set to true to get full transaction data
	})

	count := 0
	for event := range events {
		if event.Error != nil {
			fmt.Printf("Error: %v\n", event.Error)
			continue
		}

		block := event.Block
		fmt.Printf("Block %d:\n", block.Number)
		fmt.Printf("  Hash:         %s\n", block.Hash.Hex())
		fmt.Printf("  Timestamp:    %d\n", block.Timestamp)
		fmt.Printf("  Gas Used:     %d\n", block.GasUsed)
		fmt.Printf("  Transactions: %d\n", len(block.Transactions))

		count++
		if count >= 3 {
			fmt.Println("Received 3 blocks, stopping...")
			return
		}
	}
}

// watchPendingTransactionsExample demonstrates WatchPendingTransactions
func watchPendingTransactionsExample(ctx context.Context, c *client.PublicClient) {
	fmt.Println("Watching for pending transactions...")
	fmt.Println("Note: Some RPC providers don't support pending transaction filters")
	fmt.Println("Press Ctrl+C to stop")

	// Watch pending transactions with batching
	events := c.WatchPendingTransactions(ctx, public.WatchPendingTransactionsParameters{
		Batch: true,
	})

	totalTx := 0
	for event := range events {
		if event.Error != nil {
			fmt.Printf("Error: %v\n", event.Error)
			// Many providers don't support pending tx filters
			return
		}

		if len(event.Hashes) > 0 {
			fmt.Printf("Received %d pending transaction(s):\n", len(event.Hashes))
			for i, hash := range event.Hashes {
				if i < 5 { // Only show first 5
					fmt.Printf("  - %s\n", hash.Hex())
				}
			}
			if len(event.Hashes) > 5 {
				fmt.Printf("  ... and %d more\n", len(event.Hashes)-5)
			}

			totalTx += len(event.Hashes)
			if totalTx >= 20 {
				fmt.Printf("Received %d total pending transactions, stopping...\n", totalTx)
				return
			}
		}
	}
}

// watchEventExample demonstrates WatchEvent
func watchEventExample(ctx context.Context, c *client.PublicClient) {
	fmt.Println("Watching for Transfer events from USDC contract...")
	fmt.Println("Press Ctrl+C to stop")

	// USDC contract address on Ethereum mainnet
	usdcAddress := common.HexToAddress("0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359")

	// Transfer event signature
	transferTopic := common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")

	// Watch events with batching
	events := c.WatchEvent(ctx, public.WatchEventParameters{
		Address: usdcAddress,
		Event: &viemabi.Event{
			Name:  "Transfer",
			Topic: transferTopic,
		},
		Batch: true,
	})

	count := 0
	for event := range events {
		if event.Error != nil {
			fmt.Printf("Error: %v\n", event.Error)
			continue
		}

		if len(event.Logs) > 0 {
			fmt.Printf("Received %d Transfer event(s):\n", len(event.Logs))
			for i, log := range event.Logs {
				if i < 3 { // Only show first 3
					txHash := ""
					if log.TransactionHash != nil {
						txHash = *log.TransactionHash
					}
					fmt.Printf("  Block %d, Tx: %s\n", log.BlockNumber, txHash)
					if len(log.Topics) >= 3 {
						from := common.HexToAddress(log.Topics[1])
						to := common.HexToAddress(log.Topics[2])
						fmt.Printf("    From: %s\n", from.Hex())
						fmt.Printf("    To:   %s\n", to.Hex())
					}
				}
			}
			if len(event.Logs) > 3 {
				fmt.Printf("  ... and %d more\n", len(event.Logs)-3)
			}

			count += len(event.Logs)
			if count >= 100 {
				fmt.Printf("Received %d total events, stopping...\n", count)
				return
			}
		}
	}
}

// watchContractEventExample demonstrates WatchContractEvent with ABI decoding
func watchContractEventExample(ctx context.Context, c *client.PublicClient) {
	fmt.Println("Watching for Transfer events with ABI decoding...")
	fmt.Println("Press Ctrl+C to stop")

	// Parse the ERC20 ABI
	abi, err := viemabi.Parse([]byte(erc20TransferEventABI))
	if err != nil {
		fmt.Printf("Failed to parse ABI: %v\n", err)
		return
	}

	// USDC contract address on Ethereum mainnet
	usdcAddress := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")

	// Watch contract events with ABI decoding
	events := c.WatchContractEvent(ctx, public.WatchContractEventParameters{
		Address:   usdcAddress,
		ABI:       abi,
		EventName: "Transfer",
		Batch:     true,
	})

	count := 0
	for event := range events {
		if event.Error != nil {
			fmt.Printf("Error: %v\n", event.Error)
			continue
		}

		if len(event.Logs) > 0 {
			fmt.Printf("Received %d decoded Transfer event(s):\n", len(event.Logs))
			for i, log := range event.Logs {
				if i < 3 { // Only show first 3
					fmt.Printf("  Block %d:\n", log.BlockNumber)
					if log.EventName != "" {
						fmt.Printf("    Event: %s\n", log.EventName)
					}
					if log.Args != nil {
						if args, ok := log.Args.(map[string]any); ok {
							if from, ok := args["from"]; ok {
								fmt.Printf("    From: %v\n", from)
							}
							if to, ok := args["to"]; ok {
								fmt.Printf("    To: %v\n", to)
							}
							if value, ok := args["value"]; ok {
								fmt.Printf("    Value: %v\n", value)
							}
						}
					}
				}
			}
			if len(event.Logs) > 3 {
				fmt.Printf("  ... and %d more\n", len(event.Logs)-3)
			}

			count += len(event.Logs)
			if count >= 10 {
				fmt.Printf("Received %d total events, stopping...\n", count)
				return
			}
		}
	}
}
