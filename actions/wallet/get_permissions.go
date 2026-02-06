package wallet

import (
	"context"
	"encoding/json"
	"fmt"
)

// WalletPermissionCaveat represents a caveat on a wallet permission.
// This mirrors viem's WalletPermissionCaveat type from EIP-2255.
type WalletPermissionCaveat struct {
	Type  string `json:"type"`
	Value any    `json:"value"`
}

// WalletPermission represents a wallet permission.
// This mirrors viem's WalletPermission type from EIP-2255.
type WalletPermission struct {
	Caveats          []WalletPermissionCaveat `json:"caveats"`
	Date             int64                    `json:"date"`
	ID               string                   `json:"id"`
	Invoker          string                   `json:"invoker"`
	ParentCapability string                   `json:"parentCapability"`
}

// GetPermissionsReturnType is the return type for the GetPermissions action.
type GetPermissionsReturnType = []WalletPermission

// GetPermissions gets the wallet's current permissions.
//
// This is equivalent to viem's `getPermissions` action.
//
// JSON-RPC Method: wallet_getPermissions (EIP-2255)
//
// Example:
//
//	permissions, err := wallet.GetPermissions(ctx, client)
func GetPermissions(ctx context.Context, client Client) (GetPermissionsReturnType, error) {
	resp, err := client.Request(ctx, "wallet_getPermissions")
	if err != nil {
		return nil, fmt.Errorf("wallet_getPermissions failed: %w", err)
	}

	var permissions []WalletPermission
	if unmarshalErr := json.Unmarshal(resp.Result, &permissions); unmarshalErr != nil {
		return nil, fmt.Errorf("failed to unmarshal permissions: %w", unmarshalErr)
	}

	return permissions, nil
}
