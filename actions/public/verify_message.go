package public

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/utils/signature"
)

// VerifyMessageParameters contains the parameters for the VerifyMessage action.
// This mirrors viem's VerifyMessageParameters type with full feature support including:
//   - ERC-6492 (counterfactual signature verification for undeployed smart accounts)
//   - ERC-1271 (smart contract signature verification)
//   - ECDSA recovery fallback (for EOA accounts)
type VerifyMessageParameters struct {
	// Address is the address that is expected to have signed the message.
	Address common.Address

	// Message is the message that was signed.
	//
	// For convenience, use signature.NewSignableMessage, NewSignableMessageRaw,
	// or NewSignableMessageRawHex to construct this value.
	Message signature.SignableMessage

	// Signature is the signature produced by signing the message.
	//
	// Accepts:
	//   - string hex-encoded signature
	//   - []byte raw signature bytes
	//   - *signature.Signature (r, s, v, yParity)
	Signature any

	// BlockNumber is the block number to verify at.
	// Mutually exclusive with BlockTag.
	BlockNumber *uint64

	// BlockTag is the block tag to verify at (e.g., "latest", "pending").
	// Mutually exclusive with BlockNumber.
	BlockTag BlockTag

	// Factory is the ERC-4337 Account Factory address for counterfactual verification.
	// Used with FactoryData for undeployed smart accounts.
	Factory *common.Address

	// FactoryData is the calldata to deploy the account via Factory.
	// Used with Factory for undeployed smart accounts.
	FactoryData []byte

	// ERC6492VerifierAddress is the address of a deployed ERC-6492 signature verifier contract.
	// If provided, uses this contract instead of deployless verification.
	ERC6492VerifierAddress *common.Address
}

// VerifyMessageReturnType is the return type for the VerifyMessage action.
// It indicates whether the signature is valid.
type VerifyMessageReturnType = bool

// VerifyMessage verifies that a message was signed by the provided address.
//
// Compatible with Smart Contract Accounts & Externally Owned Accounts via ERC-6492.
// https://eips.ethereum.org/EIPS/eip-6492
//
// This is equivalent to viem's `verifyMessage` action with full feature support:
//   - ERC-6492 verification for counterfactual (undeployed) smart accounts
//   - ERC-1271 verification for deployed smart contracts
//   - ECDSA recovery fallback for EOA accounts
//
// Example:
//
//	valid, err := public.VerifyMessage(ctx, client, public.VerifyMessageParameters{
//	    Address: common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
//	    Message: signature.NewSignableMessage("hello world"),
//	    Signature: "0x...",
//	})
//
// Example with counterfactual verification:
//
//	valid, err := public.VerifyMessage(ctx, client, public.VerifyMessageParameters{
//	    Address:     common.HexToAddress("0x..."),
//	    Message:     signature.NewSignableMessage("hello"),
//	    Signature:   "0x...",
//	    Factory:     &factoryAddress,
//	    FactoryData: factoryCalldata,
//	})
func VerifyMessage(ctx context.Context, client Client, params VerifyMessageParameters) (VerifyMessageReturnType, error) {
	// Hash the message using Ethereum Signed Message prefix
	hash := signature.HashMessage(params.Message)

	// Delegate to VerifyHash with full verification support
	return VerifyHash(ctx, client, VerifyHashParameters{
		Address:                params.Address,
		Hash:                   hash,
		Signature:              params.Signature,
		BlockNumber:            params.BlockNumber,
		BlockTag:               params.BlockTag,
		Factory:                params.Factory,
		FactoryData:            params.FactoryData,
		ERC6492VerifierAddress: params.ERC6492VerifierAddress,
	})
}

// VerifyMessageLocal verifies a message locally using ECDSA recovery only.
// This is a convenience function for simple EOA verification without network calls.
//
// Example:
//
//	valid, err := public.VerifyMessageLocal(public.VerifyMessageLocalParameters{
//	    Address:   common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
//	    Message:   signature.NewSignableMessage("hello world"),
//	    Signature: "0x...",
//	})
func VerifyMessageLocal(params VerifyMessageParameters) (VerifyMessageReturnType, error) {
	return signature.VerifyMessage(
		params.Address.Hex(),
		params.Message,
		params.Signature,
	)
}
