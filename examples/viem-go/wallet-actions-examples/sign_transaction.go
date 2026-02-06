package main

import (
	"fmt"
	"math/big"

	"github.com/ChefBingbong/viem-go/actions/wallet"
	"github.com/ChefBingbong/viem-go/utils/formatters"
	utiltx "github.com/ChefBingbong/viem-go/utils/transaction"
)

// RunSignTransaction demonstrates signing transactions without broadcasting via the WalletClient.
func RunSignTransaction() {
	PrintSection("signTransaction - Sign without broadcasting")

	nonce := 0

	// EIP-1559 via walletClient (delegates to wallet.SignTransaction)
	signed1559, err := WalletCl.SignTransaction(ctx, wallet.SignTransactionParameters{
		To:                   TargetAddr.Hex(),
		Value:                MustParseEther("0.01"),
		Gas:                  big.NewInt(21000),
		MaxFeePerGas:         MustParseGwei("50"),
		MaxPriorityFeePerGas: MustParseGwei("2"),
		Type:                 formatters.TransactionTypeEIP1559,
		Nonce:                &nonce,
	})
	if err != nil {
		// Expected if RPC doesn't manage the key — fall back to local signing
		fmt.Printf("WalletClient.SignTransaction: %v\n", Truncate(err.Error(), 100))
		fmt.Println("Falling back to local account signing...")

		tx1559 := &utiltx.Transaction{
			Type:                 utiltx.TransactionTypeEIP1559,
			ChainId:              137,
			To:                   TargetAddr.Hex(),
			Value:                MustParseEther("0.01"),
			Gas:                  big.NewInt(21000),
			MaxFeePerGas:         MustParseGwei("50"),
			MaxPriorityFeePerGas: MustParseGwei("2"),
			Nonce:                0,
		}
		signed1559, err = LocalAccount.SignTransaction(tx1559)
		if err != nil {
			fmt.Printf("Local EIP-1559 error: %v\n", err)
		}
	}
	if err == nil {
		fmt.Printf("Signed tx (EIP-1559): %s...\n", Truncate(signed1559, 42))
	}

	// Legacy via local account
	txLegacy := &utiltx.Transaction{
		Type:     utiltx.TransactionTypeLegacy,
		ChainId:  137,
		To:       TargetAddr.Hex(),
		Value:    MustParseEther("0.01"),
		Gas:      big.NewInt(21000),
		GasPrice: MustParseGwei("50"),
		Nonce:    0,
	}
	signedLegacy, err := LocalAccount.SignTransaction(txLegacy)
	if err != nil {
		fmt.Printf("Legacy error: %v\n", err)
	} else {
		fmt.Printf("Signed tx (legacy):   %s...\n", Truncate(signedLegacy, 42))
	}

	// EIP-2930 via local account
	tx2930 := &utiltx.Transaction{
		Type:     utiltx.TransactionTypeEIP2930,
		ChainId:  137,
		To:       TargetAddr.Hex(),
		Value:    MustParseEther("0.01"),
		Gas:      big.NewInt(21000),
		GasPrice: MustParseGwei("50"),
		Nonce:    0,
		AccessList: utiltx.AccessList{
			{Address: TargetAddr.Hex(), StorageKeys: []string{
				"0x0000000000000000000000000000000000000000000000000000000000000001",
			}},
		},
	}
	signed2930, err := LocalAccount.SignTransaction(tx2930)
	if err != nil {
		fmt.Printf("EIP-2930 error: %v\n", err)
	} else {
		fmt.Printf("Signed tx (EIP-2930): %s...\n", Truncate(signed2930, 42))
	}
}
