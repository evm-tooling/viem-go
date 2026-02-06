package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/actions/public"
	"github.com/ChefBingbong/viem-go/actions/wallet"
	"github.com/ChefBingbong/viem-go/utils/unit"
)

// RunSendRawTransaction demonstrates local sign + broadcast via walletClient.SendRawTransaction.
func RunSendRawTransaction() {
	PrintSection("sendRawTransaction - Local sign + broadcast")

	value := MustParseEther("0.00005")
	to := TargetAddr

	// Step 1: Simulate
	fmt.Printf("Simulating: %s -> %s (%s POL)\n", SourceAddr.Hex(), to.Hex(), unit.FormatEther(value))
	_, simErr := public.Call(ctx, PublicClient, public.CallParameters{
		Account: &SourceAddr,
		To:      &to,
		Value:   value,
	})
	if simErr != nil {
		fmt.Printf("Simulation note: %s\n", Truncate(simErr.Error(), 120))
	} else {
		fmt.Println("Simulation passed.")
	}

	// Step 2: Prepare the transaction (fills nonce, gas, fees, chainId, type)
	fmt.Println("\nPreparing transaction request...")
	prepared, err := WalletCl.PrepareTransactionRequest(ctx, wallet.PrepareTransactionRequestParameters{
		To:    to.Hex(),
		Value: value,
	})
	if err != nil {
		fmt.Printf("Prepare error: %v\n", err)
		return
	}
	fmt.Printf("Prepared: nonce=%v gas=%v type=%s chainId=%v\n",
		prepared.Nonce, prepared.Gas, prepared.Type, prepared.ChainID)

	// Step 3: Sign the prepared transaction locally
	fmt.Println("\nSigning prepared transaction...")
	serialized, err := WalletCl.SignPreparedTransaction(ctx, prepared)
	if err != nil {
		fmt.Printf("Sign error: %v\n", err)
		return
	}
	fmt.Printf("Signed tx: %s...\n", Truncate(serialized, 42))

	// Step 4: Broadcast via walletClient.SendRawTransaction
	fmt.Println("\nBroadcasting via walletClient.SendRawTransaction...")
	hash, err := WalletCl.SendRawTransaction(ctx, wallet.SendRawTransactionParameters{
		SerializedTransaction: serialized,
	})
	if err != nil {
		fmt.Printf("Broadcast error: %s\n", Truncate(err.Error(), 120))
		fmt.Println("(Expected on mainnet if account has insufficient funds for gas)")
		return
	}
	fmt.Printf("Tx hash: %s\n", hash)

	// Step 5: Wait for receipt
	receipt, err := PublicClient.WaitForTransactionReceipt(ctx, common.HexToHash(hash))
	if err != nil {
		fmt.Printf("Wait for receipt error: %s\n", Truncate(err.Error(), 120))
		return
	}
	fmt.Printf("Receipt: status=%d gasUsed=%d\n", receipt.Status, receipt.GasUsed)
}
