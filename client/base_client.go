package client

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"

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

// Chain represents chain configuration for the client.
type Chain struct {
	// ID is the chain ID.
	ID int
	// Name is the human-readable chain name.
	Name string
	// BlockTime is the block time in milliseconds.
	BlockTime int
	// NativeCurrency contains the native currency information.
	NativeCurrency *NativeCurrency
	// RPCUrls contains the RPC endpoints.
	RPCUrls *ChainRPCUrls
	// Contracts contains well-known contract addresses.
	Contracts *ChainContracts
	// Testnet indicates if this is a test network.
	Testnet bool
}

// NativeCurrency represents the native currency of a chain.
type NativeCurrency struct {
	Name     string
	Symbol   string
	Decimals int
}

// ChainRPCUrls contains RPC URLs for a chain.
type ChainRPCUrls struct {
	Default ChainRPCEndpoints
	Public  ChainRPCEndpoints
}

// ChainRPCEndpoints contains HTTP and WebSocket endpoints.
type ChainRPCEndpoints struct {
	HTTP      []string
	WebSocket []string
}

// ChainContracts contains well-known contract addresses.
type ChainContracts struct {
	Multicall3 *ChainContract
}

// ChainContract represents a contract address with optional creation block.
type ChainContract struct {
	Address      common.Address
	BlockCreated uint64
}

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
	blockTime := 12000 // Default 12 seconds
	if config.Chain != nil && config.Chain.BlockTime > 0 {
		blockTime = config.Chain.BlockTime
	}

	defaultPollingInterval := time.Duration(min(max(blockTime/2, 500), 4000)) * time.Millisecond
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

	// Convert chain to transport chain
	var transportChain *transport.Chain
	if config.Chain != nil {
		transportChain = &transport.Chain{
			ID:        config.Chain.ID,
			Name:      config.Chain.Name,
			BlockTime: config.Chain.BlockTime,
		}
		if config.Chain.RPCUrls != nil {
			transportChain.RPCUrls = transport.ChainRPCUrls{
				Default: transport.ChainRPCEndpoints{
					HTTP:      config.Chain.RPCUrls.Default.HTTP,
					WebSocket: config.Chain.RPCUrls.Default.WebSocket,
				},
			}
		}
	}

	tr, err := config.Transport(transport.TransportParams{
		Chain:           transportChain,
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

// Request sends a JSON-RPC request.
func (c *BaseClient) Request(ctx context.Context, method string, params ...any) (*transport.RPCResponse, error) {
	req := transport.RPCRequest{
		Method: method,
		Params: params,
	}
	return c.transport.Request(ctx, req)
}

// Close closes the client.
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

// ---- RPC Methods ----

// GetBlockNumber returns the current block number.
func (c *BaseClient) GetBlockNumber(ctx context.Context) (uint64, error) {
	resp, err := c.Request(ctx, "eth_blockNumber")
	if err != nil {
		return 0, err
	}

	var hexNumber string
	if err := json.Unmarshal(resp.Result, &hexNumber); err != nil {
		return 0, err
	}

	return parseHexUint64(hexNumber)
}

// GetChainID returns the chain ID.
func (c *BaseClient) GetChainID(ctx context.Context) (uint64, error) {
	resp, err := c.Request(ctx, "eth_chainId")
	if err != nil {
		return 0, err
	}

	var hexChainID string
	if err := json.Unmarshal(resp.Result, &hexChainID); err != nil {
		return 0, err
	}

	return parseHexUint64(hexChainID)
}

// GetGasPrice returns the current gas price.
func (c *BaseClient) GetGasPrice(ctx context.Context) (uint64, error) {
	resp, err := c.Request(ctx, "eth_gasPrice")
	if err != nil {
		return 0, err
	}

	var hexGasPrice string
	if err := json.Unmarshal(resp.Result, &hexGasPrice); err != nil {
		return 0, err
	}

	return parseHexUint64(hexGasPrice)
}

// GetBalance returns the balance of an address.
func (c *BaseClient) GetBalance(ctx context.Context, address common.Address, blockTag ...BlockTag) (json.RawMessage, error) {
	tag := BlockTagLatest
	if len(blockTag) > 0 {
		tag = blockTag[0]
	}
	if c.experimentalBlockTag != "" && len(blockTag) == 0 {
		tag = c.experimentalBlockTag
	}

	resp, err := c.Request(ctx, "eth_getBalance", address.Hex(), string(tag))
	if err != nil {
		return nil, err
	}

	return resp.Result, nil
}

// GetTransactionCount returns the nonce for an address.
func (c *BaseClient) GetTransactionCount(ctx context.Context, address common.Address, blockTag ...BlockTag) (uint64, error) {
	tag := BlockTagLatest
	if len(blockTag) > 0 {
		tag = blockTag[0]
	}
	if c.experimentalBlockTag != "" && len(blockTag) == 0 {
		tag = c.experimentalBlockTag
	}

	resp, err := c.Request(ctx, "eth_getTransactionCount", address.Hex(), string(tag))
	if err != nil {
		return 0, err
	}

	var hexNonce string
	if err := json.Unmarshal(resp.Result, &hexNonce); err != nil {
		return 0, err
	}

	return parseHexUint64(hexNonce)
}

// GetCode returns the code at an address.
func (c *BaseClient) GetCode(ctx context.Context, address common.Address, blockTag ...BlockTag) ([]byte, error) {
	tag := BlockTagLatest
	if len(blockTag) > 0 {
		tag = blockTag[0]
	}
	if c.experimentalBlockTag != "" && len(blockTag) == 0 {
		tag = c.experimentalBlockTag
	}

	resp, err := c.Request(ctx, "eth_getCode", address.Hex(), string(tag))
	if err != nil {
		return nil, err
	}

	var hexCode string
	if err := json.Unmarshal(resp.Result, &hexCode); err != nil {
		return nil, err
	}

	return parseHexBytes(hexCode)
}

// Call performs an eth_call.
func (c *BaseClient) Call(ctx context.Context, callData map[string]any, blockTag ...BlockTag) ([]byte, error) {
	tag := BlockTagLatest
	if len(blockTag) > 0 {
		tag = blockTag[0]
	}
	if c.experimentalBlockTag != "" && len(blockTag) == 0 {
		tag = c.experimentalBlockTag
	}

	resp, err := c.Request(ctx, "eth_call", callData, string(tag))
	if err != nil {
		return nil, err
	}

	var hexResult string
	if err := json.Unmarshal(resp.Result, &hexResult); err != nil {
		return nil, err
	}

	return parseHexBytes(hexResult)
}

// EstimateGas estimates gas for a call.
func (c *BaseClient) EstimateGas(ctx context.Context, callData map[string]any) (uint64, error) {
	resp, err := c.Request(ctx, "eth_estimateGas", callData)
	if err != nil {
		return 0, err
	}

	var hexGas string
	if err := json.Unmarshal(resp.Result, &hexGas); err != nil {
		return 0, err
	}

	return parseHexUint64(hexGas)
}

// SendRawTransaction sends a signed raw transaction.
func (c *BaseClient) SendRawTransaction(ctx context.Context, signedTx []byte) (common.Hash, error) {
	hexTx := "0x" + hex.EncodeToString(signedTx)

	resp, err := c.Request(ctx, "eth_sendRawTransaction", hexTx)
	if err != nil {
		return common.Hash{}, err
	}

	var hashHex string
	if err := json.Unmarshal(resp.Result, &hashHex); err != nil {
		return common.Hash{}, err
	}

	return common.HexToHash(hashHex), nil
}

// GetTransaction returns a transaction by hash.
func (c *BaseClient) GetTransaction(ctx context.Context, hash common.Hash) (json.RawMessage, error) {
	resp, err := c.Request(ctx, "eth_getTransactionByHash", hash.Hex())
	if err != nil {
		return nil, err
	}

	return resp.Result, nil
}

// GetTransactionReceipt returns a transaction receipt.
func (c *BaseClient) GetTransactionReceipt(ctx context.Context, hash common.Hash) (json.RawMessage, error) {
	resp, err := c.Request(ctx, "eth_getTransactionReceipt", hash.Hex())
	if err != nil {
		return nil, err
	}

	return resp.Result, nil
}

// GetBlock returns a block by number or tag.
func (c *BaseClient) GetBlock(ctx context.Context, blockTag BlockTag, includeTransactions bool) (json.RawMessage, error) {
	resp, err := c.Request(ctx, "eth_getBlockByNumber", string(blockTag), includeTransactions)
	if err != nil {
		return nil, err
	}

	return resp.Result, nil
}

// GetBlockByHash returns a block by hash.
func (c *BaseClient) GetBlockByHash(ctx context.Context, hash common.Hash, includeTransactions bool) (json.RawMessage, error) {
	resp, err := c.Request(ctx, "eth_getBlockByHash", hash.Hex(), includeTransactions)
	if err != nil {
		return nil, err
	}

	return resp.Result, nil
}

// GetLogs returns logs matching the filter.
func (c *BaseClient) GetLogs(ctx context.Context, filter map[string]any) (json.RawMessage, error) {
	resp, err := c.Request(ctx, "eth_getLogs", filter)
	if err != nil {
		return nil, err
	}

	return resp.Result, nil
}

// ---- Helper functions ----

// parseHexUint64 parses a hex string to uint64.
func parseHexUint64(hexStr string) (uint64, error) {
	if len(hexStr) >= 2 && hexStr[:2] == "0x" {
		hexStr = hexStr[2:]
	}
	if hexStr == "" {
		return 0, nil
	}

	var result uint64
	for _, c := range hexStr {
		result *= 16
		switch {
		case c >= '0' && c <= '9':
			result += uint64(c - '0')
		case c >= 'a' && c <= 'f':
			result += uint64(c - 'a' + 10)
		case c >= 'A' && c <= 'F':
			result += uint64(c - 'A' + 10)
		}
	}
	return result, nil
}

// parseHexBytes parses a hex string to bytes.
func parseHexBytes(hexStr string) ([]byte, error) {
	if len(hexStr) >= 2 && hexStr[:2] == "0x" {
		hexStr = hexStr[2:]
	}
	if hexStr == "" {
		return []byte{}, nil
	}
	return hex.DecodeString(hexStr)
}
