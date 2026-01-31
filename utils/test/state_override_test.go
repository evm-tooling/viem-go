package utils_test

import (
	"math/big"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ChefBingbong/viem-go/utils"
)

var _ = Describe("StateOverride", func() {
	Describe("SerializeStateMapping", func() {
		It("should serialize valid state mapping", func() {
			mapping := []utils.StateMapping{
				{
					Slot:  "0x0000000000000000000000000000000000000000000000000000000000000001",
					Value: "0x0000000000000000000000000000000000000000000000000000000000000002",
				},
			}

			result, err := utils.SerializeStateMapping(mapping)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(HaveKey("0x0000000000000000000000000000000000000000000000000000000000000001"))
			Expect(result["0x0000000000000000000000000000000000000000000000000000000000000001"]).To(Equal("0x0000000000000000000000000000000000000000000000000000000000000002"))
		})

		It("should return nil for empty mapping", func() {
			result, err := utils.SerializeStateMapping(nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(BeNil())

			result, err = utils.SerializeStateMapping([]utils.StateMapping{})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(BeNil())
		})

		It("should error on invalid slot length", func() {
			mapping := []utils.StateMapping{
				{
					Slot:  "0x01", // Too short
					Value: "0x0000000000000000000000000000000000000000000000000000000000000002",
				},
			}

			_, err := utils.SerializeStateMapping(mapping)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid bytes length"))
		})

		It("should error on invalid value length", func() {
			mapping := []utils.StateMapping{
				{
					Slot:  "0x0000000000000000000000000000000000000000000000000000000000000001",
					Value: "0x02", // Too short
				},
			}

			_, err := utils.SerializeStateMapping(mapping)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid bytes length"))
		})
	})

	Describe("SerializeAccountStateOverride", func() {
		It("should serialize balance", func() {
			override := utils.AccountStateOverride{
				Balance: big.NewInt(1000000000000000000), // 1 ETH
			}

			result, err := utils.SerializeAccountStateOverride(override)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Balance).To(Equal("0xde0b6b3a7640000"))
		})

		It("should serialize nonce", func() {
			nonce := uint64(5)
			override := utils.AccountStateOverride{
				Nonce: &nonce,
			}

			result, err := utils.SerializeAccountStateOverride(override)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Nonce).To(Equal("0x5"))
		})

		It("should serialize code", func() {
			override := utils.AccountStateOverride{
				Code: "0x608060405234801561001057600080fd5b50",
			}

			result, err := utils.SerializeAccountStateOverride(override)
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Code).To(Equal("0x608060405234801561001057600080fd5b50"))
		})

		It("should error when both state and stateDiff are set", func() {
			override := utils.AccountStateOverride{
				State: []utils.StateMapping{
					{
						Slot:  "0x0000000000000000000000000000000000000000000000000000000000000001",
						Value: "0x0000000000000000000000000000000000000000000000000000000000000001",
					},
				},
				StateDiff: []utils.StateMapping{
					{
						Slot:  "0x0000000000000000000000000000000000000000000000000000000000000002",
						Value: "0x0000000000000000000000000000000000000000000000000000000000000002",
					},
				},
			}

			_, err := utils.SerializeAccountStateOverride(override)
			Expect(err).To(Equal(utils.ErrStateAssignmentConflict))
		})
	})

	Describe("SerializeStateOverride", func() {
		It("should serialize valid state override", func() {
			overrides := utils.StateOverride{
				{
					Address: "0x1234567890123456789012345678901234567890",
					Balance: big.NewInt(1000),
				},
			}

			result, err := utils.SerializeStateOverride(overrides)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(HaveKey("0x1234567890123456789012345678901234567890"))
		})

		It("should return nil for empty override", func() {
			result, err := utils.SerializeStateOverride(nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(BeNil())
		})

		It("should error on invalid address", func() {
			overrides := utils.StateOverride{
				{
					Address: "not-an-address",
					Balance: big.NewInt(1000),
				},
			}

			_, err := utils.SerializeStateOverride(overrides)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid address"))
		})

		It("should error on duplicate address", func() {
			overrides := utils.StateOverride{
				{
					Address: "0x1234567890123456789012345678901234567890",
					Balance: big.NewInt(1000),
				},
				{
					Address: "0x1234567890123456789012345678901234567890",
					Balance: big.NewInt(2000),
				},
			}

			_, err := utils.SerializeStateOverride(overrides)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("already set"))
		})
	})

	Describe("Helper functions", func() {
		It("should create state override with NewStateOverride", func() {
			override := utils.NewStateOverride(
				utils.AccountStateOverride{
					Address: "0x1234567890123456789012345678901234567890",
					Balance: big.NewInt(1000),
				},
			)

			Expect(override).To(HaveLen(1))
			Expect(override[0].Address).To(Equal("0x1234567890123456789012345678901234567890"))
		})

		It("should create state mapping with NewStateMapping", func() {
			mapping := utils.NewStateMapping(
				"0x0000000000000000000000000000000000000000000000000000000000000001",
				"0x0000000000000000000000000000000000000000000000000000000000000002",
			)

			Expect(mapping.Slot).To(Equal("0x0000000000000000000000000000000000000000000000000000000000000001"))
			Expect(mapping.Value).To(Equal("0x0000000000000000000000000000000000000000000000000000000000000002"))
		})
	})
})
