package utils

import "math/big"

// AccountType represents the type of account.
type AccountType string

const (
	// AccountTypeLocal represents a local account with a private key.
	AccountTypeLocal AccountType = "local"
	// AccountTypeJSONRPC represents a JSON-RPC account (address only).
	AccountTypeJSONRPC AccountType = "json-rpc"
)

// Account represents an Ethereum account.
type Account struct {
	// Address is the account's Ethereum address.
	Address string `json:"address"`
	// Type indicates the account type (local or json-rpc).
	Type AccountType `json:"type"`
}

// LocalAccount represents a local account with signing capabilities.
type LocalAccount struct {
	Account
	// PublicKey is the account's public key (hex encoded).
	PublicKey string `json:"publicKey,omitempty"`
}

// Signature represents an ECDSA signature.
type Signature struct {
	// R component of the signature (32 bytes as hex string).
	R string `json:"r"`
	// S component of the signature (32 bytes as hex string).
	S string `json:"s"`
	// V value (27 or 28 for legacy, optional for EIP-2930/1559).
	V *big.Int `json:"v,omitempty"`
	// YParity is the parity of the y-coordinate of the curve point (0 or 1).
	YParity int `json:"yParity"`
}

// SignedAuthorization represents an EIP-7702 signed authorization.
type SignedAuthorization struct {
	// Address is the contract address being authorized.
	Address string `json:"address"`
	// ChainId is the chain ID for this authorization.
	ChainId int `json:"chainId"`
	// Nonce is the account nonce for this authorization.
	Nonce int `json:"nonce"`
	// R component of the signature.
	R string `json:"r"`
	// S component of the signature.
	S string `json:"s"`
	// V value.
	V *big.Int `json:"v,omitempty"`
	// YParity is the parity of the y-coordinate.
	YParity int `json:"yParity"`
}

// SignReturnFormat specifies the output format for signing operations.
type SignReturnFormat string

const (
	// SignReturnFormatObject returns a Signature struct.
	SignReturnFormatObject SignReturnFormat = "object"
	// SignReturnFormatHex returns a hex string.
	SignReturnFormatHex SignReturnFormat = "hex"
	// SignReturnFormatBytes returns raw bytes.
	SignReturnFormatBytes SignReturnFormat = "bytes"
)
