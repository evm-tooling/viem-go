package bench

import (
	"math/big"
	"testing"

	"github.com/ChefBingbong/viem-go/utils/unit"
)

// --- ParseEther benchmarks ---

func BenchmarkUnit_ParseEther(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := unit.ParseEther("1.5")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnit_ParseEtherLarge(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := unit.ParseEther("123456789.123456789012345678")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// --- FormatEther benchmarks ---

func BenchmarkUnit_FormatEther(b *testing.B) {
	wei := new(big.Int)
	wei.SetString("1500000000000000000", 10) // 1.5 ETH
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = unit.FormatEther(wei)
	}
}

// --- ParseUnits benchmarks ---

func BenchmarkUnit_ParseUnits6(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := unit.ParseUnits("100.50", 6) // USDC-style
		if err != nil {
			b.Fatal(err)
		}
	}
}

// --- ParseGwei benchmarks ---

func BenchmarkUnit_ParseGwei(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := unit.ParseGwei("20.5")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// --- FormatUnits benchmarks ---

func BenchmarkUnit_FormatUnits(b *testing.B) {
	val := big.NewInt(100500000) // 100.5 USDC
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = unit.FormatUnits(val, 6)
	}
}
