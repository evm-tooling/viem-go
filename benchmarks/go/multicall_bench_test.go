// Package bench provides cross-language benchmarks for viem-go.
//
// These benchmarks mirror the TypeScript viem benchmarks for fair comparison.
//
// IMPORTANT: All benchmarks use BatchSize: 0 to disable chunking,
// ensuring a single RPC call for fair comparison with TypeScript.
package bench

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/_typed/templates/erc200"
	"github.com/ChefBingbong/viem-go/actions/public"
)

// ERC20 ABI for multicall benchmarks
var multicallERC20ABI = erc200.MustParsedABI()

// Additional token addresses for multicall benchmarks (Mainnet)
var (
	wethAddress = common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")
	daiAddress  = common.HexToAddress("0x6B175474E89094C44Da98b954EecdB131e560f7")
)

// multicallBenchConfig returns the standard config for fair benchmark comparison.
// BatchSize: 0 disables chunking, ensuring a single RPC call.
func multicallBenchConfig(contracts []public.MulticallContract) public.MulticallParameters {
	return public.MulticallParameters{
		Contracts:           contracts,
		BatchSize:           0, // Disable chunking - ensures single RPC call
		MaxConcurrentChunks: 1, // Safety: only 1 chunk anyway with BatchSize=0
	}
}

// BenchmarkMulticall_Basic benchmarks a simple multicall with 3 calls.
// This tests the basic overhead of using multicall vs individual calls.
func BenchmarkMulticall_Basic(b *testing.B) {
	params := multicallBenchConfig([]public.MulticallContract{
		{Address: usdcAddress, ABI: multicallERC20ABI, FunctionName: "name"},
		{Address: usdcAddress, ABI: multicallERC20ABI, FunctionName: "symbol"},
		{Address: usdcAddress, ABI: multicallERC20ABI, FunctionName: "decimals"},
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := public.Multicall(benchCtx, benchClient, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMulticall_WithArgs benchmarks multicall with function arguments.
// This tests calls that require parameter encoding (balanceOf).
func BenchmarkMulticall_WithArgs(b *testing.B) {
	params := multicallBenchConfig([]public.MulticallContract{
		{Address: usdcAddress, ABI: multicallERC20ABI, FunctionName: "balanceOf", Args: []any{vitalikAddress}},
		{Address: usdcAddress, ABI: multicallERC20ABI, FunctionName: "balanceOf", Args: []any{anvilAccount0}},
		{Address: usdcAddress, ABI: multicallERC20ABI, FunctionName: "balanceOf", Args: []any{usdcAddress}},
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := public.Multicall(benchCtx, benchClient, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMulticall_MultiContract benchmarks multicall across multiple contracts.
// This tests the common use case of querying multiple tokens in one call.
func BenchmarkMulticall_MultiContract(b *testing.B) {
	params := multicallBenchConfig([]public.MulticallContract{
		{Address: usdcAddress, ABI: multicallERC20ABI, FunctionName: "name"},
		{Address: wethAddress, ABI: multicallERC20ABI, FunctionName: "name"},
		{Address: usdcAddress, ABI: multicallERC20ABI, FunctionName: "balanceOf", Args: []any{vitalikAddress}},
		{Address: wethAddress, ABI: multicallERC20ABI, FunctionName: "balanceOf", Args: []any{vitalikAddress}},
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := public.Multicall(benchCtx, benchClient, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMulticall_10Calls benchmarks multicall with 10 calls.
// Tests medium-sized multicall batches.
func BenchmarkMulticall_10Calls(b *testing.B) {
	contracts := make([]public.MulticallContract, 10)
	for i := 0; i < 10; i++ {
		contracts[i] = public.MulticallContract{
			Address:      usdcAddress,
			ABI:          multicallERC20ABI,
			FunctionName: "balanceOf",
			Args:         []any{vitalikAddress},
		}
	}

	params := multicallBenchConfig(contracts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := public.Multicall(benchCtx, benchClient, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMulticall_30Calls benchmarks multicall with 30 calls.
// Tests larger multicall batches.
func BenchmarkMulticall_30Calls(b *testing.B) {
	contracts := make([]public.MulticallContract, 30)
	for i := 0; i < 30; i++ {
		contracts[i] = public.MulticallContract{
			Address:      usdcAddress,
			ABI:          multicallERC20ABI,
			FunctionName: "balanceOf",
			Args:         []any{vitalikAddress},
		}
	}

	params := multicallBenchConfig(contracts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := public.Multicall(benchCtx, benchClient, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMulticall_Deployless benchmarks deployless multicall.
// This tests multicall without requiring a deployed multicall3 contract.
func BenchmarkMulticall_Deployless(b *testing.B) {
	params := multicallBenchConfig([]public.MulticallContract{
		{Address: usdcAddress, ABI: multicallERC20ABI, FunctionName: "name"},
		{Address: usdcAddress, ABI: multicallERC20ABI, FunctionName: "symbol"},
		{Address: usdcAddress, ABI: multicallERC20ABI, FunctionName: "decimals"},
	})
	params.Deployless = true

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := public.Multicall(benchCtx, benchClient, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMulticall_TokenMetadata benchmarks the common pattern of
// fetching complete token metadata in a single call.
func BenchmarkMulticall_TokenMetadata(b *testing.B) {
	params := multicallBenchConfig([]public.MulticallContract{
		{Address: usdcAddress, ABI: multicallERC20ABI, FunctionName: "name"},
		{Address: usdcAddress, ABI: multicallERC20ABI, FunctionName: "symbol"},
		{Address: usdcAddress, ABI: multicallERC20ABI, FunctionName: "decimals"},
		{Address: usdcAddress, ABI: multicallERC20ABI, FunctionName: "totalSupply"},
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := public.Multicall(benchCtx, benchClient, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// ============================================================
// STRESS TESTS - Large batch sizes to test performance at scale
// ============================================================

// BenchmarkMulticall_50Calls stress tests multicall with 50 calls.
func BenchmarkMulticall_50Calls(b *testing.B) {
	contracts := make([]public.MulticallContract, 50)
	for i := 0; i < 50; i++ {
		contracts[i] = public.MulticallContract{
			Address:      usdcAddress,
			ABI:          multicallERC20ABI,
			FunctionName: "balanceOf",
			Args:         []any{vitalikAddress},
		}
	}

	params := multicallBenchConfig(contracts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := public.Multicall(benchCtx, benchClient, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMulticall_100Calls stress tests multicall with 100 calls.
func BenchmarkMulticall_100Calls(b *testing.B) {
	contracts := make([]public.MulticallContract, 100)
	for i := 0; i < 100; i++ {
		contracts[i] = public.MulticallContract{
			Address:      usdcAddress,
			ABI:          multicallERC20ABI,
			FunctionName: "balanceOf",
			Args:         []any{vitalikAddress},
		}
	}

	params := multicallBenchConfig(contracts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := public.Multicall(benchCtx, benchClient, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMulticall_200Calls stress tests multicall with 200 calls.
func BenchmarkMulticall_200Calls(b *testing.B) {
	contracts := make([]public.MulticallContract, 200)
	for i := 0; i < 200; i++ {
		contracts[i] = public.MulticallContract{
			Address:      usdcAddress,
			ABI:          multicallERC20ABI,
			FunctionName: "balanceOf",
			Args:         []any{vitalikAddress},
		}
	}

	params := multicallBenchConfig(contracts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := public.Multicall(benchCtx, benchClient, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMulticall_500Calls extreme stress test with 500 calls.
func BenchmarkMulticall_500Calls(b *testing.B) {
	contracts := make([]public.MulticallContract, 500)
	for i := 0; i < 500; i++ {
		contracts[i] = public.MulticallContract{
			Address:      usdcAddress,
			ABI:          multicallERC20ABI,
			FunctionName: "balanceOf",
			Args:         []any{vitalikAddress},
		}
	}

	params := multicallBenchConfig(contracts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := public.Multicall(benchCtx, benchClient, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMulticall_MixedContracts_100 tests 100 calls across multiple contracts.
// More realistic workload with varied targets.
func BenchmarkMulticall_MixedContracts_100(b *testing.B) {
	contracts := make([]public.MulticallContract, 100)
	for i := 0; i < 100; i++ {
		addr := usdcAddress
		if i%2 != 0 {
			addr = wethAddress
		}
		contracts[i] = public.MulticallContract{
			Address:      addr,
			ABI:          multicallERC20ABI,
			FunctionName: "balanceOf",
			Args:         []any{vitalikAddress},
		}
	}

	params := multicallBenchConfig(contracts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := public.Multicall(benchCtx, benchClient, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// ============================================================
// EXTREME STRESS TESTS - Maximum throughput benchmarks
// ============================================================

// BenchmarkMulticall_1000Calls tests 1000 calls in a single RPC.
func BenchmarkMulticall_1000Calls(b *testing.B) {
	contracts := make([]public.MulticallContract, 1000)
	for i := 0; i < 1000; i++ {
		contracts[i] = public.MulticallContract{
			Address:      usdcAddress,
			ABI:          multicallERC20ABI,
			FunctionName: "balanceOf",
			Args:         []any{vitalikAddress},
		}
	}

	params := multicallBenchConfig(contracts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := public.Multicall(benchCtx, benchClient, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMulticall_10000Calls_SingleRPC tests 10,000 calls in a single massive RPC.
// Tests maximum payload size handling.
func BenchmarkMulticall_10000Calls_SingleRPC(b *testing.B) {
	contracts := make([]public.MulticallContract, 10000)
	for i := 0; i < 10000; i++ {
		contracts[i] = public.MulticallContract{
			Address:      usdcAddress,
			ABI:          multicallERC20ABI,
			FunctionName: "balanceOf",
			Args:         []any{vitalikAddress},
		}
	}

	// Single RPC - no chunking
	params := public.MulticallParameters{
		Contracts:           contracts,
		BatchSize:           0, // Disable chunking
		MaxConcurrentChunks: 1,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := public.Multicall(benchCtx, benchClient, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMulticall_10000Calls_Chunked tests 10,000 calls with optimized chunking.
// Uses large batches for parallel RPC execution.
// batchSize: 8192 bytes (~228 calls per chunk) = ~44 chunks
func BenchmarkMulticall_10000Calls_Chunked(b *testing.B) {
	contracts := make([]public.MulticallContract, 10000)
	for i := 0; i < 10000; i++ {
		contracts[i] = public.MulticallContract{
			Address:      usdcAddress,
			ABI:          multicallERC20ABI,
			FunctionName: "balanceOf",
			Args:         []any{vitalikAddress},
		}
	}

	// Optimized chunking for parallel execution
	params := public.MulticallParameters{
		Contracts:           contracts,
		BatchSize:           8192, // Large batches
		MaxConcurrentChunks: 20,   // Parallel RPC calls
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := public.Multicall(benchCtx, benchClient, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMulticall_10000Calls_AggressiveChunking tests 10,000 calls with aggressive chunking.
// Uses smaller batches for maximum parallelism.
// batchSize: 2048 bytes (~57 calls per chunk) = ~175 chunks
func BenchmarkMulticall_10000Calls_AggressiveChunking(b *testing.B) {
	contracts := make([]public.MulticallContract, 10000)
	for i := 0; i < 10000; i++ {
		contracts[i] = public.MulticallContract{
			Address:      usdcAddress,
			ABI:          multicallERC20ABI,
			FunctionName: "balanceOf",
			Args:         []any{vitalikAddress},
		}
	}

	// Aggressive chunking for maximum parallelism
	params := public.MulticallParameters{
		Contracts:           contracts,
		BatchSize:           2048, // Smaller batches = more chunks
		MaxConcurrentChunks: 10,   // High parallelism
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := public.Multicall(benchCtx, benchClient, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}
