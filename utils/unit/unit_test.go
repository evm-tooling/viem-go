package unit_test

import (
	"math/big"
	"testing"

	"github.com/ChefBingbong/viem-go/utils/unit"
)

func TestFormatUnits(t *testing.T) {
	tests := []struct {
		name     string
		value    *big.Int
		decimals int
		expected string
	}{
		{
			"420 with 9 decimals",
			big.NewInt(420000000000),
			9,
			"420",
		},
		{
			"1 with 18 decimals",
			new(big.Int).SetUint64(1000000000000000000),
			18,
			"1",
		},
		{
			"0.5 with 18 decimals",
			big.NewInt(500000000000000000),
			18,
			"0.5",
		},
		{
			"0.123456789 with 18 decimals",
			big.NewInt(123456789000000000),
			18,
			"0.123456789",
		},
		{
			"negative value",
			big.NewInt(-1000000000000000000),
			18,
			"-1",
		},
		{
			"zero",
			big.NewInt(0),
			18,
			"0",
		},
		{
			"small value",
			big.NewInt(1),
			18,
			"0.000000000000000001",
		},
		{
			"nil value",
			nil,
			18,
			"0",
		},
		{
			"69420.1234 with 4 decimals",
			big.NewInt(694201234),
			4,
			"69420.1234",
		},
		{
			"trailing zeros trimmed",
			big.NewInt(69420123400),
			6,
			"69420.1234",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := unit.FormatUnits(tt.value, tt.decimals)
			if result != tt.expected {
				t.Errorf("FormatUnits(%v, %d) = %s, want %s", tt.value, tt.decimals, result, tt.expected)
			}
		})
	}
}

func TestParseUnits(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		decimals int
		expected string
		hasError bool
	}{
		{
			"420 with 9 decimals",
			"420",
			9,
			"420000000000",
			false,
		},
		{
			"1 with 18 decimals",
			"1",
			18,
			"1000000000000000000",
			false,
		},
		{
			"1.5 with 18 decimals",
			"1.5",
			18,
			"1500000000000000000",
			false,
		},
		{
			"0.1 with 18 decimals",
			"0.1",
			18,
			"100000000000000000",
			false,
		},
		{
			"0.123456789 with 18 decimals",
			"0.123456789",
			18,
			"123456789000000000",
			false,
		},
		{
			"negative value",
			"-1",
			18,
			"-1000000000000000000",
			false,
		},
		{
			"zero",
			"0",
			18,
			"0",
			false,
		},
		{
			"invalid - letters",
			"abc",
			18,
			"",
			true,
		},
		{
			"invalid - multiple decimals",
			"1.2.3",
			18,
			"",
			true,
		},
		{
			"69420.1234 with 4 decimals",
			"69420.1234",
			4,
			"694201234",
			false,
		},
		{
			"rounding up",
			"1.99999999999999999999",
			18,
			"2000000000000000000",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := unit.ParseUnits(tt.value, tt.decimals)
			if tt.hasError && err == nil {
				t.Errorf("ParseUnits(%q, %d) expected error", tt.value, tt.decimals)
				return
			}
			if !tt.hasError && err != nil {
				t.Errorf("ParseUnits(%q, %d) unexpected error: %v", tt.value, tt.decimals, err)
				return
			}
			if !tt.hasError && result.String() != tt.expected {
				t.Errorf("ParseUnits(%q, %d) = %s, want %s", tt.value, tt.decimals, result.String(), tt.expected)
			}
		})
	}
}

func TestFormatEther(t *testing.T) {
	tests := []struct {
		name     string
		wei      *big.Int
		expected string
	}{
		{
			"1 ether",
			new(big.Int).SetUint64(1000000000000000000),
			"1",
		},
		{
			"1.5 ether",
			big.NewInt(1500000000000000000),
			"1.5",
		},
		{
			"0.1 ether",
			big.NewInt(100000000000000000),
			"0.1",
		},
		{
			"420 ether",
			func() *big.Int {
				v, _ := new(big.Int).SetString("420000000000000000000", 10)
				return v
			}(),
			"420",
		},
		{
			"0 ether",
			big.NewInt(0),
			"0",
		},
		{
			"1 wei",
			big.NewInt(1),
			"0.000000000000000001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := unit.FormatEther(tt.wei)
			if result != tt.expected {
				t.Errorf("FormatEther(%v) = %s, want %s", tt.wei, result, tt.expected)
			}
		})
	}
}

func TestParseEther(t *testing.T) {
	tests := []struct {
		name     string
		ether    string
		expected string
		hasError bool
	}{
		{
			"1 ether",
			"1",
			"1000000000000000000",
			false,
		},
		{
			"1.5 ether",
			"1.5",
			"1500000000000000000",
			false,
		},
		{
			"0.1 ether",
			"0.1",
			"100000000000000000",
			false,
		},
		{
			"420 ether",
			"420",
			"420000000000000000000",
			false,
		},
		{
			"0 ether",
			"0",
			"0",
			false,
		},
		{
			"invalid",
			"abc",
			"",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := unit.ParseEther(tt.ether)
			if tt.hasError && err == nil {
				t.Errorf("ParseEther(%q) expected error", tt.ether)
				return
			}
			if !tt.hasError && err != nil {
				t.Errorf("ParseEther(%q) unexpected error: %v", tt.ether, err)
				return
			}
			if !tt.hasError && result.String() != tt.expected {
				t.Errorf("ParseEther(%q) = %s, want %s", tt.ether, result.String(), tt.expected)
			}
		})
	}
}

func TestFormatGwei(t *testing.T) {
	tests := []struct {
		name     string
		wei      *big.Int
		expected string
	}{
		{
			"1 gwei",
			big.NewInt(1000000000),
			"1",
		},
		{
			"1.5 gwei",
			big.NewInt(1500000000),
			"1.5",
		},
		{
			"0.1 gwei",
			big.NewInt(100000000),
			"0.1",
		},
		{
			"420 gwei",
			big.NewInt(420000000000),
			"420",
		},
		{
			"0 gwei",
			big.NewInt(0),
			"0",
		},
		{
			"1 wei",
			big.NewInt(1),
			"0.000000001",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := unit.FormatGwei(tt.wei)
			if result != tt.expected {
				t.Errorf("FormatGwei(%v) = %s, want %s", tt.wei, result, tt.expected)
			}
		})
	}
}

func TestParseGwei(t *testing.T) {
	tests := []struct {
		name     string
		gwei     string
		expected string
		hasError bool
	}{
		{
			"1 gwei",
			"1",
			"1000000000",
			false,
		},
		{
			"1.5 gwei",
			"1.5",
			"1500000000",
			false,
		},
		{
			"0.1 gwei",
			"0.1",
			"100000000",
			false,
		},
		{
			"420 gwei",
			"420",
			"420000000000",
			false,
		},
		{
			"0 gwei",
			"0",
			"0",
			false,
		},
		{
			"invalid",
			"abc",
			"",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := unit.ParseGwei(tt.gwei)
			if tt.hasError && err == nil {
				t.Errorf("ParseGwei(%q) expected error", tt.gwei)
				return
			}
			if !tt.hasError && err != nil {
				t.Errorf("ParseGwei(%q) unexpected error: %v", tt.gwei, err)
				return
			}
			if !tt.hasError && result.String() != tt.expected {
				t.Errorf("ParseGwei(%q) = %s, want %s", tt.gwei, result.String(), tt.expected)
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	// Test that parse -> format -> parse gives same result
	tests := []struct {
		name     string
		value    string
		decimals int
	}{
		{"ether", "1.5", 18},
		{"gwei", "123.456", 9},
		{"small decimals", "69420.1234", 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := unit.ParseUnits(tt.value, tt.decimals)
			if err != nil {
				t.Errorf("ParseUnits(%q, %d) error: %v", tt.value, tt.decimals, err)
				return
			}

			formatted := unit.FormatUnits(parsed, tt.decimals)
			reparsed, err := unit.ParseUnits(formatted, tt.decimals)
			if err != nil {
				t.Errorf("ParseUnits(%q, %d) error: %v", formatted, tt.decimals, err)
				return
			}

			if parsed.Cmp(reparsed) != 0 {
				t.Errorf("Round trip failed: %s -> %s -> %s", tt.value, formatted, reparsed.String())
			}
		})
	}
}

// Benchmark tests

func BenchmarkFormatUnits(b *testing.B) {
	value := new(big.Int).SetUint64(1234567890123456789)
	for i := 0; i < b.N; i++ {
		unit.FormatUnits(value, 18)
	}
}

func BenchmarkParseUnits(b *testing.B) {
	for i := 0; i < b.N; i++ {
		unit.ParseUnits("1234567.890123456789", 18)
	}
}

func BenchmarkFormatEther(b *testing.B) {
	value := new(big.Int).SetUint64(1234567890123456789)
	for i := 0; i < b.N; i++ {
		unit.FormatEther(value)
	}
}

func BenchmarkParseEther(b *testing.B) {
	for i := 0; i < b.N; i++ {
		unit.ParseEther("1234567.890123456789")
	}
}
