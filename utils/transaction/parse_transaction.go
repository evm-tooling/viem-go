package transaction

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ChefBingbong/viem-go/utils/encoding"
)

// ParseTransaction parses a serialized transaction.
//
// Example:
//
//	tx, err := ParseTransaction("0x02f8...")
func ParseTransaction(serializedTx string) (*Transaction, error) {
	txType, err := GetSerializedTransactionType(serializedTx)
	if err != nil {
		return nil, err
	}

	switch txType {
	case TransactionTypeEIP7702:
		return parseTransactionEIP7702(serializedTx)
	case TransactionTypeEIP4844:
		return parseTransactionEIP4844(serializedTx)
	case TransactionTypeEIP1559:
		return parseTransactionEIP1559(serializedTx)
	case TransactionTypeEIP2930:
		return parseTransactionEIP2930(serializedTx)
	case TransactionTypeLegacy:
		return parseTransactionLegacy(serializedTx)
	default:
		return nil, fmt.Errorf("%w: unknown transaction type", ErrInvalidSerializedTransaction)
	}
}

func parseTransactionEIP7702(serializedTx string) (*Transaction, error) {
	// Remove type prefix (0x04) and decode RLP
	data, err := decodeTransactionRlp(serializedTx)
	if err != nil {
		return nil, err
	}

	items, ok := data.([]any)
	if !ok {
		return nil, ErrInvalidSerializedTransaction
	}

	// EIP-7702: [chainId, nonce, maxPriorityFeePerGas, maxFeePerGas, gas, to, value, data, accessList, authorizationList, v, r, s]
	if len(items) != 10 && len(items) != 13 {
		return nil, fmt.Errorf("%w: expected 10 or 13 fields, got %d", ErrInvalidSerializedTransaction, len(items))
	}

	tx := &Transaction{Type: TransactionTypeEIP7702}

	tx.ChainId = hexToNumber(getHexString(items[0]))
	tx.Nonce = hexToNumber(getHexString(items[1]))
	tx.MaxPriorityFeePerGas = hexToBigInt(getHexString(items[2]))
	tx.MaxFeePerGas = hexToBigInt(getHexString(items[3]))
	tx.Gas = hexToBigInt(getHexString(items[4]))
	tx.To = getNonEmptyHex(getHexString(items[5]))
	tx.Value = hexToBigInt(getHexString(items[6]))
	tx.Data = getNonEmptyHex(getHexString(items[7]))

	if accessList, ok := items[8].([]any); ok && len(accessList) > 0 {
		tx.AccessList, _ = ParseAccessList(accessList)
	}

	if authList, ok := items[9].([]any); ok && len(authList) > 0 {
		tx.AuthorizationList = parseAuthorizationList(authList)
	}

	// Parse signature if present
	if len(items) == 13 {
		parseEIP155Signature(tx, items[10:])
	}

	return tx, nil
}

func parseTransactionEIP4844(serializedTx string) (*Transaction, error) {
	// Remove type prefix (0x03) and decode RLP
	data, err := decodeTransactionRlp(serializedTx)
	if err != nil {
		return nil, err
	}

	items, ok := data.([]any)
	if !ok {
		return nil, ErrInvalidSerializedTransaction
	}

	// Check if it's a wrapper format (4 items: transaction, blobs, commitments, proofs)
	hasWrapper := len(items) == 4
	var txItems []any

	if hasWrapper {
		txItems, ok = items[0].([]any)
		if !ok {
			return nil, ErrInvalidSerializedTransaction
		}
	} else {
		txItems = items
	}

	// EIP-4844: [chainId, nonce, maxPriorityFeePerGas, maxFeePerGas, gas, to, value, data, accessList, maxFeePerBlobGas, blobVersionedHashes, v, r, s]
	if len(txItems) != 11 && len(txItems) != 14 {
		return nil, fmt.Errorf("%w: expected 11 or 14 fields, got %d", ErrInvalidSerializedTransaction, len(txItems))
	}

	tx := &Transaction{Type: TransactionTypeEIP4844}

	tx.ChainId = hexToNumber(getHexString(txItems[0]))
	tx.Nonce = hexToNumber(getHexString(txItems[1]))
	tx.MaxPriorityFeePerGas = hexToBigInt(getHexString(txItems[2]))
	tx.MaxFeePerGas = hexToBigInt(getHexString(txItems[3]))
	tx.Gas = hexToBigInt(getHexString(txItems[4]))
	tx.To = getNonEmptyHex(getHexString(txItems[5]))
	tx.Value = hexToBigInt(getHexString(txItems[6]))
	tx.Data = getNonEmptyHex(getHexString(txItems[7]))

	if accessList, ok := txItems[8].([]any); ok && len(accessList) > 0 {
		tx.AccessList, _ = ParseAccessList(accessList)
	}

	tx.MaxFeePerBlobGas = hexToBigInt(getHexString(txItems[9]))

	if hashes, ok := txItems[10].([]any); ok {
		tx.BlobVersionedHashes = make([]string, len(hashes))
		for i, h := range hashes {
			tx.BlobVersionedHashes[i] = getHexString(h)
		}
	}

	// Parse signature if present
	if len(txItems) == 14 {
		parseEIP155Signature(tx, txItems[11:])
	}

	// Parse wrapper (sidecars)
	if hasWrapper {
		blobs, _ := items[1].([]any)
		commitments, _ := items[2].([]any)
		proofs, _ := items[3].([]any)

		if len(blobs) > 0 {
			tx.Sidecars = make([]BlobSidecar, len(blobs))
			for i := range blobs {
				tx.Sidecars[i] = BlobSidecar{
					Blob:       getHexString(blobs[i]),
					Commitment: getHexString(commitments[i]),
					Proof:      getHexString(proofs[i]),
				}
			}
		}
	}

	return tx, nil
}

func parseTransactionEIP1559(serializedTx string) (*Transaction, error) {
	// Remove type prefix (0x02) and decode RLP
	data, err := decodeTransactionRlp(serializedTx)
	if err != nil {
		return nil, err
	}

	items, ok := data.([]any)
	if !ok {
		return nil, ErrInvalidSerializedTransaction
	}

	// EIP-1559: [chainId, nonce, maxPriorityFeePerGas, maxFeePerGas, gas, to, value, data, accessList, v, r, s]
	if len(items) != 9 && len(items) != 12 {
		return nil, fmt.Errorf("%w: expected 9 or 12 fields, got %d", ErrInvalidSerializedTransaction, len(items))
	}

	tx := &Transaction{Type: TransactionTypeEIP1559}

	tx.ChainId = hexToNumber(getHexString(items[0]))
	tx.Nonce = hexToNumber(getHexString(items[1]))
	tx.MaxPriorityFeePerGas = hexToBigInt(getHexString(items[2]))
	tx.MaxFeePerGas = hexToBigInt(getHexString(items[3]))
	tx.Gas = hexToBigInt(getHexString(items[4]))
	tx.To = getNonEmptyHex(getHexString(items[5]))
	tx.Value = hexToBigInt(getHexString(items[6]))
	tx.Data = getNonEmptyHex(getHexString(items[7]))

	if accessList, ok := items[8].([]any); ok && len(accessList) > 0 {
		tx.AccessList, _ = ParseAccessList(accessList)
	}

	// Parse signature if present
	if len(items) == 12 {
		parseEIP155Signature(tx, items[9:])
	}

	return tx, nil
}

func parseTransactionEIP2930(serializedTx string) (*Transaction, error) {
	// Remove type prefix (0x01) and decode RLP
	data, err := decodeTransactionRlp(serializedTx)
	if err != nil {
		return nil, err
	}

	items, ok := data.([]any)
	if !ok {
		return nil, ErrInvalidSerializedTransaction
	}

	// EIP-2930: [chainId, nonce, gasPrice, gas, to, value, data, accessList, v, r, s]
	if len(items) != 8 && len(items) != 11 {
		return nil, fmt.Errorf("%w: expected 8 or 11 fields, got %d", ErrInvalidSerializedTransaction, len(items))
	}

	tx := &Transaction{Type: TransactionTypeEIP2930}

	tx.ChainId = hexToNumber(getHexString(items[0]))
	tx.Nonce = hexToNumber(getHexString(items[1]))
	tx.GasPrice = hexToBigInt(getHexString(items[2]))
	tx.Gas = hexToBigInt(getHexString(items[3]))
	tx.To = getNonEmptyHex(getHexString(items[4]))
	tx.Value = hexToBigInt(getHexString(items[5]))
	tx.Data = getNonEmptyHex(getHexString(items[6]))

	if accessList, ok := items[7].([]any); ok && len(accessList) > 0 {
		tx.AccessList, _ = ParseAccessList(accessList)
	}

	// Parse signature if present
	if len(items) == 11 {
		parseEIP155Signature(tx, items[8:])
	}

	return tx, nil
}

func parseTransactionLegacy(serializedTx string) (*Transaction, error) {
	// Legacy transactions are RLP encoded directly (no type prefix)
	data, err := encoding.RlpDecodeHex(serializedTx)
	if err != nil {
		return nil, err
	}

	items, ok := data.([]any)
	if !ok {
		return nil, ErrInvalidSerializedTransaction
	}

	// Legacy: [nonce, gasPrice, gas, to, value, data, v, r, s] or [nonce, gasPrice, gas, to, value, data]
	if len(items) != 6 && len(items) != 9 {
		return nil, fmt.Errorf("%w: expected 6 or 9 fields, got %d", ErrInvalidSerializedTransaction, len(items))
	}

	tx := &Transaction{Type: TransactionTypeLegacy}

	tx.Nonce = hexToNumber(getHexString(items[0]))
	tx.GasPrice = hexToBigInt(getHexString(items[1]))
	tx.Gas = hexToBigInt(getHexString(items[2]))
	tx.To = getNonEmptyHex(getHexString(items[3]))
	tx.Value = hexToBigInt(getHexString(items[4]))
	tx.Data = getNonEmptyHex(getHexString(items[5]))

	// Parse signature if present
	if len(items) == 9 {
		parseLegacySignature(tx, items[6:])
	}

	return tx, nil
}

func decodeTransactionRlp(serializedTx string) (any, error) {
	// Remove type prefix (first 2 bytes after 0x)
	rlpData := "0x" + serializedTx[4:]
	return encoding.RlpDecodeHex(rlpData)
}

func parseEIP155Signature(tx *Transaction, sigItems []any) {
	if len(sigItems) < 3 {
		return
	}

	vHex := getHexString(sigItems[0])
	v := hexToBigInt(vHex)

	// Determine actual v value (0 or 1 becomes 27 or 28)
	if v == nil || v.Sign() == 0 {
		tx.V = big.NewInt(27)
		tx.YParity = 0
	} else {
		tx.V = big.NewInt(28)
		tx.YParity = 1
	}

	tx.R = padHex(getHexString(sigItems[1]), 32)
	tx.S = padHex(getHexString(sigItems[2]), 32)
}

func parseLegacySignature(tx *Transaction, sigItems []any) {
	if len(sigItems) < 3 {
		return
	}

	vHex := getHexString(sigItems[0])
	rHex := getHexString(sigItems[1])
	sHex := getHexString(sigItems[2])

	v := hexToBigInt(vHex)
	if v == nil {
		v = big.NewInt(0)
	}

	// Check if it's an unsigned transaction with chainId
	if rHex == "0x" && sHex == "0x" {
		if v.Sign() > 0 {
			tx.ChainId = int(v.Int64())
		}
		return
	}

	tx.V = v
	tx.R = rHex
	tx.S = sHex

	// Derive chainId from v (EIP-155)
	vInt := v.Int64()
	if vInt > 28 {
		chainId := (vInt - 35) / 2
		if chainId > 0 {
			tx.ChainId = int(chainId)
		}
		tx.YParity = int(vInt % 2)
	} else if vInt == 27 {
		tx.YParity = 0
	} else if vInt == 28 {
		tx.YParity = 1
	}
}

func parseAuthorizationList(items []any) []SignedAuthorization {
	result := make([]SignedAuthorization, 0, len(items))

	for _, item := range items {
		authItems, ok := item.([]any)
		if !ok || len(authItems) < 6 {
			continue
		}

		vHex := getHexString(authItems[3])
		yParity := 0
		if hexToBigInt(vHex) != nil && hexToBigInt(vHex).Sign() > 0 {
			yParity = 1
		}

		result = append(result, SignedAuthorization{
			Authorization: Authorization{
				ChainId: hexToNumber(getHexString(authItems[0])),
				Address: getHexString(authItems[1]),
				Nonce:   hexToNumber(getHexString(authItems[2])),
			},
			YParity: yParity,
			R:       getHexString(authItems[4]),
			S:       getHexString(authItems[5]),
		})
	}

	return result
}

// Helper functions

func getHexString(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case []byte:
		return BytesToHex(val)
	default:
		return "0x"
	}
}

func getNonEmptyHex(h string) string {
	if h == "0x" || h == "" {
		return ""
	}
	return h
}

func hexToNumber(h string) int {
	if h == "" || h == "0x" {
		return 0
	}
	h = strings.TrimPrefix(h, "0x")
	h = strings.TrimPrefix(h, "0X")

	n := new(big.Int)
	n.SetString(h, 16)
	return int(n.Int64())
}

func hexToBigInt(h string) *big.Int {
	if h == "" || h == "0x" {
		return nil
	}
	h = strings.TrimPrefix(h, "0x")
	h = strings.TrimPrefix(h, "0X")

	n := new(big.Int)
	n.SetString(h, 16)
	return n
}

func padHex(h string, size int) string {
	h = strings.TrimPrefix(h, "0x")
	h = strings.TrimPrefix(h, "0X")

	targetLen := size * 2
	if len(h) >= targetLen {
		return "0x" + h
	}

	return "0x" + strings.Repeat("0", targetLen-len(h)) + h
}
