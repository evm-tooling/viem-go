// Package erc721 provides bindings for the ERC721 NFT standard.
package erc721

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/client"
	"github.com/ChefBingbong/viem-go/contract"
	"github.com/ChefBingbong/viem-go/types"
)

// ContractABI is the ABI of the ERC721 contract.
var ContractABI = `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"type":"function"},{"constant":true,"inputs":[{"name":"tokenId","type":"uint256"}],"name":"tokenURI","outputs":[{"name":"","type":"string"}],"type":"function"},{"constant":true,"inputs":[{"name":"owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":true,"inputs":[{"name":"tokenId","type":"uint256"}],"name":"ownerOf","outputs":[{"name":"","type":"address"}],"type":"function"},{"constant":true,"inputs":[{"name":"tokenId","type":"uint256"}],"name":"getApproved","outputs":[{"name":"","type":"address"}],"type":"function"},{"constant":true,"inputs":[{"name":"owner","type":"address"},{"name":"operator","type":"address"}],"name":"isApprovedForAll","outputs":[{"name":"","type":"bool"}],"type":"function"},{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"tokenId","type":"uint256"}],"name":"approve","outputs":[],"type":"function"},{"constant":false,"inputs":[{"name":"operator","type":"address"},{"name":"approved","type":"bool"}],"name":"setApprovalForAll","outputs":[],"type":"function"},{"constant":false,"inputs":[{"name":"from","type":"address"},{"name":"to","type":"address"},{"name":"tokenId","type":"uint256"}],"name":"transferFrom","outputs":[],"type":"function"},{"constant":false,"inputs":[{"name":"from","type":"address"},{"name":"to","type":"address"},{"name":"tokenId","type":"uint256"}],"name":"safeTransferFrom","outputs":[],"type":"function"},{"constant":false,"inputs":[{"name":"from","type":"address"},{"name":"to","type":"address"},{"name":"tokenId","type":"uint256"},{"name":"data","type":"bytes"}],"name":"safeTransferFrom","outputs":[],"type":"function"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":true,"name":"tokenId","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"},{"indexed":true,"name":"approved","type":"address"},{"indexed":true,"name":"tokenId","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"},{"indexed":true,"name":"operator","type":"address"},{"indexed":false,"name":"approved","type":"bool"}],"name":"ApprovalForAll","type":"event"}]`

// ERC721 is a binding to an ERC721 NFT contract.
type ERC721 struct {
	contract *contract.Contract
}

// New creates a new ERC721 contract binding.
func New(address common.Address, c *client.PublicClient) (*ERC721, error) {
	cont, err := contract.NewContract(address, []byte(ContractABI), c)
	if err != nil {
		return nil, err
	}
	return &ERC721{contract: cont}, nil
}

// MustNew creates a new ERC721 contract binding, panicking on error.
func MustNew(address common.Address, c *client.PublicClient) *ERC721 {
	cont, err := New(address, c)
	if err != nil {
		panic(err)
	}
	return cont
}

// Address returns the contract address.
func (e *ERC721) Address() common.Address {
	return e.contract.Address()
}

// Contract returns the underlying contract instance.
func (e *ERC721) Contract() *contract.Contract {
	return e.contract
}

// ---- Read Methods (Public Actions) ----

// Name returns the token name.
func (e *ERC721) Name(ctx context.Context) (string, error) {
	return e.contract.ReadString(ctx, "name")
}

// Symbol returns the token symbol.
func (e *ERC721) Symbol(ctx context.Context) (string, error) {
	return e.contract.ReadString(ctx, "symbol")
}

// TokenURI returns the URI for a token.
func (e *ERC721) TokenURI(ctx context.Context, tokenId *big.Int) (string, error) {
	return e.contract.ReadString(ctx, "tokenURI", tokenId)
}

// BalanceOf returns the number of NFTs owned by an address.
func (e *ERC721) BalanceOf(ctx context.Context, owner common.Address) (*big.Int, error) {
	return e.contract.ReadBigInt(ctx, "balanceOf", owner)
}

// OwnerOf returns the owner of a token.
func (e *ERC721) OwnerOf(ctx context.Context, tokenId *big.Int) (common.Address, error) {
	return e.contract.ReadAddress(ctx, "ownerOf", tokenId)
}

// GetApproved returns the approved address for a token.
func (e *ERC721) GetApproved(ctx context.Context, tokenId *big.Int) (common.Address, error) {
	return e.contract.ReadAddress(ctx, "getApproved", tokenId)
}

// IsApprovedForAll returns if an operator is approved for all tokens of an owner.
func (e *ERC721) IsApprovedForAll(ctx context.Context, owner, operator common.Address) (bool, error) {
	return e.contract.ReadBool(ctx, "isApprovedForAll", owner, operator)
}

// ---- Write Methods (Prepare Transaction for Signing) ----

// PrepareApprove prepares an approve transaction for signing.
func (e *ERC721) PrepareApprove(ctx context.Context, opts contract.WriteOptions, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return e.contract.PrepareTransaction(ctx, opts, "approve", to, tokenId)
}

// PrepareSetApprovalForAll prepares a setApprovalForAll transaction for signing.
func (e *ERC721) PrepareSetApprovalForAll(ctx context.Context, opts contract.WriteOptions, operator common.Address, approved bool) (*types.Transaction, error) {
	return e.contract.PrepareTransaction(ctx, opts, "setApprovalForAll", operator, approved)
}

// PrepareTransferFrom prepares a transferFrom transaction for signing.
func (e *ERC721) PrepareTransferFrom(ctx context.Context, opts contract.WriteOptions, from, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return e.contract.PrepareTransaction(ctx, opts, "transferFrom", from, to, tokenId)
}

// PrepareSafeTransferFrom prepares a safeTransferFrom transaction for signing.
func (e *ERC721) PrepareSafeTransferFrom(ctx context.Context, opts contract.WriteOptions, from, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return e.contract.PrepareTransaction(ctx, opts, "safeTransferFrom", from, to, tokenId)
}

// PrepareSafeTransferFromWithData prepares a safeTransferFrom transaction with data for signing.
func (e *ERC721) PrepareSafeTransferFromWithData(ctx context.Context, opts contract.WriteOptions, from, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return e.contract.PrepareTransaction(ctx, opts, "safeTransferFrom", from, to, tokenId, data)
}

// ---- Gas Estimation ----

// EstimateApprove estimates gas for an approve transaction.
func (e *ERC721) EstimateApprove(ctx context.Context, opts contract.WriteOptions, to common.Address, tokenId *big.Int) (uint64, error) {
	return e.contract.EstimateGas(ctx, opts, "approve", to, tokenId)
}

// EstimateSetApprovalForAll estimates gas for a setApprovalForAll transaction.
func (e *ERC721) EstimateSetApprovalForAll(ctx context.Context, opts contract.WriteOptions, operator common.Address, approved bool) (uint64, error) {
	return e.contract.EstimateGas(ctx, opts, "setApprovalForAll", operator, approved)
}

// EstimateTransferFrom estimates gas for a transferFrom transaction.
func (e *ERC721) EstimateTransferFrom(ctx context.Context, opts contract.WriteOptions, from, to common.Address, tokenId *big.Int) (uint64, error) {
	return e.contract.EstimateGas(ctx, opts, "transferFrom", from, to, tokenId)
}

// ---- Events ----

// TransferEvent represents a Transfer event.
type TransferEvent struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
}

// ApprovalEvent represents an Approval event.
type ApprovalEvent struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
}

// ApprovalForAllEvent represents an ApprovalForAll event.
type ApprovalForAllEvent struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
}

// ParseTransfer parses a Transfer event from a log.
func (e *ERC721) ParseTransfer(log types.Log) (*TransferEvent, error) {
	event, err := e.contract.DecodeEvent("Transfer", log.Topics, log.Data)
	if err != nil {
		return nil, err
	}

	return &TransferEvent{
		From:    event["from"].(common.Address),
		To:      event["to"].(common.Address),
		TokenId: event["tokenId"].(*big.Int),
	}, nil
}

// ParseApproval parses an Approval event from a log.
func (e *ERC721) ParseApproval(log types.Log) (*ApprovalEvent, error) {
	event, err := e.contract.DecodeEvent("Approval", log.Topics, log.Data)
	if err != nil {
		return nil, err
	}

	return &ApprovalEvent{
		Owner:    event["owner"].(common.Address),
		Approved: event["approved"].(common.Address),
		TokenId:  event["tokenId"].(*big.Int),
	}, nil
}

// ParseApprovalForAll parses an ApprovalForAll event from a log.
func (e *ERC721) ParseApprovalForAll(log types.Log) (*ApprovalForAllEvent, error) {
	event, err := e.contract.DecodeEvent("ApprovalForAll", log.Topics, log.Data)
	if err != nil {
		return nil, err
	}

	return &ApprovalForAllEvent{
		Owner:    event["owner"].(common.Address),
		Operator: event["operator"].(common.Address),
		Approved: event["approved"].(bool),
	}, nil
}
