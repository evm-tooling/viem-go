package main

import (
	"fmt"

	"github.com/ChefBingbong/viem-go/actions/wallet"
	"github.com/ChefBingbong/viem-go/utils/signature"
)

// RunSignMessage demonstrates EIP-191 personal message signing via the WalletClient.
func RunSignMessage() {
	PrintSection("signMessage - Sign a plain text message (EIP-191)")

	// Use the walletClient instance method (delegates to wallet.SignMessage)
	sig, err := WalletCl.SignMessage(ctx, wallet.SignMessageParameters{
		Message: signature.NewSignableMessage("hello world"),
	})
	if err != nil {
		// Expected on Polygon — the RPC node doesn't manage our private key
		fmt.Printf("JSON-RPC sign: %v (expected — node doesn't manage this key)\n", Truncate(err.Error(), 100))
	} else {
		fmt.Printf("JSON-RPC signature: %s...\n", Truncate(sig, 42))
	}

	// Local account path: sign without any RPC call
	localSig, err := LocalAccount.SignMessage(signature.NewSignableMessage("hello world"))
	if err != nil {
		fmt.Printf("Local sign error: %v\n", err)
	} else {
		fmt.Printf("Local signature:    %s...\n", Truncate(localSig, 42))
	}

	// Raw hex variant
	rawSig, err := LocalAccount.SignMessage(signature.NewSignableMessageRawHex("0x68656c6c6f20776f726c64"))
	if err != nil {
		fmt.Printf("Raw sign error: %v\n", err)
	} else {
		fmt.Printf("Raw hex signature:  %s...\n", Truncate(rawSig, 42))
	}

	// Raw bytes variant
	bytesSig, err := LocalAccount.SignMessage(signature.NewSignableMessageRaw([]byte{104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100}))
	if err != nil {
		fmt.Printf("Bytes sign error: %v\n", err)
	} else {
		fmt.Printf("Raw bytes sig:      %s...\n", Truncate(bytesSig, 42))
	}

	fmt.Printf("\nAll local sigs match: %v\n", localSig == rawSig && localSig == bytesSig)
}
