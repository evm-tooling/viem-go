package bench

import (
	"testing"

	"github.com/ChefBingbong/viem-go/utils/ens"
)

// --- Namehash benchmarks ---

func BenchmarkEns_Namehash(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ens.Namehash("vitalik.eth")
	}
}

func BenchmarkEns_NamehashDeep(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ens.Namehash("sub.domain.vitalik.eth")
	}
}

// --- Labelhash benchmarks ---

func BenchmarkEns_Labelhash(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ens.Labelhash("vitalik")
	}
}

// --- Normalize benchmarks ---

func BenchmarkEns_Normalize(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ens.Normalize("Vitalik.ETH")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEns_NormalizeLong(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ens.Normalize("My.Long.SubDomain.Name.vitalik.eth")
		if err != nil {
			b.Fatal(err)
		}
	}
}
