package main

import (
	"fmt"
	"math/big"

	"github.com/ChefBingbong/viem-go/actions/public"
	"github.com/ChefBingbong/viem-go/actions/wallet"
	viemclient "github.com/ChefBingbong/viem-go/client"
)

// ERC-20 transfer ABI
const erc20TransferABI = `[
	{"name":"transfer","type":"function","stateMutability":"nonpayable",
	 "inputs":[{"name":"to","type":"address"},{"name":"amount","type":"uint256"}],
	 "outputs":[{"name":"","type":"bool"}]}
]`

// RunWriteContract demonstrates the simulate-then-write pattern via walletClient.WriteContract.
func RunWriteContract() {
	PrintSection("writeContract - Simulate then call contract write")

	amount := new(big.Int).Mul(big.NewInt(10), new(big.Int).Exp(big.NewInt(10), big.NewInt(6), nil)) // 10 USDC
	to := USDCAddr

	// Encode the calldata for simulation
	data, err := viemclient.EncodeFunctionData(viemclient.EncodeFunctionDataOptions{
		ABI:          erc20TransferABI,
		FunctionName: "transfer",
		Args:         []any{TargetAddr, amount},
	})
	if err != nil {
		fmt.Printf("Encode error: %v\n", err)
		return
	}
	fmt.Printf("Encoded transfer(%s, 10 USDC)\n", TargetAddr.Hex())
	fmt.Printf("Calldata: 0x%s...\n", Truncate(fmt.Sprintf("%x", data), 40))

	// Step 1: Simulate the write via eth_call
	fmt.Println("\nSimulating contract write...")
	result, simErr := public.Call(ctx, PublicClient, public.CallParameters{
		Account: &SourceAddr,
		To:      &to,
		Data:    data,
	})
	if simErr != nil {
		fmt.Printf("Simulation note: %s\n", Truncate(simErr.Error(), 120))
		fmt.Println("(Expected — account has no USDC on Polygon)")
	} else {
		fmt.Printf("Simulation result: 0x%x\n", result.Data)
	}

	// Step 2: writeContract via walletClient (delegates to wallet.WriteContract)
	fmt.Println("\nSending writeContract via walletClient...")
	hash, writeErr := WalletCl.WriteContract(ctx, wallet.WriteContractParameters{
		Address:      USDCAddr.Hex(),
		ABI:          erc20TransferABI,
		FunctionName: "transfer",
		Args:         []any{TargetAddr, amount},
	})
	if writeErr != nil {
		fmt.Printf("writeContract error: %s\n", Truncate(writeErr.Error(), 120))
		fmt.Println("(Expected on mainnet with unfunded account)")
	} else {
		fmt.Printf("Tx hash: %s\n", hash)
	}
}
