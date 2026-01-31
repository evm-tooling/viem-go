package test

import (
	"math/big"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ChefBingbong/viem-go/utils/formatters"
)

var _ = Describe("Formatters", func() {
	Describe("FormatLog", func() {
		It("should format a basic log", func() {
			rpcLog := formatters.RpcLog{
				Address:          "0x1234567890123456789012345678901234567890",
				BlockHash:        "0xabcd",
				BlockNumber:      "0x100",
				Data:             "0x",
				LogIndex:         "0x0",
				TransactionHash:  "0xdef0",
				TransactionIndex: "0x1",
				Topics:           []string{"0xtopic1"},
			}

			log := formatters.FormatLog(rpcLog, nil)

			Expect(log.Address).To(Equal("0x1234567890123456789012345678901234567890"))
			Expect(*log.BlockHash).To(Equal("0xabcd"))
			Expect(log.BlockNumber.Cmp(big.NewInt(256))).To(Equal(0))
			Expect(*log.LogIndex).To(Equal(0))
			Expect(*log.TransactionHash).To(Equal("0xdef0"))
			Expect(*log.TransactionIndex).To(Equal(1))
		})

		It("should handle event name and args", func() {
			rpcLog := formatters.RpcLog{
				Address:     "0x1234567890123456789012345678901234567890",
				BlockNumber: "0x1",
			}

			opts := &formatters.LogFormatOptions{
				EventName: "Transfer",
				Args:      map[string]any{"from": "0x1", "to": "0x2"},
			}

			log := formatters.FormatLog(rpcLog, opts)

			Expect(log.EventName).To(Equal("Transfer"))
			Expect(log.Args).NotTo(BeNil())
		})
	})

	Describe("FormatTransaction", func() {
		It("should format a basic transaction", func() {
			rpcTx := formatters.RpcTransaction{
				BlockHash:   "0xabcd",
				BlockNumber: "0x100",
				ChainID:     "0x1",
				From:        "0xfrom",
				Gas:         "0x5208",
				GasPrice:    "0x3b9aca00",
				Hash:        "0xhash",
				Nonce:       "0xa",
				To:          "0xto",
				Value:       "0xde0b6b3a7640000",
				Type:        "0x0",
			}

			tx := formatters.FormatTransaction(rpcTx)

			Expect(*tx.BlockHash).To(Equal("0xabcd"))
			Expect(tx.BlockNumber.Cmp(big.NewInt(256))).To(Equal(0))
			Expect(*tx.ChainID).To(Equal(1))
			Expect(tx.From).To(Equal("0xfrom"))
			Expect(tx.Gas.Cmp(big.NewInt(21000))).To(Equal(0))
			Expect(tx.GasPrice.Cmp(big.NewInt(1000000000))).To(Equal(0))
			Expect(*tx.Nonce).To(Equal(10))
			Expect(*tx.To).To(Equal("0xto"))
			Expect(tx.Type).To(Equal(formatters.TransactionTypeLegacy))
		})

		It("should handle EIP-1559 transaction", func() {
			rpcTx := formatters.RpcTransaction{
				Type:                 "0x2",
				MaxFeePerGas:         "0x3b9aca00",
				MaxPriorityFeePerGas: "0x77359400",
			}

			tx := formatters.FormatTransaction(rpcTx)

			Expect(tx.Type).To(Equal(formatters.TransactionTypeEIP1559))
			Expect(tx.MaxFeePerGas.Cmp(big.NewInt(1000000000))).To(Equal(0))
			Expect(tx.MaxPriorityFeePerGas.Cmp(big.NewInt(2000000000))).To(Equal(0))
		})

		It("should derive yParity from v", func() {
			// v = 27 -> yParity = 0
			rpcTx := formatters.RpcTransaction{
				V: "0x1b", // 27
			}
			tx := formatters.FormatTransaction(rpcTx)
			Expect(*tx.YParity).To(Equal(0))

			// v = 28 -> yParity = 1
			rpcTx = formatters.RpcTransaction{
				V: "0x1c", // 28
			}
			tx = formatters.FormatTransaction(rpcTx)
			Expect(*tx.YParity).To(Equal(1))
		})
	})

	Describe("FormatBlock", func() {
		It("should format a basic block", func() {
			rpcBlock := formatters.RpcBlock{
				Number:          "0x100",
				Hash:            "0xblockhash",
				Timestamp:       "0x5f5e100",
				GasLimit:        "0x1c9c380",
				GasUsed:         "0x5208",
				BaseFeePerGas:   "0x3b9aca00",
				Nonce:           "0x0",
				Miner:           "0xminer",
			}

			block := formatters.FormatBlock(rpcBlock)

			Expect(block.Number.Cmp(big.NewInt(256))).To(Equal(0))
			Expect(*block.Hash).To(Equal("0xblockhash"))
			Expect(block.Timestamp.Cmp(big.NewInt(100000000))).To(Equal(0))
			Expect(block.GasLimit.Cmp(big.NewInt(30000000))).To(Equal(0))
			Expect(block.GasUsed.Cmp(big.NewInt(21000))).To(Equal(0))
			Expect(block.BaseFeePerGas.Cmp(big.NewInt(1000000000))).To(Equal(0))
			Expect(*block.Nonce).To(Equal("0x0"))
			Expect(block.Miner).To(Equal("0xminer"))
		})

		It("should handle transactions as hashes", func() {
			rpcBlock := formatters.RpcBlock{
				Number:       "0x1",
				Transactions: []any{"0xtx1", "0xtx2"},
			}

			block := formatters.FormatBlock(rpcBlock)

			Expect(len(block.Transactions)).To(Equal(2))
			Expect(block.Transactions[0]).To(Equal("0xtx1"))
		})
	})

	Describe("FormatTransactionReceipt", func() {
		It("should format a basic receipt", func() {
			rpcReceipt := formatters.RpcTransactionReceipt{
				BlockNumber:       "0x100",
				GasUsed:           "0x5208",
				Status:            "0x1",
				TransactionIndex:  "0x0",
				ContractAddress:   "0xcontract",
				CumulativeGasUsed: "0x5208",
				EffectiveGasPrice: "0x3b9aca00",
				Type:              "0x2",
			}

			receipt := formatters.FormatTransactionReceipt(rpcReceipt)

			Expect(receipt.BlockNumber.Cmp(big.NewInt(256))).To(Equal(0))
			Expect(receipt.GasUsed.Cmp(big.NewInt(21000))).To(Equal(0))
			Expect(receipt.Status).To(Equal(formatters.ReceiptStatusSuccess))
			Expect(*receipt.TransactionIndex).To(Equal(0))
			Expect(*receipt.ContractAddress).To(Equal("0xcontract"))
			Expect(receipt.Type).To(Equal(formatters.TransactionTypeEIP1559))
		})

		It("should format receipt with reverted status", func() {
			rpcReceipt := formatters.RpcTransactionReceipt{
				Status: "0x0",
			}

			receipt := formatters.FormatTransactionReceipt(rpcReceipt)

			Expect(receipt.Status).To(Equal(formatters.ReceiptStatusReverted))
		})
	})

	Describe("FormatTransactionRequest", func() {
		It("should format a basic request", func() {
			request := formatters.TransactionRequest{
				From:  "0xfrom",
				To:    "0xto",
				Value: big.NewInt(1000000000000000000),
				Gas:   big.NewInt(21000),
				Nonce: intPtr(5),
				Type:  formatters.TransactionTypeEIP1559,
			}

			rpcRequest := formatters.FormatTransactionRequest(request)

			Expect(rpcRequest.From).To(Equal("0xfrom"))
			Expect(rpcRequest.To).To(Equal("0xto"))
			Expect(rpcRequest.Value).To(Equal("0xde0b6b3a7640000"))
			Expect(rpcRequest.Gas).To(Equal("0x5208"))
			Expect(rpcRequest.Nonce).To(Equal("0x5"))
			Expect(rpcRequest.Type).To(Equal("0x2"))
		})

		It("should format EIP-1559 fields", func() {
			request := formatters.TransactionRequest{
				MaxFeePerGas:         big.NewInt(1000000000),
				MaxPriorityFeePerGas: big.NewInt(2000000000),
			}

			rpcRequest := formatters.FormatTransactionRequest(request)

			Expect(rpcRequest.MaxFeePerGas).To(Equal("0x3b9aca00"))
			Expect(rpcRequest.MaxPriorityFeePerGas).To(Equal("0x77359400"))
		})
	})

	Describe("FormatFeeHistory", func() {
		It("should format fee history", func() {
			rpcFeeHistory := formatters.RpcFeeHistory{
				BaseFeePerGas: []string{"0x3b9aca00", "0x3b9aca01"},
				GasUsedRatio:  []float64{0.5, 0.6},
				OldestBlock:   "0x100",
				Reward:        [][]string{{"0x1", "0x2"}, {"0x3", "0x4"}},
			}

			feeHistory := formatters.FormatFeeHistory(rpcFeeHistory)

			Expect(len(feeHistory.BaseFeePerGas)).To(Equal(2))
			Expect(feeHistory.BaseFeePerGas[0].Cmp(big.NewInt(1000000000))).To(Equal(0))
			Expect(feeHistory.GasUsedRatio).To(Equal([]float64{0.5, 0.6}))
			Expect(feeHistory.OldestBlock.Cmp(big.NewInt(256))).To(Equal(0))
			Expect(len(feeHistory.Reward)).To(Equal(2))
		})
	})

	Describe("FormatProof", func() {
		It("should format proof", func() {
			rpcProof := formatters.RpcProof{
				Address:      "0xaddress",
				AccountProof: []string{"0xproof1", "0xproof2"},
				Balance:      "0xde0b6b3a7640000",
				CodeHash:     "0xcodehash",
				Nonce:        "0x5",
				StorageHash:  "0xstoragehash",
				StorageProof: []formatters.RpcStorageProof{
					{
						Key:   "0xkey",
						Proof: []string{"0xp1"},
						Value: "0x100",
					},
				},
			}

			proof := formatters.FormatProof(rpcProof)

			Expect(proof.Address).To(Equal("0xaddress"))
			Expect(proof.Balance.Cmp(big.NewInt(1000000000000000000))).To(Equal(0))
			Expect(*proof.Nonce).To(Equal(5))
			Expect(len(proof.StorageProof)).To(Equal(1))
			Expect(proof.StorageProof[0].Value.Cmp(big.NewInt(256))).To(Equal(0))
		})
	})

	Describe("Extract", func() {
		It("should extract keys from value based on formatter", func() {
			value := map[string]any{
				"gasPrice":     "0x1",
				"maxFeePerGas": "0x2",
				"customField":  "ignored",
			}

			formatter := func(v map[string]any) map[string]any {
				return map[string]any{
					"gasPrice": v["gasPrice"],
				}
			}

			extracted := formatters.Extract(value, formatter)

			Expect(extracted).To(HaveKey("gasPrice"))
			Expect(extracted).NotTo(HaveKey("customField"))
		})

		It("should return empty map when formatter is nil", func() {
			value := map[string]any{"key": "value"}
			extracted := formatters.Extract(value, nil)
			Expect(extracted).To(BeEmpty())
		})
	})

	Describe("ExtractKeys", func() {
		It("should extract only specified keys", func() {
			value := map[string]any{
				"a": 1,
				"b": 2,
				"c": 3,
			}

			extracted := formatters.ExtractKeys(value, []string{"a", "c"})

			Expect(extracted).To(HaveKey("a"))
			Expect(extracted).To(HaveKey("c"))
			Expect(extracted).NotTo(HaveKey("b"))
		})
	})

	Describe("OmitKeys", func() {
		It("should omit specified keys", func() {
			value := map[string]any{
				"a": 1,
				"b": 2,
				"c": 3,
			}

			result := formatters.OmitKeys(value, []string{"b"})

			Expect(result).To(HaveKey("a"))
			Expect(result).To(HaveKey("c"))
			Expect(result).NotTo(HaveKey("b"))
		})
	})
})

func intPtr(n int) *int {
	return &n
}
