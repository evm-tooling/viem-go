package accounts

import (
	"fmt"
	"regexp"

	"github.com/ChefBingbong/viem-go/utils/signature"
	"github.com/ChefBingbong/viem-go/utils/transaction"
)

var addressRegex = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)

// CustomSource represents a custom account source with signing implementations.
type CustomSource struct {
	// Address is the account's Ethereum address.
	Address string
	// Sign signs a hash (optional).
	Sign SignHashFunc
	// SignMessage signs a message.
	SignMessage SignMessageFunc
	// SignTransaction signs a transaction.
	SignTransaction SignTransactionFunc
	// SignTypedData signs EIP-712 typed data.
	SignTypedData SignTypedDataFunc
	// SignAuthorization signs an EIP-7702 authorization (optional).
	SignAuthorization SignAuthorizationFunc
}

// ToAccountFromAddress creates a JSON-RPC account from an address string.
//
// Example:
//
//	account, err := ToAccountFromAddress("0x1234567890123456789012345678901234567890")
func ToAccountFromAddress(address string) (*JsonRpcAccount, error) {
	if !isValidAddress(address) {
		return nil, fmt.Errorf("%w: %s", ErrInvalidAddress, address)
	}

	return &JsonRpcAccount{
		Address: address,
		Type:    AccountTypeJSONRPC,
	}, nil
}

// ToAccount creates a local account from a custom source with signing implementations.
//
// Example:
//
//	account, err := ToAccount(CustomSource{
//		Address: "0x1234567890123456789012345678901234567890",
//		SignMessage: func(message signature.SignableMessage) (string, error) {
//			// Custom signing implementation
//			return "0x...", nil
//		},
//		SignTransaction: func(tx *transaction.Transaction) (string, error) {
//			// Custom signing implementation
//			return "0x...", nil
//		},
//		SignTypedData: func(data signature.TypedDataDefinition) (string, error) {
//			// Custom signing implementation
//			return "0x...", nil
//		},
//	})
func ToAccount(source CustomSource) (*LocalAccount, error) {
	if !isValidAddress(source.Address) {
		return nil, fmt.Errorf("%w: %s", ErrInvalidAddress, source.Address)
	}

	return &LocalAccount{
		Addr:              source.Address,
		Source:            AccountSourceCustom,
		Type:              AccountTypeLocal,
		sign:              source.Sign,
		signMessage:       source.SignMessage,
		signTransaction:   source.SignTransaction,
		signTypedData:     source.SignTypedData,
		signAuthorization: source.SignAuthorization,
	}, nil
}

// ToAccountGeneric creates an account from either an address string or a CustomSource.
// This mirrors the TypeScript overloaded function behavior.
//
// Example:
//
//	// From address string - returns JsonRpcAccount
//	account, err := ToAccountGeneric("0x1234...")
//
//	// From CustomSource - returns LocalAccount
//	account, err := ToAccountGeneric(CustomSource{...})
func ToAccountGeneric(source any) (Account, error) {
	switch v := source.(type) {
	case string:
		return ToAccountFromAddress(v)
	case CustomSource:
		return ToAccount(v)
	case *CustomSource:
		if v != nil {
			return ToAccount(*v)
		}
		return nil, fmt.Errorf("%w: nil source", ErrInvalidAddress)
	default:
		return nil, fmt.Errorf("%w: unsupported source type", ErrInvalidAddress)
	}
}

// isValidAddress checks if a string is a valid Ethereum address format.
func isValidAddress(address string) bool {
	return addressRegex.MatchString(address)
}

// Helper to create a signing account with all methods wired up
func createLocalAccount(
	address string,
	publicKey string,
	source AccountSource,
	sign SignHashFunc,
	signMessage SignMessageFunc,
	signTransaction SignTransactionFunc,
	signTypedData SignTypedDataFunc,
	signAuthorization SignAuthorizationFunc,
) *LocalAccount {
	return &LocalAccount{
		Addr:              address,
		PublicKey:         publicKey,
		Source:            source,
		Type:              AccountTypeLocal,
		sign:              sign,
		signMessage:       signMessage,
		signTransaction:   signTransaction,
		signTypedData:     signTypedData,
		signAuthorization: signAuthorization,
	}
}

// Ensure we use the imports
var (
	_ signature.SignableMessage
	_ *transaction.Transaction
)
