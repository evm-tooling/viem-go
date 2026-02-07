package public

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"sync"

	gethabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ChefBingbong/viem-go/abi"
	"github.com/ChefBingbong/viem-go/constants"
	"github.com/ChefBingbong/viem-go/utils/deployless"
)

// Cached aggregate3 ABI parameters and pre-parsed go-ethereum Arguments.
// The Arguments are parsed once and reused for all multicalls, avoiding
// the overhead of rebuilding type definitions on every call.
var (
	aggregate3EncodeArgs gethabi.Arguments
	aggregate3DecodeArgs gethabi.Arguments
	aggregate3ArgsOnce   sync.Once
	aggregate3ArgsErr    error
)

// initAggregate3Args initializes the cached, pre-parsed go-ethereum Arguments
// for aggregate3 encoding and decoding. This is the key optimization: we parse
// the complex tuple[] type definition once and reuse the Arguments forever.
func initAggregate3Args() {
	aggregate3ArgsOnce.Do(func() {
		// Build encode args: tuple(address target, bool allowFailure, bytes callData)[]
		encodeType, err := gethabi.NewType("tuple[]", "", []gethabi.ArgumentMarshaling{
			{Name: "target", Type: "address"},
			{Name: "allowFailure", Type: "bool"},
			{Name: "callData", Type: "bytes"},
		})
		if err != nil {
			aggregate3ArgsErr = fmt.Errorf("failed to parse aggregate3 encode type: %w", err)
			return
		}
		aggregate3EncodeArgs = gethabi.Arguments{{Type: encodeType}}

		// Build decode args: tuple(bool success, bytes returnData)[]
		decodeType, err := gethabi.NewType("tuple[]", "", []gethabi.ArgumentMarshaling{
			{Name: "success", Type: "bool"},
			{Name: "returnData", Type: "bytes"},
		})
		if err != nil {
			aggregate3ArgsErr = fmt.Errorf("failed to parse aggregate3 decode type: %w", err)
			return
		}
		aggregate3DecodeArgs = gethabi.Arguments{{Type: decodeType}}
	})
}

// Cached aggregate3 selector - parsed once
var (
	aggregate3Selector     []byte
	aggregate3SelectorOnce sync.Once
)

func getAggregate3Selector() []byte {
	aggregate3SelectorOnce.Do(func() {
		aggregate3Selector = common.FromHex(constants.Aggregate3Signature)
	})
	return aggregate3Selector
}

// MulticallContract defines a contract call for multicall.
// This mirrors viem's ContractFunctionParameters type.
type MulticallContract struct {
	// Address is the contract address to call.
	Address common.Address

	// ABI is the contract ABI as JSON bytes, string, or *abi.ABI.
	ABI *abi.ABI

	// FunctionName is the name of the function to call.
	FunctionName string

	// Args are the function arguments.
	Args []any
}

// MulticallParameters contains the parameters for the Multicall action.
// This mirrors viem's MulticallParameters type.
type MulticallParameters struct {
	// Contracts is the list of contract calls to execute.
	Contracts []MulticallContract

	// AllowFailure determines whether to continue if individual calls fail.
	// If true, failed calls will be marked with status "failure" but won't
	// stop the entire multicall. Default is true.
	AllowFailure *bool

	// BatchSize is the maximum size in bytes for each batch of calls.
	// Calls are chunked into batches based on their calldata size.
	// Default is 1024 bytes.
	BatchSize int

	// Deployless enables deployless multicall using bytecode execution.
	// This allows multicall on chains without a deployed multicall3 contract.
	Deployless bool

	// MulticallAddress overrides the default multicall3 contract address.
	MulticallAddress *common.Address

	// BlockNumber is the block number to execute the calls at.
	// Mutually exclusive with BlockTag.
	BlockNumber *uint64

	// BlockTag is the block tag to execute the calls at.
	// Mutually exclusive with BlockNumber.
	BlockTag BlockTag

	// MaxConcurrentChunks limits the number of concurrent chunk executions.
	// This prevents overwhelming RPC endpoints. Default is 4.
	// Set to 0 or negative for unlimited concurrency.
	MaxConcurrentChunks int
}

// MulticallResult represents the result of a single contract call in a multicall.
type MulticallResult struct {
	// Status is either "success" or "failure".
	Status string

	// Result contains the decoded return value(s) if Status is "success".
	Result any

	// Error contains the error if Status is "failure".
	Error error
}

// MulticallReturnType is the return type for the Multicall action.
type MulticallReturnType = []MulticallResult

// Call3 represents a single call in the aggregate3 function.
// The struct tags match the Multicall3 ABI parameter names.
type Call3 struct {
	Target       common.Address `abi:"target"`
	AllowFailure bool           `abi:"allowFailure"`
	CallData     []byte         `abi:"callData"`
}

// aggregate3Result represents the result from aggregate3.
type aggregate3Result struct {
	Success    bool
	ReturnData []byte
}

// chunkResult holds the result of executing a chunk.
type chunkResult struct {
	Results []aggregate3Result
	Err     error
}

// encodeJob represents a contract to encode.
type encodeJob struct {
	index    int
	contract MulticallContract
}

// encodeResult represents the result of encoding a contract call.
type encodeResult struct {
	index     int
	call      Call3
	parsedABI *abi.ABI
	err       error
}

// decodeJob represents a result to decode.
type decodeJob struct {
	index       int
	aggResult   aggregate3Result
	contract    MulticallContract
	parsedABI   *abi.ABI
	encodeError error
	callData    []byte
}

// decodeResult represents the result of decoding.
type decodeResult struct {
	index  int
	result MulticallResult
}

// chunkJob represents a chunk to execute.
type chunkJob struct {
	chunkIndex int
	chunk      []Call3
}

// getNumWorkers returns the number of workers to use based on job count.
// Uses GOMAXPROCS but caps at job count to avoid idle workers.
func getNumWorkers(numJobs int) int {
	workers := runtime.GOMAXPROCS(1)
	if workers > numJobs {
		workers = numJobs
	}
	if workers < 1 {
		workers = 1
	}
	return workers
}

// Multicall batches multiple contract function calls into a single RPC call
// using the multicall3 contract.
//
// When the client has Batch.Multicall configured, concurrent Multicall() calls
// within the Wait window are automatically aggregated into a single larger
// multicall RPC call. This mirrors viem's `batch: { multicall: { ... } }` behavior.
//
// This is equivalent to viem's `multicall` action.
//
// Example:
//
//	results, err := public.Multicall(ctx, client, public.MulticallParameters{
//	    Contracts: []public.MulticallContract{
//	        {
//	            Address:      tokenAddress,
//	            ABI:          erc20ABI,
//	            FunctionName: "balanceOf",
//	            Args:         []any{userAddress},
//	        },
//	        {
//	            Address:      tokenAddress,
//	            ABI:          erc20ABI,
//	            FunctionName: "totalSupply",
//	        },
//	    },
//	})
func Multicall(ctx context.Context, client Client, params MulticallParameters) (MulticallReturnType, error) {
	// Check if client has multicall batch aggregation enabled
	if batch := client.Batch(); batch != nil && batch.Multicall != nil {
		batcher := getMulticallBatcher(client, batch.Multicall)
		if batcher != nil {
			return batcher.Schedule(ctx, params)
		}
	}

	return multicallDirect(ctx, client, params)
}

// MulticallConcurrent is like Multicall but optimized for fan-out patterns where
// many goroutines call it simultaneously. When the client has Batch.Multicall
// configured, all concurrent calls within the wait window are merged into a single
// large aggregate3 RPC call.
//
// Use this instead of Multicall when you know multiple goroutines will call it
// concurrently (e.g., resolving N tokens in parallel).
func MulticallConcurrent(ctx context.Context, client Client, params MulticallParameters) (MulticallReturnType, error) {
	if batch := client.Batch(); batch != nil && batch.Multicall != nil {
		batcher := getMulticallBatcher(client, batch.Multicall)
		if batcher != nil {
			return batcher.ScheduleConcurrent(ctx, params)
		}
	}

	return multicallDirect(ctx, client, params)
}

// multicallDirect is the actual multicall implementation that executes immediately
// without batching. This is called directly by Multicall when batching is not
// enabled, and by the MulticallBatcher when flushing a batch.
func multicallDirect(ctx context.Context, client Client, params MulticallParameters) (MulticallReturnType, error) {
	// Set defaults
	allowFailure := true
	if params.AllowFailure != nil {
		allowFailure = *params.AllowFailure
	}

	batchSize := params.BatchSize
	if batchSize <= 0 {
		batchSize = 8192
	}

	maxConcurrent := params.MaxConcurrentChunks
	if maxConcurrent <= 0 {
		maxConcurrent = 10
	}

	// Resolve multicall address
	multicallAddress, err := resolveMulticallAddress(client, params)
	if err != nil && !params.Deployless {
		return nil, err
	}

	contracts := params.Contracts
	numContracts := len(contracts)

	// ============================================================
	// PHASE 1: Parallel Encoding with Workers
	// ============================================================
	encodedCalls := make([]Call3, numContracts)
	parsedABIs := make([]*abi.ABI, numContracts)
	encodeErrors := make([]error, numContracts)

	// For small batches, skip worker overhead
	if numContracts <= 100000 {
		for i, contract := range contracts {
			parsedABIs[i] = contract.ABI
			callData, encodeErr := contract.ABI.EncodeFunctionData(contract.FunctionName, contract.Args...)
			if encodeErr != nil {
				encodeErrors[i] = fmt.Errorf("failed to encode call for %q: %w", contract.FunctionName, encodeErr)
				encodedCalls[i] = Call3{Target: contract.Address, AllowFailure: true}
			} else {
				encodedCalls[i] = Call3{Target: contract.Address, AllowFailure: true, CallData: callData}
			}
		}
	} else {
		// Use worker pool for parallel encoding
		encodeJobs := make(chan encodeJob, numContracts)
		encodeResults := make(chan encodeResult, numContracts)

		numEncodeWorkers := getNumWorkers(numContracts)
		var encodeWg sync.WaitGroup
		encodeWg.Add(numEncodeWorkers)

		// Start encode workers
		for w := 0; w < numEncodeWorkers; w++ {
			go func() {
				defer encodeWg.Done()
				for job := range encodeJobs {
					parsedABI := job.contract.ABI
					callData, encodeErr := parsedABI.EncodeFunctionData(job.contract.FunctionName, job.contract.Args...)
					if encodeErr != nil {
						encodeResults <- encodeResult{
							index:     job.index,
							call:      Call3{Target: job.contract.Address, AllowFailure: true},
							parsedABI: parsedABI,
							err:       fmt.Errorf("failed to encode call for %q: %w", job.contract.FunctionName, encodeErr),
						}
					} else {
						encodeResults <- encodeResult{
							index:     job.index,
							call:      Call3{Target: job.contract.Address, AllowFailure: true, CallData: callData},
							parsedABI: parsedABI,
						}
					}
				}
			}()
		}

		// Send encode jobs
		for i, contract := range contracts {
			encodeJobs <- encodeJob{index: i, contract: contract}
		}
		close(encodeJobs)

		// Collect encode results in background
		go func() {
			encodeWg.Wait()
			close(encodeResults)
		}()

		for result := range encodeResults {
			encodedCalls[result.index] = result.call
			parsedABIs[result.index] = result.parsedABI
			encodeErrors[result.index] = result.err
		}
	}

	// ============================================================
	// PHASE 2: Chunk Calls and Execute with Workers
	// ============================================================
	chunkedCalls := chunkCalls(encodedCalls, batchSize)
	numChunks := len(chunkedCalls)
	chunkResults := make([]*chunkResult, numChunks)

	if numChunks == 1 {
		// Single chunk - no need for workers
		result, execErr := executeChunk(ctx, client, chunkedCalls[0], multicallAddress, params)
		chunkResults[0] = &chunkResult{Results: result, Err: execErr}
	} else {
		// Use worker pool for parallel RPC execution
		chunkJobs := make(chan chunkJob, numChunks)
		chunkResultsChan := make(chan struct {
			index  int
			result *chunkResult
		}, numChunks)

		numChunkWorkers := maxConcurrent
		if numChunkWorkers > numChunks {
			numChunkWorkers = numChunks
		}

		var chunkWg sync.WaitGroup
		chunkWg.Add(numChunkWorkers)

		// Start RPC execution workers
		for w := 0; w < numChunkWorkers; w++ {
			go func() {
				defer chunkWg.Done()
				for job := range chunkJobs {
					result, execErr := executeChunk(ctx, client, job.chunk, multicallAddress, params)
					chunkResultsChan <- struct {
						index  int
						result *chunkResult
					}{job.chunkIndex, &chunkResult{Results: result, Err: execErr}}
				}
			}()
		}

		// Send chunk jobs
		for i, chunk := range chunkedCalls {
			chunkJobs <- chunkJob{chunkIndex: i, chunk: chunk}
		}
		close(chunkJobs)

		// Collect results
		go func() {
			chunkWg.Wait()
			close(chunkResultsChan)
		}()

		for res := range chunkResultsChan {
			chunkResults[res.index] = res.result
		}
	}

	// ============================================================
	// PHASE 3: Build Decode Jobs from Chunk Results
	// ============================================================
	decodeJobs := make([]decodeJob, 0, numContracts)
	resultIndex := 0

	for chunkIdx, chunkRes := range chunkResults {
		chunkLen := len(chunkedCalls[chunkIdx])

		if chunkRes.Err != nil {
			// Chunk-level error - create failure jobs for all calls in chunk
			for j := 0; j < chunkLen; j++ {
				decodeJobs = append(decodeJobs, decodeJob{
					index:       resultIndex,
					aggResult:   aggregate3Result{Success: false},
					contract:    contracts[resultIndex],
					encodeError: chunkRes.Err,
				})
				resultIndex++
			}
			continue
		}

		// Process individual results
		for j, aggResult := range chunkRes.Results {
			decodeJobs = append(decodeJobs, decodeJob{
				index:       resultIndex,
				aggResult:   aggResult,
				contract:    contracts[resultIndex],
				parsedABI:   parsedABIs[resultIndex],
				encodeError: encodeErrors[resultIndex],
				callData:    chunkedCalls[chunkIdx][j].CallData,
			})
			resultIndex++
		}
	}

	// ============================================================
	// PHASE 4: Parallel Decoding with Workers
	// ============================================================
	results := make(MulticallReturnType, numContracts)

	if numContracts <= 10000000 {
		// Small batch - decode sequentially
		for _, job := range decodeJobs {
			results[job.index] = decodeOneResult(job, allowFailure)
		}
	} else {
		// Use worker pool for parallel decoding
		decodeJobsChan := make(chan decodeJob, len(decodeJobs))
		decodeResultsChan := make(chan decodeResult, len(decodeJobs))

		numDecodeWorkers := getNumWorkers(len(decodeJobs))
		var decodeWg sync.WaitGroup
		decodeWg.Add(numDecodeWorkers)

		// Start decode workers
		for w := 0; w < numDecodeWorkers; w++ {
			go func() {
				defer decodeWg.Done()
				for job := range decodeJobsChan {
					decodeResultsChan <- decodeResult{
						index:  job.index,
						result: decodeOneResult(job, allowFailure),
					}
				}
			}()
		}

		// Send decode jobs
		for _, job := range decodeJobs {
			decodeJobsChan <- job
		}
		close(decodeJobsChan)

		// Collect decode results
		go func() {
			decodeWg.Wait()
			close(decodeResultsChan)
		}()

		for res := range decodeResultsChan {
			results[res.index] = res.result
		}
	}

	// Check for early failure if allowFailure is false
	if !allowFailure {
		for _, r := range results {
			if r.Status == "failure" {
				return nil, r.Error
			}
		}
	}

	return results, nil
}

// decodeOneResult decodes a single multicall result.
func decodeOneResult(job decodeJob, allowFailure bool) MulticallResult {
	// Check for encode errors first
	if job.encodeError != nil {
		return MulticallResult{Status: "failure", Error: job.encodeError}
	}

	// Check if the call succeeded
	if !job.aggResult.Success {
		return MulticallResult{Status: "failure", Error: &RawContractError{Data: job.aggResult.ReturnData}}
	}

	// Check for empty calldata
	if len(job.callData) == 0 {
		return MulticallResult{Status: "failure", Error: &AbiDecodingZeroDataError{}}
	}

	// Decode the result
	decoded, decodeErr := job.parsedABI.DecodeFunctionResult(job.contract.FunctionName, job.aggResult.ReturnData)
	if decodeErr != nil {
		return MulticallResult{
			Status: "failure",
			Error:  fmt.Errorf("failed to decode result for %q: %w", job.contract.FunctionName, decodeErr),
		}
	}

	// Unwrap single return value
	var result any
	if len(decoded) == 1 {
		result = decoded[0]
	} else {
		result = decoded
	}

	return MulticallResult{Status: "success", Result: result}
}

// chunkCalls splits calls into chunks based on batch size.
// Pre-allocates slices for efficiency.
func chunkCalls(calls []Call3, batchSize int) [][]Call3 {
	if len(calls) == 0 {
		return nil
	}

	// If batchSize is 0 or negative, return all calls in a single chunk
	if batchSize <= 0 {
		return [][]Call3{calls}
	}

	// Estimate number of chunks (avg call ~36 bytes for balanceOf)
	estimatedChunks := (len(calls)*36)/batchSize + 1
	chunks := make([][]Call3, 0, estimatedChunks)

	// Pre-allocate current chunk with reasonable capacity
	currentChunk := make([]Call3, 0, min(len(calls), batchSize/36+1))
	currentSize := 0

	for _, call := range calls {
		callSize := len(call.CallData)
		if callSize == 0 {
			callSize = 2 // "0x" placeholder
		}

		// Check if we need a new chunk
		if currentSize+callSize > batchSize && len(currentChunk) > 0 {
			chunks = append(chunks, currentChunk)
			currentChunk = make([]Call3, 0, min(len(calls)-len(chunks)*len(currentChunk), batchSize/36+1))
			currentSize = 0
		}

		currentChunk = append(currentChunk, call)
		currentSize += callSize
	}

	// Add final chunk
	if len(currentChunk) > 0 {
		chunks = append(chunks, currentChunk)
	}

	return chunks
}

// executeChunk executes a single chunk of calls via multicall3.
func executeChunk(ctx context.Context, client Client, calls []Call3, multicallAddress *common.Address, params MulticallParameters) ([]aggregate3Result, error) {
	// Encode aggregate3 call
	calldata, err := encodeAggregate3(calls)
	if err != nil {
		return nil, fmt.Errorf("failed to encode aggregate3: %w", err)
	}

	// Build call request
	blockTag := resolveBlockTag(client, params.BlockNumber, params.BlockTag)

	var req callRequest
	var rpcParams []any

	if params.Deployless || multicallAddress == nil {
		// Deployless multicall - wrap in deployless bytecode
		deploylessData, deploylessErr := deployless.ToDeploylessCallViaBytecodeData(
			common.FromHex(constants.Multicall3Bytecode),
			calldata,
		)
		if deploylessErr != nil {
			return nil, fmt.Errorf("failed to encode deployless multicall: %w", deploylessErr)
		}
		req = callRequest{Data: hexutil.Encode(deploylessData)}
	} else {
		req = callRequest{
			To:   multicallAddress.Hex(),
			Data: hexutil.Encode(calldata),
		}
	}

	rpcParams = []any{req, blockTag}

	// Execute call
	resp, requestErr := client.Request(ctx, "eth_call", rpcParams...)
	if requestErr != nil {
		return nil, fmt.Errorf("eth_call failed: %w", requestErr)
	}

	var hexResult string
	if unmarshalErr := json.Unmarshal(resp.Result, &hexResult); unmarshalErr != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %w", unmarshalErr)
	}

	// Decode aggregate3 result
	resultData := common.FromHex(hexResult)
	return decodeAggregate3Result(resultData)
}

// encodeAggregate3 encodes calls for the aggregate3 function.
// Uses pre-parsed go-ethereum Arguments for direct packing, bypassing
// the generic EncodeAbiParameters path which re-parses types on every call.
func encodeAggregate3(calls []Call3) ([]byte, error) {
	initAggregate3Args()
	if aggregate3ArgsErr != nil {
		return nil, aggregate3ArgsErr
	}

	// Convert Call3 structs to the format go-ethereum expects for tuple[]
	tuples := make([]struct {
		Target       common.Address `abi:"target"`
		AllowFailure bool           `abi:"allowFailure"`
		CallData     []byte         `abi:"callData"`
	}, len(calls))
	for i, c := range calls {
		tuples[i].Target = c.Target
		tuples[i].AllowFailure = c.AllowFailure
		tuples[i].CallData = c.CallData
	}

	// Pack directly using cached Arguments
	encoded, err := aggregate3EncodeArgs.Pack(tuples)
	if err != nil {
		return nil, fmt.Errorf("failed to pack aggregate3: %w", err)
	}

	// Prepend cached aggregate3 selector (0x82ad56cb)
	selector := getAggregate3Selector()
	result := make([]byte, len(selector)+len(encoded))
	copy(result, selector)
	copy(result[len(selector):], encoded)

	return result, nil
}

// decodeAggregate3Result decodes the result from aggregate3.
// Uses pre-parsed go-ethereum Arguments and direct struct unpacking for speed,
// bypassing the generic DecodeAbiParameters path entirely.
func decodeAggregate3Result(data []byte) ([]aggregate3Result, error) {
	initAggregate3Args()
	if aggregate3ArgsErr != nil {
		return nil, aggregate3ArgsErr
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("empty aggregate3 result data")
	}

	// Unpack directly using cached Arguments
	unpacked, err := aggregate3DecodeArgs.Unpack(data)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack aggregate3 result: %w", err)
	}

	if len(unpacked) == 0 {
		return nil, fmt.Errorf("empty aggregate3 result")
	}

	// The unpacked result is a []struct{Success bool, ReturnData []byte}
	// go-ethereum returns this as a slice of anonymous structs
	// We need to use reflect to extract the values since the exact type is generated at runtime
	type resultStruct = struct {
		Success    bool
		ReturnData []byte
	}

	// Try direct type assertion for the common case
	if tuples, ok := unpacked[0].([]resultStruct); ok {
		results := make([]aggregate3Result, len(tuples))
		for i, t := range tuples {
			results[i] = aggregate3Result{Success: t.Success, ReturnData: t.ReturnData}
		}
		return results, nil
	}

	// Fallback: use reflect for go-ethereum's generated anonymous struct types
	return decodeAggregate3Reflect(unpacked[0])
}

// decodeAggregate3Reflect handles the case where go-ethereum returns an anonymous struct
// type that can't be directly asserted. Uses reflect to extract fields by name.
func decodeAggregate3Reflect(raw any) ([]aggregate3Result, error) {
	rv := reflect.ValueOf(raw)
	if rv.Kind() != reflect.Slice {
		return nil, fmt.Errorf("unexpected aggregate3 result type: %T (expected slice)", raw)
	}

	results := make([]aggregate3Result, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		elem := rv.Index(i)
		if elem.Kind() == reflect.Struct {
			successField := elem.FieldByName("Success")
			returnDataField := elem.FieldByName("ReturnData")
			if successField.IsValid() && returnDataField.IsValid() {
				results[i] = aggregate3Result{
					Success:    successField.Bool(),
					ReturnData: returnDataField.Bytes(),
				}
				continue
			}
		}
		return nil, fmt.Errorf("invalid aggregate3 result tuple at index %d: %T", i, elem.Interface())
	}

	return results, nil
}

// resolveMulticallAddress determines the multicall3 contract address.
func resolveMulticallAddress(client Client, params MulticallParameters) (*common.Address, error) {
	// Use provided address if specified
	if params.MulticallAddress != nil {
		return params.MulticallAddress, nil
	}

	// Deployless doesn't need an address
	if params.Deployless {
		return nil, nil
	}

	// Get from chain config
	chain := client.Chain()
	if chain == nil {
		return nil, &ChainNotConfiguredError{}
	}

	if chain.Contracts == nil || chain.Contracts.Multicall3 == nil {
		return nil, &ChainDoesNotSupportContractError{
			ChainID:      chain.ID,
			ContractName: "multicall3",
		}
	}

	// Check block number constraint
	if params.BlockNumber != nil && chain.Contracts.Multicall3.BlockCreated != nil {
		if *params.BlockNumber < *chain.Contracts.Multicall3.BlockCreated {
			return nil, &ChainDoesNotSupportContractError{
				ChainID:      chain.ID,
				ContractName: "multicall3",
				BlockNumber:  params.BlockNumber,
			}
		}
	}

	return &chain.Contracts.Multicall3.Address, nil
}

// parseABIParam parses the ABI parameter which can be []byte, string, or *abi.ABI.
func parseABIParam(abiParam any) (*abi.ABI, error) {
	switch v := abiParam.(type) {
	case *abi.ABI:
		return v, nil
	case []byte:
		return abi.Parse(v)
	case string:
		return abi.Parse([]byte(v))
	default:
		return nil, fmt.Errorf("ABI must be []byte, string, or *abi.ABI, got %T", abiParam)
	}
}

// AbiDecodingZeroDataError is returned when trying to decode zero data.
type AbiDecodingZeroDataError struct{}

func (e *AbiDecodingZeroDataError) Error() string {
	return "cannot decode zero data (0x) - the function may have reverted"
}
