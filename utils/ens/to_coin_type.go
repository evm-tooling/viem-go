package ens

import (
	"errors"
	"fmt"
)

const slip44MSB = 0x80000000

// ErrInvalidChainId is returned when a chain ID is invalid for coin type conversion.
var ErrInvalidChainId = errors.New("invalid chain ID")

// ToCoinType converts a chain ID to an ENSIP-9 compliant coin type.
//
// For Ethereum mainnet (chainId 1), returns 60 (ETH's SLIP-44 coin type).
// For other chains, applies the ENSIP-9 formula: 0x80000000 | chainId
//
// Example:
//
//	coinType, _ := ToCoinType(1)
//	// 60 (Ethereum)
//
//	coinType, _ := ToCoinType(10)
//	// 2147483658 (Optimism)
//
//	coinType, _ := ToCoinType(137)
//	// 2147483785 (Polygon)
//
// @see https://docs.ens.domains/ensip/9
func ToCoinType(chainId int) (uint64, error) {
	// Ethereum mainnet uses SLIP-44 coin type 60
	if chainId == 1 {
		return 60, nil
	}

	// Validate chain ID
	if chainId >= slip44MSB || chainId < 0 {
		return 0, fmt.Errorf("%w: %d", ErrInvalidChainId, chainId)
	}

	// Apply ENSIP-9 formula
	return uint64((slip44MSB | chainId) & 0xFFFFFFFF), nil
}

// MustToCoinType converts a chain ID to a coin type, panicking on error.
func MustToCoinType(chainId int) uint64 {
	coinType, err := ToCoinType(chainId)
	if err != nil {
		panic(err)
	}
	return coinType
}
