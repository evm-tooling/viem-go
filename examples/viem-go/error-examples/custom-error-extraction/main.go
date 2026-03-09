// Custom Error Extraction (viem-go)
//
// Showcases the errors.As pattern for extracting ContractFunctionRevertedError
// and accessing ErrorName / Data for custom error handling.
//
// Parity with viem TypeScript:
//
//	viem:    err.walk(e => e instanceof ContractFunctionRevertedError)
//	viem-go: errors.As(err, &revertErr)
//
// Both expose: errorName, reason, args/data for branching on custom errors.
//
// Run: go run ./custom-error-extraction
package main

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/abi"
	"github.com/ChefBingbong/viem-go/actions/public"
	"github.com/ChefBingbong/viem-go/chain/definitions"
	"github.com/ChefBingbong/viem-go/client"
	"github.com/ChefBingbong/viem-go/client/transport"
	"github.com/ChefBingbong/viem-go/contracts/erc20"
	pkgerrors "github.com/ChefBingbong/viem-go/errors"
)

var (
	usdcAddress = common.HexToAddress("0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359")
	testAddress = common.HexToAddress("0x1234567890123456789012345678901234567890")
	zeroAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")
)

func main() {
	ctx := context.Background()

	erc20ABI, err := abi.Parse([]byte(erc20.ContractABI))
	if err != nil {
		fmt.Printf("Failed to parse ERC20 ABI: %v\n", err)
		return
	}

	fmt.Println("============================================================")
	fmt.Println("  Custom Error Extraction (viem-go)")
	fmt.Println("  Pattern: errors.As(err, &revertErr)")
	fmt.Println("============================================================\n")

	publicClient, err := client.CreatePublicClient(client.PublicClientConfig{
		Chain:     &definitions.Polygon,
		Transport: transport.HTTP("https://rpc-mainnet.matic.quiknode.pro"),
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}
	defer publicClient.Close()

	// Simulate a reverting transfer (from zero address -> Error(string))
	_, err = public.SimulateContract(ctx, publicClient, public.SimulateContractParameters{
		Account:      &zeroAddress,
		Address:      usdcAddress,
		ABI:          erc20ABI,
		FunctionName: "transfer",
		Args:         []any{testAddress, new(big.Int).Exp(big.NewInt(10), big.NewInt(60), nil)},
	})
	if err != nil {
		var revertErr *pkgerrors.ContractFunctionRevertedError
		if errors.As(err, &revertErr) {
			errorName := revertErr.ErrorName
			reason := revertErr.Reason
			data := revertErr.Data

			fmt.Println("✓ Extracted ContractFunctionRevertedError via errors.As()")
			fmt.Printf("  ErrorName: %q\n", errorName)
			fmt.Printf("  Reason:   %q\n", reason)
			if data != nil {
				fmt.Printf("  Data:     %v\n", data)
			}
			fmt.Println()

			// Example: branch on error type (parity with viem switch on errorName)
			switch errorName {
			case "Error":
				fmt.Println("  → Standard Error(string) - handle generic revert message")
			case "Panic":
				fmt.Println("  → Standard Panic(uint256) - handle assertion failure")
			case "ERC20InsufficientBalance", "ERC20InsufficientAllowance":
				fmt.Println("  → Custom ERC20 error - handle insufficient balance/allowance")
			default:
				if errorName != "" {
					fmt.Printf("  → Custom error %q - handle as needed\n", errorName)
				} else {
					fmt.Println("  → Unknown/undecoded error")
				}
			}
		} else {
			// ContractFunctionRevertedError is wrapped inside ContractFunctionExecutionError.Cause
			var execErr *pkgerrors.ContractFunctionExecutionError
			if errors.As(err, &execErr) && execErr.Cause != nil {
				if errors.As(execErr.Cause, &revertErr) {
					errorName := revertErr.ErrorName
					fmt.Println("✓ Extracted via execErr.Cause (nested in ContractFunctionExecutionError)")
					fmt.Printf("  ErrorName: %q\n", errorName)
					fmt.Printf("  Reason:   %q\n", revertErr.Reason)
				}
			}
		}
		if err != nil {
			fmt.Printf("\n  Full error: %v\n", err)
		}
	}

	fmt.Println("\n--- Summary ---")
	fmt.Println("  viem-go supports: errors.As(err, &revertErr) to find ContractFunctionRevertedError,")
	fmt.Println("  then access revertErr.ErrorName, revertErr.Reason, revertErr.Data")
	fmt.Println("  (ErrorName added for parity with viem's data.errorName)")
	fmt.Println()
}
