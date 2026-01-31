package client

import (
	"context"
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
)

// ClientDecorator is a function that extends a client with additional functionality.
type ClientDecorator func(*BaseClient) map[string]any

// WithPublicActions adds public (read) actions to the client.
// This mirrors viem's publicActions decorator.
func WithPublicActions() ClientDecorator {
	return func(c *BaseClient) map[string]any {
		return map[string]any{
			"getBlockNumber":        c.GetBlockNumber,
			"getChainId":            c.GetChainID,
			"getGasPrice":           c.GetGasPrice,
			"getBalance":            c.GetBalance,
			"getTransactionCount":   c.GetTransactionCount,
			"getCode":               c.GetCode,
			"call":                  c.Call,
			"estimateGas":           c.EstimateGas,
			"getTransaction":        c.GetTransaction,
			"getTransactionReceipt": c.GetTransactionReceipt,
			"getBlock":              c.GetBlock,
			"getBlockByHash":        c.GetBlockByHash,
			"getLogs":               c.GetLogs,
		}
	}
}

// WithWalletActions adds wallet (write) actions to the client.
// This mirrors viem's walletActions decorator.
func WithWalletActions() ClientDecorator {
	return func(c *BaseClient) map[string]any {
		return map[string]any{
			"sendRawTransaction": c.SendRawTransaction,
		}
	}
}

// ExtendClient applies decorators to a client.
func ExtendClient(client *BaseClient, decorators ...ClientDecorator) *BaseClient {
	for _, decorator := range decorators {
		extensions := decorator(client)
		for key, value := range extensions {
			client.Extend(key, value)
		}
	}
	return client
}

// PublicClient is a client with public actions pre-configured.
type PublicClient struct {
	*BaseClient
}

// CreatePublicClient creates a new public client.
func CreatePublicClient(config ClientConfig) (*PublicClient, error) {
	if config.Type == "" {
		config.Type = "public"
	}
	if config.Name == "" {
		config.Name = "Public Client"
	}

	base, err := CreateClient(config)
	if err != nil {
		return nil, err
	}

	ExtendClient(base, WithPublicActions())

	return &PublicClient{BaseClient: base}, nil
}

// WalletClient is a client with wallet actions pre-configured.
type WalletClient struct {
	*BaseClient
}

// CreateWalletClient creates a new wallet client.
func CreateWalletClient(config ClientConfig) (*WalletClient, error) {
	if config.Type == "" {
		config.Type = "wallet"
	}
	if config.Name == "" {
		config.Name = "Wallet Client"
	}

	base, err := CreateClient(config)
	if err != nil {
		return nil, err
	}

	ExtendClient(base, WithWalletActions())

	return &WalletClient{BaseClient: base}, nil
}

// ---- Type-safe action wrappers ----

// PublicActions provides type-safe access to public client actions.
type PublicActions struct {
	client *BaseClient
}

// NewPublicActions creates a new PublicActions wrapper.
func NewPublicActions(client *BaseClient) *PublicActions {
	return &PublicActions{client: client}
}

// GetBlockNumber returns the current block number.
func (a *PublicActions) GetBlockNumber(ctx context.Context) (uint64, error) {
	return a.client.GetBlockNumber(ctx)
}

// GetChainID returns the chain ID.
func (a *PublicActions) GetChainID(ctx context.Context) (uint64, error) {
	return a.client.GetChainID(ctx)
}

// GetGasPrice returns the current gas price.
func (a *PublicActions) GetGasPrice(ctx context.Context) (uint64, error) {
	return a.client.GetGasPrice(ctx)
}

// GetBalance returns the balance of an address.
func (a *PublicActions) GetBalance(ctx context.Context, address common.Address, blockTag ...BlockTag) (json.RawMessage, error) {
	return a.client.GetBalance(ctx, address, blockTag...)
}

// GetTransactionCount returns the nonce for an address.
func (a *PublicActions) GetTransactionCount(ctx context.Context, address common.Address, blockTag ...BlockTag) (uint64, error) {
	return a.client.GetTransactionCount(ctx, address, blockTag...)
}

// GetCode returns the code at an address.
func (a *PublicActions) GetCode(ctx context.Context, address common.Address, blockTag ...BlockTag) ([]byte, error) {
	return a.client.GetCode(ctx, address, blockTag...)
}

// Call performs an eth_call.
func (a *PublicActions) Call(ctx context.Context, callData map[string]any, blockTag ...BlockTag) ([]byte, error) {
	return a.client.Call(ctx, callData, blockTag...)
}

// EstimateGas estimates gas for a call.
func (a *PublicActions) EstimateGas(ctx context.Context, callData map[string]any) (uint64, error) {
	return a.client.EstimateGas(ctx, callData)
}

// GetTransaction returns a transaction by hash.
func (a *PublicActions) GetTransaction(ctx context.Context, hash common.Hash) (json.RawMessage, error) {
	return a.client.GetTransaction(ctx, hash)
}

// GetTransactionReceipt returns a transaction receipt.
func (a *PublicActions) GetTransactionReceipt(ctx context.Context, hash common.Hash) (json.RawMessage, error) {
	return a.client.GetTransactionReceipt(ctx, hash)
}

// GetBlock returns a block by number or tag.
func (a *PublicActions) GetBlock(ctx context.Context, blockTag BlockTag, includeTransactions bool) (json.RawMessage, error) {
	return a.client.GetBlock(ctx, blockTag, includeTransactions)
}

// GetBlockByHash returns a block by hash.
func (a *PublicActions) GetBlockByHash(ctx context.Context, hash common.Hash, includeTransactions bool) (json.RawMessage, error) {
	return a.client.GetBlockByHash(ctx, hash, includeTransactions)
}

// GetLogs returns logs matching the filter.
func (a *PublicActions) GetLogs(ctx context.Context, filter map[string]any) (json.RawMessage, error) {
	return a.client.GetLogs(ctx, filter)
}

// WalletActions provides type-safe access to wallet client actions.
type WalletActions struct {
	client *BaseClient
}

// NewWalletActions creates a new WalletActions wrapper.
func NewWalletActions(client *BaseClient) *WalletActions {
	return &WalletActions{client: client}
}

// SendRawTransaction sends a signed raw transaction.
func (a *WalletActions) SendRawTransaction(ctx context.Context, signedTx []byte) (common.Hash, error) {
	return a.client.SendRawTransaction(ctx, signedTx)
}
