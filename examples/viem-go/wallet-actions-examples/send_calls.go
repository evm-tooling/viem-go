package main

import (
	"fmt"

	"github.com/ChefBingbong/viem-go/actions/public"
	"github.com/ChefBingbong/viem-go/actions/wallet"
	"github.com/ChefBingbong/viem-go/utils/unit"
)

// RunSendCalls demonstrates EIP-5792 batch call sending via walletClient.SendCalls.
func RunSendCalls() {
	PrintSection("sendCalls - EIP-5792 batch calls")

	// Simulate each call individually
	values := []string{"0.001", "0.002", "0.003"}
	to := TargetAddr
	for i, v := range values {
		val := MustParseEther(v)
		fmt.Printf("Simulating call %d: %s -> %s (%s POL)\n", i+1, SourceAddr.Hex(), to.Hex(), unit.FormatEther(val))
		_, simErr := public.Call(ctx, PublicClient, public.CallParameters{
			Account: &SourceAddr,
			To:      &to,
			Value:   val,
		})
		if simErr != nil {
			fmt.Printf("  Simulation note: %s\n", Truncate(simErr.Error(), 100))
		} else {
			fmt.Printf("  Simulation passed.\n")
		}
	}

	// Send batch via walletClient.SendCalls (delegates to wallet.SendCalls)
	fmt.Println("\nSending batch via walletClient.SendCalls...")
	result, err := WalletCl.SendCalls(ctx, wallet.SendCallsParameters{
		Calls: []wallet.Call{
			{To: TargetAddr.Hex(), Value: MustParseEther("0.001")},
			{To: TargetAddr.Hex(), Value: MustParseEther("0.002")},
			{To: TargetAddr.Hex(), Value: MustParseEther("0.003")},
		},
	})
	if err != nil {
		fmt.Printf("sendCalls error: %s\n", Truncate(err.Error(), 120))
		fmt.Println("(Expected — most public RPCs do not support wallet_sendCalls)")
	} else {
		fmt.Printf("Batch ID: %s\n", result.ID)
	}
}
