// Wallet Actions Examples — Entry Point
//
// Usage:
//
//	go run ./examples/viem-go/wallet-actions-examples/                                # Run all
//	go run ./examples/viem-go/wallet-actions-examples/ signMessage                    # Run one
//	go run ./examples/viem-go/wallet-actions-examples/ signMessage signTransaction    # Run multiple
//
// Available examples:
//
//	signMessage, signTypedData, signTransaction, signAuthorization,
//	sendTransaction, sendRawTransaction, writeContract, deployContract,
//	sendCalls
package main

import (
	"context"
	"fmt"
	"os"
	"strings"
)

var ctx = context.Background()

// registry maps CLI names to runner functions.
var registry = map[string]func(){
	"signMessage":        RunSignMessage,
	"signTypedData":      RunSignTypedData,
	"signTransaction":    RunSignTransaction,
	"signAuthorization":  RunSignAuthorization,
	"sendTransaction":    RunSendTransaction,
	"sendRawTransaction": RunSendRawTransaction,
	"writeContract":      RunWriteContract,
	"deployContract":     RunDeployContract,
	"sendCalls":          RunSendCalls,
}

// ordered preserves the display / execution order.
var ordered = []string{
	"signMessage",
	"signTypedData",
	"signTransaction",
	"signAuthorization",
	"sendTransaction",
	"sendRawTransaction",
	"writeContract",
	"deployContract",
	"sendCalls",
}

func main() {
	args := os.Args[1:]

	// Validate args
	var selected []string
	if len(args) == 0 {
		selected = ordered
	} else {
		for _, a := range args {
			if _, ok := registry[a]; !ok {
				fmt.Fprintf(os.Stderr, "Unknown example: %s\n", a)
				fmt.Fprintf(os.Stderr, "Available: %s\n", strings.Join(ordered, ", "))
				os.Exit(1)
			}
		}
		selected = args
	}

	PrintHeader("Wallet Actions Examples (viem-go) — Polygon Mainnet")

	// Initialise clients
	if err := Setup(); err != nil {
		fmt.Fprintf(os.Stderr, "Setup error: %v\n", err)
		os.Exit(1)
	}
	defer Teardown()

	PrintAccountInfo()
	fmt.Printf("\nRunning %d example(s): %s\n", len(selected), strings.Join(selected, ", "))

	for _, name := range selected {
		fn := registry[name]
		fn()
	}

	PrintHeader("Done")
}
