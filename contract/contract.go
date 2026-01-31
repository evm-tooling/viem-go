package contract

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/abi"
	"github.com/ChefBingbong/viem-go/client"
)

// Contract represents a smart contract instance with ABI and client.
type Contract struct {
	address common.Address
	abi     *abi.ABI
	client  *client.PublicClient
}

// NewContract creates a new Contract instance.
func NewContract(address common.Address, abiJSON []byte, c *client.PublicClient) (*Contract, error) {
	parsedABI, err := abi.Parse(abiJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	return &Contract{
		address: address,
		abi:     parsedABI,
		client:  c,
	}, nil
}

// NewContractWithABI creates a new Contract instance with a pre-parsed ABI.
func NewContractWithABI(address common.Address, parsedABI *abi.ABI, c *client.PublicClient) *Contract {
	return &Contract{
		address: address,
		abi:     parsedABI,
		client:  c,
	}
}

// MustNewContract creates a new Contract instance, panicking on error.
func MustNewContract(address common.Address, abiJSON []byte, c *client.PublicClient) *Contract {
	contract, err := NewContract(address, abiJSON, c)
	if err != nil {
		panic(err)
	}
	return contract
}

// Address returns the contract address.
func (c *Contract) Address() common.Address {
	return c.address
}

// ABI returns the contract's ABI.
func (c *Contract) ABI() *abi.ABI {
	return c.abi
}

// Client returns the underlying RPC client.
func (c *Contract) Client() *client.PublicClient {
	return c.client
}

// HasFunction returns true if the contract ABI contains the specified function.
func (c *Contract) HasFunction(name string) bool {
	return c.abi.HasFunction(name)
}

// HasEvent returns true if the contract ABI contains the specified event.
func (c *Contract) HasEvent(name string) bool {
	return c.abi.HasEvent(name)
}

// FunctionNames returns the names of all functions in the contract.
func (c *Contract) FunctionNames() []string {
	return c.abi.FunctionNames()
}

// EventNames returns the names of all events in the contract.
func (c *Contract) EventNames() []string {
	return c.abi.EventNames()
}

// GetFunction returns information about a function by name.
func (c *Contract) GetFunction(name string) (*abi.Function, error) {
	return c.abi.GetFunction(name)
}

// GetEvent returns information about an event by name.
func (c *Contract) GetEvent(name string) (*abi.Event, error) {
	return c.abi.GetEvent(name)
}

// EncodeCall encodes a function call for this contract.
func (c *Contract) EncodeCall(method string, args ...any) ([]byte, error) {
	return c.abi.EncodeCall(method, args...)
}

// DecodeReturn decodes the return data from a function call.
func (c *Contract) DecodeReturn(method string, data []byte) ([]any, error) {
	return c.abi.DecodeReturn(method, data)
}

// DecodeEvent decodes event data using the contract's ABI.
func (c *Contract) DecodeEvent(name string, topics []common.Hash, data []byte) (map[string]any, error) {
	return c.abi.DecodeEvent(name, topics, data)
}

// Clone creates a copy of the contract with a different address.
func (c *Contract) Clone(newAddress common.Address) *Contract {
	return &Contract{
		address: newAddress,
		abi:     c.abi,
		client:  c.client,
	}
}

// WithClient creates a copy of the contract with a different client.
func (c *Contract) WithClient(newClient *client.PublicClient) *Contract {
	return &Contract{
		address: c.address,
		abi:     c.abi,
		client:  newClient,
	}
}
