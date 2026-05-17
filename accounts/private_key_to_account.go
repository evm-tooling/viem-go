package accounts

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"

	accountUtils "github.com/ChefBingbong/viem-go/accounts/utils"
	"github.com/ChefBingbong/viem-go/utils/signature"
	"github.com/ChefBingbong/viem-go/utils/transaction"
)

// PrivateKeyToAccountOptions contains options for creating an account from a private key.
type PrivateKeyToAccountOptions struct {
	// Add any options here (e.g., nonceManager in the future)
}

// PrivateKeyToAccount creates an Account from a private key.
//
// The account can sign messages, transactions, typed data, and EIP-7702 authorizations.
//
// Example:
//
//	account, err := PrivateKeyToAccount("0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
//	// account.Address = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
//
//	// Sign a message
//	sig, err := account.SignMessage(signature.NewSignableMessage("hello world"))
func PrivateKeyToAccount(privateKey string, opts ...PrivateKeyToAccountOptions) (*PrivateKeyAccount, error) {
	// Validate and parse the private key
	keyHex := strings.TrimPrefix(privateKey, "0x")
	keyHex = strings.TrimPrefix(keyHex, "0X")

	ecdsaKey, err := crypto.HexToECDSA(keyHex)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidPrivateKey, err)
	}

	// Derive public key (uncompressed format)
	publicKeyBytes := crypto.FromECDSAPub(&ecdsaKey.PublicKey)
	publicKey := "0x" + hex.EncodeToString(publicKeyBytes)

	// Derive address
	address, err := accountUtils.PublicKeyToAddress(publicKey)
	if err != nil {
		return nil, err
	}

	// Create signing functions that capture the private key
	signFunc := func(hash string) (string, error) {
		return accountUtils.SignToHex(hash, privateKey)
	}

	signMessageFunc := func(message signature.SignableMessage) (string, error) {
		return accountUtils.SignMessage(accountUtils.SignMessageParameters{
			Message:    message,
			PrivateKey: privateKey,
		})
	}

	signTransactionFunc := func(tx *transaction.Transaction) (string, error) {
		return accountUtils.SignTransaction(accountUtils.SignTransactionParameters{
			PrivateKey:  privateKey,
			Transaction: tx,
		})
	}

	signTypedDataFunc := func(data signature.TypedDataDefinition) (string, error) {
		return accountUtils.SignTypedData(accountUtils.SignTypedDataParameters{
			Domain:      data.Domain,
			Types:       data.Types,
			PrimaryType: data.PrimaryType,
			Message:     data.Message,
			PrivateKey:  privateKey,
		})
	}

	signAuthorizationFunc := func(auth AuthorizationRequest) (*SignedAuthorization, error) {
		result, err := accountUtils.SignAuthorization(accountUtils.SignAuthorizationParameters{
			Address:    auth.GetAddress(),
			ChainId:    auth.ChainId,
			Nonce:      auth.Nonce,
			PrivateKey: privateKey,
		})
		if err != nil {
			return nil, err
		}

		return &SignedAuthorization{
			Address: result.SignedAuthorization.Address,
			ChainId: result.SignedAuthorization.ChainId,
			Nonce:   result.SignedAuthorization.Nonce,
			R:       result.SignedAuthorization.R,
			S:       result.SignedAuthorization.S,
			V:       result.SignedAuthorization.V,
			YParity: result.SignedAuthorization.YParity,
		}, nil
	}

	localAccount := createLocalAccount(
		address,
		publicKey,
		AccountSourcePrivateKey,
		signFunc,
		signMessageFunc,
		signTransactionFunc,
		signTypedDataFunc,
		signAuthorizationFunc,
	)

	return &PrivateKeyAccount{
		LocalAccount: localAccount,
	}, nil
}

// MustPrivateKeyToAccount creates an account from a private key or panics on error.
func MustPrivateKeyToAccount(privateKey string, opts ...PrivateKeyToAccountOptions) *PrivateKeyAccount {
	account, err := PrivateKeyToAccount(privateKey, opts...)
	if err != nil {
		panic(err)
	}
	return account
}
