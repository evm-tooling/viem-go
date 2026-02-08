// Package bench provides cross-language benchmarks for viem-go.
//
// These benchmarks are designed to run against a shared Anvil instance
// for fair comparison with the TypeScript viem benchmarks.
package bench

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/chain/definitions"
	"github.com/ChefBingbong/viem-go/client"
	"github.com/ChefBingbong/viem-go/client/transport"
)

// Shared test fixtures
var (
	// benchClient is the shared PublicClient for all benchmarks
	benchClient *client.PublicClient

	// benchCtx is the shared context for all benchmarks
	benchCtx context.Context

	// Test addresses
	usdcAddress    = common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	vitalikAddress = common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")

	// Anvil's default funded account (account 0)
	anvilAccount0 = common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")

	// ERC20 function selectors
	nameSelector      = common.Hex2Bytes("06fdde03")                                                         // name()
	decimalsSelector  = common.Hex2Bytes("313ce567")                                                         // decimals()
	balanceOfSelector = common.Hex2Bytes("70a08231")                                                         // balanceOf(address)
	symbolSelector    = common.Hex2Bytes("95d89b41")                                                         // symbol()

	// Pre-encoded calldata for balanceOf(vitalikAddress)
	balanceOfVitalikData []byte
)

var rpcURL = os.Getenv("ANVIL_RPC_URL")


// Create the public client
var cl, err = client.CreatePublicClient(client.PublicClientConfig{
	Chain:     &definitions.Mainnet,
	Transport: transport.HTTP(rpcURL),
	Batch: &client.BatchOptions{
		Multicall: &client.MulticallBatchOptions{
			BatchSize: 8192,
			Wait:      16 * time.Millisecond,
		},
	},
})

func init() {
	// Pre-encode balanceOf(vitalikAddress) calldata
	balanceOfVitalikData = append(
		append([]byte{}, balanceOfSelector...),
		common.LeftPadBytes(vitalikAddress.Bytes(), 32)...,
	)
}

// TestMain sets up the shared PublicClient for all benchmarks.
// It reads the RPC URL from the ANVIL_RPC_URL environment variable,
// falling back to localhost:8545 if not set.
func TestMain(m *testing.M) {
	// Get RPC URL from environment or use default

if rpcURL == "" {
	rpcURL = "http://127.0.0.1:8545"
}
if err != nil {
	panic("failed to create benchmark client: " + err.Error())
}

benchClient = cl

	// Create shared context
	benchCtx = context.Background()

	// Verify connection by getting block number.
	// NOTE: We intentionally do not do any "warmup" calls here. Anvil warmup is
	// handled centrally by benchmarks/bench.sh for fairness across runtimes.
	// blockNum, err := benchClient.GetBlockNumber(benchCtx)
	// if err != nil {
	// 	panic("failed to connect to Anvil: " + err.Error())
	// }

	// Log connection info (visible with -v flag)
	// println("Benchmark client connected to:", rpcURL)
	// println("Current block number:", blockNum)

	// Run benchmarks
	code := m.Run()

	// Cleanup
	if benchClient != nil {
		benchClient.Close()
	}

	os.Exit(code)
}

// Helper function to create balanceOf calldata for any address
func encodeBalanceOf(addr common.Address) []byte {
	return append(
		append([]byte{}, balanceOfSelector...),
		common.LeftPadBytes(addr.Bytes(), 32)...,
	)
}
