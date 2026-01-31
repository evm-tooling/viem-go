package client

import (
	"context"
	"encoding/json"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/chain"
	"github.com/ChefBingbong/viem-go/client/transport"
	"github.com/ChefBingbong/viem-go/types"
)

// WalletClientConfig contains configuration for creating a wallet client.
// It mirrors viem's WalletClientConfig, picking relevant fields from ClientConfig.
type WalletClientConfig struct {
	// Account is the account to use for the client.
	Account Account
	// Chain is the chain configuration.
	Chain *chain.Chain
	// CacheTime is the time (in ms) that cached data will remain in memory.
	CacheTime time.Duration
	// Key is a key for the client (default: "wallet").
	Key string
	// Name is a name for the client (default: "Wallet Client").
	Name string
	// PollingInterval is the frequency (in ms) for polling enabled actions & events.
	PollingInterval time.Duration
	// Transport is the transport factory to use.
	Transport transport.TransportFactory
}

// WalletClient is a client with wallet (write) actions.
// It wraps BaseClient and provides typed methods for wallet JSON-RPC calls.
// This mirrors viem's createWalletClient.
type WalletClient struct {
	*BaseClient
}

// CreateWalletClient creates a new wallet client with the given configuration.
// A Wallet Client is an interface to wallet JSON-RPC API methods such as
// sending transactions, signing messages, etc.
//
// Example:
//
//	client, err := CreateWalletClient(WalletClientConfig{
//	    Account:   account,
//	    Chain:     mainnet,
//	    Transport: transport.HTTP("https://eth.merkle.io"),
//	})
func CreateWalletClient(config WalletClientConfig) (*WalletClient, error) {
	// Set defaults
	key := config.Key
	if key == "" {
		key = "wallet"
	}
	name := config.Name
	if name == "" {
		name = "Wallet Client"
	}

	// Create the base client
	baseConfig := ClientConfig{
		Account:         config.Account,
		CacheTime:       config.CacheTime,
		Chain:           config.Chain,
		Key:             key,
		Name:            name,
		PollingInterval: config.PollingInterval,
		Transport:       config.Transport,
		Type:            "walletClient",
	}

	base, err := CreateClient(baseConfig)
	if err != nil {
		return nil, err
	}

	return &WalletClient{BaseClient: base}, nil
}

// ---- Wallet Actions (Write Methods) ----

// SendRawTransaction sends a signed raw transaction.
func (c *WalletClient) SendRawTransaction(ctx context.Context, signedTx []byte) (common.Hash, error) {
	hexTx := "0x" + common.Bytes2Hex(signedTx)

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

// SendTransaction sends a transaction using eth_sendTransaction.
// Requires the node to manage the account (e.g., via personal_unlockAccount).
func (c *WalletClient) SendTransaction(ctx context.Context, tx types.Transaction) (common.Hash, error) {
	resp, err := c.Request(ctx, "eth_sendTransaction", tx)
	if err != nil {
		return common.Hash{}, err
	}

	var hashHex string
	if err := json.Unmarshal(resp.Result, &hashHex); err != nil {
		return common.Hash{}, err
	}

	return common.HexToHash(hashHex), nil
}

// SignMessage signs a message using eth_sign.
func (c *WalletClient) SignMessage(ctx context.Context, address common.Address, message []byte) ([]byte, error) {
	hexMessage := "0x" + common.Bytes2Hex(message)

	resp, err := c.Request(ctx, "eth_sign", address.Hex(), hexMessage)
	if err != nil {
		return nil, err
	}

	var hexSig string
	if err := json.Unmarshal(resp.Result, &hexSig); err != nil {
		return nil, err
	}

	return parseHexBytes(hexSig)
}

// SignTypedData signs typed data using eth_signTypedData_v4.
func (c *WalletClient) SignTypedData(ctx context.Context, address common.Address, typedData any) ([]byte, error) {
	resp, err := c.Request(ctx, "eth_signTypedData_v4", address.Hex(), typedData)
	if err != nil {
		return nil, err
	}

	var hexSig string
	if err := json.Unmarshal(resp.Result, &hexSig); err != nil {
		return nil, err
	}

	return parseHexBytes(hexSig)
}

// GetAccounts returns the accounts managed by the wallet.
func (c *WalletClient) GetAccounts(ctx context.Context) ([]common.Address, error) {
	resp, err := c.Request(ctx, "eth_accounts")
	if err != nil {
		return nil, err
	}

	var hexAddresses []string
	if err := json.Unmarshal(resp.Result, &hexAddresses); err != nil {
		return nil, err
	}

	addresses := make([]common.Address, len(hexAddresses))
	for i, hex := range hexAddresses {
		addresses[i] = common.HexToAddress(hex)
	}

	return addresses, nil
}

// RequestAccounts requests accounts using eth_requestAccounts.
func (c *WalletClient) RequestAccounts(ctx context.Context) ([]common.Address, error) {
	resp, err := c.Request(ctx, "eth_requestAccounts")
	if err != nil {
		return nil, err
	}

	var hexAddresses []string
	if err := json.Unmarshal(resp.Result, &hexAddresses); err != nil {
		return nil, err
	}

	addresses := make([]common.Address, len(hexAddresses))
	for i, hex := range hexAddresses {
		addresses[i] = common.HexToAddress(hex)
	}

	return addresses, nil
}

// SwitchChain switches to a different chain using wallet_switchEthereumChain.
func (c *WalletClient) SwitchChain(ctx context.Context, chainID *big.Int) error {
	params := map[string]string{
		"chainId": "0x" + chainID.Text(16),
	}

	_, err := c.Request(ctx, "wallet_switchEthereumChain", params)
	return err
}

// AddChain adds a chain using wallet_addEthereumChain.
func (c *WalletClient) AddChain(ctx context.Context, ch *chain.Chain) error {
	params := map[string]any{
		"chainId":   "0x" + big.NewInt(ch.ID).Text(16),
		"chainName": ch.Name,
	}

	params["nativeCurrency"] = map[string]any{
		"name":     ch.NativeCurrency.Name,
		"symbol":   ch.NativeCurrency.Symbol,
		"decimals": ch.NativeCurrency.Decimals,
	}

	if urls, ok := ch.RpcUrls["default"]; ok && len(urls.HTTP) > 0 {
		params["rpcUrls"] = urls.HTTP
	}

	_, err := c.Request(ctx, "wallet_addEthereumChain", params)
	return err
}

// WatchAsset adds a token to the wallet's asset list using wallet_watchAsset.
func (c *WalletClient) WatchAsset(ctx context.Context, tokenType string, address common.Address, symbol string, decimals uint8, image string) (bool, error) {
	params := map[string]any{
		"type": tokenType,
		"options": map[string]any{
			"address":  address.Hex(),
			"symbol":   symbol,
			"decimals": decimals,
		},
	}

	if image != "" {
		params["options"].(map[string]any)["image"] = image
	}

	resp, err := c.Request(ctx, "wallet_watchAsset", params)
	if err != nil {
		return false, err
	}

	var success bool
	if err := json.Unmarshal(resp.Result, &success); err != nil {
		return false, err
	}

	return success, nil
}

// GetPermissions returns the wallet's permissions using wallet_getPermissions.
func (c *WalletClient) GetPermissions(ctx context.Context) (json.RawMessage, error) {
	resp, err := c.Request(ctx, "wallet_getPermissions")
	if err != nil {
		return nil, err
	}

	return resp.Result, nil
}

// RequestPermissions requests permissions using wallet_requestPermissions.
func (c *WalletClient) RequestPermissions(ctx context.Context, permissions []map[string]any) (json.RawMessage, error) {
	resp, err := c.Request(ctx, "wallet_requestPermissions", permissions)
	if err != nil {
		return nil, err
	}

	return resp.Result, nil
}
