package wallet

import (
	"context"
	"fmt"

	json "github.com/goccy/go-json"

	"github.com/ChefBingbong/viem-go/utils/address"
)

// GetAddressesReturnType is the return type for the GetAddresses action.
// It represents a list of checksummed Ethereum addresses.
type GetAddressesReturnType = []address.Address

// GetAddresses returns a list of account addresses owned by the wallet or client.
//
// If the client has a local account attached (implements LocalAccount),
// it returns the local address directly without making an RPC call,
// mirroring viem's `client.account?.type === 'local'` check.
//
// This is equivalent to viem's `getAddresses` action.
//
// JSON-RPC Method: eth_accounts
//
// Example:
//
//	addresses, err := wallet.GetAddresses(ctx, client)
//	// addresses == ["0xa5cc3c03994DB5b0d9A5eEdD10CabaB0813678AC"]
func GetAddresses(ctx context.Context, client Client) (GetAddressesReturnType, error) {
	// If the client has a local account, return it directly (mirrors viem's local account check)
	if acct := client.Account(); acct != nil {
		if _, isLocal := acct.(LocalAccount); isLocal {
			checksummed := address.ChecksumAddress(acct.Address().Hex())
			return []address.Address{checksummed}, nil
		}
	}

	resp, err := client.Request(ctx, "eth_accounts")
	if err != nil {
		return nil, fmt.Errorf("eth_accounts failed: %w", err)
	}

	var hexAddresses []string
	if unmarshalErr := json.Unmarshal(resp.Result, &hexAddresses); unmarshalErr != nil {
		return nil, fmt.Errorf("failed to unmarshal addresses: %w", unmarshalErr)
	}

	addresses := make([]address.Address, len(hexAddresses))
	for i, addr := range hexAddresses {
		addresses[i] = address.ChecksumAddress(addr)
	}

	return addresses, nil
}
