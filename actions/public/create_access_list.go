package public

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ChefBingbong/viem-go/types"
	"github.com/ChefBingbong/viem-go/utils/transaction"
)

// CreateAccessListParameters contains the parameters for the CreateAccessList action.
// This mirrors viem's CreateAccessListParameters type.
type CreateAccessListParameters struct {
	// Account is the account attached to the call (msg.sender).
	Account *common.Address

	// To is the contract address to call.
	To *common.Address

	// Data is the calldata to send.
	Data []byte

	// Value is the amount of wei to send with the call.
	Value *big.Int

	// Gas is the gas limit for the call.
	Gas *uint64

	// GasPrice is the gas price for the call (legacy).
	GasPrice *big.Int

	// MaxFeePerGas is the max fee per gas (EIP-1559).
	MaxFeePerGas *big.Int

	// MaxPriorityFeePerGas is the max priority fee per gas (EIP-1559).
	MaxPriorityFeePerGas *big.Int

	// MaxFeePerBlobGas is the EIP-4844 max fee per blob gas.
	MaxFeePerBlobGas *big.Int

	// BlockNumber is the block number to create access list at.
	// Mutually exclusive with BlockTag.
	BlockNumber *uint64

	// BlockTag is the block tag to create access list at (e.g., "latest", "pending").
	// Mutually exclusive with BlockNumber.
	BlockTag BlockTag

	// Blobs is the EIP-4844 blob data.
	Blobs [][]byte
}

// CreateAccessListReturnType is the return type for the CreateAccessList action.
type CreateAccessListReturnType struct {
	// AccessList is the generated access list.
	AccessList types.AccessList

	// GasUsed is the estimated gas used with this access list.
	GasUsed *big.Int
}

// accessListRequest is the internal request format for eth_createAccessList.
type accessListRequest struct {
	From                 string `json:"from,omitempty"`
	To                   string `json:"to,omitempty"`
	Data                 string `json:"data,omitempty"`
	Value                string `json:"value,omitempty"`
	Gas                  string `json:"gas,omitempty"`
	GasPrice             string `json:"gasPrice,omitempty"`
	MaxFeePerGas         string `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas,omitempty"`
	MaxFeePerBlobGas     string `json:"maxFeePerBlobGas,omitempty"`
}

// rpcAccessListResponse is the RPC response format for eth_createAccessList.
type rpcAccessListResponse struct {
	AccessList []rpcAccessListItem `json:"accessList"`
	GasUsed    string              `json:"gasUsed"`
}

// rpcAccessListItem is the RPC format for an access list entry.
type rpcAccessListItem struct {
	Address     string   `json:"address"`
	StorageKeys []string `json:"storageKeys"`
}

// CreateAccessList creates an EIP-2930 access list.
//
// An access list contains a list of addresses and storage keys that the transaction
// plans to access. This can reduce gas costs by pre-warming the accessed state.
//
// This is equivalent to viem's `createAccessList` action.
//
// JSON-RPC Method: eth_createAccessList
//
// Example:
//
//	result, err := public.CreateAccessList(ctx, client, public.CreateAccessListParameters{
//	    Account: &senderAddr,
//	    To:      &contractAddr,
//	    Data:    calldata,
//	})
//	// result.AccessList contains the generated access list
//	// result.GasUsed contains the estimated gas with this access list
func CreateAccessList(ctx context.Context, client Client, params CreateAccessListParameters) (*CreateAccessListReturnType, error) {
	// Validate request
	accountAddr := ""
	if params.Account != nil {
		accountAddr = params.Account.Hex()
	}
	toAddr := ""
	if params.To != nil {
		toAddr = params.To.Hex()
	}

	if err := transaction.AssertRequest(transaction.AssertRequestParams{
		Account:              accountAddr,
		To:                   toAddr,
		MaxFeePerGas:         params.MaxFeePerGas,
		MaxPriorityFeePerGas: params.MaxPriorityFeePerGas,
	}); err != nil {
		return nil, err
	}

	// Determine block tag
	blockTag := resolveBlockTag(client, params.BlockNumber, params.BlockTag)

	// Build the request
	req := accessListRequest{}

	if params.Account != nil {
		req.From = params.Account.Hex()
	}
	if params.To != nil {
		req.To = params.To.Hex()
	}
	if len(params.Data) > 0 {
		req.Data = hexutil.Encode(params.Data)
	}
	if params.Value != nil {
		req.Value = hexutil.EncodeBig(params.Value)
	}
	if params.Gas != nil {
		req.Gas = hexutil.EncodeUint64(*params.Gas)
	}
	if params.GasPrice != nil {
		req.GasPrice = hexutil.EncodeBig(params.GasPrice)
	}
	if params.MaxFeePerGas != nil {
		req.MaxFeePerGas = hexutil.EncodeBig(params.MaxFeePerGas)
	}
	if params.MaxPriorityFeePerGas != nil {
		req.MaxPriorityFeePerGas = hexutil.EncodeBig(params.MaxPriorityFeePerGas)
	}
	if params.MaxFeePerBlobGas != nil {
		req.MaxFeePerBlobGas = hexutil.EncodeBig(params.MaxFeePerBlobGas)
	}

	// Execute the call
	resp, err := client.Request(ctx, "eth_createAccessList", req, blockTag)
	if err != nil {
		return nil, &CreateAccessListError{Cause: err, To: params.To, Data: params.Data}
	}

	// Parse the response
	var rpcResult rpcAccessListResponse
	if err := json.Unmarshal(resp.Result, &rpcResult); err != nil {
		return nil, fmt.Errorf("failed to unmarshal access list result: %w", err)
	}

	// Convert access list
	accessList := make(types.AccessList, len(rpcResult.AccessList))
	for i, item := range rpcResult.AccessList {
		storageKeys := make([]common.Hash, len(item.StorageKeys))
		for j, key := range item.StorageKeys {
			storageKeys[j] = common.HexToHash(key)
		}
		accessList[i] = types.AccessTuple{
			Address:     common.HexToAddress(item.Address),
			StorageKeys: storageKeys,
		}
	}

	// Parse gas used
	gasUsed := new(big.Int)
	if rpcResult.GasUsed != "" {
		gasUsed.SetString(rpcResult.GasUsed[2:], 16) // Remove 0x prefix
	}

	return &CreateAccessListReturnType{
		AccessList: accessList,
		GasUsed:    gasUsed,
	}, nil
}

// CreateAccessListError is returned when access list creation fails.
type CreateAccessListError struct {
	Cause error
	To    *common.Address
	Data  []byte
}

func (e *CreateAccessListError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("access list creation failed: %v", e.Cause)
	}
	return "access list creation failed"
}

func (e *CreateAccessListError) Unwrap() error {
	return e.Cause
}
