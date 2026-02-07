package wallet

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	json "github.com/goccy/go-json"

	"github.com/ChefBingbong/viem-go/actions/public"
	viemchain "github.com/ChefBingbong/viem-go/chain"
	"github.com/ChefBingbong/viem-go/utils"
	"github.com/ChefBingbong/viem-go/utils/authorization"
	"github.com/ChefBingbong/viem-go/utils/data"
	"github.com/ChefBingbong/viem-go/utils/encoding"
	"github.com/ChefBingbong/viem-go/utils/formatters"
	"github.com/ChefBingbong/viem-go/utils/signature"
	"github.com/ChefBingbong/viem-go/utils/transaction"
)

// supportsWalletNamespace is an LRU cache keyed by client.UID() that tracks
// whether a given transport supports `wallet_sendTransaction` over the standard
// `eth_sendTransaction`. This mirrors viem's `supportsWalletNamespace` LruMap.
var supportsWalletNamespace = utils.NewLruMap[bool](128)

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

	// Resolve `to` — infer from authorizationList if not provided.
	// Mirrors viem's: if no `to` and authorizationList present, recover address from first auth.
	to := params.To
	if to == "" && len(params.AuthorizationList) > 0 {
		recovered, recoverErr := recoverAuthorizationAddr(params.AuthorizationList[0])
		if recoverErr != nil {
			return "", fmt.Errorf("`to` is required. Could not infer from `authorizationList`: %w", recoverErr)
		}
		to = recovered
	}

	// Check if this is a local account that can sign transactions
	signable, isLocal := account.(TransactionSignableAccount)

	if !isLocal {
		// ---- JSON-RPC Account path (eth_sendTransaction) ----
		return sendTransactionViaRPC(ctx, client, account, params, txData, to)
	}

	// ---- Local Account path (prepare + sign + sendRawTransaction) ----
	return sendTransactionViaLocalSign(ctx, client, account, signable, params, txData, to)
}

// sendTransactionViaRPC handles the JSON-RPC account path using eth_sendTransaction.
// This mirrors viem's `account?.type === 'json-rpc'` branch, including the
// wallet_sendTransaction namespace fallback with LRU caching.
func sendTransactionViaRPC(
	ctx context.Context,
	client Client,
	account Account,
	params SendTransactionParameters,
	txData string,
	to string,
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
		To:                   to,
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

	// Send with wallet_sendTransaction namespace fallback.
	// Mirrors viem's: try eth_sendTransaction first, on certain RPC errors
	// retry with wallet_sendTransaction, cache the result per client.UID().
	return sendWithNamespaceFallback(ctx, client, rpcReq)
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
	to string,
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
		To:                   to,
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

// sendWithNamespaceFallback sends a transaction via eth_sendTransaction, falling back
// to wallet_sendTransaction if the transport doesn't support eth_sendTransaction.
// Caches the result per client.UID() in the supportsWalletNamespace LRU.
//
// This mirrors viem's namespace fallback logic in sendTransaction.ts lines 256-308.
func sendWithNamespaceFallback(ctx context.Context, client Client, rpcReq any) (string, error) {
	uid := client.UID()

	// Check LRU cache for known wallet namespace support
	isWalletSupported, hasCachedValue := supportsWalletNamespace.Get(uid)

	// Determine the primary method to try
	method := "eth_sendTransaction"
	if hasCachedValue && isWalletSupported {
		method = "wallet_sendTransaction"
	}

	// Try the primary method
	resp, err := client.Request(ctx, method, rpcReq)
	if err == nil {
		var hash string
		if unmarshalErr := json.Unmarshal(resp.Result, &hash); unmarshalErr != nil {
			return "", fmt.Errorf("failed to unmarshal transaction hash: %w", unmarshalErr)
		}
		return hash, nil
	}

	// If we already know wallet namespace is NOT supported, don't retry
	if hasCachedValue && !isWalletSupported {
		return "", fmt.Errorf("%s failed: %w", method, err)
	}

	// Check if the error is a type that warrants trying the wallet_ namespace
	// Mirrors viem's: InvalidInputRpcError, InvalidParamsRpcError, MethodNotFoundRpcError, MethodNotSupportedRpcError
	originalErr := err
	if !isNamespaceFallbackError(err) {
		return "", fmt.Errorf("%s failed: %w", method, err)
	}

	// Attempt wallet_sendTransaction as fallback
	walletResp, walletErr := client.Request(ctx, "wallet_sendTransaction", rpcReq)
	if walletErr == nil {
		// wallet_sendTransaction succeeded — cache for future calls
		supportsWalletNamespace.Set(uid, true)

		var hash string
		if unmarshalErr := json.Unmarshal(walletResp.Result, &hash); unmarshalErr != nil {
			return "", fmt.Errorf("failed to unmarshal transaction hash: %w", unmarshalErr)
		}
		return hash, nil
	}

	// wallet_sendTransaction also failed — check if it's MethodNotFound/NotSupported/NotImplemented
	walletErrStr := strings.ToLower(walletErr.Error())
	if strings.Contains(walletErrStr, "method not found") ||
		strings.Contains(walletErrStr, "method not supported") ||
		strings.Contains(walletErrStr, "not implemented") {
		// The wallet namespace is confirmed not supported; cache and throw the original error
		supportsWalletNamespace.Set(uid, false)
		return "", fmt.Errorf("eth_sendTransaction failed: %w", originalErr)
	}

	// Different wallet_sendTransaction error — throw it
	return "", fmt.Errorf("wallet_sendTransaction failed: %w", walletErr)
}

// isNamespaceFallbackError checks if an error should trigger a wallet_ namespace fallback.
// Mirrors viem's check for InvalidInputRpcError, InvalidParamsRpcError,
// MethodNotFoundRpcError, MethodNotSupportedRpcError.
// Also covers "not implemented" which some RPC providers (e.g. Polygon public RPC) return.
func isNamespaceFallbackError(err error) bool {
	if err == nil {
		return false
	}
	lower := strings.ToLower(err.Error())
	return strings.Contains(lower, "invalid input") ||
		strings.Contains(lower, "invalid params") ||
		strings.Contains(lower, "method not found") ||
		strings.Contains(lower, "method not supported") ||
		strings.Contains(lower, "not implemented")
}

// recoverAuthorizationAddr recovers the signer address from a signed EIP-7702 authorization.
// Mirrors viem's recoverAuthorizationAddress utility.
func recoverAuthorizationAddr(auth transaction.SignedAuthorization) (string, error) {
	// Hash the authorization
	authHash, err := authorization.HashAuthorizationHex(authorization.AuthorizationRequest{
		Address: auth.Address,
		ChainId: auth.ChainId,
		Nonce:   auth.Nonce,
	})
	if err != nil {
		return "", fmt.Errorf("failed to hash authorization: %w", err)
	}

	// Build compact signature hex: r + s + yParity
	yParityHex := "0x00"
	if auth.YParity == 1 {
		yParityHex = "0x01"
	}
	sig := data.ConcatHex(auth.R, auth.S, yParityHex)

	// Recover the address
	addr, recoverErr := signature.RecoverAddress(authHash, sig)
	if recoverErr != nil {
		return "", fmt.Errorf("failed to recover authorization address: %w", recoverErr)
	}

	return addr, nil
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
