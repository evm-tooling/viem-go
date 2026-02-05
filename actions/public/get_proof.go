package public

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ChefBingbong/viem-go/utils/formatters"
)

// GetProofParameters contains the parameters for the GetProof action.
// This mirrors viem's GetProofParameters type.
type GetProofParameters struct {
	// Address is the account address to retrieve proof for.
	Address common.Address

	// StorageKeys is the list of storage keys to include in the proof.
	StorageKeys []common.Hash

	// BlockNumber is the block number to get the proof at.
	// Mutually exclusive with BlockTag.
	BlockNumber *uint64

	// BlockTag is the block tag to get the proof at (e.g., "latest", "pending").
	// Mutually exclusive with BlockNumber.
	// Default: "latest"
	BlockTag BlockTag
}

// GetProofReturnType is the return type for the GetProof action.
// It represents an account proof with formatted numeric fields.
type GetProofReturnType = formatters.Proof

// GetProofError is returned when getProof fails.
type GetProofError struct {
	Cause error
}

func (e *GetProofError) Error() string {
	return fmt.Sprintf("eth_getProof failed: %v", e.Cause)
}

func (e *GetProofError) Unwrap() error {
	return e.Cause
}

// GetProof returns the account and storage values of the specified account
// including the Merkle-proof.
//
// This is equivalent to viem's `getProof` action.
//
// JSON-RPC Method: eth_getProof (EIP-1186)
//
// Example:
//
//	proof, err := public.GetProof(ctx, client, public.GetProofParameters{
//	    Address:     common.HexToAddress("0x..."),
//	    StorageKeys: []common.Hash{common.HexToHash("0x0")},
//	})
func GetProof(ctx context.Context, client Client, params GetProofParameters) (GetProofReturnType, error) {
	// Determine block tag/number
	blockTag := resolveBlockTag(client, params.BlockNumber, params.BlockTag)

	// Execute the request
	resp, err := client.Request(ctx, "eth_getProof", params.Address.Hex(), params.StorageKeys, blockTag)
	if err != nil {
		return formatters.Proof{}, &GetProofError{Cause: err}
	}

	var rpcProof formatters.RpcProof
	if unmarshalErr := json.Unmarshal(resp.Result, &rpcProof); unmarshalErr != nil {
		return formatters.Proof{}, fmt.Errorf("failed to unmarshal proof: %w", unmarshalErr)
	}

	// Format RPC proof into typed Proof (similar to viem's formatProof).
	proof, formatErr := formatProof(rpcProof)
	if formatErr != nil {
		return formatters.Proof{}, formatErr
	}

	return proof, nil
}

// formatProof formats an RpcProof into a Proof by decoding numeric fields.
func formatProof(rpc formatters.RpcProof) (formatters.Proof, error) {
	var (
		balance *big.Int
		nonce   *int
		err     error
	)

	// Decode balance (hex string) to *big.Int
	if rpc.Balance != "" {
		var b *big.Int
		b, err = hexutil.DecodeBig(rpc.Balance)
		if err != nil {
			return formatters.Proof{}, fmt.Errorf("invalid balance in proof: %w", err)
		}
		balance = b
	}

	// Decode nonce (hex string) to *int
	if rpc.Nonce != "" {
		var n uint64
		n, err = hexutil.DecodeUint64(rpc.Nonce)
		if err != nil {
			return formatters.Proof{}, fmt.Errorf("invalid nonce in proof: %w", err)
		}
		tmp := int(n)
		nonce = &tmp
	}

	// Decode storage proofs
	storageProofs := make([]formatters.StorageProof, 0, len(rpc.StorageProof))
	for _, sp := range rpc.StorageProof {
		var value *big.Int
		if sp.Value != "" {
			v, decErr := hexutil.DecodeBig(sp.Value)
			if decErr != nil {
				return formatters.Proof{}, fmt.Errorf("invalid storage value in proof: %w", decErr)
			}
			value = v
		}

		storageProofs = append(storageProofs, formatters.StorageProof{
			Key:   sp.Key,
			Proof: sp.Proof,
			Value: value,
		})
	}

	return formatters.Proof{
		Address:      rpc.Address,
		AccountProof: rpc.AccountProof,
		Balance:      balance,
		CodeHash:     rpc.CodeHash,
		Nonce:        nonce,
		StorageHash:  rpc.StorageHash,
		StorageProof: storageProofs,
	}, nil
}
