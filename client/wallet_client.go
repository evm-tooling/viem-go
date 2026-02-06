package client

import (
	"context"
	"math/big"
	"time"

	"github.com/ChefBingbong/viem-go/actions/wallet"
	"github.com/ChefBingbong/viem-go/chain"
	"github.com/ChefBingbong/viem-go/client/transport"
	"github.com/ChefBingbong/viem-go/types"
	"github.com/ChefBingbong/viem-go/utils/address"
	"github.com/ChefBingbong/viem-go/utils/signature"
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
// It wraps BaseClient and delegates to the standalone action functions
// in actions/wallet, mirroring viem's createWalletClient.
//
// WalletClient implements wallet.Client so the standalone action functions
// can be used with it directly.
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

// ---------------------------------------------------------------------------
// wallet.Client interface implementation
// ---------------------------------------------------------------------------

// Account returns the account as a wallet.Account.
// This bridges client.Account -> wallet.Account so WalletClient satisfies wallet.Client.
func (c *WalletClient) Account() wallet.Account {
	acct := c.BaseClient.Account()
	if acct == nil {
		return nil
	}
	// client.Account and wallet.Account are structurally identical —
	// any concrete type satisfying client.Account also satisfies wallet.Account.
	if wa, ok := acct.(wallet.Account); ok {
		return wa
	}
	return nil
}

// ---------------------------------------------------------------------------
// Wallet Actions — Signing
// ---------------------------------------------------------------------------

// SignMessage calculates an Ethereum-specific EIP-191 signature.
// Delegates to wallet.SignMessage.
func (c *WalletClient) SignMessage(ctx context.Context, params wallet.SignMessageParameters) (string, error) {
	return wallet.SignMessage(ctx, c, params)
}

// SignTypedData signs EIP-712 typed structured data.
// Delegates to wallet.SignTypedData.
func (c *WalletClient) SignTypedData(ctx context.Context, params wallet.SignTypedDataParameters) (string, error) {
	return wallet.SignTypedData(ctx, c, params)
}

// SignTransaction signs a transaction without broadcasting.
// Delegates to wallet.SignTransaction.
func (c *WalletClient) SignTransaction(ctx context.Context, params wallet.SignTransactionParameters) (string, error) {
	return wallet.SignTransaction(ctx, c, params)
}

// SignPreparedTransaction converts a prepared transaction request into signing
// parameters, then signs it. This is a convenience method that bridges
// PrepareTransactionRequest -> SignTransaction in a single call.
//
// Example:
//
//	prepared, _ := client.PrepareTransactionRequest(ctx, wallet.PrepareTransactionRequestParameters{
//	    To: "0x...", Value: big.NewInt(1000),
//	})
//	signed, _ := client.SignPreparedTransaction(ctx, prepared)
func (c *WalletClient) SignPreparedTransaction(ctx context.Context, prepared wallet.PrepareTransactionRequestReturnType) (string, error) {
	return wallet.SignTransaction(ctx, c, wallet.PreparedToSignParams(prepared))
}

// SignAuthorization signs an EIP-7702 authorization.
// Delegates to wallet.SignAuthorization.
func (c *WalletClient) SignAuthorization(ctx context.Context, params wallet.SignAuthorizationParameters) (*types.SignedAuthorization, error) {
	return wallet.SignAuthorization(ctx, c, params)
}

// PrepareAuthorization prepares an EIP-7702 authorization for signing.
// Delegates to wallet.PrepareAuthorization.
func (c *WalletClient) PrepareAuthorization(ctx context.Context, params wallet.PrepareAuthorizationParameters) (wallet.PrepareAuthorizationReturnType, error) {
	return wallet.PrepareAuthorization(ctx, c, params)
}

// ---------------------------------------------------------------------------
// Wallet Actions — Transactions
// ---------------------------------------------------------------------------

// SendTransaction creates, signs, and sends a new transaction.
// Delegates to wallet.SendTransaction.
func (c *WalletClient) SendTransaction(ctx context.Context, params wallet.SendTransactionParameters) (string, error) {
	return wallet.SendTransaction(ctx, c, params)
}

// SendTransactionSync sends a transaction and waits for the receipt.
// Delegates to wallet.SendTransactionSync.
func (c *WalletClient) SendTransactionSync(ctx context.Context, params wallet.SendTransactionSyncParameters) (wallet.SendTransactionSyncReturnType, error) {
	return wallet.SendTransactionSync(ctx, c, params)
}

// SendRawTransaction sends a signed serialized transaction.
// Delegates to wallet.SendRawTransaction.
func (c *WalletClient) SendRawTransaction(ctx context.Context, params wallet.SendRawTransactionParameters) (string, error) {
	return wallet.SendRawTransaction(ctx, c, params)
}

// SendRawTransactionSync sends a signed transaction and waits for the receipt.
// Delegates to wallet.SendRawTransactionSync.
func (c *WalletClient) SendRawTransactionSync(ctx context.Context, params wallet.SendRawTransactionSyncParameters) (wallet.SendRawTransactionSyncReturnType, error) {
	return wallet.SendRawTransactionSync(ctx, c, params)
}

// PrepareTransactionRequest fills in missing transaction fields (nonce, gas, fees, type, chainId).
// Delegates to wallet.PrepareTransactionRequest.
func (c *WalletClient) PrepareTransactionRequest(ctx context.Context, params wallet.PrepareTransactionRequestParameters) (wallet.PrepareTransactionRequestReturnType, error) {
	return wallet.PrepareTransactionRequest(ctx, c, params)
}

// ---------------------------------------------------------------------------
// Wallet Actions — Contracts
// ---------------------------------------------------------------------------

// WriteContract executes a write function on a contract.
// Delegates to wallet.WriteContract.
func (c *WalletClient) WriteContract(ctx context.Context, params wallet.WriteContractParameters) (string, error) {
	return wallet.WriteContract(ctx, c, params)
}

// WriteContractSync executes a contract write and waits for the receipt.
// Delegates to wallet.WriteContractSync.
func (c *WalletClient) WriteContractSync(ctx context.Context, params wallet.WriteContractSyncParameters) (wallet.WriteContractSyncReturnType, error) {
	return wallet.WriteContractSync(ctx, c, params)
}

// DeployContract deploys a contract to the network.
// Delegates to wallet.DeployContract.
func (c *WalletClient) DeployContract(ctx context.Context, params wallet.DeployContractParameters) (string, error) {
	return wallet.DeployContract(ctx, c, params)
}

// ---------------------------------------------------------------------------
// Wallet Actions — Account Management
// ---------------------------------------------------------------------------

// GetAddresses returns a list of account addresses owned by the wallet.
// Delegates to wallet.GetAddresses.
func (c *WalletClient) GetAddresses(ctx context.Context) ([]address.Address, error) {
	return wallet.GetAddresses(ctx, c)
}

// RequestAddresses requests a list of accounts managed by the wallet.
// Delegates to wallet.RequestAddresses.
func (c *WalletClient) RequestAddresses(ctx context.Context) ([]address.Address, error) {
	return wallet.RequestAddresses(ctx, c)
}

// AddChain adds an EVM chain to the wallet.
// Delegates to wallet.AddChain.
func (c *WalletClient) AddChain(ctx context.Context, params wallet.AddChainParameters) error {
	return wallet.AddChain(ctx, c, params)
}

// SwitchChain switches the target chain in the wallet.
// Delegates to wallet.SwitchChain.
func (c *WalletClient) SwitchChain(ctx context.Context, params wallet.SwitchChainParameters) error {
	return wallet.SwitchChain(ctx, c, params)
}

// ---------------------------------------------------------------------------
// Wallet Actions — Permissions & Assets
// ---------------------------------------------------------------------------

// GetPermissions gets the wallet's current permissions.
// Delegates to wallet.GetPermissions.
func (c *WalletClient) GetPermissions(ctx context.Context) ([]wallet.WalletPermission, error) {
	return wallet.GetPermissions(ctx, c)
}

// RequestPermissions requests permissions for the wallet.
// Delegates to wallet.RequestPermissions.
func (c *WalletClient) RequestPermissions(ctx context.Context, permissions wallet.RequestPermissionsParameters) ([]wallet.WalletPermission, error) {
	return wallet.RequestPermissions(ctx, c, permissions)
}

// WatchAsset requests that the wallet tracks a specified token.
// Delegates to wallet.WatchAsset.
func (c *WalletClient) WatchAsset(ctx context.Context, params wallet.WatchAssetParameters) (bool, error) {
	return wallet.WatchAsset(ctx, c, params)
}

// ---------------------------------------------------------------------------
// Wallet Actions — EIP-5792 Batch Calls
// ---------------------------------------------------------------------------

// GetCapabilities extracts capabilities that a connected wallet supports.
// Delegates to wallet.GetCapabilities.
func (c *WalletClient) GetCapabilities(ctx context.Context, params wallet.GetCapabilitiesParameters) (wallet.GetCapabilitiesReturnType, error) {
	return wallet.GetCapabilities(ctx, c, params)
}

// SendCalls sends a batch of calls via wallet_sendCalls (EIP-5792).
// Delegates to wallet.SendCalls.
func (c *WalletClient) SendCalls(ctx context.Context, params wallet.SendCallsParameters) (*wallet.SendCallsReturnType, error) {
	return wallet.SendCalls(ctx, c, params)
}

// SendCallsSync sends a batch of calls and waits for inclusion.
// Delegates to wallet.SendCallsSync.
func (c *WalletClient) SendCallsSync(ctx context.Context, params wallet.SendCallsSyncParameters) (wallet.SendCallsSyncReturnType, error) {
	return wallet.SendCallsSync(ctx, c, params)
}

// GetCallsStatus returns the status of a call batch.
// Delegates to wallet.GetCallsStatus.
func (c *WalletClient) GetCallsStatus(ctx context.Context, params wallet.GetCallsStatusParameters) (*wallet.GetCallsStatusReturnType, error) {
	return wallet.GetCallsStatus(ctx, c, params)
}

// WaitForCallsStatus waits for a call batch to be confirmed.
// Delegates to wallet.WaitForCallsStatus.
func (c *WalletClient) WaitForCallsStatus(ctx context.Context, params wallet.WaitForCallsStatusParameters) (wallet.WaitForCallsStatusReturnType, error) {
	return wallet.WaitForCallsStatus(ctx, c, params)
}

// ShowCallsStatus requests the wallet to show call batch status.
// Delegates to wallet.ShowCallsStatus.
func (c *WalletClient) ShowCallsStatus(ctx context.Context, params wallet.ShowCallsStatusParameters) error {
	return wallet.ShowCallsStatus(ctx, c, params)
}

// ---------------------------------------------------------------------------
// Deprecated / backward-compat helpers
// ---------------------------------------------------------------------------

// SignMessageRaw signs a message using personal_sign RPC (JSON-RPC account path).
// Deprecated: Use SignMessage with a SignMessageParameters instead.
func (c *WalletClient) SignMessageRaw(ctx context.Context, params wallet.SignMessageParameters) (string, error) {
	return wallet.SignMessage(ctx, c, params)
}

// SignTypedDataRaw signs typed data using eth_signTypedData_v4 RPC.
// Deprecated: Use SignTypedData with SignTypedDataParameters instead.
func (c *WalletClient) SignTypedDataRaw(ctx context.Context, domain signature.TypedDataDomain, types map[string][]signature.TypedDataField, primaryType string, message map[string]any) (string, error) {
	return wallet.SignTypedData(ctx, c, wallet.SignTypedDataParameters{
		Domain:      domain,
		Types:       types,
		PrimaryType: primaryType,
		Message:     message,
	})
}

// SwitchChainByID switches chain by numeric ID.
// Deprecated: Use SwitchChain with SwitchChainParameters instead.
func (c *WalletClient) SwitchChainByID(ctx context.Context, chainID *big.Int) error {
	return wallet.SwitchChain(ctx, c, wallet.SwitchChainParameters{ID: chainID.Int64()})
}
