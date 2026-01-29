package contract_test

import (
	"context"
	"encoding/json"
	"math/big"

	erc200 "github.com/ChefBingbong/viem-go/_typed/templates/erc200"
	"github.com/ChefBingbong/viem-go/client"
	"github.com/ChefBingbong/viem-go/contract"
	"github.com/ethereum/go-ethereum/common"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// Newerc200.Erc200 creates a new typed ERC20 contract.
func NewErc200(address common.Address, c *client.Client) (*erc200.Erc200, error) {
	return erc200.New(address, c)
}

// ============================================================================
// TESTS
// ============================================================================

var erc20ABI = `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"owner","type":"address"},{"name":"spender","type":"address"}],"name":"allowance","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"type":"function"},{"constant":false,"inputs":[{"name":"spender","type":"address"},{"name":"value","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"type":"function"},{"constant":false,"inputs":[{"name":"from","type":"address"},{"name":"to","type":"address"},{"name":"value","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"","type":"bool"}],"type":"function"}]`

// MockTransport for testing
type mockTransport struct {
	responses map[string]json.RawMessage
}

func newMockTransport() *mockTransport {
	return &mockTransport{responses: make(map[string]json.RawMessage)}
}

func (m *mockTransport) SetResponse(method string, result any) {
	data, _ := json.Marshal(result)
	m.responses[method] = data
}

func (m *mockTransport) Call(ctx context.Context, method string, params ...any) (json.RawMessage, error) {
	if resp, ok := m.responses[method]; ok {
		return resp, nil
	}
	return json.RawMessage(`"0x"`), nil
}

func (m *mockTransport) Close() error { return nil }

var _ = Describe("Typed Contract Methods", func() {
	var (
		mock      *mockTransport
		rpcClient *client.Client
		ctx       context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		mock = newMockTransport()
		mock.SetResponse("eth_chainId", "0x1")
		var err error
		rpcClient, err = client.NewClientWithTransport(mock)
		Expect(err).ToNot(HaveOccurred())
	})

	Context("using ReadTyped directly with method descriptors", func() {
		It("should return typed values", func() {
			// Setup mock response for name()
			nameEncoded := "0x" +
				"0000000000000000000000000000000000000000000000000000000000000020" +
				"0000000000000000000000000000000000000000000000000000000000000004" +
				"5553444300000000000000000000000000000000000000000000000000000000"
			mock.SetResponse("eth_call", nameEncoded)

			// Create base contract
			addr := common.HexToAddress("0x1234567890123456789012345678901234567890")
			c, err := contract.NewContract(addr, []byte(erc20ABI), rpcClient)
			Expect(err).ToNot(HaveOccurred())

			// Use ReadTyped with the method descriptor - return type is inferred!
			name, err := contract.ReadTyped(c, ctx, erc200.Methods.Name)
			Expect(err).ToNot(HaveOccurred())
			Expect(name).To(Equal("USDC")) // name is string, not any!
		})

		It("should return typed *big.Int for balanceOf", func() {
			balanceEncoded := "0x00000000000000000000000000000000000000000000000000000000000f4240"
			mock.SetResponse("eth_call", balanceEncoded)

			addr := common.HexToAddress("0x1234567890123456789012345678901234567890")
			c, err := contract.NewContract(addr, []byte(erc20ABI), rpcClient)
			Expect(err).ToNot(HaveOccurred())

			owner := common.HexToAddress("0xaaaa")

			// Return type *big.Int is inferred from ERC20.BalanceOf
			balance, err := contract.ReadTyped(c, ctx, erc200.Methods.BalanceOf, owner)
			Expect(err).ToNot(HaveOccurred())
			Expect(balance.Cmp(big.NewInt(1000000))).To(Equal(0)) // balance is *big.Int!
		})
	})

	Context("using erc200.Erc200 wrapper", func() {
		It("should provide typed methods", func() {
			nameEncoded := "0x" +
				"0000000000000000000000000000000000000000000000000000000000000020" +
				"0000000000000000000000000000000000000000000000000000000000000004" +
				"5553444300000000000000000000000000000000000000000000000000000000"
			mock.SetResponse("eth_call", nameEncoded)

			addr := common.HexToAddress("0x1234567890123456789012345678901234567890")
			token, err := NewErc200(addr, rpcClient)
			Expect(err).ToNot(HaveOccurred())

			// Clean typed API
			name, err := token.Name(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(name).To(Equal("USDC"))
		})
	})
})
