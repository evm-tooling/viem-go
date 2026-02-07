package wallet

import (
	"context"
	"fmt"

	json "github.com/goccy/go-json"

	"github.com/ChefBingbong/viem-go/utils/address"
)

// RequestAddressesReturnType is the return type for the RequestAddresses action.
// It represents a list of checksummed Ethereum addresses.
type RequestAddressesReturnType = []address.Address

// RequestAddresses requests a list of accounts managed by a wallet.
//
// Sends a request to the wallet, asking for permission to access the user's
// accounts. After the user accepts the request, it will return a list of
// accounts (addresses).
//
// This API can be useful for dapps that need to access the user's accounts
// in order to execute transactions or interact with smart contracts.
//
// This is equivalent to viem's `requestAddresses` action.
//
// JSON-RPC Method: eth_requestAccounts (EIP-1102)
//
// Example:
//
//	addresses, err := wallet.RequestAddresses(ctx, client)
//	// addresses == ["0xa5cc3c03994DB5b0d9A5eEdD10CabaB0813678AC"]
func RequestAddresses(ctx context.Context, client Client) (RequestAddressesReturnType, error) {
	resp, err := client.Request(ctx, "eth_requestAccounts")
	if err != nil {
		return nil, fmt.Errorf("eth_requestAccounts failed: %w", err)
	}

	var hexAddresses []string
	if unmarshalErr := json.Unmarshal(resp.Result, &hexAddresses); unmarshalErr != nil {
		return nil, fmt.Errorf("failed to unmarshal addresses: %w", unmarshalErr)
	}

	addresses := make([]address.Address, len(hexAddresses))
	for i, addr := range hexAddresses {
		checksummed, getErr := address.GetAddress(addr)
		if getErr != nil {
			return nil, fmt.Errorf("failed to checksum address %s: %w", addr, getErr)
		}
		addresses[i] = checksummed
	}

	return addresses, nil
}
