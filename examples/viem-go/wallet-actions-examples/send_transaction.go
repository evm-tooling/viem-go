package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/actions/public"
	"github.com/ChefBingbong/viem-go/actions/wallet"
	"github.com/ChefBingbong/viem-go/chain/definitions"
	"github.com/ChefBingbong/viem-go/utils/unit"
)

// RunSendTransaction demonstrates the simulate-then-send pattern via the WalletClient.
func RunSendTransaction() {
	PrintSection("sendTransaction - Simulate then send ETH transfer")

	value := MustParseEther("0.0001")
	to := TargetAddr

	// Step 1: Simulate via eth_call
	fmt.Printf("Simulating: %s -> %s (%s POL)\n", SourceAddr.Hex(), to.Hex(), unit.FormatEther(value))
	_, simErr := public.Call(ctx, PublicClient, public.CallParameters{
		Account: &SourceAddr,
		To:      &to,
		Value:   value,
	})
	if simErr != nil {
		fmt.Printf("Simulation note: %s\n", Truncate(simErr.Error(), 120))
		fmt.Println("(Expected on mainnet with unfunded account)")
	} else {
		fmt.Println("Simulation passed.")
	}

	// Step 2: Send using walletClient.SendTransaction (delegates to wallet.SendTransaction)
	fmt.Println("\nSending transaction via walletClient.SendTransaction...")
	hash, err := WalletCl.SendTransaction(ctx, wallet.SendTransactionParameters{
		Account: WalletCl.Account(),
		Chain:   &definitions.Polygon,
		To:      to.Hex(),
		Value:   value,
	})

	if err != nil {
		fmt.Printf("Send error: %s\n", Truncate(err.Error(), 120))
		fmt.Println("(Expected — public RPCs don't support eth_sendTransaction/wallet_sendTransaction;")
		fmt.Println(" they don't manage private keys. Use a local account or wallet-connected RPC.)")
	} else {
		fmt.Printf("Tx hash: %s\n", hash)
	}
	receipt, err := PublicClient.WaitForTransactionReceipt(ctx, common.HexToHash(hash))
	if err != nil {
		fmt.Printf("Wait for Receipt error: %s\n", Truncate(err.Error(), 120))

	}
	fmt.Println("Tx hash: \n", receipt)

}
