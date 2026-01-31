// Package erc1155 provides bindings for the ERC1155 multi-token standard.
package erc1155

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/client"
	"github.com/ChefBingbong/viem-go/contract"
	"github.com/ChefBingbong/viem-go/types"
)

// ContractABI is the ABI of the ERC1155 contract.
var ContractABI = `[{"constant":true,"inputs":[{"name":"account","type":"address"},{"name":"id","type":"uint256"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"accounts","type":"address[]"},{"name":"ids","type":"uint256[]"}],"name":"balanceOfBatch","outputs":[{"name":"","type":"uint256[]"}],"type":"function"},{"constant":true,"inputs":[{"name":"account","type":"address"},{"name":"operator","type":"address"}],"name":"isApprovedForAll","outputs":[{"name":"","type":"bool"}],"type":"function"},{"constant":true,"inputs":[{"name":"id","type":"uint256"}],"name":"uri","outputs":[{"name":"","type":"string"}],"type":"function"},{"constant":false,"inputs":[{"name":"operator","type":"address"},{"name":"approved","type":"bool"}],"name":"setApprovalForAll","outputs":[],"type":"function"},{"constant":false,"inputs":[{"name":"from","type":"address"},{"name":"to","type":"address"},{"name":"id","type":"uint256"},{"name":"amount","type":"uint256"},{"name":"data","type":"bytes"}],"name":"safeTransferFrom","outputs":[],"type":"function"},{"constant":false,"inputs":[{"name":"from","type":"address"},{"name":"to","type":"address"},{"name":"ids","type":"uint256[]"},{"name":"amounts","type":"uint256[]"},{"name":"data","type":"bytes"}],"name":"safeBatchTransferFrom","outputs":[],"type":"function"},{"anonymous":false,"inputs":[{"indexed":true,"name":"operator","type":"address"},{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"id","type":"uint256"},{"indexed":false,"name":"value","type":"uint256"}],"name":"TransferSingle","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"operator","type":"address"},{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"ids","type":"uint256[]"},{"indexed":false,"name":"values","type":"uint256[]"}],"name":"TransferBatch","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"account","type":"address"},{"indexed":true,"name":"operator","type":"address"},{"indexed":false,"name":"approved","type":"bool"}],"name":"ApprovalForAll","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"value","type":"string"},{"indexed":true,"name":"id","type":"uint256"}],"name":"URI","type":"event"}]`

// ERC1155 is a binding to an ERC1155 multi-token contract.
type ERC1155 struct {
	contract *contract.Contract
}

// New creates a new ERC1155 contract binding.
func New(address common.Address, c *client.PublicClient) (*ERC1155, error) {
	cont, err := contract.NewContract(address, []byte(ContractABI), c)
	if err != nil {
		return nil, err
	}
	return &ERC1155{contract: cont}, nil
}

// MustNew creates a new ERC1155 contract binding, panicking on error.
func MustNew(address common.Address, c *client.PublicClient) *ERC1155 {
	cont, err := New(address, c)
	if err != nil {
		panic(err)
	}
	return cont
}

// Address returns the contract address.
func (e *ERC1155) Address() common.Address {
	return e.contract.Address()
}

// Contract returns the underlying contract instance.
func (e *ERC1155) Contract() *contract.Contract {
	return e.contract
}

// ---- Read Methods (Public Actions) ----

// URI returns the URI for a token type.
func (e *ERC1155) URI(ctx context.Context, id *big.Int) (string, error) {
	return e.contract.ReadString(ctx, "uri", id)
}

// BalanceOf returns the balance of a specific token for an account.
func (e *ERC1155) BalanceOf(ctx context.Context, account common.Address, id *big.Int) (*big.Int, error) {
	return e.contract.ReadBigInt(ctx, "balanceOf", account, id)
}

// BalanceOfBatch returns the balances of multiple tokens for multiple accounts.
func (e *ERC1155) BalanceOfBatch(ctx context.Context, accounts []common.Address, ids []*big.Int) ([]*big.Int, error) {
	result, err := e.contract.Read(ctx, "balanceOfBatch", accounts, ids)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}
	// Convert result to []*big.Int
	balances, ok := result[0].([]*big.Int)
	if !ok {
		// Try to convert from []interface{}
		if arr, ok := result[0].([]interface{}); ok {
			balances = make([]*big.Int, len(arr))
			for i, v := range arr {
				if b, ok := v.(*big.Int); ok {
					balances[i] = b
				}
			}
		}
	}
	return balances, nil
}

// IsApprovedForAll returns if an operator is approved for all tokens of an account.
func (e *ERC1155) IsApprovedForAll(ctx context.Context, account, operator common.Address) (bool, error) {
	return e.contract.ReadBool(ctx, "isApprovedForAll", account, operator)
}

// ---- Write Methods (Prepare Transaction for Signing) ----

// PrepareSetApprovalForAll prepares a setApprovalForAll transaction for signing.
func (e *ERC1155) PrepareSetApprovalForAll(ctx context.Context, opts contract.WriteOptions, operator common.Address, approved bool) (*types.Transaction, error) {
	return e.contract.PrepareTransaction(ctx, opts, "setApprovalForAll", operator, approved)
}

// PrepareSafeTransferFrom prepares a safeTransferFrom transaction for signing.
func (e *ERC1155) PrepareSafeTransferFrom(ctx context.Context, opts contract.WriteOptions, from, to common.Address, id, amount *big.Int, data []byte) (*types.Transaction, error) {
	return e.contract.PrepareTransaction(ctx, opts, "safeTransferFrom", from, to, id, amount, data)
}

// PrepareSafeBatchTransferFrom prepares a safeBatchTransferFrom transaction for signing.
func (e *ERC1155) PrepareSafeBatchTransferFrom(ctx context.Context, opts contract.WriteOptions, from, to common.Address, ids, amounts []*big.Int, data []byte) (*types.Transaction, error) {
	return e.contract.PrepareTransaction(ctx, opts, "safeBatchTransferFrom", from, to, ids, amounts, data)
}

// ---- Gas Estimation ----

// EstimateSetApprovalForAll estimates gas for a setApprovalForAll transaction.
func (e *ERC1155) EstimateSetApprovalForAll(ctx context.Context, opts contract.WriteOptions, operator common.Address, approved bool) (uint64, error) {
	return e.contract.EstimateGas(ctx, opts, "setApprovalForAll", operator, approved)
}

// EstimateSafeTransferFrom estimates gas for a safeTransferFrom transaction.
func (e *ERC1155) EstimateSafeTransferFrom(ctx context.Context, opts contract.WriteOptions, from, to common.Address, id, amount *big.Int, data []byte) (uint64, error) {
	return e.contract.EstimateGas(ctx, opts, "safeTransferFrom", from, to, id, amount, data)
}

// EstimateSafeBatchTransferFrom estimates gas for a safeBatchTransferFrom transaction.
func (e *ERC1155) EstimateSafeBatchTransferFrom(ctx context.Context, opts contract.WriteOptions, from, to common.Address, ids, amounts []*big.Int, data []byte) (uint64, error) {
	return e.contract.EstimateGas(ctx, opts, "safeBatchTransferFrom", from, to, ids, amounts, data)
}

// ---- Events ----

// TransferSingleEvent represents a TransferSingle event.
type TransferSingleEvent struct {
	Operator common.Address
	From     common.Address
	To       common.Address
	Id       *big.Int
	Value    *big.Int
}

// TransferBatchEvent represents a TransferBatch event.
type TransferBatchEvent struct {
	Operator common.Address
	From     common.Address
	To       common.Address
	Ids      []*big.Int
	Values   []*big.Int
}

// ApprovalForAllEvent represents an ApprovalForAll event.
type ApprovalForAllEvent struct {
	Account  common.Address
	Operator common.Address
	Approved bool
}

// URIEvent represents a URI event.
type URIEvent struct {
	Value string
	Id    *big.Int
}

// ParseTransferSingle parses a TransferSingle event from a log.
func (e *ERC1155) ParseTransferSingle(log types.Log) (*TransferSingleEvent, error) {
	event, err := e.contract.DecodeEvent("TransferSingle", log.Topics, log.Data)
	if err != nil {
		return nil, err
	}

	operator, _ := event["operator"].(common.Address)
	from, _ := event["from"].(common.Address)
	to, _ := event["to"].(common.Address)
	id, _ := event["id"].(*big.Int)
	value, _ := event["value"].(*big.Int)

	return &TransferSingleEvent{
		Operator: operator,
		From:     from,
		To:       to,
		Id:       id,
		Value:    value,
	}, nil
}

// ParseTransferBatch parses a TransferBatch event from a log.
func (e *ERC1155) ParseTransferBatch(log types.Log) (*TransferBatchEvent, error) {
	event, err := e.contract.DecodeEvent("TransferBatch", log.Topics, log.Data)
	if err != nil {
		return nil, err
	}

	return &TransferBatchEvent{
		Operator: event["operator"].(common.Address),
		From:     event["from"].(common.Address),
		To:       event["to"].(common.Address),
		Ids:      event["ids"].([]*big.Int),
		Values:   event["values"].([]*big.Int),
	}, nil
}

// ParseApprovalForAll parses an ApprovalForAll event from a log.
func (e *ERC1155) ParseApprovalForAll(log types.Log) (*ApprovalForAllEvent, error) {
	event, err := e.contract.DecodeEvent("ApprovalForAll", log.Topics, log.Data)
	if err != nil {
		return nil, err
	}

	return &ApprovalForAllEvent{
		Account:  event["account"].(common.Address),
		Operator: event["operator"].(common.Address),
		Approved: event["approved"].(bool),
	}, nil
}

// ParseURI parses a URI event from a log.
func (e *ERC1155) ParseURI(log types.Log) (*URIEvent, error) {
	event, err := e.contract.DecodeEvent("URI", log.Topics, log.Data)
	if err != nil {
		return nil, err
	}

	return &URIEvent{
		Value: event["value"].(string),
		Id:    event["id"].(*big.Int),
	}, nil
}
