package wallet

import (
	"context"
	"fmt"

	json "github.com/goccy/go-json"
)

// RequestPermissionsParameters contains the parameters for the RequestPermissions action.
// This mirrors viem's RequestPermissionsParameters type.
// The map keys are permission names (e.g., "eth_accounts") and values are
// permission-specific configuration objects.
type RequestPermissionsParameters = map[string]map[string]any

// RequestPermissionsReturnType is the return type for the RequestPermissions action.
type RequestPermissionsReturnType = []WalletPermission

// RequestPermissions requests permissions for a wallet.
//
// This is equivalent to viem's `requestPermissions` action.
//
// JSON-RPC Method: wallet_requestPermissions (EIP-2255)
//
// Example:
//
//	permissions, err := wallet.RequestPermissions(ctx, client, wallet.RequestPermissionsParameters{
//	    "eth_accounts": {},
//	})
func RequestPermissions(ctx context.Context, client Client, permissions RequestPermissionsParameters) (RequestPermissionsReturnType, error) {
	resp, err := client.Request(ctx, "wallet_requestPermissions", permissions)
	if err != nil {
		return nil, fmt.Errorf("wallet_requestPermissions failed: %w", err)
	}

	var result []WalletPermission
	if unmarshalErr := json.Unmarshal(resp.Result, &result); unmarshalErr != nil {
		return nil, fmt.Errorf("failed to unmarshal permissions: %w", unmarshalErr)
	}

	return result, nil
}
