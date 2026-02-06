package main

import (
	"fmt"
	"math/big"

	"github.com/ChefBingbong/viem-go/actions/wallet"
	"github.com/ChefBingbong/viem-go/utils/signature"
)

// RunSignTypedData demonstrates EIP-712 structured typed data signing via the WalletClient.
func RunSignTypedData() {
	PrintSection("signTypedData - Sign EIP-712 typed data")

	sig, err := WalletCl.SignTypedData(ctx, wallet.SignTypedDataParameters{
		Domain: signature.TypedDataDomain{
			Name:              "Ether Mail",
			Version:           "1",
			ChainId:           big.NewInt(137), // Polygon
			VerifyingContract: "0xCcCCccccCCCCcCCCCCCcCcCccCcCCCcCcccccccC",
		},
		Types: map[string][]signature.TypedDataField{
			"Person": {
				{Name: "name", Type: "string"},
				{Name: "wallet", Type: "address"},
			},
			"Mail": {
				{Name: "from", Type: "Person"},
				{Name: "to", Type: "Person"},
				{Name: "contents", Type: "string"},
			},
		},
		PrimaryType: "Mail",
		Message: map[string]any{
			"from":     map[string]any{"name": "Cow", "wallet": "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826"},
			"to":       map[string]any{"name": "Bob", "wallet": "0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB"},
			"contents": "Hello, Bob!",
		},
	})
	if err != nil {
		// Expected — the public RPC doesn't manage our key, so JSON-RPC path will fail.
		// For local signing, the WalletClient will use the account's SignTypedData if available.
		fmt.Printf("signTypedData: %v\n", Truncate(err.Error(), 120))
		fmt.Println("(Expected — RPC doesn't manage this key. Use a local account for offline signing.)")

		// Demonstrate local signing as well
		typedData := signature.TypedDataDefinition{
			Domain: signature.TypedDataDomain{
				Name:              "Ether Mail",
				Version:           "1",
				ChainId:           big.NewInt(137),
				VerifyingContract: "0xCcCCccccCCCCcCCCCCCcCcCccCcCCCcCcccccccC",
			},
			Types: map[string][]signature.TypedDataField{
				"Person": {
					{Name: "name", Type: "string"},
					{Name: "wallet", Type: "address"},
				},
				"Mail": {
					{Name: "from", Type: "Person"},
					{Name: "to", Type: "Person"},
					{Name: "contents", Type: "string"},
				},
			},
			PrimaryType: "Mail",
			Message: map[string]any{
				"from":     map[string]any{"name": "Cow", "wallet": "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826"},
				"to":       map[string]any{"name": "Bob", "wallet": "0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB"},
				"contents": "Hello, Bob!",
			},
		}
		localSig, localErr := LocalAccount.SignTypedData(typedData)
		if localErr != nil {
			fmt.Printf("Local signTypedData error: %v\n", localErr)
		} else {
			fmt.Printf("Local typed data signature: %s...\n", Truncate(localSig, 42))
		}
	} else {
		fmt.Printf("Typed data signature: %s...\n", Truncate(sig, 42))
	}
}
