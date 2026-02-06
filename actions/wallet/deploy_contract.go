package wallet

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	viemabi "github.com/ChefBingbong/viem-go/abi"
	viemchain "github.com/ChefBingbong/viem-go/chain"
	"github.com/ChefBingbong/viem-go/utils/formatters"
	"github.com/ChefBingbong/viem-go/utils/transaction"
)

// DeployContractParameters contains the parameters for the DeployContract action.
// This mirrors viem's DeployContractParameters type.
type DeployContractParameters struct {
	// Account is the account to deploy from. If nil, uses the client's account.
	Account Account

	// ABI is the contract ABI as JSON bytes, string, or a pre-parsed *abi.ABI.
	ABI any // []byte, string, or *abi.ABI

	// Bytecode is the contract bytecode as a hex string (with or without 0x prefix).
	Bytecode string

	// Args are the constructor arguments.
	Args []any

	// Chain optionally overrides the client's chain for chain ID validation.
	Chain *viemchain.Chain

	// AssertChainID when true, asserts the chain ID matches. Default: true.
	AssertChainID *bool

	// DataSuffix is data to append to the end of the calldata.
	DataSuffix string

	// Value is the amount of ETH to send with the deployment transaction.
	Value *big.Int

	// Transaction fields
	AuthorizationList    []transaction.SignedAuthorization `json:"authorizationList,omitempty"`
	BlobVersionedHashes  []string                          `json:"blobVersionedHashes,omitempty"`
	Blobs                []string                          `json:"blobs,omitempty"`
	Gas                  *big.Int                          `json:"gas,omitempty"`
	GasPrice             *big.Int                          `json:"gasPrice,omitempty"`
	MaxFeePerBlobGas     *big.Int                          `json:"maxFeePerBlobGas,omitempty"`
	MaxFeePerGas         *big.Int                          `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas *big.Int                          `json:"maxPriorityFeePerGas,omitempty"`
	Nonce                *int                              `json:"nonce,omitempty"`
	Type                 formatters.TransactionType        `json:"type,omitempty"`
}

// DeployContractReturnType is the return type for the DeployContract action.
// It is the transaction hash as a hex string.
type DeployContractReturnType = SendTransactionReturnType

// DeployContract deploys a contract to the network, given bytecode and constructor arguments.
//
// Internally, encodes the deploy data (bytecode + ABI-encoded constructor args) and delegates
// to SendTransaction with no `to` address (contract creation).
//
// This is equivalent to viem's `deployContract` action.
//
// Example:
//
//	hash, err := wallet.DeployContract(ctx, client, wallet.DeployContractParameters{
//	    ABI:      contractABI,
//	    Bytecode: "0x608060405260405161083e38038061083e833981016040819052610...",
//	})
//
// Example with constructor arguments:
//
//	hash, err := wallet.DeployContract(ctx, client, wallet.DeployContractParameters{
//	    ABI:      contractABI,
//	    Bytecode: "0x608060405260405161083e38038061083e833981016040819052610...",
//	    Args:     []any{"MyToken", "MTK", uint8(18)},
//	})
func DeployContract(ctx context.Context, client Client, params DeployContractParameters) (DeployContractReturnType, error) {
	// Parse the ABI
	parsedABI, err := parseABIParam(params.ABI)
	if err != nil {
		return "", fmt.Errorf("failed to parse ABI: %w", err)
	}

	// Encode deploy data: bytecode + constructor args
	// This mirrors viem's encodeDeployData({ abi, args, bytecode })
	calldata, err := encodeDeployData(parsedABI, params.Bytecode, params.Args)
	if err != nil {
		return "", fmt.Errorf("failed to encode deploy data: %w", err)
	}

	// Convert to hex string
	calldataHex := "0x" + hex.EncodeToString(calldata)

	// Delegate to SendTransaction with no `to` address (contract creation)
	// This mirrors viem's: sendTransaction(walletClient, { ...request, data: calldata })
	// Note: `to` is intentionally empty for contract deployment
	return SendTransaction(ctx, client, SendTransactionParameters{
		Account:              params.Account,
		Chain:                params.Chain,
		AssertChainID:        params.AssertChainID,
		DataSuffix:           params.DataSuffix,
		Data:                 calldataHex,
		Value:                params.Value,
		AuthorizationList:    params.AuthorizationList,
		BlobVersionedHashes:  params.BlobVersionedHashes,
		Blobs:                params.Blobs,
		Gas:                  params.Gas,
		GasPrice:             params.GasPrice,
		MaxFeePerBlobGas:     params.MaxFeePerBlobGas,
		MaxFeePerGas:         params.MaxFeePerGas,
		MaxPriorityFeePerGas: params.MaxPriorityFeePerGas,
		Nonce:                params.Nonce,
		Type:                 params.Type,
		// To is intentionally empty (nil) — this is a contract creation transaction
	})
}

// encodeDeployData encodes bytecode + ABI-encoded constructor arguments.
// This mirrors viem's encodeDeployData utility.
func encodeDeployData(parsedABI *viemabi.ABI, bytecode string, args []any) ([]byte, error) {
	// Parse bytecode hex string to bytes
	bytecodeHex := strings.TrimPrefix(bytecode, "0x")
	bytecodeHex = strings.TrimPrefix(bytecodeHex, "0X")
	bytecodeBytes, err := hex.DecodeString(bytecodeHex)
	if err != nil {
		return nil, fmt.Errorf("invalid bytecode hex: %w", err)
	}

	// If there are no constructor args, just return the bytecode
	if len(args) == 0 {
		return bytecodeBytes, nil
	}

	// Encode constructor arguments using the ABI
	constructorArgs, err := parsedABI.EncodeConstructor(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to encode constructor arguments: %w", err)
	}

	// Concatenate bytecode + constructor args
	result := make([]byte, len(bytecodeBytes)+len(constructorArgs))
	copy(result, bytecodeBytes)
	copy(result[len(bytecodeBytes):], constructorArgs)

	return result, nil
}
