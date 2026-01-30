package abi_test

import (
	"math/big"

	"github.com/ChefBingbong/viem-go/abi"
	"github.com/ethereum/go-ethereum/common"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ParseEventLogs", func() {
	var erc20ABI *abi.ABI
	var transferTopic common.Hash

	BeforeEach(func() {
		jsonABI := []byte(`[
			{"type":"event","name":"Transfer","inputs":[{"name":"from","type":"address","indexed":true},{"name":"to","type":"address","indexed":true},{"name":"value","type":"uint256","indexed":false}]},
			{"type":"event","name":"Approval","inputs":[{"name":"owner","type":"address","indexed":true},{"name":"spender","type":"address","indexed":true},{"name":"value","type":"uint256","indexed":false}]}
		]`)
		var err error
		erc20ABI, err = abi.Parse(jsonABI)
		Expect(err).ToNot(HaveOccurred())

		// Transfer event topic
		transferTopic = common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
	})

	Context("basic parsing", func() {
		It("should parse a single Transfer event", func() {
			from := common.HexToAddress("0x1111111111111111111111111111111111111111")
			to := common.HexToAddress("0x2222222222222222222222222222222222222222")

			// Encode the value (non-indexed)
			value := big.NewInt(1000)
			valueBytes := common.LeftPadBytes(value.Bytes(), 32)

			logs := []abi.RawLog{
				{
					Topics: []common.Hash{
						transferTopic,
						common.BytesToHash(from.Bytes()),
						common.BytesToHash(to.Bytes()),
					},
					Data: valueBytes,
				},
			}

			parsed := erc20ABI.ParseEventLogs(logs, nil)
			Expect(parsed).To(HaveLen(1))
			Expect(parsed[0].EventName).To(Equal("Transfer"))
		})

		It("should parse multiple logs", func() {
			from := common.HexToAddress("0x1111111111111111111111111111111111111111")
			to := common.HexToAddress("0x2222222222222222222222222222222222222222")
			value := big.NewInt(1000)
			valueBytes := common.LeftPadBytes(value.Bytes(), 32)

			logs := []abi.RawLog{
				{
					Topics: []common.Hash{
						transferTopic,
						common.BytesToHash(from.Bytes()),
						common.BytesToHash(to.Bytes()),
					},
					Data: valueBytes,
				},
				{
					Topics: []common.Hash{
						transferTopic,
						common.BytesToHash(to.Bytes()),
						common.BytesToHash(from.Bytes()),
					},
					Data: valueBytes,
				},
			}

			parsed := erc20ABI.ParseEventLogs(logs, nil)
			Expect(parsed).To(HaveLen(2))
			Expect(parsed[0].EventName).To(Equal("Transfer"))
			Expect(parsed[1].EventName).To(Equal("Transfer"))
		})

		It("should skip logs with unknown topics", func() {
			unknownTopic := common.HexToHash("0xdeadbeef00000000000000000000000000000000000000000000000000000000")

			logs := []abi.RawLog{
				{
					Topics: []common.Hash{unknownTopic},
					Data:   []byte{},
				},
			}

			parsed := erc20ABI.ParseEventLogs(logs, nil)
			Expect(parsed).To(HaveLen(0))
		})
	})

	Context("filtering by event name", func() {
		It("should filter by single event name", func() {
			from := common.HexToAddress("0x1111111111111111111111111111111111111111")
			to := common.HexToAddress("0x2222222222222222222222222222222222222222")
			value := big.NewInt(1000)
			valueBytes := common.LeftPadBytes(value.Bytes(), 32)

			approvalTopic := common.HexToHash("0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925")

			logs := []abi.RawLog{
				{
					Topics: []common.Hash{
						transferTopic,
						common.BytesToHash(from.Bytes()),
						common.BytesToHash(to.Bytes()),
					},
					Data: valueBytes,
				},
				{
					Topics: []common.Hash{
						approvalTopic,
						common.BytesToHash(from.Bytes()),
						common.BytesToHash(to.Bytes()),
					},
					Data: valueBytes,
				},
			}

			opts := &abi.ParseEventLogsOptions{
				EventName: []string{"Transfer"},
				Strict:    true,
			}
			parsed := erc20ABI.ParseEventLogs(logs, opts)
			Expect(parsed).To(HaveLen(1))
			Expect(parsed[0].EventName).To(Equal("Transfer"))
		})

		It("should filter by multiple event names", func() {
			from := common.HexToAddress("0x1111111111111111111111111111111111111111")
			to := common.HexToAddress("0x2222222222222222222222222222222222222222")
			value := big.NewInt(1000)
			valueBytes := common.LeftPadBytes(value.Bytes(), 32)

			approvalTopic := common.HexToHash("0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925")

			logs := []abi.RawLog{
				{
					Topics: []common.Hash{
						transferTopic,
						common.BytesToHash(from.Bytes()),
						common.BytesToHash(to.Bytes()),
					},
					Data: valueBytes,
				},
				{
					Topics: []common.Hash{
						approvalTopic,
						common.BytesToHash(from.Bytes()),
						common.BytesToHash(to.Bytes()),
					},
					Data: valueBytes,
				},
			}

			opts := &abi.ParseEventLogsOptions{
				EventName: []string{"Transfer", "Approval"},
				Strict:    true,
			}
			parsed := erc20ABI.ParseEventLogs(logs, opts)
			Expect(parsed).To(HaveLen(2))
		})
	})
})
