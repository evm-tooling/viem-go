package accounts

import (
	"math/big"

	"github.com/ChefBingbong/viem-go/utils/signature"
	"github.com/ChefBingbong/viem-go/utils/transaction"
)

// AccountType represents the type of account.
type AccountType string

const (
	// AccountTypeLocal represents a local account with signing capabilities.
	AccountTypeLocal AccountType = "local"
	// AccountTypeJSONRPC represents a JSON-RPC account (address only).
	AccountTypeJSONRPC AccountType = "json-rpc"
)

// AccountSource represents the source of a local account.
type AccountSource string

const (
	// AccountSourcePrivateKey indicates the account was created from a private key.
	AccountSourcePrivateKey AccountSource = "privateKey"
	// AccountSourceHD indicates the account was created from an HD wallet.
	AccountSourceHD AccountSource = "hd"
	// AccountSourceMnemonic indicates the account was created from a mnemonic.
	AccountSourceMnemonic AccountSource = "mnemonic"
	// AccountSourceCustom indicates the account has custom signing implementations.
	AccountSourceCustom AccountSource = "custom"
)

// Account is the base interface for all account types.
type Account interface {
	// GetAddress returns the account's Ethereum address.
	GetAddress() string
	// GetType returns the account type.
	GetType() AccountType
}

// JsonRpcAccount represents an account that signs via JSON-RPC.
type JsonRpcAccount struct {
	Address string      `json:"address"`
	Type    AccountType `json:"type"`
}

// GetAddress returns the account's address.
func (a *JsonRpcAccount) GetAddress() string { return a.Address }

// GetType returns the account type.
func (a *JsonRpcAccount) GetType() AccountType { return a.Type }

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

// AuthorizationRequest represents an EIP-7702 authorization request.
type AuthorizationRequest struct {
	Address         string `json:"address,omitempty"`
	ContractAddress string `json:"contractAddress,omitempty"`
	ChainId         int    `json:"chainId"`
	Nonce           int    `json:"nonce"`
}

// GetAddress returns the address, preferring ContractAddress if set.
func (a *AuthorizationRequest) GetAddress() string {
	if a.ContractAddress != "" {
		return a.ContractAddress
	}
	return a.Address
}

// SignedAuthorization represents a signed EIP-7702 authorization.
type SignedAuthorization struct {
	Address string   `json:"address"`
	ChainId int      `json:"chainId"`
	Nonce   int      `json:"nonce"`
	R       string   `json:"r"`
	S       string   `json:"s"`
	V       *big.Int `json:"v,omitempty"`
	YParity int      `json:"yParity"`
}

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

// HDKey represents a hierarchical deterministic key.
type HDKey interface {
	// Derive derives a child key at the given path.
	Derive(path string) (HDKey, error)
	// PrivateKey returns the private key bytes.
	PrivateKey() []byte
	// PublicKey returns the public key bytes.
	PublicKey() []byte
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

// HDOptions contains options for deriving HD accounts.
type HDOptions struct {
	// AccountIndex is the account index in the path (m/44'/60'/{accountIndex}'/0/0).
	AccountIndex int
	// AddressIndex is the address index in the path (m/44'/60'/0'/0/{addressIndex}).
	AddressIndex int
	// ChangeIndex is the change index in the path (m/44'/60'/0'/{changeIndex}/0).
	ChangeIndex int
	// Path is a custom derivation path. If set, overrides the index options.
	Path string
}

// DefaultHDPath returns the default Ethereum HD derivation path.
func DefaultHDPath(accountIndex, changeIndex, addressIndex int) string {
	return "m/44'/60'/" + itoa(accountIndex) + "'/" + itoa(changeIndex) + "/" + itoa(addressIndex)
}

// itoa converts int to string without importing strconv.
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	if i < 0 {
		return "-" + uitoa(uint(-i))
	}
	return uitoa(uint(i))
}

func uitoa(u uint) string {
	var buf [20]byte
	i := len(buf)
	for u >= 10 {
		i--
		q := u / 10
		buf[i] = byte('0' + u - q*10)
		u = q
	}
	i--
	buf[i] = byte('0' + u)
	return string(buf[i:])
}
