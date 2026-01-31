package formatters

import (
	"math/big"
)

// LogFormatOptions contains optional formatting options for logs.
type LogFormatOptions struct {
	Args      any
	EventName string
}

// FormatLog formats an RPC log into a Log struct.
//
// Example:
//
//	rpcLog := RpcLog{
//		Address:     "0x...",
//		BlockNumber: "0x1234",
//	}
//	log := FormatLog(rpcLog, nil)
func FormatLog(log RpcLog, opts *LogFormatOptions) Log {
	result := Log{
		Address: log.Address,
		Data:    log.Data,
		Removed: log.Removed,
		Topics:  log.Topics,
	}

	// Block hash
	if log.BlockHash != "" {
		result.BlockHash = &log.BlockHash
	}

	// Block number
	if log.BlockNumber != "" {
		result.BlockNumber = hexToBigInt(log.BlockNumber)
	}

	// Block timestamp
	if log.BlockTimestamp != nil {
		if *log.BlockTimestamp != "" {
			result.BlockTimestamp = hexToBigInt(*log.BlockTimestamp)
		}
	}

	// Log index
	if log.LogIndex != "" {
		idx := hexToInt(log.LogIndex)
		result.LogIndex = &idx
	}

	// Transaction hash
	if log.TransactionHash != "" {
		result.TransactionHash = &log.TransactionHash
	}

	// Transaction index
	if log.TransactionIndex != "" {
		idx := hexToInt(log.TransactionIndex)
		result.TransactionIndex = &idx
	}

	// Optional args and event name
	if opts != nil {
		if opts.EventName != "" {
			result.Args = opts.Args
			result.EventName = opts.EventName
		}
	}

	return result
}

// FormatLogs formats multiple RPC logs.
func FormatLogs(logs []RpcLog) []Log {
	result := make([]Log, len(logs))
	for i, log := range logs {
		result[i] = FormatLog(log, nil)
	}
	return result
}

// hexToBigInt converts a hex string to *big.Int.
func hexToBigInt(hex string) *big.Int {
	if hex == "" {
		return nil
	}
	// Remove 0x prefix
	if len(hex) >= 2 && hex[:2] == "0x" {
		hex = hex[2:]
	}
	n := new(big.Int)
	n.SetString(hex, 16)
	return n
}

// hexToInt converts a hex string to int.
func hexToInt(hex string) int {
	if hex == "" {
		return 0
	}
	// Remove 0x prefix
	if len(hex) >= 2 && hex[:2] == "0x" {
		hex = hex[2:]
	}
	n := new(big.Int)
	n.SetString(hex, 16)
	return int(n.Int64())
}
