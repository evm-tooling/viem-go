package contract_test

import (
	"context"
	"encoding/json"
	"math/big"

	"github.com/ChefBingbong/viem-go/client"
	"github.com/ChefBingbong/viem-go/contract"
	"github.com/ethereum/go-ethereum/common"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// ============================================================================
// STEP 1: Define your contract's method template
// ============================================================================

// ERC20Methods defines the typed methods for an ERC20 contract.
// Each field is a method descriptor with its return type encoded in the generic parameter.
type ERC20Methods struct {
	// Read methods (view/pure functions)
	Name        contract.ReadString   // returns string
	Symbol      contract.ReadString   // returns string
	Decimals    contract.ReadUint8    // returns uint8
	TotalSupply contract.ReadBigInt   // returns *big.Int
	BalanceOf   contract.ReadBigInt   // returns *big.Int
	Allowance   contract.ReadBigInt   // returns *big.Int

	// Write methods (state-changing functions)
	Transfer     contract.WriteMethod
	Approve      contract.WriteMethod
	TransferFrom contract.WriteMethod
}

// ERC20 is the method template instance with all method names defined.
// This acts as a "schema" for the contract.
var ERC20 = ERC20Methods{
	Name:         contract.ReadString{Name: "name"},
	Symbol:       contract.ReadString{Name: "symbol"},
	Decimals:     contract.ReadUint8{Name: "decimals"},
	TotalSupply:  contract.ReadBigInt{Name: "totalSupply"},
	BalanceOf:    contract.ReadBigInt{Name: "balanceOf"},
	Allowance:    contract.ReadBigInt{Name: "allowance"},
	Transfer:     contract.WriteMethod{Name: "transfer"},
	Approve:      contract.WriteMethod{Name: "approve"},
	TransferFrom: contract.WriteMethod{Name: "transferFrom"},
}

// ============================================================================
// STEP 2: Create a typed contract wrapper (optional but cleaner)
// ============================================================================

// TypedERC20 wraps a Contract with typed method access.
type TypedERC20 struct {
	*contract.Contract
	M ERC20Methods // Method descriptors
}

// NewTypedERC20 creates a new typed ERC20 contract.
func NewTypedERC20(address common.Address, c *client.Client) (*TypedERC20, error) {
	cont, err := contract.NewContract(address, []byte(erc20ABI), c)
	if err != nil {
		return nil, err
	}
	return &TypedERC20{Contract: cont, M: ERC20}, nil
}

// Name returns the token name with typed return.
func (t *TypedERC20) Name(ctx context.Context) (string, error) {
	return contract.ReadTyped(t.Contract, ctx, t.M.Name)
}

// BalanceOf returns the balance with typed return.
func (t *TypedERC20) BalanceOf(ctx context.Context, owner common.Address) (*big.Int, error) {
	return contract.ReadTyped(t.Contract, ctx, t.M.BalanceOf, owner)
}

// Transfer sends a transfer transaction.
func (t *TypedERC20) Transfer(ctx context.Context, opts contract.WriteOptions, to common.Address, amount *big.Int) (common.Hash, error) {
	return contract.WriteTyped(t.Contract, ctx, opts, t.M.Transfer, to, amount)
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
			name, err := contract.ReadTyped(c, ctx, ERC20.Name)
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
			balance, err := contract.ReadTyped(c, ctx, ERC20.BalanceOf, owner)
			Expect(err).ToNot(HaveOccurred())
			Expect(balance.Cmp(big.NewInt(1000000))).To(Equal(0)) // balance is *big.Int!
		})
	})

	Context("using TypedERC20 wrapper", func() {
		It("should provide typed methods", func() {
			nameEncoded := "0x" +
				"0000000000000000000000000000000000000000000000000000000000000020" +
				"0000000000000000000000000000000000000000000000000000000000000004" +
				"5553444300000000000000000000000000000000000000000000000000000000"
			mock.SetResponse("eth_call", nameEncoded)

			addr := common.HexToAddress("0x1234567890123456789012345678901234567890")
			token, err := NewTypedERC20(addr, rpcClient)
			Expect(err).ToNot(HaveOccurred())

			// Clean typed API
			name, err := token.Name(ctx)
			Expect(err).ToNot(HaveOccurred())
			Expect(name).To(Equal("USDC"))
		})
	})
})
