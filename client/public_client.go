package client

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"

	viemabi "github.com/ChefBingbong/viem-go/abi"
	"github.com/ChefBingbong/viem-go/actions/public"
	"github.com/ChefBingbong/viem-go/chain"
	"github.com/ChefBingbong/viem-go/client/transport"
	"github.com/ChefBingbong/viem-go/types"
)

// PublicClientConfig contains configuration for creating a public client.
// It mirrors viem's PublicClientConfig, picking relevant fields from ClientConfig.
type PublicClientConfig struct {
	// Batch contains batch settings.
	Batch *BatchOptions
	// CacheTime is the time (in ms) that cached data will remain in memory.
	CacheTime time.Duration
	// Chain is the chain configuration.
	Chain *chain.Chain
	// ExperimentalBlockTag is the default block tag for RPC requests.
	ExperimentalBlockTag BlockTag
	// Key is a key for the client (default: "public").
	Key string
	// Name is a name for the client (default: "Public Client").
	Name string
	// PollingInterval is the frequency (in ms) for polling enabled actions & events.
	PollingInterval time.Duration
	// Transport is the transport factory to use.
	Transport transport.TransportFactory
}

// PublicClient is a client with public (read) actions.
// It wraps BaseClient and provides typed methods for public JSON-RPC calls.
// This mirrors viem's createPublicClient.
type PublicClient struct {
	*BaseClient
}

// CreatePublicClient creates a new public client with the given configuration.
// A Public Client is an interface to "public" JSON-RPC API methods such as
// retrieving block numbers, transactions, reading from smart contracts, etc.
//
// Example:
//
//	client, err := CreatePublicClient(PublicClientConfig{
//	    Chain:     mainnet,
//	    Transport: transport.HTTP("https://eth.merkle.io"),
//	})
func CreatePublicClient(config PublicClientConfig) (*PublicClient, error) {
	// Set defaults
	key := config.Key
	if key == "" {
		key = "public"
	}
	name := config.Name
	if name == "" {
		name = "Public Client"
	}

	// Create the base client
	baseConfig := ClientConfig{
		Batch:                config.Batch,
		CacheTime:            config.CacheTime,
		Chain:                config.Chain,
		ExperimentalBlockTag: config.ExperimentalBlockTag,
		Key:                  key,
		Name:                 name,
		PollingInterval:      config.PollingInterval,
		Transport:            config.Transport,
		Type:                 "publicClient",
	}

	base, err := CreateClient(baseConfig)
	if err != nil {
		return nil, err
	}

	return &PublicClient{BaseClient: base}, nil
}

// ---- Public Actions (Read Methods) ----

// GetBlockNumber returns the current block number.
func (c *PublicClient) GetBlockNumber(ctx context.Context) (uint64, error) {
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
func (c *PublicClient) GetChainID(ctx context.Context) (uint64, error) {
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

// GetGasPrice returns the current gas price in wei.
func (c *PublicClient) GetGasPrice(ctx context.Context) (*big.Int, error) {
	resp, err := c.Request(ctx, "eth_gasPrice")
	if err != nil {
		return nil, err
	}

	var hexGasPrice string
	if err := json.Unmarshal(resp.Result, &hexGasPrice); err != nil {
		return nil, err
	}

	return parseHexBigInt(hexGasPrice)
}

// GetBalance returns the balance of an address in wei.
// This delegates to the standalone public.GetBalance action.
func (c *PublicClient) GetBalance(ctx context.Context, address common.Address, blockTag ...BlockTag) (*big.Int, error) {
	params := public.GetBalanceParameters{Address: address}
	if len(blockTag) > 0 {
		params.BlockTag = blockTag[0]
	}
	return public.GetBalance(ctx, c, params)
}

// GetTransactionCount returns the nonce for an address.
func (c *PublicClient) GetTransactionCount(ctx context.Context, address common.Address, blockTag ...BlockTag) (uint64, error) {
	tag := c.resolveBlockTag(blockTag)

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

// GetCode returns the bytecode at an address.
func (c *PublicClient) GetCode(ctx context.Context, address common.Address, blockTag ...BlockTag) ([]byte, error) {
	tag := c.resolveBlockTag(blockTag)

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

// GetStorageAt returns the value at a storage position.
func (c *PublicClient) GetStorageAt(ctx context.Context, address common.Address, slot common.Hash, blockTag ...BlockTag) ([]byte, error) {
	tag := c.resolveBlockTag(blockTag)

	resp, err := c.Request(ctx, "eth_getStorageAt", address.Hex(), slot.Hex(), string(tag))
	if err != nil {
		return nil, err
	}

	var hexData string
	if err := json.Unmarshal(resp.Result, &hexData); err != nil {
		return nil, err
	}

	return parseHexBytes(hexData)
}

// CallRequest represents the parameters for an eth_call request.
type CallRequest = types.CallRequest

// Call performs an eth_call.
// This delegates to the standalone public.Call action.
func (c *PublicClient) Call(ctx context.Context, call CallRequest, blockTag ...BlockTag) ([]byte, error) {
	to := call.To // copy to get addressable value
	params := public.CallParameters{
		Account:  call.From,
		To:       &to,
		Data:     call.Data,
		Value:    call.Value,
		GasPrice: call.GasPrice,
	}
	if call.Gas > 0 {
		gas := call.Gas
		params.Gas = &gas
	}
	if len(blockTag) > 0 {
		params.BlockTag = blockTag[0]
	}
	result, err := public.Call(ctx, c, params)
	if err != nil {
		return nil, err
	}
	return result.Data, nil
}

// EstimateGas estimates gas for a call.
func (c *PublicClient) EstimateGas(ctx context.Context, call CallRequest) (uint64, error) {
	resp, err := c.Request(ctx, "eth_estimateGas", call)
	if err != nil {
		return 0, err
	}

	var hexGas string
	if err := json.Unmarshal(resp.Result, &hexGas); err != nil {
		return 0, err
	}

	return parseHexUint64(hexGas)
}

// GetBlock returns a block by tag.
// This delegates to the standalone public.GetBlock action.
func (c *PublicClient) GetBlock(ctx context.Context, blockTag BlockTag, includeTransactions bool) (*types.Block, error) {
	block, err := public.GetBlock(ctx, c, public.GetBlockParameters{
		BlockTag:            blockTag,
		IncludeTransactions: includeTransactions,
	})
	if err != nil {
		// Convert BlockNotFoundError to nil for backward compatibility
		if _, ok := err.(*public.BlockNotFoundError); ok {
			return nil, nil
		}
		return nil, err
	}
	return block, nil
}

// GetBlockByNumber returns a block by number.
// This delegates to the standalone public.GetBlock action.
func (c *PublicClient) GetBlockByNumber(ctx context.Context, blockNumber uint64, includeTransactions bool) (*types.Block, error) {
	block, err := public.GetBlock(ctx, c, public.GetBlockParameters{
		BlockNumber:         &blockNumber,
		IncludeTransactions: includeTransactions,
	})
	if err != nil {
		// Convert BlockNotFoundError to nil for backward compatibility
		if _, ok := err.(*public.BlockNotFoundError); ok {
			return nil, nil
		}
		return nil, err
	}
	return block, nil
}

// GetBlockByHash returns a block by hash.
// This delegates to the standalone public.GetBlock action.
func (c *PublicClient) GetBlockByHash(ctx context.Context, hash common.Hash, includeTransactions bool) (*types.Block, error) {
	block, err := public.GetBlock(ctx, c, public.GetBlockParameters{
		BlockHash:           &hash,
		IncludeTransactions: includeTransactions,
	})
	if err != nil {
		// Convert BlockNotFoundError to nil for backward compatibility
		if _, ok := err.(*public.BlockNotFoundError); ok {
			return nil, nil
		}
		return nil, err
	}
	return block, nil
}

// GetTransaction returns a transaction by hash.
// This delegates to the standalone public.GetTransaction action.
func (c *PublicClient) GetTransaction(ctx context.Context, hash common.Hash) (*public.TransactionResponse, error) {
	tx, err := public.GetTransaction(ctx, c, public.GetTransactionParameters{
		Hash: &hash,
	})
	if err != nil {
		// Convert TransactionNotFoundError to nil for backward compatibility
		if _, ok := err.(*public.TransactionNotFoundError); ok {
			return nil, nil
		}
		return nil, err
	}
	return tx, nil
}

// GetTransactionReceipt returns a transaction receipt.
func (c *PublicClient) GetTransactionReceipt(ctx context.Context, hash common.Hash) (*types.Receipt, error) {
	resp, err := c.Request(ctx, "eth_getTransactionReceipt", hash.Hex())
	if err != nil {
		return nil, err
	}

	if resp.Result == nil || string(resp.Result) == "null" {
		return nil, nil
	}

	var receipt types.Receipt
	if err := json.Unmarshal(resp.Result, &receipt); err != nil {
		return nil, err
	}

	return &receipt, nil
}

// FilterQuery represents parameters for eth_getLogs.
type FilterQuery = types.FilterQuery

// GetLogs returns logs matching the filter.
func (c *PublicClient) GetLogs(ctx context.Context, filter FilterQuery) ([]types.Log, error) {
	resp, err := c.Request(ctx, "eth_getLogs", filter)
	if err != nil {
		return nil, err
	}

	var logs []types.Log
	if err := json.Unmarshal(resp.Result, &logs); err != nil {
		return nil, err
	}

	return logs, nil
}

// GetFeeHistory returns fee history.
func (c *PublicClient) GetFeeHistory(ctx context.Context, blockCount uint64, newestBlock BlockTag, rewardPercentiles []float64) (json.RawMessage, error) {
	resp, err := c.Request(ctx, "eth_feeHistory", fmt.Sprintf("0x%x", blockCount), string(newestBlock), rewardPercentiles)
	if err != nil {
		return nil, err
	}

	return resp.Result, nil
}

// GetMaxPriorityFeePerGas returns the max priority fee per gas.
func (c *PublicClient) GetMaxPriorityFeePerGas(ctx context.Context) (*big.Int, error) {
	resp, err := c.Request(ctx, "eth_maxPriorityFeePerGas")
	if err != nil {
		return nil, err
	}

	var hexFee string
	if err := json.Unmarshal(resp.Result, &hexFee); err != nil {
		return nil, err
	}

	return parseHexBigInt(hexFee)
}

// GetProof returns the account and storage values with Merkle proof.
func (c *PublicClient) GetProof(ctx context.Context, address common.Address, storageKeys []common.Hash, blockTag ...BlockTag) (json.RawMessage, error) {
	tag := c.resolveBlockTag(blockTag)

	keys := make([]string, len(storageKeys))
	for i, k := range storageKeys {
		keys[i] = k.Hex()
	}

	resp, err := c.Request(ctx, "eth_getProof", address.Hex(), keys, string(tag))
	if err != nil {
		return nil, err
	}

	return resp.Result, nil
}

// WaitForTransactionReceipt waits for a transaction to be mined and returns its receipt.
func (c *PublicClient) WaitForTransactionReceipt(ctx context.Context, hash common.Hash) (*types.Receipt, error) {
	ticker := time.NewTicker(c.pollingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			receipt, err := c.GetTransactionReceipt(ctx, hash)
			if err != nil {
				return nil, err
			}
			if receipt != nil {
				return receipt, nil
			}
		}
	}
}

// ---- Watch Actions ----

// TransportType returns the type of transport being used.
// Implements the WatchClient interface.
func (c *PublicClient) TransportType() string {
	return c.transport.Config().Type
}

// Subscribe creates a WebSocket subscription.
// Implements the WatchClient interface.
// Returns ErrSubscriptionNotSupported if the transport doesn't support subscriptions.
func (c *PublicClient) Subscribe(
	params transport.SubscribeParams,
	onData func(data json.RawMessage),
	onError func(err error),
) (*transport.Subscription, error) {
	// Check if transport supports subscriptions
	if wsTransport, ok := c.transport.(*transport.WebSocketTransport); ok {
		return wsTransport.Subscribe(params, onData, onError)
	}
	return nil, public.ErrSubscriptionNotSupported
}

// WatchBlockNumber watches and returns incoming block numbers.
// Returns a channel that receives block number events.
// Close the context to stop watching.
//
// Example:
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	events := client.WatchBlockNumber(ctx, public.WatchBlockNumberParameters{
//	    EmitOnBegin: true,
//	})
//
//	for event := range events {
//	    fmt.Printf("Block: %d\n", event.BlockNumber)
//	}
func (c *PublicClient) WatchBlockNumber(ctx context.Context, params public.WatchBlockNumberParameters) <-chan public.WatchBlockNumberEvent {
	return public.WatchBlockNumber(ctx, c, params)
}

// WatchBlocks watches and returns incoming blocks.
// Returns a channel that receives block events.
// Close the context to stop watching.
func (c *PublicClient) WatchBlocks(ctx context.Context, params public.WatchBlocksParameters) <-chan public.WatchBlocksEvent {
	return public.WatchBlocks(ctx, c, params)
}

// WatchPendingTransactions watches and returns pending transaction hashes.
// Returns a channel that receives pending transaction events.
// Close the context to stop watching.
func (c *PublicClient) WatchPendingTransactions(ctx context.Context, params public.WatchPendingTransactionsParameters) <-chan public.WatchPendingTransactionsEvent {
	return public.WatchPendingTransactions(ctx, c, params)
}

// WatchEvent watches and returns emitted event logs.
// Returns a channel that receives event log events.
// Close the context to stop watching.
func (c *PublicClient) WatchEvent(ctx context.Context, params public.WatchEventParameters) <-chan public.WatchEventEvent {
	return public.WatchEvent(ctx, c, params)
}

// WatchContractEvent watches and returns emitted contract event logs with ABI decoding.
// Returns a channel that receives decoded event log events.
// Close the context to stop watching.
//
// Example:
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	events := client.WatchContractEvent(ctx, public.WatchContractEventParameters{
//	    Address:   contractAddress,
//	    ABI:       erc20ABI,
//	    EventName: "Transfer",
//	})
//
//	for event := range events {
//	    for _, log := range event.Logs {
//	        fmt.Printf("Transfer: %v\n", log.Args)
//	    }
//	}
func (c *PublicClient) WatchContractEvent(ctx context.Context, params public.WatchContractEventParameters) <-chan public.WatchContractEventEvent {
	return public.WatchContractEvent(ctx, c, params)
}

// ---- Contract Read/Write shortcuts that use ABI ----

// ReadContractWithABI reads a contract function using an ABI.
func (c *PublicClient) ReadContractWithABI(ctx context.Context, address common.Address, abi *viemabi.ABI, functionName string, args ...any) ([]any, error) {
	data, err := abi.EncodeFunctionData(functionName, args...)
	if err != nil {
		return nil, err
	}

	result, err := c.Call(ctx, CallRequest{
		To:   address,
		Data: data,
	})
	if err != nil {
		return nil, err
	}

	return abi.DecodeFunctionResult(functionName, result)
}

// ---- Helper methods ----

// resolveBlockTag returns the block tag to use, considering experimental block tag.
func (c *PublicClient) resolveBlockTag(tags []BlockTag) BlockTag {
	if len(tags) > 0 {
		return tags[0]
	}
	if c.experimentalBlockTag != "" {
		return c.experimentalBlockTag
	}
	return BlockTagLatest
}

// ---- Parsing helpers ----

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

// parseHexBigInt parses a hex string to *big.Int.
func parseHexBigInt(hexStr string) (*big.Int, error) {
	if len(hexStr) >= 2 && hexStr[:2] == "0x" {
		hexStr = hexStr[2:]
	}
	if hexStr == "" {
		return big.NewInt(0), nil
	}

	n := new(big.Int)
	_, ok := n.SetString(hexStr, 16)
	if !ok {
		return nil, fmt.Errorf("invalid hex string: %s", hexStr)
	}
	return n, nil
}
