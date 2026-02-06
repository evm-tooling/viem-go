package wallet

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ChefBingbong/viem-go/actions/public"
	viemchain "github.com/ChefBingbong/viem-go/chain"
	"github.com/ChefBingbong/viem-go/utils/encoding"
	"github.com/ChefBingbong/viem-go/utils/formatters"
	"github.com/ChefBingbong/viem-go/utils/transaction"
)

// SignTransactionParameters contains the parameters for the SignTransaction action.
// This mirrors viem's SignTransactionParameters type.
type SignTransactionParameters struct {
	// Account is the account to sign with. If nil, uses the client's account.
	Account Account

	// Chain optionally overrides the client's chain for chain ID validation.
	// Set to nil to skip chain assertion (mirrors viem's `chain = client.chain`).
	Chain *viemchain.Chain

	// SkipChainAssertion when true, skips the chain ID mismatch assertion.
	// This mirrors viem's `chain !== null` check.
	SkipChainAssertion bool

	// Transaction fields
	AccessList           transaction.AccessList            `json:"accessList,omitempty"`
	AuthorizationList    []transaction.SignedAuthorization `json:"authorizationList,omitempty"`
	BlobVersionedHashes  []string                          `json:"blobVersionedHashes,omitempty"`
	Blobs                []string                          `json:"blobs,omitempty"`
	Data                 string                            `json:"data,omitempty"`
	Gas                  *big.Int                          `json:"gas,omitempty"`
	GasPrice             *big.Int                          `json:"gasPrice,omitempty"`
	MaxFeePerBlobGas     *big.Int                          `json:"maxFeePerBlobGas,omitempty"`
	MaxFeePerGas         *big.Int                          `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas *big.Int                          `json:"maxPriorityFeePerGas,omitempty"`
	Nonce                *int                              `json:"nonce,omitempty"`
	To                   string                            `json:"to,omitempty"`
	Type                 formatters.TransactionType        `json:"type,omitempty"`
	Value                *big.Int                          `json:"value,omitempty"`
}

// SignTransactionReturnType is the return type for the SignTransaction action (hex string).
type SignTransactionReturnType = string

// SignTransaction signs a transaction.
//
// - For local accounts (implementing TransactionSignableAccount), signs locally without an RPC call.
// - For JSON-RPC accounts, delegates to the `eth_signTransaction` RPC method.
//
// This is equivalent to viem's `signTransaction` action.
//
// JSON-RPC Method: eth_signTransaction
//
// Example:
//
//	signed, err := wallet.SignTransaction(ctx, client, wallet.SignTransactionParameters{
//	    To:    "0x0000000000000000000000000000000000000000",
//	    Value: big.NewInt(1),
//	})
//
// Example with account hoisting:
//
//	signed, err := wallet.SignTransaction(ctx, client, wallet.SignTransactionParameters{
//	    Account: myAccount,
//	    To:      "0x0000000000000000000000000000000000000000",
//	    Value:   big.NewInt(1),
//	})
func SignTransaction(ctx context.Context, client Client, params SignTransactionParameters) (SignTransactionReturnType, error) {
	// Resolve account: param > client
	account := params.Account
	if account == nil {
		account = client.Account()
	}
	if account == nil {
		return "", &AccountNotFoundError{DocsPath: "/docs/actions/wallet/signTransaction"}
	}

	// Validate the request (mirrors viem's assertRequest)
	if err := transaction.AssertRequest(transaction.AssertRequestParams{
		Account:              account.Address().Hex(),
		To:                   params.To,
		MaxFeePerGas:         params.MaxFeePerGas,
		MaxPriorityFeePerGas: params.MaxPriorityFeePerGas,
	}); err != nil {
		return "", err
	}

	// Get chain ID from the network (mirrors viem's getAction(client, getChainId, 'getChainId')({}))
	chainID, err := public.GetChainID(ctx, client)
	if err != nil {
		return "", fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Assert chain matches if chain is provided (mirrors viem's `if (chain !== null) assertCurrentChain(...)`)
	resolvedChain := params.Chain
	if resolvedChain == nil {
		resolvedChain = client.Chain()
	}
	if !params.SkipChainAssertion && resolvedChain != nil {
		if chainErr := viemchain.AssertCurrentChain(resolvedChain, int64(chainID)); chainErr != nil {
			return "", chainErr
		}
	}

	// If the account can sign transactions locally, use it directly
	if signable, ok := account.(TransactionSignableAccount); ok {
		tx := paramsToTransaction(params, int(chainID))
		return signable.SignTransaction(tx)
	}

	// Otherwise, format the transaction request and send via eth_signTransaction RPC
	txRequest := paramsToTransactionRequest(params)
	rpcRequest := formatters.FormatTransactionRequest(txRequest)

	// Add chainId and from (mirrors viem's override of format result)
	type signTxRpcRequest struct {
		formatters.RpcTransactionRequest
		ChainID string `json:"chainId"`
		From    string `json:"from"`
	}

	rpcReq := signTxRpcRequest{
		RpcTransactionRequest: rpcRequest,
		ChainID:               encoding.NumberToHex(new(big.Int).SetUint64(chainID)),
		From:                  account.Address().Hex(),
	}

	resp, err := client.Request(ctx, "eth_signTransaction", rpcReq)
	if err != nil {
		return "", fmt.Errorf("eth_signTransaction failed: %w", err)
	}

	var hexResult string
	if unmarshalErr := json.Unmarshal(resp.Result, &hexResult); unmarshalErr != nil {
		return "", fmt.Errorf("failed to unmarshal signed transaction: %w", unmarshalErr)
	}

	return hexResult, nil
}

// PreparedToSignParams converts a PrepareTransactionRequestReturnType (from PrepareTransactionRequest)
// into SignTransactionParameters. This bridges the two types so you can do:
//
//	prepared, _ := wallet.PrepareTransactionRequest(ctx, client, ...)
//	signed, _ := wallet.SignTransaction(ctx, client, wallet.PreparedToSignParams(prepared))
func PreparedToSignParams(p *PrepareTransactionRequestParameters) SignTransactionParameters {
	params := SignTransactionParameters{
		To:                   p.To,
		Value:                p.Value,
		Data:                 p.Data,
		Gas:                  p.Gas,
		GasPrice:             p.GasPrice,
		MaxFeePerGas:         p.MaxFeePerGas,
		MaxPriorityFeePerGas: p.MaxPriorityFeePerGas,
		MaxFeePerBlobGas:     p.MaxFeePerBlobGas,
		BlobVersionedHashes:  p.BlobVersionedHashes,
		Blobs:                p.Blobs,
		Nonce:                p.Nonce,
		Type:                 p.Type,
	}

	if len(p.AccessList) > 0 {
		al := make(transaction.AccessList, len(p.AccessList))
		for i, item := range p.AccessList {
			al[i] = transaction.AccessListItem{
				Address:     item.Address,
				StorageKeys: item.StorageKeys,
			}
		}
		params.AccessList = al
	}

	if len(p.AuthorizationList) > 0 {
		params.AuthorizationList = p.AuthorizationList
	}

	return params
}

// paramsToTransaction converts SignTransactionParameters to a transaction.Transaction for local signing.
func paramsToTransaction(params SignTransactionParameters, chainID int) *transaction.Transaction {
	tx := &transaction.Transaction{
		ChainId:              chainID,
		To:                   params.To,
		Value:                params.Value,
		Data:                 params.Data,
		Gas:                  params.Gas,
		GasPrice:             params.GasPrice,
		MaxFeePerGas:         params.MaxFeePerGas,
		MaxPriorityFeePerGas: params.MaxPriorityFeePerGas,
		MaxFeePerBlobGas:     params.MaxFeePerBlobGas,
		BlobVersionedHashes:  params.BlobVersionedHashes,
		Blobs:                params.Blobs,
		AuthorizationList:    params.AuthorizationList,
	}

	if params.Nonce != nil {
		tx.Nonce = *params.Nonce
	}

	if len(params.AccessList) > 0 {
		tx.AccessList = params.AccessList
	}

	if params.Type != "" {
		tx.Type = transaction.TransactionType(params.Type)
	}

	return tx
}

// paramsToTransactionRequest converts SignTransactionParameters to a formatters.TransactionRequest for RPC formatting.
func paramsToTransactionRequest(params SignTransactionParameters) formatters.TransactionRequest {
	req := formatters.TransactionRequest{
		To:                   params.To,
		Data:                 params.Data,
		Value:                params.Value,
		Gas:                  params.Gas,
		GasPrice:             params.GasPrice,
		MaxFeePerGas:         params.MaxFeePerGas,
		MaxPriorityFeePerGas: params.MaxPriorityFeePerGas,
		MaxFeePerBlobGas:     params.MaxFeePerBlobGas,
		BlobVersionedHashes:  params.BlobVersionedHashes,
		Nonce:                params.Nonce,
		Type:                 params.Type,
	}

	if len(params.AccessList) > 0 {
		fmtAccessList := make(formatters.AccessList, len(params.AccessList))
		for i, item := range params.AccessList {
			fmtAccessList[i] = formatters.AccessListItem{
				Address:     item.Address,
				StorageKeys: item.StorageKeys,
			}
		}
		req.AccessList = fmtAccessList
	}

	if len(params.AuthorizationList) > 0 {
		authList := make([]any, len(params.AuthorizationList))
		for i, auth := range params.AuthorizationList {
			authList[i] = map[string]any{
				"address": auth.Address,
				"chainId": auth.ChainId,
				"nonce":   auth.Nonce,
				"r":       auth.R,
				"s":       auth.S,
				"yParity": auth.YParity,
			}
		}
		req.AuthorizationList = authList
	}

	if len(params.Blobs) > 0 {
		blobs := make([]any, len(params.Blobs))
		for i, b := range params.Blobs {
			blobs[i] = b
		}
		req.Blobs = blobs
	}

	return req
}
