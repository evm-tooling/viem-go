package main

import (
	"fmt"

	"github.com/ChefBingbong/viem-go/actions/wallet"
)

// RunSignAuthorization demonstrates EIP-7702 authorization signing via the WalletClient.
func RunSignAuthorization() {
	PrintSection("signAuthorization - EIP-7702 authorization signing")

	// prepareAuthorization fills in chainId + nonce from the network
	fmt.Println("Preparing authorization (fetches chainId + nonce from RPC)...")
	auth, err := WalletCl.PrepareAuthorization(ctx, wallet.PrepareAuthorizationParameters{
		ContractAddress: TargetAddr.Hex(),
	})
	if err != nil {
		fmt.Printf("prepareAuthorization: %v\n", Truncate(err.Error(), 120))
		fmt.Println("(Expected — account not found or RPC doesn't support pending nonce)")
	} else {
		fmt.Printf("Prepared authorization: address=%s chainId=%d nonce=%d\n", auth.Address, auth.ChainId, auth.Nonce)
	}

	// signAuthorization signs the authorization locally
	// (requires a local account that implements AuthorizationSignableAccount)
	fmt.Println("\nSigning authorization locally...")
	signed, signErr := WalletCl.SignAuthorization(ctx, wallet.SignAuthorizationParameters{
		ContractAddress: TargetAddr.Hex(),
	})
	if signErr != nil {
		fmt.Printf("signAuthorization: %v\n", Truncate(signErr.Error(), 120))
		fmt.Println("(Expected — the JSON-RPC account doesn't support local signing)")
		fmt.Println("Use a PrivateKeyAccount for local EIP-7702 signing.")
	} else {
		fmt.Printf("Authorization address: %s\n", signed.Address)
		fmt.Printf("Authorization chainId: %d\n", signed.ChainId)
		fmt.Printf("Authorization nonce:   %d\n", signed.Nonce)
		fmt.Printf("Signature r: %s...\n", Truncate(signed.R, 22))
	}
}
