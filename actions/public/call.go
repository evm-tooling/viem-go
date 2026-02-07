package public

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ChefBingbong/viem-go/abi"
	"github.com/ChefBingbong/viem-go/constants"
	"github.com/ChefBingbong/viem-go/types"
	blockoverride "github.com/ChefBingbong/viem-go/utils/block_override"
	"github.com/ChefBingbong/viem-go/utils/ccip"
	"github.com/ChefBingbong/viem-go/utils/deployless"
	stateoverride "github.com/ChefBingbong/viem-go/utils/state_override"
	"github.com/ChefBingbong/viem-go/utils/transaction"
)

// CallParameters contains the parameters for the Call action.
// This mirrors viem's CallParameters type with full feature support.
type CallParameters struct {
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

	// Nonce is the nonce for the call.
	Nonce *uint64

	// BlockNumber is the block number to execute the call at.
	// Mutually exclusive with BlockTag.
	BlockNumber *uint64

	// BlockTag is the block tag to execute the call at (e.g., "latest", "pending").
	// Mutually exclusive with BlockNumber.
	BlockTag BlockTag

	// Batch enables multicall batching for this call.
	// If nil, uses client's batch setting.
	Batch *bool

	// BlockOverrides contains block-level overrides (baseFeePerGas, gasLimit, etc.)
	BlockOverrides *types.BlockOverrides

	// Code is bytecode for deployless calls (call code without deploying).
	// Mutually exclusive with Factory/FactoryData.
	Code []byte

	// Factory is the contract deployment factory address (e.g., Create2 factory).
	// Used with FactoryData for deployless calls via factory.
	Factory *common.Address

	// FactoryData is the calldata to execute on the factory to deploy the contract.
	// Used with Factory for deployless calls via factory.
	FactoryData []byte

	// StateOverride contains state overrides for the call.
	StateOverride types.StateOverride

	// AccessList is the EIP-2930 access list.
	AccessList types.AccessList

	// AuthorizationList is the EIP-7702 authorization list.
	AuthorizationList []types.SignedAuthorization

	// Blobs is the EIP-4844 blob data.
	Blobs [][]byte

	// MaxFeePerBlobGas is the EIP-4844 max fee per blob gas.
	MaxFeePerBlobGas *big.Int

	// BlobVersionedHashes is the EIP-4844 blob versioned hashes.
	BlobVersionedHashes []common.Hash
}

// CallReturnType is the return type for the Call action.
type CallReturnType struct {
	// Data is the return data from the call, or nil if the call returned empty.
	Data []byte
}

// callRequest is the internal request format for eth_call.
type callRequest struct {
	From                 string           `json:"from,omitempty"`
	To                   string           `json:"to,omitempty"`
	Data                 string           `json:"data,omitempty"`
	Value                string           `json:"value,omitempty"`
	Gas                  string           `json:"gas,omitempty"`
	GasPrice             string           `json:"gasPrice,omitempty"`
	MaxFeePerGas         string           `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas string           `json:"maxPriorityFeePerGas,omitempty"`
	MaxFeePerBlobGas     string           `json:"maxFeePerBlobGas,omitempty"`
	Nonce                string           `json:"nonce,omitempty"`
	AccessList           types.AccessList `json:"accessList,omitempty"`
}

// Call executes a new message call immediately without submitting a transaction to the network.
//
// This is equivalent to viem's `call` action with full feature support including:
//   - Deployless calls (via bytecode or factory)
//   - State and block overrides
//   - Multicall batching
//   - CCIP-Read support
//   - Request validation
//
// JSON-RPC Method: eth_call
//
// Example:
//
//	result, err := public.Call(ctx, client, public.CallParameters{
//	    To:   &contractAddress,
//	    Data: calldata,
//	})
//
// Deployless call example:
//
//	result, err := public.Call(ctx, client, public.CallParameters{
//	    Code: contractBytecode,
//	    Data: calldata,
//	})
func Call(ctx context.Context, client Client, params CallParameters) (*CallReturnType, error) {
	// Validate mutually exclusive parameters
	if len(params.Code) > 0 && (params.Factory != nil || len(params.FactoryData) > 0) {
		return nil, &InvalidCallParamsError{
			Message: "cannot provide both 'code' and 'factory'/'factoryData' as parameters",
		}
	}
	if len(params.Code) > 0 && params.To != nil {
		return nil, &InvalidCallParamsError{
			Message: "cannot provide both 'code' and 'to' as parameters",
		}
	}

	// Check for deployless call types
	deploylessCallViaBytecode := len(params.Code) > 0 && len(params.Data) > 0
	deploylessCallViaFactory := params.Factory != nil && len(params.FactoryData) > 0 && params.To != nil && len(params.Data) > 0
	isDeploylessCall := deploylessCallViaBytecode || deploylessCallViaFactory

	// Build the calldata (potentially wrapping for deployless calls)
	data := params.Data
	if deploylessCallViaBytecode {
		var err error
		data, err = deployless.ToDeploylessCallViaBytecodeData(params.Code, params.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to encode deployless call: %w", err)
		}
	} else if deploylessCallViaFactory {
		var err error
		data, err = deployless.ToDeploylessCallViaFactoryData(*params.To, params.Data, *params.Factory, params.FactoryData)
		if err != nil {
			return nil, fmt.Errorf("failed to encode deployless factory call: %w", err)
		}
	}

	// Validate request
	accountAddr := ""
	if params.Account != nil {
		accountAddr = params.Account.Hex()
	}
	toAddr := ""
	if params.To != nil && !isDeploylessCall {
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

	// Serialize overrides
	rpcBlockOverrides := blockoverride.SerializeBlockOverrides(params.BlockOverrides)
	rpcStateOverride, err := stateoverride.SerializeStateOverride(params.StateOverride)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize state override: %w", err)
	}

	// Build the call request
	req := callRequest{}

	if params.Account != nil {
		req.From = params.Account.Hex()
	}
	if params.To != nil && !isDeploylessCall {
		req.To = params.To.Hex()
	}
	if len(data) > 0 {
		req.Data = hexutil.Encode(data)
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
	if params.Nonce != nil {
		req.Nonce = hexutil.EncodeUint64(*params.Nonce)
	}
	if len(params.AccessList) > 0 {
		req.AccessList = params.AccessList
	}

	// Check if we should batch via multicall
	batch := params.Batch
	if batch == nil && client.Batch() != nil && client.Batch().Multicall != nil {
		b := true
		batch = &b
	}

	if batch != nil && *batch && shouldPerformMulticall(req) && rpcStateOverride == nil && rpcBlockOverrides == nil {
		result, multicallErr := scheduleMulticall(ctx, client, req, params.BlockNumber, params.BlockTag)
		if multicallErr != nil {
			// Fall through to regular call if multicall fails due to chain not supporting it
			if _, ok := multicallErr.(*ChainNotConfiguredError); !ok {
				if _, ok := multicallErr.(*ChainDoesNotSupportContractError); !ok {
					return nil, multicallErr
				}
			}
		} else {
			return result, nil
		}
	}

	// Build params array
	rpcParams := []any{req, blockTag}
	if rpcStateOverride != nil && rpcBlockOverrides != nil {
		rpcParams = append(rpcParams, rpcStateOverride, rpcBlockOverrides)
	} else if rpcStateOverride != nil {
		rpcParams = append(rpcParams, rpcStateOverride)
	} else if rpcBlockOverrides != nil {
		rpcParams = append(rpcParams, map[string]any{}, rpcBlockOverrides)
	}

	// Execute the call
	resp, err := client.Request(ctx, "eth_call", rpcParams...)
	if err != nil {
		// Handle CCIP-Read
		revertData := getRevertErrorData(err)
		if len(revertData) >= 4 {
			selector := hexutil.Encode(revertData[:4])

			// Check for CCIP-Read offchain lookup
			ccipReadConfig := client.CCIPRead()
			if ccipReadConfig != nil && selector == ccip.OffchainLookupSignature && params.To != nil {
				result, ccipErr := handleCCIPRead(ctx, client, params, revertData)
				if ccipErr == nil {
					return result, nil
				}
				// If CCIP-Read fails, fall through to return original error
			}

			// Check for counterfactual deployment failure
			if isDeploylessCall && selector == constants.CounterfactualDeploymentFailedSignature {
				return nil, &CounterfactualDeploymentFailedError{Factory: params.Factory}
			}
		}

		return nil, &CallExecutionError{Cause: err, To: params.To, Data: data}
	}

	var hexResult string
	if unmarshalErr := json.Unmarshal(resp.Result, &hexResult); unmarshalErr != nil {
		return nil, fmt.Errorf("failed to unmarshal call result: %w", unmarshalErr)
	}

	// Parse the result
	resultData, parseErr := parseHexBytes(hexResult)
	if parseErr != nil {
		return nil, fmt.Errorf("failed to parse call result: %w", parseErr)
	}

	// Return nil data if the result is empty (0x)
	if len(resultData) == 0 {
		return &CallReturnType{Data: nil}, nil
	}

	return &CallReturnType{Data: resultData}, nil
}

// shouldPerformMulticall determines if a call should be batched via multicall.
// Returns true if the call has data, has a target, is not already a multicall,
// and has no extra parameters that can't be multicalled.
func shouldPerformMulticall(req callRequest) bool {
	// Must have data
	if req.Data == "" {
		return false
	}

	// Must have target
	if req.To == "" {
		return false
	}

	// Must not already be a multicall
	if strings.HasPrefix(strings.ToLower(req.Data), strings.ToLower(constants.Aggregate3Signature)) {
		return false
	}

	// Must not have extra parameters that can't be multicalled
	if req.Gas != "" || req.GasPrice != "" || req.MaxFeePerGas != "" ||
		req.MaxPriorityFeePerGas != "" || req.Value != "" || req.Nonce != "" {
		return false
	}

	return true
}

// scheduleMulticall batches the call via multicall3.
func scheduleMulticall(ctx context.Context, client Client, req callRequest, blockNumber *uint64, blockTag BlockTag) (*CallReturnType, error) {
	chain := client.Chain()
	if chain == nil {
		return nil, &ChainNotConfiguredError{}
	}

	batchOpts := client.Batch()
	if batchOpts == nil || batchOpts.Multicall == nil {
		return nil, &ChainNotConfiguredError{}
	}

	// Check if chain has multicall3
	var multicallAddress *common.Address
	if !batchOpts.Multicall.Deployless {
		if chain.Contracts == nil || chain.Contracts.Multicall3 == nil {
			return nil, &ChainDoesNotSupportContractError{
				ChainID:      chain.ID,
				ContractName: "multicall3",
				BlockNumber:  blockNumber,
			}
		}

		// Check block number constraint
		if blockNumber != nil && chain.Contracts.Multicall3.BlockCreated != nil {
			if *blockNumber < *chain.Contracts.Multicall3.BlockCreated {
				return nil, &ChainDoesNotSupportContractError{
					ChainID:      chain.ID,
					ContractName: "multicall3",
					BlockNumber:  blockNumber,
				}
			}
		}

		multicallAddress = &chain.Contracts.Multicall3.Address
	}

	// Encode the aggregate3 call using low-level ABI encoding
	// aggregate3 signature: aggregate3((address target, bool allowFailure, bytes callData)[])
	// Function selector: 0x82ad56cb

	// Encode the calls tuple array
	target := common.HexToAddress(req.To)
	callData := common.FromHex(req.Data)

	calls := make([]Call3, 1)
	calls = append(calls, Call3{target, true, callData})

	// Encode a single Call3 struct: (address, bool, bytes)
	callEncoded, err := abi.EncodeAbiParameters(
		[]abi.AbiParam{
			{
				Type: "tuple[]",
				Components: []abi.AbiParam{
					{Name: "target", Type: "address"},
					{Name: "allowFailure", Type: "bool"},
					{Name: "callData", Type: "bytes"},
				},
			},
		},
		[]any{calls},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to encode multicall args: %w", err)
	}

	// Prepend function selector
	selector := common.FromHex(constants.Aggregate3Signature)
	calldata := make([]byte, len(selector)+len(callEncoded))
	copy(calldata, selector)
	copy(calldata[len(selector):], callEncoded)

	// Build multicall request
	multicallReq := callRequest{Data: hexutil.Encode(calldata)}
	if multicallAddress != nil {
		multicallReq.To = multicallAddress.Hex()
	} else {
		// Deployless multicall - wrap in deployless bytecode
		deploylessData, deploylessErr := deployless.ToDeploylessCallViaBytecodeData(
			common.FromHex(constants.Multicall3Bytecode),
			calldata,
		)
		if deploylessErr != nil {
			return nil, fmt.Errorf("failed to encode deployless multicall: %w", deploylessErr)
		}
		multicallReq.Data = hexutil.Encode(deploylessData)
	}

	// Execute multicall
	block := resolveBlockTag(client, blockNumber, blockTag)
	resp, requestErr := client.Request(ctx, "eth_call", multicallReq, block)
	if requestErr != nil {
		return nil, requestErr
	}

	var hexResult string
	if unmarshalErr := json.Unmarshal(resp.Result, &hexResult); unmarshalErr != nil {
		return nil, fmt.Errorf("failed to unmarshal multicall result: %w", unmarshalErr)
	}

	// Decode multicall result
	// aggregate3 returns: (bool success, bytes returnData)[]
	resultData := common.FromHex(hexResult)
	decodedResults, decodeErr := decodeAggregate3Result(resultData)
	if decodeErr != nil {
		return nil, fmt.Errorf("failed to decode multicall result: %w", decodeErr)
	}

	callReturn := CallReturnType{Data: decodedResults[0].ReturnData}
	return &callReturn, nil
}

// handleCCIPRead handles CCIP-Read offchain lookup.
func handleCCIPRead(ctx context.Context, client Client, params CallParameters, revertData []byte) (*CallReturnType, error) {
	// Decode the offchain lookup error
	lookup, err := ccip.DecodeOffchainLookupError(revertData)
	if err != nil {
		return nil, err
	}

	// Verify sender matches
	if params.To != nil && lookup.Sender != *params.To {
		return nil, &ccip.ErrOffchainLookupSenderMismatch{
			Sender: lookup.Sender,
			To:     *params.To,
		}
	}

	// Make CCIP request
	ccipResult, err := ccip.CCIPRequest(ctx, ccip.CCIPRequestParams{
		Data:   lookup.CallData,
		Sender: lookup.Sender,
		URLs:   lookup.URLs,
	})
	if err != nil {
		return nil, &ccip.ErrOffchainLookup{
			CallbackSelector: lookup.CallbackFunction,
			Data:             revertData,
			ExtraData:        lookup.ExtraData,
			Sender:           lookup.Sender,
			URLs:             lookup.URLs,
			Cause:            err,
		}
	}

	// Build callback calldata
	callbackData, err := ccip.BuildCallbackData(lookup.CallbackFunction, ccipResult, lookup.ExtraData)
	if err != nil {
		return nil, err
	}

	// Execute callback
	return Call(ctx, client, CallParameters{
		To:          params.To,
		Data:        callbackData,
		BlockNumber: params.BlockNumber,
		BlockTag:    params.BlockTag,
	})
}

// getRevertErrorData extracts revert data from an error.
func getRevertErrorData(err error) []byte {
	if err == nil {
		return nil
	}

	// Try to extract from error message (common RPC error format)
	errStr := err.Error()

	// Look for hex data in the error
	if idx := strings.Index(errStr, "0x"); idx >= 0 {
		hexStr := errStr[idx:]
		// Find end of hex string
		end := len(hexStr)
		for i, c := range hexStr[2:] {
			isHexDigit := (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
			if !isHexDigit {
				end = i + 2
				break
			}
		}
		if end > 2 {
			return common.FromHex(hexStr[:end])
		}
	}

	return nil
}

// resolveBlockTag determines the block tag to use for a request.
func resolveBlockTag(client Client, blockNumber *uint64, blockTag BlockTag) string {
	if blockNumber != nil {
		return hexutil.EncodeUint64(*blockNumber)
	}
	if blockTag != "" {
		return string(blockTag)
	}
	if experimentalTag := client.ExperimentalBlockTag(); experimentalTag != "" {
		return string(experimentalTag)
	}
	return string(BlockTagLatest)
}

// parseHexBytes parses a hex string to bytes.
func parseHexBytes(hexStr string) ([]byte, error) {
	if hexStr == "" || hexStr == "0x" {
		return []byte{}, nil
	}
	hexStr = strings.TrimPrefix(hexStr, "0x")
	// Handle odd-length hex strings by padding with leading zero
	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}
	return hex.DecodeString(hexStr)
}
