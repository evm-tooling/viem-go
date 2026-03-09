// Error Examples - viem-go
//
// Full coverage of custom errors integrated across 4 phases:
//
// Phase 1: GetNodeError, GetTransactionError, GetEstimateGasError
//   - EstimateGas with reverting call -> EstimateGasExecutionError
//
// Phase 2: GetCallError, GetContractError
//   - CreateAccessList with reverting call -> CallExecutionError
//   - EstimateContractGas with reverting transfer -> ContractFunctionExecutionError
//   - SimulateContract with reverting transfer -> ContractFunctionExecutionError
//
// Phase 3: ContractFunctionZeroDataError, AbiDecodingZeroDataError
//   - ReadContract on EOA (returns 0x) -> ContractFunctionZeroDataError
//
// Phase 4: RPC Error Code Mapping
//   - Non-existent RPC method -> MethodNotFoundRpcError
//
// Also: Invalid params, contract reverts, CallExecutionError
package main

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/abi"
	"github.com/ChefBingbong/viem-go/actions/public"
	"github.com/ChefBingbong/viem-go/chain/definitions"
	"github.com/ChefBingbong/viem-go/client"
	"github.com/ChefBingbong/viem-go/client/transport"
	"github.com/ChefBingbong/viem-go/contracts/erc20"
	pkgerrors "github.com/ChefBingbong/viem-go/errors"
)

var (
	usdcAddress = common.HexToAddress("0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359")
	testAddress = common.HexToAddress("0x1234567890123456789012345678901234567890")
	zeroAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")
)

func main() {
	ctx := context.Background()

	erc20ABI, err := abi.Parse([]byte(erc20.ContractABI))
	if err != nil {
		fmt.Printf("Failed to parse ERC20 ABI: %v\n", err)
		return
	}

	fmt.Println("============================================================")
	fmt.Println("  Error Path Examples (viem-go)")
	fmt.Println("============================================================")

	publicClient, err := client.CreatePublicClient(client.PublicClientConfig{
		Chain:     &definitions.Polygon,
		Transport: transport.HTTP("https://rpc-mainnet.matic.quiknode.pro"),
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}
	defer publicClient.Close()

	// 1. Invalid params: code + to (mutually exclusive)
	fmt.Println("\n--- 1. Invalid Params: code + to (mutually exclusive) ---")
	_, err = public.Call(ctx, publicClient, public.CallParameters{
		Code: []byte{0x60, 0x2a, 0x60, 0x00, 0x52, 0x60, 0x20, 0x60, 0x00, 0xf3},
		To:   &usdcAddress,
		Data: []byte{0x06, 0xfd, 0xde, 0x03},
	})
	if err != nil {
		var invalidErr *public.InvalidCallParamsError
		if errors.As(err, &invalidErr) {
			fmt.Printf("  ✓ Caught InvalidCallParamsError: %v\n", err)
		} else {
			fmt.Printf("  Error (unexpected type): %v\n", err)
		}
	}

	// 2. Invalid params: code + factory (mutually exclusive)
	fmt.Println("\n--- 2. Invalid Params: code + factory (mutually exclusive) ---")
	factoryAddr := common.HexToAddress("0x4e59b44847b379578588920cA78FbF26c0B4956C")
	_, err = public.Call(ctx, publicClient, public.CallParameters{
		Code:        []byte{0x60, 0x2a, 0x60, 0x00, 0x52},
		Factory:     &factoryAddr,
		FactoryData: []byte{0x12, 0x34},
	})
	if err != nil {
		var invalidErr *public.InvalidCallParamsError
		if errors.As(err, &invalidErr) {
			fmt.Printf("  ✓ Caught InvalidCallParamsError: %v\n", err)
		} else {
			fmt.Printf("  Error (unexpected type): %v\n", err)
		}
	}

	// 3. Contract revert: transfer exceeds balance (ERC20)
	fmt.Println("\n--- 3. Contract Revert: transfer exceeds balance ---")
	_, err = publicClient.ReadContract(ctx, client.ReadContractOptions{
		Address:      usdcAddress,
		ABI:          erc20.ContractABI,
		FunctionName: "transfer",
		Args:         []any{testAddress, new(big.Int).Exp(big.NewInt(10), big.NewInt(60), nil)},
	})
	if err != nil {
		var execErr *pkgerrors.ContractFunctionExecutionError
		if errors.As(err, &execErr) {
			fmt.Printf("  ✓ Caught ContractFunctionExecutionError for %q\n", execErr.FunctionName)
		}
		var revertedErr *pkgerrors.ContractFunctionRevertedError
		if errors.As(err, &revertedErr) {
			fmt.Printf("  ✓ Decoded revert - Reason: %q\n", revertedErr.Reason)
		}
		fmt.Printf("  Full error: %v\n", err)
	}

	// 4. Contract revert: transfer from zero address (simulate as 0x0...0)
	fmt.Println("\n--- 4. Contract Revert: transfer from zero address ---")
	transferData, _ := client.EncodeFunctionData(client.EncodeFunctionDataOptions{
		ABI:          erc20.ContractABI,
		FunctionName: "transfer",
		Args:         []any{common.HexToAddress("0x0000000000000000000000000000000000000001"), big.NewInt(1)},
	})
	_, err = public.Call(ctx, publicClient, public.CallParameters{
		Account: &zeroAddress,
		To:      &usdcAddress,
		Data:    transferData,
	})
	if err != nil {
		var execErr *pkgerrors.ContractFunctionExecutionError
		if errors.As(err, &execErr) {
			fmt.Printf("  ✓ Caught ContractFunctionExecutionError for %q\n", execErr.FunctionName)
		}
		var revertedErr *pkgerrors.ContractFunctionRevertedError
		if errors.As(err, &revertedErr) {
			fmt.Printf("  ✓ Decoded revert - Reason: %q, Signature: %s\n", revertedErr.Reason, revertedErr.Signature)
		}
		fmt.Printf("  Full error: %v\n", err)
	}

	// 5. Call to non-contract address (EOA with random data)
	fmt.Println("\n--- 5. Call to non-contract address (EOA) ---")
	_, err = public.Call(ctx, publicClient, public.CallParameters{
		To:   &testAddress,
		Data: common.Hex2Bytes("0x12345678"),
	})
	if err != nil {
		var callErr *pkgerrors.CallExecutionError
		if errors.As(err, &callErr) {
			fmt.Printf("  ✓ Caught CallExecutionError (wraps underlying error)\n")
		}
		fmt.Printf("  Error: %v\n", err)
	}

	// 6. Invalid function selector (non-existent function on contract)
	fmt.Println("\n--- 6. Invalid selector: non-existent function ---")
	_, err = public.Call(ctx, publicClient, public.CallParameters{
		To:   &usdcAddress,
		Data: common.Hex2Bytes("0xdeadbeef"), // Random selector
	})
	if err != nil {
		fmt.Printf("  Error: %v\n", err)
	}

	// --- Phase 1: GetEstimateGasError ---
	fmt.Println("\n--- Phase 1: EstimateGasExecutionError (GetEstimateGasError) ---")
	transferRevertData, _ := client.EncodeFunctionData(client.EncodeFunctionDataOptions{
		ABI:          erc20.ContractABI,
		FunctionName: "transfer",
		Args:         []any{testAddress, new(big.Int).Exp(big.NewInt(10), big.NewInt(60), nil)},
	})
	_, err = public.EstimateGas(ctx, publicClient, public.EstimateGasParameters{
		Account: &zeroAddress,
		To:      &usdcAddress,
		Data:    transferRevertData,
		Value:   big.NewInt(0),
	})
	if err != nil {
		var estErr *pkgerrors.EstimateGasExecutionError
		if errors.As(err, &estErr) {
			fmt.Printf("  ✓ Caught EstimateGasExecutionError\n")
		}
		var execReverted *pkgerrors.ExecutionRevertedError
		if errors.As(err, &execReverted) {
			fmt.Printf("  ✓ Cause: ExecutionRevertedError\n")
		}
		fmt.Printf("  Error: %v\n", err)
	}

	// --- Phase 2: CreateAccessList + GetCallError ---
	fmt.Println("\n--- Phase 2a: CreateAccessList -> CallExecutionError (GetCallError) ---")
	_, err = public.CreateAccessList(ctx, publicClient, public.CreateAccessListParameters{
		Account: &zeroAddress,
		To:      &usdcAddress,
		Data:    transferRevertData,
	})
	if err != nil {
		var callErr *pkgerrors.CallExecutionError
		if errors.As(err, &callErr) {
			fmt.Printf("  ✓ Caught CallExecutionError from CreateAccessList\n")
		}
		fmt.Printf("  Error: %v\n", err)
	}

	// --- Phase 2: EstimateContractGas + GetContractError ---
	fmt.Println("\n--- Phase 2b: EstimateContractGas -> ContractFunctionExecutionError (GetContractError) ---")
	_, err = public.EstimateContractGas(ctx, publicClient, public.EstimateContractGasParameters{
		Account:      &zeroAddress,
		Address:      usdcAddress,
		ABI:          erc20ABI,
		FunctionName: "transfer",
		Args:         []any{testAddress, new(big.Int).Exp(big.NewInt(10), big.NewInt(60), nil)},
	})
	if err != nil {
		var execErr *pkgerrors.ContractFunctionExecutionError
		if errors.As(err, &execErr) {
			fmt.Printf("  ✓ Caught ContractFunctionExecutionError for %q\n", execErr.FunctionName)
		}
		var revertedErr *pkgerrors.ContractFunctionRevertedError
		if errors.As(err, &revertedErr) {
			fmt.Printf("  ✓ Decoded revert - Reason: %q\n", revertedErr.Reason)
		}
		fmt.Printf("  Error: %v\n", err)
	}

	// --- Phase 2: SimulateContract + GetContractError ---
	fmt.Println("\n--- Phase 2c: SimulateContract -> ContractFunctionExecutionError (GetContractError) ---")
	_, err = public.SimulateContract(ctx, publicClient, public.SimulateContractParameters{
		Account:      &zeroAddress,
		Address:      usdcAddress,
		ABI:          erc20ABI,
		FunctionName: "transfer",
		Args:         []any{testAddress, new(big.Int).Exp(big.NewInt(10), big.NewInt(60), nil)},
	})
	if err != nil {
		var execErr *pkgerrors.ContractFunctionExecutionError
		if errors.As(err, &execErr) {
			fmt.Printf("  ✓ Caught ContractFunctionExecutionError for %q\n", execErr.FunctionName)
		}
		fmt.Printf("  Error: %v\n", err)
	}

	// --- Phase 3: ContractFunctionZeroDataError (ReadContract on EOA returns 0x) ---
	fmt.Println("\n--- Phase 3: ContractFunctionZeroDataError (ReadContract on EOA) ---")
	_, err = publicClient.ReadContract(ctx, client.ReadContractOptions{
		Address:      testAddress, // EOA, not a contract
		ABI:          erc20.ContractABI,
		FunctionName: "balanceOf",
		Args:         []any{testAddress},
	})
	if err != nil {
		var zeroErr *pkgerrors.ContractFunctionZeroDataError
		if errors.As(err, &zeroErr) {
			fmt.Printf("  ✓ Caught ContractFunctionZeroDataError for %q\n", zeroErr.FunctionName)
		}
		var execErr *pkgerrors.ContractFunctionExecutionError
		if errors.As(err, &execErr) {
			fmt.Printf("  ✓ Caught ContractFunctionExecutionError (EOA may revert instead)\n")
		}
		fmt.Printf("  Error: %v\n", err)
	}

	// --- Phase 4: RPC Error Code Mapping (MethodNotFoundRpcError) ---
	fmt.Println("\n--- Phase 4: MethodNotFoundRpcError (non-existent RPC method) ---")
	_, err = publicClient.Request(ctx, "eth_nonexistent", "latest")
	if err != nil {
		var rpcErr *pkgerrors.RpcError
		if errors.As(err, &rpcErr) && rpcErr.Code == pkgerrors.RpcCodeMethodNotFound {
			fmt.Printf("  ✓ Caught RpcError (MethodNotFound, code %d)\n", rpcErr.Code)
		}
		fmt.Printf("  Error: %v\n", err)
	}

	// --- Summary ---
	fmt.Println("\n--- Error type classification summary ---")
	fmt.Println("  Phase 1: EstimateGasExecutionError, ExecutionRevertedError")
	fmt.Println("  Phase 2: CallExecutionError (CreateAccessList), ContractFunctionExecutionError (EstimateContractGas, SimulateContract)")
	fmt.Println("  Phase 3: ContractFunctionZeroDataError (ReadContract on EOA)")
	fmt.Println("  Phase 4: RpcError (MethodNotFound, InvalidParams, etc.)")
	fmt.Println("  Also: InvalidCallParamsError, CallExecutionError, ContractFunctionRevertedError")
	fmt.Println()
}
