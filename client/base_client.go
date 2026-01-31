package client

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/chain"
	"github.com/ChefBingbong/viem-go/client/transport"
	"github.com/ChefBingbong/viem-go/types"
)

// BlockTag is an alias for types.BlockTag.
type BlockTag = types.BlockTag

const (
	BlockTagLatest    = types.BlockTagLatest
	BlockTagPending   = types.BlockTagPending
	BlockTagEarliest  = types.BlockTagEarliest
	BlockTagSafe      = types.BlockTagSafe
	BlockTagFinalized = types.BlockTagFinalized
)

// MulticallBatchOptions contains options for multicall batching.
type MulticallBatchOptions struct {
	// BatchSize is the maximum size (in bytes) for each calldata chunk.
	BatchSize int
	// Deployless enables deployless multicall.
	Deployless bool
	// Wait is the maximum number of milliseconds to wait before sending a batch.
	Wait time.Duration
}

// BatchOptions contains batch settings.
type BatchOptions struct {
	// Multicall enables eth_call multicall aggregation.
	Multicall *MulticallBatchOptions
}

// Account represents an account that can be used with the client.
type Account interface {
	// Address returns the account address.
	Address() common.Address
}

// AddressAccount is a simple account that only has an address (JSON-RPC account).
type AddressAccount struct {
	addr common.Address
}

// NewAddressAccount creates a new address-only account.
func NewAddressAccount(addr common.Address) *AddressAccount {
	return &AddressAccount{addr: addr}
}

// Address returns the account address.
func (a *AddressAccount) Address() common.Address {
	return a.addr
}

// Chain is an alias for chain.Chain.
type Chain = chain.Chain

// ClientConfig contains configuration for creating a client.
type ClientConfig struct {
	// Account is the account to use for the client.
	Account Account
	// Batch contains batch settings.
	Batch *BatchOptions
	// CacheTime is the time (in ms) that cached data will remain in memory.
	CacheTime time.Duration
	// Chain is the chain configuration.
	Chain *Chain
	// DataSuffix is the data suffix to append to transaction data.
	DataSuffix []byte
	// ExperimentalBlockTag is the default block tag for RPC requests.
	ExperimentalBlockTag BlockTag
	// Key is a key for the client.
	Key string
	// Name is a name for the client.
	Name string
	// PollingInterval is the frequency (in ms) for polling enabled actions & events.
	PollingInterval time.Duration
	// Transport is the transport factory to use.
	Transport transport.TransportFactory
	// Type is the type of client.
	Type string
}

// DefaultClientConfig returns default client configuration.
func DefaultClientConfig() ClientConfig {
	return ClientConfig{
		Key:             "base",
		Name:            "Base Client",
		Type:            "base",
		CacheTime:       4000 * time.Millisecond,
		PollingInterval: 4000 * time.Millisecond,
	}
}

// BaseClient is the base JSON-RPC client that mirrors viem's createClient.
// It only provides raw request capabilities - use PublicClient or WalletClient
// for typed RPC method wrappers.
type BaseClient struct {
	// Account is the account of the client.
	account Account
	// Batch contains batch settings.
	batch *BatchOptions
	// CacheTime is the time (in ms) that cached data will remain in memory.
	cacheTime time.Duration
	// Chain is the chain configuration.
	chain *Chain
	// DataSuffix is the data suffix to append to transaction data.
	dataSuffix []byte
	// ExperimentalBlockTag is the default block tag for RPC requests.
	experimentalBlockTag BlockTag
	// Key is a key for the client.
	key string
	// Name is a name for the client.
	name string
	// PollingInterval is the frequency for polling.
	pollingInterval time.Duration
	// Transport is the underlying transport.
	transport transport.Transport
	// Type is the type of client.
	clientType string
	// UID is a unique identifier for the client.
	uid string

	// extensions holds extended functionality
	extensions map[string]any
	extMu      sync.RWMutex
}

// CreateClient creates a new base client with the given configuration.
// This is the low-level client factory - prefer CreatePublicClient or
// CreateWalletClient for most use cases.
func CreateClient(config ClientConfig) (*BaseClient, error) {
	// Apply defaults
	if config.Key == "" {
		config.Key = "base"
	}
	if config.Name == "" {
		config.Name = "Base Client"
	}
	if config.Type == "" {
		config.Type = "base"
	}

	// Calculate default polling interval based on chain block time
	blockTime := int64(12000) // Default 12 seconds
	if config.Chain != nil && config.Chain.BlockTime != nil && *config.Chain.BlockTime > 0 {
		blockTime = *config.Chain.BlockTime
	}

	defaultPollingInterval := time.Duration(min(max(int(blockTime/2), 500), 4000)) * time.Millisecond
	if config.PollingInterval == 0 {
		config.PollingInterval = defaultPollingInterval
	}
	if config.CacheTime == 0 {
		config.CacheTime = config.PollingInterval
	}

	// Determine experimental block tag
	experimentalBlockTag := config.ExperimentalBlockTag

	// Create transport
	if config.Transport == nil {
		return nil, transport.ErrURLRequired
	}

	tr, err := config.Transport(transport.TransportParams{
		Chain:           config.Chain,
		PollingInterval: config.PollingInterval,
	})
	if err != nil {
		return nil, err
	}

	// Generate UID
	uid := generateUID(11)

	client := &BaseClient{
		account:              config.Account,
		batch:                config.Batch,
		cacheTime:            config.CacheTime,
		chain:                config.Chain,
		dataSuffix:           config.DataSuffix,
		experimentalBlockTag: experimentalBlockTag,
		key:                  config.Key,
		name:                 config.Name,
		pollingInterval:      config.PollingInterval,
		transport:            tr,
		clientType:           config.Type,
		uid:                  uid,
		extensions:           make(map[string]any),
	}

	return client, nil
}

// Account returns the client account.
func (c *BaseClient) Account() Account {
	return c.account
}

// Batch returns the batch options.
func (c *BaseClient) Batch() *BatchOptions {
	return c.batch
}

// CacheTime returns the cache time.
func (c *BaseClient) CacheTime() time.Duration {
	return c.cacheTime
}

// Chain returns the chain configuration.
func (c *BaseClient) Chain() *Chain {
	return c.chain
}

// DataSuffix returns the data suffix.
func (c *BaseClient) DataSuffix() []byte {
	return c.dataSuffix
}

// ExperimentalBlockTag returns the experimental block tag.
func (c *BaseClient) ExperimentalBlockTag() BlockTag {
	return c.experimentalBlockTag
}

// Key returns the client key.
func (c *BaseClient) Key() string {
	return c.key
}

// Name returns the client name.
func (c *BaseClient) Name() string {
	return c.name
}

// PollingInterval returns the polling interval.
func (c *BaseClient) PollingInterval() time.Duration {
	return c.pollingInterval
}

// Transport returns the underlying transport.
func (c *BaseClient) Transport() transport.Transport {
	return c.transport
}

// Type returns the client type.
func (c *BaseClient) Type() string {
	return c.clientType
}

// UID returns the unique client identifier.
func (c *BaseClient) UID() string {
	return c.uid
}

// Request sends a raw JSON-RPC request.
// This is the only RPC method on BaseClient - use PublicClient or WalletClient
// for typed method wrappers.
func (c *BaseClient) Request(ctx context.Context, method string, params ...any) (*transport.RPCResponse, error) {
	req := transport.RPCRequest{
		Method: method,
		Params: params,
	}
	return c.transport.Request(ctx, req)
}

// Close closes the client and its underlying transport.
func (c *BaseClient) Close() error {
	return c.transport.Close()
}

// Extend adds extended functionality to the client.
// This mirrors viem's extend pattern for adding decorators.
func (c *BaseClient) Extend(key string, value any) *BaseClient {
	c.extMu.Lock()
	defer c.extMu.Unlock()
	c.extensions[key] = value
	return c
}

// GetExtension retrieves an extended value by key.
func (c *BaseClient) GetExtension(key string) (any, bool) {
	c.extMu.RLock()
	defer c.extMu.RUnlock()
	val, ok := c.extensions[key]
	return val, ok
}

// Extensions returns all extensions.
func (c *BaseClient) Extensions() map[string]any {
	c.extMu.RLock()
	defer c.extMu.RUnlock()
	result := make(map[string]any, len(c.extensions))
	for k, v := range c.extensions {
		result[k] = v
	}
	return result
}

// generateUID generates a unique identifier.
func generateUID(length int) string {
	bytes := make([]byte, (length+1)/2)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based ID
		return hex.EncodeToString([]byte(time.Now().String()))[:length]
	}
	return hex.EncodeToString(bytes)[:length]
}

// Helper functions for min/max
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
