package accounts

import (
	"github.com/ChefBingbong/viem-go/types"
	"github.com/ChefBingbong/viem-go/utils/signature"
	"github.com/ChefBingbong/viem-go/utils/transaction"
)

// Re-export types from the types package for convenience
type (
	AccountType          = types.AccountType
	AccountSource        = types.AccountSource
	Account              = types.Account
	JsonRpcAccount       = types.JsonRpcAccount
	AuthorizationRequest = types.AuthorizationRequest
	SignedAuthorization  = types.SignedAuthorization
	HDKey                = types.HDKey
	HDOptions            = types.HDOptions
)

// Re-export constants
const (
	AccountTypeLocal   = types.AccountTypeLocal
	AccountTypeJSONRPC = types.AccountTypeJSONRPC

	AccountSourcePrivateKey = types.AccountSourcePrivateKey
	AccountSourceHD         = types.AccountSourceHD
	AccountSourceMnemonic   = types.AccountSourceMnemonic
	AccountSourceCustom     = types.AccountSourceCustom
)

// Re-export functions
var DefaultHDPath = types.DefaultHDPath

// SignMessageFunc is the function signature for signing messages.
type SignMessageFunc func(message signature.SignableMessage) (string, error)

// SignTransactionFunc is the function signature for signing transactions.
type SignTransactionFunc func(tx *transaction.Transaction) (string, error)

// SignTypedDataFunc is the function signature for signing typed data.
type SignTypedDataFunc func(data signature.TypedDataDefinition) (string, error)

// SignHashFunc is the function signature for signing a hash.
type SignHashFunc func(hash string) (string, error)

// SignAuthorizationFunc is the function signature for signing EIP-7702 authorizations.
type SignAuthorizationFunc func(auth AuthorizationRequest) (*SignedAuthorization, error)

// LocalAccount represents a local account with signing capabilities.
type LocalAccount struct {
	Address   string        `json:"address"`
	PublicKey string        `json:"publicKey"`
	Source    AccountSource `json:"source"`
	Type      AccountType   `json:"type"`

	// Signing functions
	sign              SignHashFunc
	signMessage       SignMessageFunc
	signTransaction   SignTransactionFunc
	signTypedData     SignTypedDataFunc
	signAuthorization SignAuthorizationFunc
}

// GetAddress returns the account's address.
func (a *LocalAccount) GetAddress() string { return a.Address }

// GetType returns the account type.
func (a *LocalAccount) GetType() AccountType { return a.Type }

// GetPublicKey returns the account's public key.
func (a *LocalAccount) GetPublicKey() string { return a.PublicKey }

// GetSource returns the account's source.
func (a *LocalAccount) GetSource() AccountSource { return a.Source }

// Sign signs a hash and returns the signature as hex.
func (a *LocalAccount) Sign(hash string) (string, error) {
	if a.sign == nil {
		return "", ErrSigningNotSupported
	}
	return a.sign(hash)
}

// SignMessage signs a message and returns the signature as hex.
func (a *LocalAccount) SignMessage(message signature.SignableMessage) (string, error) {
	if a.signMessage == nil {
		return "", ErrSigningNotSupported
	}
	return a.signMessage(message)
}

// SignTransaction signs a transaction and returns the serialized signed transaction.
func (a *LocalAccount) SignTransaction(tx *transaction.Transaction) (string, error) {
	if a.signTransaction == nil {
		return "", ErrSigningNotSupported
	}
	return a.signTransaction(tx)
}

// SignTypedData signs EIP-712 typed data and returns the signature as hex.
func (a *LocalAccount) SignTypedData(data signature.TypedDataDefinition) (string, error) {
	if a.signTypedData == nil {
		return "", ErrSigningNotSupported
	}
	return a.signTypedData(data)
}

// SignAuthorization signs an EIP-7702 authorization.
func (a *LocalAccount) SignAuthorization(auth AuthorizationRequest) (*SignedAuthorization, error) {
	if a.signAuthorization == nil {
		return nil, ErrSigningNotSupported
	}
	return a.signAuthorization(auth)
}

// PrivateKeyAccount extends LocalAccount for accounts created from private keys.
type PrivateKeyAccount struct {
	*LocalAccount
}

// HDAccount extends LocalAccount for HD wallet accounts.
type HDAccount struct {
	*LocalAccount
	hdKey HDKey
}

// GetHdKey returns the underlying HD key.
func (a *HDAccount) GetHdKey() HDKey {
	return a.hdKey
}
