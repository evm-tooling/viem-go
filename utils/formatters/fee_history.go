package formatters

import (
	"math/big"
)

// FormatFeeHistory formats an RPC fee history into a FeeHistory struct.
//
// Example:
//
//	rpcFeeHistory := RpcFeeHistory{
//		BaseFeePerGas: []string{"0x1", "0x2"},
//		GasUsedRatio:  []float64{0.5, 0.6},
//		OldestBlock:   "0x100",
//	}
//	feeHistory := FormatFeeHistory(rpcFeeHistory)
func FormatFeeHistory(feeHistory RpcFeeHistory) FeeHistory {
	result := FeeHistory{
		GasUsedRatio: feeHistory.GasUsedRatio,
	}

	// Base fee per gas
	if len(feeHistory.BaseFeePerGas) > 0 {
		result.BaseFeePerGas = make([]*big.Int, len(feeHistory.BaseFeePerGas))
		for i, fee := range feeHistory.BaseFeePerGas {
			result.BaseFeePerGas[i] = hexToBigInt(fee)
		}
	}

	// Oldest block
	if feeHistory.OldestBlock != "" {
		result.OldestBlock = hexToBigInt(feeHistory.OldestBlock)
	}

	// Reward
	if len(feeHistory.Reward) > 0 {
		result.Reward = make([][]*big.Int, len(feeHistory.Reward))
		for i, rewards := range feeHistory.Reward {
			result.Reward[i] = make([]*big.Int, len(rewards))
			for j, reward := range rewards {
				result.Reward[i][j] = hexToBigInt(reward)
			}
		}
	}

	return result
}
