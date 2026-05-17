package transaction

import (
	"fmt"
	"strings"
)

// AssertTransactionEIP7702 validates an EIP-7702 transaction.
func AssertTransactionEIP7702(tx *Transaction) error {
	// Validate authorization list
	for _, auth := range tx.AuthorizationList {
		if !isValidAddress(auth.Address) {
			return fmt.Errorf("%w: %s", ErrInvalidAddress, auth.Address)
		}
		if auth.ChainId < 0 {
			return fmt.Errorf("%w: %d", ErrInvalidChainId, auth.ChainId)
		}
	}

	// Also validate as EIP-1559
	return AssertTransactionEIP1559(tx)
}

// AssertTransactionEIP4844 validates an EIP-4844 (blob) transaction.
func AssertTransactionEIP4844(tx *Transaction) error {
	// Validate blob versioned hashes - EIP-4844 transactions must have at least one blob
	if len(tx.BlobVersionedHashes) == 0 {
		return ErrEmptyBlob
	}

	for _, hash := range tx.BlobVersionedHashes {
		// Check size (should be 32 bytes = 66 chars with 0x prefix)
		hashClean := strings.TrimPrefix(hash, "0x")
		if len(hashClean) != 64 {
			return fmt.Errorf("%w: %s (size: %d)", ErrInvalidVersionedHashSize, hash, len(hashClean)/2)
		}

		// Check version (first byte should be 0x01 for KZG)
		if len(hashClean) >= 2 {
			version := hexToInt(hashClean[0:2])
			if version != VersionedHashVersionKzg {
				return fmt.Errorf("%w: %s (version: %d)", ErrInvalidVersionedHashVersion, hash, version)
			}
		}
	}

	// Also validate as EIP-1559
	return AssertTransactionEIP1559(tx)
}

// AssertTransactionEIP1559 validates an EIP-1559 transaction.
func AssertTransactionEIP1559(tx *Transaction) error {
	// Validate chain ID
	if tx.ChainId <= 0 {
		return fmt.Errorf("%w: %d", ErrInvalidChainId, tx.ChainId)
	}

	// Validate to address if present
	if tx.To != "" && !isValidAddress(tx.To) {
		return fmt.Errorf("%w: %s", ErrInvalidAddress, tx.To)
	}

	// Check maxFeePerGas doesn't exceed max uint256
	if tx.MaxFeePerGas != nil && tx.MaxFeePerGas.Cmp(MaxUint256) > 0 {
		return fmt.Errorf("%w: maxFeePerGas exceeds maximum value", ErrFeeCapTooHigh)
	}

	// Check tip doesn't exceed fee cap
	if tx.MaxPriorityFeePerGas != nil && tx.MaxFeePerGas != nil {
		if tx.MaxPriorityFeePerGas.Cmp(tx.MaxFeePerGas) > 0 {
			return fmt.Errorf("%w: maxPriorityFeePerGas (%s) > maxFeePerGas (%s)",
				ErrTipAboveFeeCap, tx.MaxPriorityFeePerGas.String(), tx.MaxFeePerGas.String())
		}
	}

	return nil
}

// AssertTransactionEIP2930 validates an EIP-2930 transaction.
func AssertTransactionEIP2930(tx *Transaction) error {
	// Validate chain ID
	if tx.ChainId <= 0 {
		return fmt.Errorf("%w: %d", ErrInvalidChainId, tx.ChainId)
	}

	// Validate to address if present
	if tx.To != "" && !isValidAddress(tx.To) {
		return fmt.Errorf("%w: %s", ErrInvalidAddress, tx.To)
	}

	// EIP-2930 shouldn't have EIP-1559 fee fields
	if tx.MaxPriorityFeePerGas != nil || tx.MaxFeePerGas != nil {
		return ErrMaxFeePerGasNotAllowed
	}

	// Check gasPrice doesn't exceed max uint256
	if tx.GasPrice != nil && tx.GasPrice.Cmp(MaxUint256) > 0 {
		return fmt.Errorf("%w: gasPrice exceeds maximum value", ErrFeeCapTooHigh)
	}

	return nil
}

// AssertTransactionLegacy validates a legacy transaction.
func AssertTransactionLegacy(tx *Transaction) error {
	// Validate to address if present
	if tx.To != "" && !isValidAddress(tx.To) {
		return fmt.Errorf("%w: %s", ErrInvalidAddress, tx.To)
	}

	// Validate chain ID if present
	if tx.ChainId != 0 && tx.ChainId <= 0 {
		return fmt.Errorf("%w: %d", ErrInvalidChainId, tx.ChainId)
	}

	// Legacy shouldn't have EIP-1559 fee fields
	if tx.MaxPriorityFeePerGas != nil || tx.MaxFeePerGas != nil {
		return ErrMaxFeePerGasNotAllowed
	}

	// Check gasPrice doesn't exceed max uint256
	if tx.GasPrice != nil && tx.GasPrice.Cmp(MaxUint256) > 0 {
		return fmt.Errorf("%w: gasPrice exceeds maximum value", ErrFeeCapTooHigh)
	}

	return nil
}

// AssertTransaction validates a transaction based on its type.
func AssertTransaction(tx *Transaction) error {
	txType, err := GetTransactionType(tx)
	if err != nil {
		// If we can't determine type, try basic validation
		if tx.To != "" && !isValidAddress(tx.To) {
			return fmt.Errorf("%w: %s", ErrInvalidAddress, tx.To)
		}
		return nil
	}

	switch txType {
	case TransactionTypeEIP7702:
		return AssertTransactionEIP7702(tx)
	case TransactionTypeEIP4844:
		return AssertTransactionEIP4844(tx)
	case TransactionTypeEIP1559:
		return AssertTransactionEIP1559(tx)
	case TransactionTypeEIP2930:
		return AssertTransactionEIP2930(tx)
	case TransactionTypeLegacy:
		return AssertTransactionLegacy(tx)
	default:
		return nil
	}
}
