package bench

import (
	"testing"

	"github.com/ChefBingbong/viem-go/utils/hash"
)

// --- Keccak256 benchmarks ---

func BenchmarkHash_Keccak256Short(b *testing.B) {
	data := []byte("hello world")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = hash.Keccak256(data)
	}
}

func BenchmarkHash_Keccak256Long(b *testing.B) {
	// 1KB of data
	data := make([]byte, 1024)
	for i := range data {
		data[i] = byte(i % 256)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = hash.Keccak256(data)
	}
}

func BenchmarkHash_Keccak256Hex(b *testing.B) {
	hexData := "0x68656c6c6f20776f726c64"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = hash.Keccak256(hexData)
	}
}

// --- SHA256 benchmarks ---

func BenchmarkHash_Sha256Short(b *testing.B) {
	data := []byte("hello world")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = hash.Sha256(data)
	}
}

func BenchmarkHash_Sha256Long(b *testing.B) {
	data := make([]byte, 1024)
	for i := range data {
		data[i] = byte(i % 256)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = hash.Sha256(data)
	}
}

// --- Function selector benchmarks ---

func BenchmarkHash_FunctionSelector(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := hash.ToFunctionSelector("function transfer(address to, uint256 amount)")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHash_EventSelector(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := hash.ToEventSelector("event Transfer(address indexed from, address indexed to, uint256 amount)")
		if err != nil {
			b.Fatal(err)
		}
	}
}
