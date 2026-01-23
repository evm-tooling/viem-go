package utils

import (
	"math/big"
	"strconv"
)

func IntToHex(n int64) string {
	if n < 0 {
		return "-0x" + strconv.FormatInt(-n, 16)
	}
	return "0x" + strconv.FormatInt(n, 16)
}

func UintToHex(n uint64) string {
	return "0x" + strconv.FormatUint(n, 16)
}

func BigIntToHex(n *big.Int) string {
	if n == nil {
		return "0x0"
	}
	if n.Sign() < 0 {
		return "-0x" + n.Abs(n).Text(16)
	}
	return "0x" + n.Text(16)
}

func BoolToHex(v bool) string {
	if v {
		return "0x1"
	}
	return "0x0"
}
