package bench

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/abi"
)

var (
	eventABI *abi.ABI

	// Pre-built Transfer event log data
	transferTopics []common.Hash
	transferData   []byte
)

func init() {
	var err error
	eventABI, err = abi.Parse([]byte(`[
		{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"},
		{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"},{"indexed":true,"name":"spender","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Approval","type":"event"}
	]`))
	if err != nil {
		panic("failed to parse event ABI: " + err.Error())
	}

	// Transfer event signature: Transfer(address,address,uint256)
	transferTopics = []common.Hash{
		common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
		common.HexToHash("0x000000000000000000000000f39Fd6e51aad88F6F4ce6aB8827279cffFb92266"), // from
		common.HexToHash("0x000000000000000000000000d8dA6BF26964aF9D7eEd9e03E53415D37aA96045"), // to
	}
	// ABI-encoded value: 1000000 (uint256)
	transferData = common.LeftPadBytes(big.NewInt(1000000).Bytes(), 32)
}

// --- DecodeEventLog benchmarks ---

func BenchmarkEvent_DecodeTransfer(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := eventABI.DecodeEventLog(transferTopics, transferData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEvent_DecodeBatch10(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 10; j++ {
			_, err := eventABI.DecodeEventLog(transferTopics, transferData)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkEvent_DecodeBatch100(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			_, err := eventABI.DecodeEventLog(transferTopics, transferData)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}
