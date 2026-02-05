package public

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/utils/signature"
)

// VerifyTypedDataParameters contains the parameters for the VerifyTypedData action.
// This mirrors viem's VerifyTypedDataParameters type with full feature support including:
//   - ERC-6492 (counterfactual signature verification for undeployed smart accounts)
//   - ERC-1271 (smart contract signature verification)
//   - ECDSA recovery fallback (for EOA accounts)
type VerifyTypedDataParameters struct {
	// Address is the address that is expected to have signed the typed data.
	Address common.Address

	// TypedData is the full EIP-712 typed data definition, including
	// domain, types, primary type, and message.
	TypedData signature.TypedDataDefinition

	// Signature is the signature produced by signing the typed data.
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

// VerifyTypedDataReturnType is the return type for the VerifyTypedData action.
// It indicates whether the signature is valid.
type VerifyTypedDataReturnType = bool

// VerifyTypedData verifies that EIP-712 typed data was signed by the provided address.
//
// Compatible with Smart Contract Accounts & Externally Owned Accounts via ERC-6492.
// https://eips.ethereum.org/EIPS/eip-6492
//
// This is equivalent to viem's `verifyTypedData` action with full feature support:
//   - ERC-6492 verification for counterfactual (undeployed) smart accounts
//   - ERC-1271 verification for deployed smart contracts
//   - ECDSA recovery fallback for EOA accounts
//
// Example:
//
//	valid, err := public.VerifyTypedData(ctx, client, public.VerifyTypedDataParameters{
//	    Address:   common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
//	    TypedData: typedDataDef,
//	    Signature: "0x...",
//	})
//
// Example with counterfactual verification:
//
//	valid, err := public.VerifyTypedData(ctx, client, public.VerifyTypedDataParameters{
//	    Address:     common.HexToAddress("0x..."),
//	    TypedData:   typedDataDef,
//	    Signature:   "0x...",
//	    Factory:     &factoryAddress,
//	    FactoryData: factoryCalldata,
//	})
func VerifyTypedData(ctx context.Context, client Client, params VerifyTypedDataParameters) (VerifyTypedDataReturnType, error) {
	// Hash the typed data using EIP-712
	hash, err := signature.HashTypedData(params.TypedData)
	if err != nil {
		return false, err
	}

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

// VerifyTypedDataLocal verifies typed data locally using ECDSA recovery only.
// This is a convenience function for simple EOA verification without network calls.
//
// Example:
//
//	valid, err := public.VerifyTypedDataLocal(public.VerifyTypedDataParameters{
//	    Address:   common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
//	    TypedData: typedDataDef,
//	    Signature: "0x...",
//	})
func VerifyTypedDataLocal(params VerifyTypedDataParameters) (VerifyTypedDataReturnType, error) {
	return signature.VerifyTypedData(
		params.Address.Hex(),
		params.TypedData,
		params.Signature,
	)
}
