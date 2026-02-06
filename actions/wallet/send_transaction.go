package wallet

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ChefBingbong/viem-go/actions/public"
	viemchain "github.com/ChefBingbong/viem-go/chain"
	"github.com/ChefBingbong/viem-go/utils/data"
	"github.com/ChefBingbong/viem-go/utils/encoding"
	"github.com/ChefBingbong/viem-go/utils/formatters"
	"github.com/ChefBingbong/viem-go/utils/transaction"
)

// SendTransactionParameters contains the parameters for the SendTransaction action.
// This mirrors viem's SendTransactionParameters type.
type SendTransactionParameters struct {
	// Account is the account to send from. If nil, uses the client's account.
	Account Account

	// Chain optionally overrides the client's chain for chain ID validation.
	Chain *viemchain.Chain

	// AssertChainID when true, asserts the chain ID matches. Default: true.
	AssertChainID *bool

	// DataSuffix is data to append to the end of the calldata.
	// Takes precedence over client.DataSuffix().
	DataSuffix string

	// Transaction fields
	AccessList           []formatters.AccessListItem       `json:"accessList,omitempty"`
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

// SendTransactionReturnType is the return type for the SendTransaction action.
// It is the transaction hash as a hex string.
type SendTransactionReturnType = string

// SendTransaction creates, signs, and sends a new transaction to the network.
//
//   - For JSON-RPC accounts (or when no local signer is available), sends via `eth_sendTransaction`.
//   - For local accounts (implementing TransactionSignableAccount), prepares, signs locally,
//     and sends via `eth_sendRawTransaction` (sendRawTransaction).
//
// This is equivalent to viem's `sendTransaction` action.
//
// JSON-RPC Methods:
//   - JSON-RPC Accounts: eth_sendTransaction
//   - Local Accounts: eth_sendRawTransaction
//
// Example:
//
//	hash, err := wallet.SendTransaction(ctx, client, wallet.SendTransactionParameters{
//	    To:    "0x70997970c51812dc3a010c7d01b50e0d17dc79c8",
//	    Value: big.NewInt(1000000000000000000),
//	})
//
// Example with account hoisting:
//
//	hash, err := wallet.SendTransaction(ctx, client, wallet.SendTransactionParameters{
//	    Account: myAccount,
//	    To:      "0x70997970c51812dc3a010c7d01b50e0d17dc79c8",
//	    Value:   big.NewInt(1000000000000000000),
//	})
func SendTransaction(ctx context.Context, client Client, params SendTransactionParameters) (SendTransactionReturnType, error) {
	// Resolve account: param > client
	account := params.Account
	if account == nil {
		account = client.Account()
	}
	if account == nil {
		return "", &AccountNotFoundError{DocsPath: "/docs/actions/wallet/sendTransaction"}
	}

	// Resolve data suffix: param > client
	dataSuffix := params.DataSuffix
	if dataSuffix == "" && len(client.DataSuffix()) > 0 {
		dataSuffix = encoding.BytesToHex(client.DataSuffix())
	}

	// Apply data suffix if data is present (mirrors viem's `data ? concat([data, dataSuffix ?? '0x']) : data`)
	txData := params.Data
	if txData != "" && dataSuffix != "" {
		txData = data.ConcatHex(txData, dataSuffix)
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

	// Check if this is a local account that can sign transactions
	signable, isLocal := account.(TransactionSignableAccount)

	if !isLocal {
		// ---- JSON-RPC Account path (eth_sendTransaction) ----
		return sendTransactionViaRPC(ctx, client, account, params, txData)
	}

	// ---- Local Account path (prepare + sign + sendRawTransaction) ----
	return sendTransactionViaLocalSign(ctx, client, account, signable, params, txData)
}

// sendTransactionViaRPC handles the JSON-RPC account path using eth_sendTransaction.
// This mirrors viem's `account?.type === 'json-rpc'` branch.
func sendTransactionViaRPC(
	ctx context.Context,
	client Client,
	account Account,
	params SendTransactionParameters,
	txData string,
) (string, error) {
	// Resolve chain
	ch := params.Chain
	if ch == nil {
		ch = client.Chain()
	}

	// Determine if we should assert chain ID (default: true)
	assertChainID := true
	if params.AssertChainID != nil {
		assertChainID = *params.AssertChainID
	}

	// Get chain ID and assert chain matches (mirrors viem's chain assertion)
	var chainID *uint64
	if ch != nil {
		cid, err := public.GetChainID(ctx, client)
		if err != nil {
			return "", fmt.Errorf("failed to get chain ID: %w", err)
		}
		chainID = &cid
		if assertChainID {
			if chainErr := viemchain.AssertCurrentChain(ch, int64(cid)); chainErr != nil {
				return "", chainErr
			}
		}
	}

	// Format the transaction request (mirrors viem's formatTransactionRequest)
	txRequest := formatters.TransactionRequest{
		Data:                 txData,
		From:                 account.Address().Hex(),
		To:                   params.To,
		Value:                params.Value,
		Gas:                  params.Gas,
		GasPrice:             params.GasPrice,
		MaxFeePerGas:         params.MaxFeePerGas,
		MaxPriorityFeePerGas: params.MaxPriorityFeePerGas,
		MaxFeePerBlobGas:     params.MaxFeePerBlobGas,
		Nonce:                params.Nonce,
		Type:                 params.Type,
	}

	if len(params.AccessList) > 0 {
		txRequest.AccessList = formatters.AccessList(params.AccessList)
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
		txRequest.AuthorizationList = authList
	}
	if len(params.Blobs) > 0 {
		blobs := make([]any, len(params.Blobs))
		for i, b := range params.Blobs {
			blobs[i] = b
		}
		txRequest.Blobs = blobs
	}

	rpcRequest := formatters.FormatTransactionRequest(txRequest)

	// Add chainId to the formatted request if available
	type rpcSendTxRequest struct {
		formatters.RpcTransactionRequest
		ChainID string `json:"chainId,omitempty"`
	}

	rpcReq := rpcSendTxRequest{
		RpcTransactionRequest: rpcRequest,
	}
	if chainID != nil {
		rpcReq.ChainID = encoding.NumberToHex(new(big.Int).SetUint64(*chainID))
	}

	// Send via eth_sendTransaction
	resp, err := client.Request(ctx, "eth_sendTransaction", rpcReq)
	if err != nil {
		return "", fmt.Errorf("eth_sendTransaction failed: %w", err)
	}

	var hash string
	if unmarshalErr := json.Unmarshal(resp.Result, &hash); unmarshalErr != nil {
		return "", fmt.Errorf("failed to unmarshal transaction hash: %w", unmarshalErr)
	}

	return hash, nil
}

// sendTransactionViaLocalSign handles the local account path: prepare + sign + sendRawTransaction.
// This mirrors viem's `account?.type === 'local'` branch.
func sendTransactionViaLocalSign(
	ctx context.Context,
	client Client,
	account Account,
	signable TransactionSignableAccount,
	params SendTransactionParameters,
	txData string,
) (string, error) {
	// Resolve chain
	ch := params.Chain
	if ch == nil {
		ch = client.Chain()
	}

	// Prepare the transaction request (fills in nonce, gas, fees, type, chainId)
	// This mirrors viem's: prepareTransactionRequest({ account, ...params, parameters: [...defaultParameters, 'sidecars'] })
	prepareParams := PrepareTransactionRequestParameters{
		Account:              account,
		Chain:                ch,
		Parameters:           append(append([]string{}, DefaultParameters...), "sidecars"),
		AccessList:           params.AccessList,
		AuthorizationList:    params.AuthorizationList,
		BlobVersionedHashes:  params.BlobVersionedHashes,
		Blobs:                params.Blobs,
		Data:                 txData,
		Gas:                  params.Gas,
		GasPrice:             params.GasPrice,
		MaxFeePerBlobGas:     params.MaxFeePerBlobGas,
		MaxFeePerGas:         params.MaxFeePerGas,
		MaxPriorityFeePerGas: params.MaxPriorityFeePerGas,
		Nonce:                params.Nonce,
		To:                   params.To,
		Type:                 params.Type,
		Value:                params.Value,
	}

	prepared, err := PrepareTransactionRequest(ctx, client, prepareParams)
	if err != nil {
		return "", fmt.Errorf("failed to prepare transaction request: %w", err)
	}

	// Convert prepared params to a Transaction for local signing
	tx := preparedParamsToTransaction(prepared)

	// Sign the transaction locally
	// This mirrors viem's: account.signTransaction(request, { serializer })
	serializedTx, signErr := signable.SignTransaction(tx)
	if signErr != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", signErr)
	}

	// Send the raw signed transaction
	// This mirrors viem's: sendRawTransaction({ serializedTransaction })
	return SendRawTransaction(ctx, client, SendRawTransactionParameters{
		SerializedTransaction: serializedTx,
	})
}

// preparedParamsToTransaction converts PrepareTransactionRequestParameters to a transaction.Transaction.
func preparedParamsToTransaction(params *PrepareTransactionRequestParameters) *transaction.Transaction {
	tx := &transaction.Transaction{
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
	}

	if params.ChainID != nil {
		tx.ChainId = int(*params.ChainID)
	}

	if params.Nonce != nil {
		tx.Nonce = *params.Nonce
	}

	if len(params.AccessList) > 0 {
		al := make(transaction.AccessList, len(params.AccessList))
		for i, item := range params.AccessList {
			al[i] = transaction.AccessListItem{
				Address:     item.Address,
				StorageKeys: item.StorageKeys,
			}
		}
		tx.AccessList = al
	}

	if len(params.AuthorizationList) > 0 {
		tx.AuthorizationList = params.AuthorizationList
	}

	if params.Type != "" {
		tx.Type = transaction.TransactionType(params.Type)
	}

	return tx
}
