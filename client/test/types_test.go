package client_test

import (
	"math/big"

	json "github.com/goccy/go-json"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client Types", func() {
	Context("BlockNumber", func() {
		It("should format block tags correctly", func() {
			Expect(client.BlockLatest.String()).To(Equal("latest"))
			Expect(client.BlockPending.String()).To(Equal("pending"))
			Expect(client.BlockEarliest.String()).To(Equal("earliest"))
			Expect(client.BlockSafe.String()).To(Equal("safe"))
			Expect(client.BlockFinalized.String()).To(Equal("finalized"))
		})

		It("should format numeric block numbers as hex", func() {
			bn := client.BlockNumberUint64(12345)
			Expect(bn.String()).To(Equal("0x3039"))
		})
	})

	Context("CallRequest", func() {
		It("should marshal to JSON correctly", func() {
			from := common.HexToAddress("0x1111111111111111111111111111111111111111")
			to := common.HexToAddress("0x2222222222222222222222222222222222222222")

			req := client.CallRequest{
				From:  &from,
				To:    to,
				Data:  []byte{0xa9, 0x05, 0x9c, 0xbb},
				Value: big.NewInt(1000),
				Gas:   21000,
			}

			data, err := json.Marshal(req)
			Expect(err).ToNot(HaveOccurred())

			var parsed map[string]interface{}
			err = json.Unmarshal(data, &parsed)
			Expect(err).ToNot(HaveOccurred())

			Expect(parsed["to"]).To(Equal("0x2222222222222222222222222222222222222222"))
			Expect(parsed["data"]).To(Equal("0xa9059cbb"))
			Expect(parsed["value"]).To(Equal("0x3e8"))
			Expect(parsed["gas"]).To(Equal("0x5208"))
		})

		It("should omit nil fields", func() {
			to := common.HexToAddress("0x2222222222222222222222222222222222222222")

			req := client.CallRequest{
				To:   to,
				Data: []byte{0x12, 0x34},
			}

			data, err := json.Marshal(req)
			Expect(err).ToNot(HaveOccurred())

			var parsed map[string]interface{}
			err = json.Unmarshal(data, &parsed)
			Expect(err).ToNot(HaveOccurred())

			_, hasFrom := parsed["from"]
			Expect(hasFrom).To(BeFalse())

			_, hasValue := parsed["value"]
			Expect(hasValue).To(BeFalse())
		})
	})

	Context("Transaction", func() {
		It("should marshal to JSON correctly", func() {
			from := common.HexToAddress("0x1111111111111111111111111111111111111111")
			to := common.HexToAddress("0x2222222222222222222222222222222222222222")
			nonce := uint64(5)

			tx := client.Transaction{
				From:     from,
				To:       &to,
				Data:     []byte{0xa9, 0x05, 0x9c, 0xbb},
				Value:    big.NewInt(1000),
				Nonce:    &nonce,
				Gas:      21000,
				GasPrice: big.NewInt(20000000000),
			}

			data, err := json.Marshal(tx)
			Expect(err).ToNot(HaveOccurred())

			var parsed map[string]interface{}
			err = json.Unmarshal(data, &parsed)
			Expect(err).ToNot(HaveOccurred())

			Expect(parsed["from"]).To(Equal("0x1111111111111111111111111111111111111111"))
			Expect(parsed["to"]).To(Equal("0x2222222222222222222222222222222222222222"))
			Expect(parsed["nonce"]).To(Equal("0x5"))
			Expect(parsed["gas"]).To(Equal("0x5208"))
		})
	})

	Context("RPCError", func() {
		It("should format error message correctly", func() {
			err := &client.RPCError{
				Code:    -32000,
				Message: "execution reverted",
			}

			Expect(err.Error()).To(ContainSubstring("-32000"))
			Expect(err.Error()).To(ContainSubstring("execution reverted"))
		})

		It("should include data in error message if present", func() {
			err := &client.RPCError{
				Code:    -32000,
				Message: "execution reverted",
				Data:    "0x08c379a0",
			}

			Expect(err.Error()).To(ContainSubstring("0x08c379a0"))
		})
	})

	Context("Receipt", func() {
		It("should unmarshal from JSON correctly", func() {
			receiptJSON := `{
				"transactionHash": "0x1234567890123456789012345678901234567890123456789012345678901234",
				"transactionIndex": "0x1",
				"blockHash": "0xabcdef1234567890123456789012345678901234567890123456789012345678",
				"blockNumber": "0x100",
				"from": "0x1111111111111111111111111111111111111111",
				"to": "0x2222222222222222222222222222222222222222",
				"cumulativeGasUsed": "0x5208",
				"gasUsed": "0x5208",
				"status": "0x1",
				"logs": []
			}`

			var receipt client.Receipt
			err := json.Unmarshal([]byte(receiptJSON), &receipt)
			Expect(err).ToNot(HaveOccurred())

			Expect(receipt.BlockNumber).To(Equal(uint64(256)))
			Expect(receipt.TransactionIndex).To(Equal(uint64(1)))
			Expect(receipt.GasUsed).To(Equal(uint64(21000)))
			Expect(receipt.Status).To(Equal(uint64(1)))
			Expect(receipt.IsSuccess()).To(BeTrue())
		})

		It("should detect failed transaction", func() {
			receipt := client.Receipt{
				Status: 0,
			}
			Expect(receipt.IsSuccess()).To(BeFalse())
		})
	})
})
