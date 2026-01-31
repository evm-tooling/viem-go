// Package erc20 provides bindings for the ERC20 token standard.
package erc20

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/client"
	"github.com/ChefBingbong/viem-go/contract"
	"github.com/ChefBingbong/viem-go/types"
)

// ContractABI is the ABI of the ERC20 contract.
var ContractABI = `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"owner","type":"address"},{"name":"spender","type":"address"}],"name":"allowance","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"type":"function"},{"constant":false,"inputs":[{"name":"spender","type":"address"},{"name":"value","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"type":"function"},{"constant":false,"inputs":[{"name":"from","type":"address"},{"name":"to","type":"address"},{"name":"value","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"","type":"bool"}],"type":"function"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"},{"indexed":true,"name":"spender","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Approval","type":"event"}]`

// ERC20 is a binding to an ERC20 token contract.
type ERC20 struct {
	contract *contract.Contract
}

// New creates a new ERC20 contract binding.
func New(address common.Address, c *client.PublicClient) (*ERC20, error) {
	cont, err := contract.NewContract(address, []byte(ContractABI), c)
	if err != nil {
		return nil, err
	}
	return &ERC20{contract: cont}, nil
}

// MustNew creates a new ERC20 contract binding, panicking on error.
func MustNew(address common.Address, c *client.PublicClient) *ERC20 {
	cont, err := New(address, c)
	if err != nil {
		panic(err)
	}
	return cont
}

// Address returns the contract address.
func (e *ERC20) Address() common.Address {
	return e.contract.Address()
}

// Contract returns the underlying contract instance.
func (e *ERC20) Contract() *contract.Contract {
	return e.contract
}

// ---- Read Methods (Public Actions) ----

// Name returns the token name.
func (e *ERC20) Name(ctx context.Context) (string, error) {
	return e.contract.ReadString(ctx, "name")
}

// Symbol returns the token symbol.
func (e *ERC20) Symbol(ctx context.Context) (string, error) {
	return e.contract.ReadString(ctx, "symbol")
}

// Decimals returns the token decimals.
func (e *ERC20) Decimals(ctx context.Context) (uint8, error) {
	return e.contract.ReadUint8(ctx, "decimals")
}

// TotalSupply returns the total token supply.
func (e *ERC20) TotalSupply(ctx context.Context) (*big.Int, error) {
	return e.contract.ReadBigInt(ctx, "totalSupply")
}

// BalanceOf returns the token balance of an address.
func (e *ERC20) BalanceOf(ctx context.Context, owner common.Address) (*big.Int, error) {
	return e.contract.ReadBigInt(ctx, "balanceOf", owner)
}

// Allowance returns the allowance of a spender for an owner.
func (e *ERC20) Allowance(ctx context.Context, owner, spender common.Address) (*big.Int, error) {
	return e.contract.ReadBigInt(ctx, "allowance", owner, spender)
}

// ---- Write Methods (Prepare Transaction for Signing) ----

// PrepareTransfer prepares a transfer transaction for signing.
// Use a WalletClient to sign and send the returned transaction.
func (e *ERC20) PrepareTransfer(ctx context.Context, opts contract.WriteOptions, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return e.contract.PrepareTransaction(ctx, opts, "transfer", to, amount)
}

// PrepareApprove prepares an approve transaction for signing.
func (e *ERC20) PrepareApprove(ctx context.Context, opts contract.WriteOptions, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return e.contract.PrepareTransaction(ctx, opts, "approve", spender, amount)
}

// PrepareTransferFrom prepares a transferFrom transaction for signing.
func (e *ERC20) PrepareTransferFrom(ctx context.Context, opts contract.WriteOptions, from, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return e.contract.PrepareTransaction(ctx, opts, "transferFrom", from, to, amount)
}

// ---- Gas Estimation ----

// EstimateTransfer estimates gas for a transfer.
func (e *ERC20) EstimateTransfer(ctx context.Context, opts contract.WriteOptions, to common.Address, amount *big.Int) (uint64, error) {
	return e.contract.EstimateGas(ctx, opts, "transfer", to, amount)
}

// EstimateApprove estimates gas for an approve.
func (e *ERC20) EstimateApprove(ctx context.Context, opts contract.WriteOptions, spender common.Address, amount *big.Int) (uint64, error) {
	return e.contract.EstimateGas(ctx, opts, "approve", spender, amount)
}

// EstimateTransferFrom estimates gas for a transferFrom.
func (e *ERC20) EstimateTransferFrom(ctx context.Context, opts contract.WriteOptions, from, to common.Address, amount *big.Int) (uint64, error) {
	return e.contract.EstimateGas(ctx, opts, "transferFrom", from, to, amount)
}

// ---- Events ----

// TransferEvent represents a Transfer event.
type TransferEvent struct {
	From  common.Address
	To    common.Address
	Value *big.Int
}

// ApprovalEvent represents an Approval event.
type ApprovalEvent struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
}

// ParseTransfer parses a Transfer event from a log.
func (e *ERC20) ParseTransfer(log types.Log) (*TransferEvent, error) {
	event, err := e.contract.DecodeEvent("Transfer", log.Topics, log.Data)
	if err != nil {
		return nil, err
	}

	return &TransferEvent{
		From:  event["from"].(common.Address),
		To:    event["to"].(common.Address),
		Value: event["value"].(*big.Int),
	}, nil
}

// ParseApproval parses an Approval event from a log.
func (e *ERC20) ParseApproval(log types.Log) (*ApprovalEvent, error) {
	event, err := e.contract.DecodeEvent("Approval", log.Topics, log.Data)
	if err != nil {
		return nil, err
	}

	return &ApprovalEvent{
		Owner:   event["owner"].(common.Address),
		Spender: event["spender"].(common.Address),
		Value:   event["value"].(*big.Int),
	}, nil
}
