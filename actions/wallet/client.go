// Package wallet provides standalone action functions for wallet (write) Ethereum JSON-RPC methods.
// These actions can be used directly or through a WalletClient.
//
// This mirrors viem's actions pattern where actions are standalone functions
// that take a client interface as their first parameter.
package wallet

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/chain"
	"github.com/ChefBingbong/viem-go/client/transport"
)

// Client is the interface that wallet actions require from a client.
// This allows actions to be used with any client implementation
// that satisfies this interface.
type Client interface {
	// Request sends a raw JSON-RPC request.
	Request(ctx context.Context, method string, params ...any) (*transport.RPCResponse, error)

	// Chain returns the chain configuration, if set.
	Chain() *chain.Chain

	// Account returns the account associated with this client, if any.
	// Returns nil if no account is set.
	Account() Account
}

// Account represents an account that can be used with the client.
// This mirrors the client package's Account interface.
type Account interface {
	// Address returns the account address.
	Address() common.Address
}

// LocalAccount extends Account to indicate a locally-managed account (private key, HD, etc.).
// When present, actions like GetAddresses can return the local address directly
// without making an RPC call, mirroring viem's `client.account?.type === 'local'` check.
type LocalAccount interface {
	Account
	// IsLocal is a marker method indicating this is a local account.
	IsLocal()
}
