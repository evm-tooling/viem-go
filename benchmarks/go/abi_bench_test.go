package bench

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/abi"
)

// Shared ABI fixtures for ABI benchmarks
var (
	benchERC20ABI *abi.ABI

	// Pre-encoded balanceOf return data (uint256 = 1000000)
	balanceOfReturnData []byte

	// Pre-encoded transfer calldata
	preEncodedTransferData []byte
)

func init() {
	var err error
	benchERC20ABI, err = abi.Parse([]byte(`[
		{"name":"name","type":"function","inputs":[],"outputs":[{"type":"string"}],"stateMutability":"view"},
		{"name":"symbol","type":"function","inputs":[],"outputs":[{"type":"string"}],"stateMutability":"view"},
		{"name":"decimals","type":"function","inputs":[],"outputs":[{"type":"uint8"}],"stateMutability":"view"},
		{"name":"totalSupply","type":"function","inputs":[],"outputs":[{"type":"uint256"}],"stateMutability":"view"},
		{"name":"balanceOf","type":"function","inputs":[{"name":"owner","type":"address"}],"outputs":[{"type":"uint256"}],"stateMutability":"view"},
		{"name":"allowance","type":"function","inputs":[{"name":"owner","type":"address"},{"name":"spender","type":"address"}],"outputs":[{"type":"uint256"}],"stateMutability":"view"},
		{"name":"transfer","type":"function","inputs":[{"name":"to","type":"address"},{"name":"amount","type":"uint256"}],"outputs":[{"type":"bool"}],"stateMutability":"nonpayable"},
		{"name":"approve","type":"function","inputs":[{"name":"spender","type":"address"},{"name":"amount","type":"uint256"}],"outputs":[{"type":"bool"}],"stateMutability":"nonpayable"},
		{"name":"transferFrom","type":"function","inputs":[{"name":"from","type":"address"},{"name":"to","type":"address"},{"name":"amount","type":"uint256"}],"outputs":[{"type":"bool"}],"stateMutability":"nonpayable"},
		{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"},
		{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"},{"indexed":true,"name":"spender","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Approval","type":"event"}
	]`))
	if err != nil {
		panic("failed to parse ERC20 ABI for benchmarks: " + err.Error())
	}

	// Pre-encode a balanceOf return value (uint256 = 1000000)
	balanceOfReturnData = common.LeftPadBytes(big.NewInt(1000000).Bytes(), 32)

	// Pre-encode transfer calldata
	preEncodedTransferData, err = benchERC20ABI.EncodeFunctionData("transfer",
		common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"),
		big.NewInt(1000000),
	)
	if err != nil {
		panic("failed to pre-encode transfer data: " + err.Error())
	}
}

// --- Encode benchmarks ---

func BenchmarkAbi_EncodeSimple(b *testing.B) {
	addr := common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := benchERC20ABI.EncodeFunctionData("balanceOf", addr)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAbi_EncodeComplex(b *testing.B) {
	to := common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
	amount := big.NewInt(1000000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := benchERC20ABI.EncodeFunctionData("transfer", to, amount)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAbi_EncodeMultiArg(b *testing.B) {
	from := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
	to := common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
	amount := big.NewInt(1000000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := benchERC20ABI.EncodeFunctionData("transferFrom", from, to, amount)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// --- Decode benchmarks ---

func BenchmarkAbi_DecodeResult(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := benchERC20ABI.DecodeFunctionResult("balanceOf", balanceOfReturnData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// --- EncodePacked benchmarks ---

func BenchmarkAbi_EncodePacked(b *testing.B) {
	types := []string{"address", "uint256"}
	values := []any{
		"0x14dC79964da2C08b23698B3D3cc7Ca32193d9955",
		big.NewInt(420),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := abi.EncodePacked(types, values)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAbi_EncodePackedMulti(b *testing.B) {
	types := []string{"address", "string", "uint256", "bool"}
	values := []any{
		"0x14dC79964da2C08b23698B3D3cc7Ca32193d9955",
		"hello world",
		big.NewInt(420),
		true,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := abi.EncodePacked(types, values)
		if err != nil {
			b.Fatal(err)
		}
	}
}
