package bench

import (
	"testing"

	"github.com/ChefBingbong/viem-go/utils/address"
)

// --- IsAddress benchmarks ---

func BenchmarkAddress_IsAddress(b *testing.B) {
	addr := "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = address.IsAddress(addr)
	}
}

func BenchmarkAddress_IsAddressLower(b *testing.B) {
	addr := "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = address.IsAddress(addr)
	}
}

// --- GetAddress (checksum) benchmarks ---

func BenchmarkAddress_Checksum(b *testing.B) {
	addr := "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := address.GetAddress(addr)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// --- GetContractAddress (CREATE) benchmarks ---

func BenchmarkAddress_Create(b *testing.B) {
	opts := address.GetCreateAddressOptions{
		From:  "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		Nonce: 1,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := address.GetCreateAddress(opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// --- GetContractAddress (CREATE2) benchmarks ---

func BenchmarkAddress_Create2(b *testing.B) {
	salt := make([]byte, 32)
	salt[31] = 1
	bytecode := []byte{0x60, 0x80, 0x60, 0x40, 0x52, 0x34, 0x80, 0x15} // minimal bytecode
	opts := address.GetCreate2AddressOptions{
		From:     "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		Salt:     salt,
		Bytecode: bytecode,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := address.GetCreate2Address(opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}
