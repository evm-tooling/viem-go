package bench

import (
	"testing"

	"github.com/ChefBingbong/viem-go/utils/signature"
)

// Test signature from Anvil account 0 signing "hello world"
// Account: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
const testSignature = "0x6e100a352ec6ad1b70802290e18aeed190704973570f3b8ed42cb9808e2ea6bf4a90a229a244495b41890987806fcbd2d5d23fc0dbe5f5256c2613c039d76db81c"
const testAddress = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"

// --- HashMessage benchmarks ---

func BenchmarkSignature_HashMessage(b *testing.B) {
	msg := signature.NewSignableMessage("hello world")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = signature.HashMessage(msg)
	}
}

func BenchmarkSignature_HashMessageLong(b *testing.B) {
	longMsg := "The quick brown fox jumps over the lazy dog. " +
		"This is a much longer message that simulates real-world signing scenarios " +
		"where users might sign terms of service, governance proposals, or other text content."
	msg := signature.NewSignableMessage(longMsg)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = signature.HashMessage(msg)
	}
}

// --- RecoverMessageAddress benchmarks ---

func BenchmarkSignature_RecoverAddress(b *testing.B) {
	msg := signature.NewSignableMessage("hello world")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := signature.RecoverMessageAddress(msg, testSignature)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// --- VerifyMessage benchmarks ---

func BenchmarkSignature_VerifyMessage(b *testing.B) {
	msg := signature.NewSignableMessage("hello world")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := signature.VerifyMessage(testAddress, msg, testSignature)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// --- ParseSignature benchmarks ---

func BenchmarkSignature_ParseSignature(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := signature.ParseSignature(testSignature)
		if err != nil {
			b.Fatal(err)
		}
	}
}
