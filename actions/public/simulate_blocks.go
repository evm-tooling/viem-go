package public

import (
	"context"
	"fmt"
	"math/big"

	json "github.com/goccy/go-json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ChefBingbong/viem-go/abi"
	"github.com/ChefBingbong/viem-go/types"
	blockoverride "github.com/ChefBingbong/viem-go/utils/block_override"
	"github.com/ChefBingbong/viem-go/utils/formatters"
	stateoverride "github.com/ChefBingbong/viem-go/utils/state_override"
)

// SimulateBlockCall represents a single call within a block simulation.
type SimulateBlockCall struct {
	// From is the sender address (account attached to the call).
	From *common.Address

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

	// AccessList is the EIP-2930 access list.
	AccessList types.AccessList

	// AuthorizationList is the EIP-7702 authorization list.
	AuthorizationList []types.SignedAuthorization

	// ABI is the contract ABI for decoding results.
	ABI *abi.ABI

	// FunctionName is the function name for ABI decoding.
	FunctionName string

	// Args are the function arguments (for encoding if needed).
	Args []any

	// DataSuffix is optional data to append to the calldata.
	DataSuffix []byte
}

// SimulateBlock represents a single block to simulate.
type SimulateBlock struct {
	// BlockOverrides contains block-level overrides.
	BlockOverrides *types.BlockOverrides

	// Calls are the calls to execute in this block.
	Calls []SimulateBlockCall

	// StateOverrides contains state overrides for this block.
	StateOverrides types.StateOverride
}

// SimulateBlocksParameters contains the parameters for the SimulateBlocks action.
type SimulateBlocksParameters struct {
	// Blocks to simulate.
	Blocks []SimulateBlock

	// BlockNumber is the block number to simulate at.
	// Mutually exclusive with BlockTag.
	BlockNumber *uint64

	// BlockTag is the block tag to simulate at (e.g., "latest", "pending").
	// Mutually exclusive with BlockNumber.
	BlockTag BlockTag

	// ReturnFullTransactions determines whether to return full transaction objects.
	ReturnFullTransactions bool

	// TraceTransfers enables transfer tracing.
	TraceTransfers bool

	// Validation enables validation mode.
	Validation bool
}

// CallResult represents the result of a single call in the simulation.
type CallResult struct {
	// Status is the result status ("success" or "failure").
	Status string

	// Data is the return data from the call.
	Data []byte

	// GasUsed is the amount of gas used by the call.
	GasUsed *big.Int

	// Logs are the logs emitted by the call.
	Logs []formatters.Log

	// Result is the decoded result (if ABI was provided).
	Result any

	// Error is the error if the call failed.
	Error error
}

// BlockResult represents the result of a block simulation.
type BlockResult struct {
	// Block is the simulated block data.
	formatters.Block

	// Calls are the results of each call in the block.
	Calls []CallResult
}

// SimulateBlocksReturnType is the return type for the SimulateBlocks action.
type SimulateBlocksReturnType = []BlockResult

// simulateAccessListItem is the RPC format for an access list entry in simulation.
type simulateAccessListItem struct {
	Address     string   `json:"address"`
	StorageKeys []string `json:"storageKeys"`
}

// rpcSimulateCall is the RPC format for a simulation call.
type rpcSimulateCall struct {
	From                 string                   `json:"from,omitempty"`
	To                   string                   `json:"to,omitempty"`
	Data                 string                   `json:"data,omitempty"`
	Value                string                   `json:"value,omitempty"`
	Gas                  string                   `json:"gas,omitempty"`
	GasPrice             string                   `json:"gasPrice,omitempty"`
	MaxFeePerGas         string                   `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas string                   `json:"maxPriorityFeePerGas,omitempty"`
	Nonce                string                   `json:"nonce,omitempty"`
	AccessList           []simulateAccessListItem `json:"accessList,omitempty"`
}

// rpcBlockStateCall is the RPC format for a block state call.
type rpcBlockStateCall struct {
	BlockOverrides *types.RpcBlockOverrides `json:"blockOverrides,omitempty"`
	Calls          []rpcSimulateCall        `json:"calls"`
	StateOverrides types.RpcStateOverride   `json:"stateOverrides,omitempty"`
}

// rpcSimulateParams is the RPC format for eth_simulateV1 params.
type rpcSimulateParams struct {
	BlockStateCalls        []rpcBlockStateCall `json:"blockStateCalls"`
	ReturnFullTransactions *bool               `json:"returnFullTransactions,omitempty"`
	TraceTransfers         *bool               `json:"traceTransfers,omitempty"`
	Validation             *bool               `json:"validation,omitempty"`
}

// rpcSimulateCallResult is the RPC response format for a call result.
type rpcSimulateCallResult struct {
	Status     string              `json:"status"`
	ReturnData string              `json:"returnData"`
	GasUsed    string              `json:"gasUsed"`
	Logs       []formatters.RpcLog `json:"logs,omitempty"`
	Error      *rpcSimulateError   `json:"error,omitempty"`
}

// rpcSimulateError is the RPC format for a simulation error.
type rpcSimulateError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

// rpcSimulateBlockResult is the RPC response format for a block result.
type rpcSimulateBlockResult struct {
	formatters.RpcBlock
	Calls []rpcSimulateCallResult `json:"calls"`
}

// SimulateBlocks simulates a set of calls on block(s) with optional block and state overrides.
//
// This is equivalent to viem's `simulateBlocks` action.
//
// JSON-RPC Method: eth_simulateV1
//
// Example:
//
//	results, err := public.SimulateBlocks(ctx, client, public.SimulateBlocksParameters{
//	    Blocks: []public.SimulateBlock{{
//	        BlockOverrides: &types.BlockOverrides{Number: ptr(uint64(69420))},
//	        Calls: []public.SimulateBlockCall{{
//	            From: &senderAddr,
//	            To:   &contractAddr,
//	            Data: calldata,
//	        }},
//	    }},
//	})
func SimulateBlocks(ctx context.Context, client Client, params SimulateBlocksParameters) (SimulateBlocksReturnType, error) {
	// Build block state calls
	blockStateCalls := make([]rpcBlockStateCall, 0, len(params.Blocks))

	for _, block := range params.Blocks {
		// Serialize block overrides
		rpcBlockOverrides := blockoverride.SerializeBlockOverrides(block.BlockOverrides)

		// Build calls
		calls := make([]rpcSimulateCall, 0, len(block.Calls))
		for _, call := range block.Calls {
			rpcCall := rpcSimulateCall{}

			if call.From != nil {
				rpcCall.From = call.From.Hex()
			}
			if call.To != nil {
				rpcCall.To = call.To.Hex()
			}

			// Handle data - either use provided data or encode from ABI
			data := call.Data
			if len(data) == 0 && call.ABI != nil && call.FunctionName != "" {
				encoded, err := call.ABI.EncodeFunctionData(call.FunctionName, call.Args...)
				if err != nil {
					return nil, fmt.Errorf("failed to encode function data: %w", err)
				}
				data = encoded
			}

			// Append data suffix if provided
			if len(call.DataSuffix) > 0 {
				data = append(data, call.DataSuffix...)
			}

			if len(data) > 0 {
				rpcCall.Data = hexutil.Encode(data)
			}

			if call.Value != nil {
				rpcCall.Value = hexutil.EncodeBig(call.Value)
			}
			if call.Gas != nil {
				rpcCall.Gas = hexutil.EncodeUint64(*call.Gas)
			}
			if call.GasPrice != nil {
				rpcCall.GasPrice = hexutil.EncodeBig(call.GasPrice)
			}
			if call.MaxFeePerGas != nil {
				rpcCall.MaxFeePerGas = hexutil.EncodeBig(call.MaxFeePerGas)
			}
			if call.MaxPriorityFeePerGas != nil {
				rpcCall.MaxPriorityFeePerGas = hexutil.EncodeBig(call.MaxPriorityFeePerGas)
			}
			if call.Nonce != nil {
				rpcCall.Nonce = hexutil.EncodeUint64(*call.Nonce)
			}
			if len(call.AccessList) > 0 {
				rpcAccessList := make([]simulateAccessListItem, len(call.AccessList))
				for i, item := range call.AccessList {
					storageKeys := make([]string, len(item.StorageKeys))
					for j, key := range item.StorageKeys {
						storageKeys[j] = key.Hex()
					}
					rpcAccessList[i] = simulateAccessListItem{
						Address:     item.Address.Hex(),
						StorageKeys: storageKeys,
					}
				}
				rpcCall.AccessList = rpcAccessList
			}

			calls = append(calls, rpcCall)
		}

		// Serialize state overrides
		rpcStateOverride, err := stateoverride.SerializeStateOverride(block.StateOverrides)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize state override: %w", err)
		}

		blockStateCalls = append(blockStateCalls, rpcBlockStateCall{
			BlockOverrides: rpcBlockOverrides,
			Calls:          calls,
			StateOverrides: rpcStateOverride,
		})
	}

	// Build params - only include optional booleans when true
	rpcParams := rpcSimulateParams{
		BlockStateCalls: blockStateCalls,
	}
	if params.ReturnFullTransactions {
		rpcParams.ReturnFullTransactions = &params.ReturnFullTransactions
	}
	if params.TraceTransfers {
		rpcParams.TraceTransfers = &params.TraceTransfers
	}
	if params.Validation {
		rpcParams.Validation = &params.Validation
	}

	// Determine block tag
	blockTag := resolveBlockTag(client, params.BlockNumber, params.BlockTag)

	// Execute the call
	resp, err := client.Request(ctx, "eth_simulateV1", rpcParams, blockTag)
	if err != nil {
		return nil, &SimulateBlocksError{Cause: err}
	}

	// Parse the response
	var rpcResults []rpcSimulateBlockResult
	if err := json.Unmarshal(resp.Result, &rpcResults); err != nil {
		return nil, fmt.Errorf("failed to unmarshal simulate result: %w", err)
	}

	// Format results
	results := make([]BlockResult, len(rpcResults))
	for i, rpcBlock := range rpcResults {
		// Format block
		block := formatters.FormatBlock(rpcBlock.RpcBlock)

		// Format call results
		callResults := make([]CallResult, len(rpcBlock.Calls))
		for j, rpcCall := range rpcBlock.Calls {
			status := "failure"
			if rpcCall.Status == "0x1" {
				status = "success"
			}

			// Parse return data
			var data []byte
			if rpcCall.ReturnData != "" && rpcCall.ReturnData != "0x" {
				data = common.FromHex(rpcCall.ReturnData)
			}

			// Parse gas used
			var gasUsed *big.Int
			if rpcCall.GasUsed != "" {
				gasUsed = new(big.Int)
				gasUsed.SetString(rpcCall.GasUsed[2:], 16)
			}

			// Format logs
			var logs []formatters.Log
			if len(rpcCall.Logs) > 0 {
				logs = formatters.FormatLogs(rpcCall.Logs)
			}

			// Decode result if ABI provided
			var result any
			var callError error

			if status == "success" && len(data) > 0 {
				// Try to decode using the original call's ABI if available
				if j < len(params.Blocks[i].Calls) {
					originalCall := params.Blocks[i].Calls[j]
					if originalCall.ABI != nil && originalCall.FunctionName != "" {
						decoded, decodeErr := originalCall.ABI.DecodeFunctionResult(originalCall.FunctionName, data)
						if decodeErr == nil && len(decoded) > 0 {
							if len(decoded) == 1 {
								result = decoded[0]
							} else {
								result = decoded
							}
						}
					}
				}
			} else if status == "failure" {
				// Create error from revert data
				if len(data) > 0 {
					callError = &RawContractError{Data: data}
				} else if rpcCall.Error != nil {
					callError = fmt.Errorf("%s", rpcCall.Error.Message)
				} else {
					callError = fmt.Errorf("call failed")
				}
			}

			callResults[j] = CallResult{
				Status:  status,
				Data:    data,
				GasUsed: gasUsed,
				Logs:    logs,
				Result:  result,
				Error:   callError,
			}
		}

		results[i] = BlockResult{
			Block: block,
			Calls: callResults,
		}
	}

	return results, nil
}

// SimulateBlocksError is returned when block simulation fails.
type SimulateBlocksError struct {
	Cause error
}

func (e *SimulateBlocksError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("block simulation failed: %v", e.Cause)
	}
	return "block simulation failed"
}

func (e *SimulateBlocksError) Unwrap() error {
	return e.Cause
}
