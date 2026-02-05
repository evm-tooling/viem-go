// Package decorators provides decorator functions for extending clients with actions.
// This mirrors viem's decorators pattern.
package decorators

import (
	"github.com/ChefBingbong/viem-go/client"
)

// PublicActions returns public action methods as a map.
// This mirrors viem's publicActions decorator for extension purposes.
//
// Example:
//
//	client := client.CreatePublicClient(config)
//	actions := decorators.PublicActions(client)
func PublicActions(c *client.PublicClient) map[string]any {
	return map[string]any{
		// Read actions
		"getBlockNumber":            c.GetBlockNumber,
		"getChainId":                c.GetChainID,
		"getGasPrice":               c.GetGasPrice,
		"getBalance":                c.GetBalance,
		"getTransactionCount":       c.GetTransactionCount,
		"getCode":                   c.GetCode,
		"getStorageAt":              c.GetStorageAt,
		"call":                      c.Call,
		"estimateGas":               c.EstimateGas,
		"getBlock":                  c.GetBlock,
		"getBlockByNumber":          c.GetBlockByNumber,
		"getBlockByHash":            c.GetBlockByHash,
		"getTransaction":            c.GetTransaction,
		"getTransactionReceipt":     c.GetTransactionReceipt,
		"getLogs":                   c.GetLogs,
		"getFeeHistory":             c.GetFeeHistory,
		"getMaxPriorityFeePerGas":   c.GetMaxPriorityFeePerGas,
		"getProof":                  c.GetProof,
		"waitForTransactionReceipt": c.WaitForTransactionReceipt,
		"readContract":              c.ReadContract,
		"simulateContract":          c.SimulateContract,
		"prepareContractWrite":      c.PrepareContractWrite,

		// Watch actions
		"watchBlockNumber":         c.WatchBlockNumber,
		"watchBlocks":              c.WatchBlocks,
		"watchPendingTransactions": c.WatchPendingTransactions,
		"watchEvent":               c.WatchEvent,
		"watchContractEvent":       c.WatchContractEvent,
	}
}
