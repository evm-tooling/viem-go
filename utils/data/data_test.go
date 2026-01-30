package data_test

import (
	"testing"

	"github.com/ChefBingbong/viem-go/utils/data"
)

func TestConcatBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    [][]byte
		expected []byte
	}{
		{
			"two arrays",
			[][]byte{{0x01, 0x02}, {0x03, 0x04}},
			[]byte{0x01, 0x02, 0x03, 0x04},
		},
		{
			"three arrays",
			[][]byte{{0x01}, {0x02}, {0x03}},
			[]byte{0x01, 0x02, 0x03},
		},
		{
			"empty arrays",
			[][]byte{{}, {0x01}},
			[]byte{0x01},
		},
		{
			"single array",
			[][]byte{{0x01, 0x02, 0x03}},
			[]byte{0x01, 0x02, 0x03},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := data.ConcatBytes(tt.input...)
			if len(result) != len(tt.expected) {
				t.Errorf("ConcatBytes length = %d, want %d", len(result), len(tt.expected))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("ConcatBytes[%d] = %d, want %d", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestConcatHex(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{
			"two hex strings",
			[]string{"0x0102", "0x0304"},
			"0x01020304",
		},
		{
			"three hex strings",
			[]string{"0x01", "0x02", "0x03"},
			"0x010203",
		},
		{
			"empty and non-empty",
			[]string{"0x", "0x0102"},
			"0x0102",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := data.ConcatHex(tt.input...)
			if result != tt.expected {
				t.Errorf("ConcatHex = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestIsBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected bool
	}{
		{"byte slice", []byte{0x01, 0x02}, true},
		{"empty byte slice", []byte{}, true},
		{"string", "hello", false},
		{"nil", nil, false},
		{"int", 123, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := data.IsBytes(tt.input)
			if result != tt.expected {
				t.Errorf("IsBytes(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsHex(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		strict   bool
		expected bool
	}{
		{"valid hex", "0x0102", true, true},
		{"valid empty hex", "0x", true, true},
		{"valid hex uppercase", "0xABCD", true, true},
		{"invalid chars strict", "0xgg", true, false},
		{"invalid chars non-strict", "0xgg", false, true},
		{"no prefix", "0102", true, false},
		{"empty", "", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := data.IsHex(tt.input, data.IsHexOptions{Strict: tt.strict})
			if result != tt.expected {
				t.Errorf("IsHex(%q, strict=%v) = %v, want %v", tt.input, tt.strict, result, tt.expected)
			}
		})
	}
}

func TestPadBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		dir      data.PadDirection
		size     int
		expected []byte
		hasError bool
	}{
		{
			"left pad",
			[]byte{0x01, 0x02},
			data.PadLeft,
			4,
			[]byte{0x00, 0x00, 0x01, 0x02},
			false,
		},
		{
			"right pad",
			[]byte{0x01, 0x02},
			data.PadRight,
			4,
			[]byte{0x01, 0x02, 0x00, 0x00},
			false,
		},
		{
			"exact size",
			[]byte{0x01, 0x02},
			data.PadLeft,
			2,
			[]byte{0x01, 0x02},
			false,
		},
		{
			"exceeds size",
			[]byte{0x01, 0x02, 0x03},
			data.PadLeft,
			2,
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := data.PadBytes(tt.input, tt.dir, tt.size)
			if tt.hasError && err == nil {
				t.Error("expected error")
			}
			if !tt.hasError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.hasError {
				if len(result) != len(tt.expected) {
					t.Errorf("PadBytes length = %d, want %d", len(result), len(tt.expected))
					return
				}
				for i := range result {
					if result[i] != tt.expected[i] {
						t.Errorf("PadBytes[%d] = %d, want %d", i, result[i], tt.expected[i])
					}
				}
			}
		})
	}
}

func TestPadHex(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		dir      data.PadDirection
		size     int
		expected string
		hasError bool
	}{
		{
			"left pad",
			"0x0102",
			data.PadLeft,
			4,
			"0x00000102",
			false,
		},
		{
			"right pad",
			"0x0102",
			data.PadRight,
			4,
			"0x01020000",
			false,
		},
		{
			"exceeds size",
			"0x01020304",
			data.PadLeft,
			2,
			"",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := data.PadHex(tt.input, tt.dir, tt.size)
			if tt.hasError && err == nil {
				t.Error("expected error")
			}
			if !tt.hasError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.hasError && result != tt.expected {
				t.Errorf("PadHex = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestSize(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected int
	}{
		{"hex string", "0x0102", 2},
		{"empty hex", "0x", 0},
		{"byte slice", []byte{0x01, 0x02, 0x03}, 3},
		{"empty bytes", []byte{}, 0},
		{"odd hex", "0x1", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := data.Size(tt.input)
			if result != tt.expected {
				t.Errorf("Size(%v) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSliceBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		start    int
		end      int
		expected []byte
		hasError bool
	}{
		{
			"middle slice",
			[]byte{0x01, 0x02, 0x03, 0x04},
			1,
			3,
			[]byte{0x02, 0x03},
			false,
		},
		{
			"from start",
			[]byte{0x01, 0x02, 0x03},
			0,
			2,
			[]byte{0x01, 0x02},
			false,
		},
		{
			"to end",
			[]byte{0x01, 0x02, 0x03},
			1,
			3,
			[]byte{0x02, 0x03},
			false,
		},
		{
			"start out of bounds",
			[]byte{0x01, 0x02},
			5,
			6,
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := data.SliceBytes(tt.input, tt.start, tt.end)
			if tt.hasError && err == nil {
				t.Error("expected error")
			}
			if !tt.hasError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.hasError {
				if len(result) != len(tt.expected) {
					t.Errorf("SliceBytes length = %d, want %d", len(result), len(tt.expected))
					return
				}
				for i := range result {
					if result[i] != tt.expected[i] {
						t.Errorf("SliceBytes[%d] = %d, want %d", i, result[i], tt.expected[i])
					}
				}
			}
		})
	}
}

func TestSliceHex(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		start    int
		end      int
		expected string
		hasError bool
	}{
		{
			"middle slice",
			"0x01020304",
			1,
			3,
			"0x0203",
			false,
		},
		{
			"from start",
			"0x010203",
			0,
			2,
			"0x0102",
			false,
		},
		{
			"to end",
			"0x010203",
			1,
			3,
			"0x0203",
			false,
		},
		{
			"start out of bounds",
			"0x0102",
			5,
			6,
			"",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := data.SliceHex(tt.input, tt.start, tt.end)
			if tt.hasError && err == nil {
				t.Error("expected error")
			}
			if !tt.hasError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.hasError && result != tt.expected {
				t.Errorf("SliceHex = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestTrimBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		dir      data.TrimDirection
		expected []byte
	}{
		{
			"trim left",
			[]byte{0x00, 0x00, 0x01, 0x02},
			data.TrimLeft,
			[]byte{0x01, 0x02},
		},
		{
			"trim right",
			[]byte{0x01, 0x02, 0x00, 0x00},
			data.TrimRight,
			[]byte{0x01, 0x02},
		},
		{
			"no zeros left",
			[]byte{0x01, 0x02},
			data.TrimLeft,
			[]byte{0x01, 0x02},
		},
		{
			"all zeros except last",
			[]byte{0x00, 0x00, 0x00, 0x01},
			data.TrimLeft,
			[]byte{0x01},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := data.TrimBytes(tt.input, tt.dir)
			if len(result) != len(tt.expected) {
				t.Errorf("TrimBytes length = %d, want %d", len(result), len(tt.expected))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("TrimBytes[%d] = %d, want %d", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestTrimHex(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		dir      data.TrimDirection
		expected string
	}{
		{
			"trim left",
			"0x00000102",
			data.TrimLeft,
			"0x0102",
		},
		{
			"trim right",
			"0x01020000",
			data.TrimRight,
			"0x0102",
		},
		{
			"no zeros",
			"0x0102",
			data.TrimLeft,
			"0x0102",
		},
		{
			"all zeros except one",
			"0x00000001",
			data.TrimLeft,
			"0x01",
		},
		{
			"empty",
			"0x",
			data.TrimLeft,
			"0x",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := data.TrimHex(tt.input, tt.dir)
			if result != tt.expected {
				t.Errorf("TrimHex(%q, %s) = %q, want %q", tt.input, tt.dir, result, tt.expected)
			}
		})
	}
}

func TestPadDefault(t *testing.T) {
	// Test default padding (left, 32 bytes)
	result, err := data.Pad([]byte{0x01})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(result) != 32 {
		t.Errorf("Pad default length = %d, want 32", len(result))
	}
	if result[31] != 0x01 {
		t.Errorf("Pad default last byte = %d, want 1", result[31])
	}
	if result[0] != 0x00 {
		t.Errorf("Pad default first byte = %d, want 0", result[0])
	}
}

// Benchmark tests

func BenchmarkConcatBytes(b *testing.B) {
	a := []byte{0x01, 0x02, 0x03, 0x04}
	c := []byte{0x05, 0x06, 0x07, 0x08}
	for i := 0; i < b.N; i++ {
		data.ConcatBytes(a, c)
	}
}

func BenchmarkPadBytes(b *testing.B) {
	bytes := []byte{0x01, 0x02}
	for i := 0; i < b.N; i++ {
		data.PadBytes(bytes, data.PadLeft, 32)
	}
}

func BenchmarkTrimBytes(b *testing.B) {
	bytes := []byte{0x00, 0x00, 0x00, 0x01, 0x02}
	for i := 0; i < b.N; i++ {
		data.TrimBytes(bytes, data.TrimLeft)
	}
}
