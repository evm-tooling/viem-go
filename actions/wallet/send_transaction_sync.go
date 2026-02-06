package wallet

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/actions/public"
	viemchain "github.com/ChefBingbong/viem-go/chain"
	"github.com/ChefBingbong/viem-go/utils/data"
	"github.com/ChefBingbong/viem-go/utils/encoding"
	"github.com/ChefBingbong/viem-go/utils/formatters"
	"github.com/ChefBingbong/viem-go/utils/transaction"
)

// SendTransactionSyncParameters contains the parameters for the SendTransactionSync action.
// This mirrors viem's SendTransactionSyncParameters type.
type SendTransactionSyncParameters struct {
	// Account is the account to send from. If nil, uses the client's account.
	Account Account

	// Chain optionally overrides the client's chain for chain ID validation.
	Chain *viemchain.Chain

	// AssertChainID when true, asserts the chain ID matches. Default: true.
	AssertChainID *bool

	// DataSuffix is data to append to the end of the calldata.
	// Takes precedence over client.DataSuffix().
	DataSuffix string

	// PollingInterval is the polling interval to poll for the transaction receipt.
	// Defaults to client.PollingInterval().
	PollingInterval time.Duration

	// ThrowOnReceiptRevert when true, throws an error if the transaction was detected as reverted.
	// Default: true.
	ThrowOnReceiptRevert *bool

	// Timeout is the timeout to wait for a response.
	// Default: max(chain.blockTime * 3, 5000)ms.
	Timeout *time.Duration

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

// SendTransactionSyncReturnType is the return type for SendTransactionSync.
// For JSON-RPC accounts it returns a transaction receipt from waitForTransactionReceipt.
// For local accounts it returns a transaction receipt from sendRawTransactionSync.
type SendTransactionSyncReturnType = *formatters.TransactionReceipt

// SendTransactionSync creates, signs, and sends a new transaction to the network synchronously.
// Returns the transaction receipt after the transaction is included in a block.
//
//   - For JSON-RPC accounts, sends via eth_sendTransaction then waits for receipt via
//     WaitForTransactionReceipt.
//   - For local accounts, prepares, signs locally, and sends via sendRawTransactionSync (EIP-7966).
//
// This is equivalent to viem's `sendTransactionSync` action.
//
// Example:
//
//	receipt, err := wallet.SendTransactionSync(ctx, client, wallet.SendTransactionSyncParameters{
//	    To:    "0x70997970c51812dc3a010c7d01b50e0d17dc79c8",
//	    Value: big.NewInt(1000000000000000000),
//	})
//
// Example with account hoisting:
//
//	receipt, err := wallet.SendTransactionSync(ctx, client, wallet.SendTransactionSyncParameters{
//	    Account: myAccount,
//	    To:      "0x70997970c51812dc3a010c7d01b50e0d17dc79c8",
//	    Value:   big.NewInt(1000000000000000000),
//	})
func SendTransactionSync(ctx context.Context, client Client, params SendTransactionSyncParameters) (SendTransactionSyncReturnType, error) {
	// Resolve account: param > client
	account := params.Account
	if account == nil {
		account = client.Account()
	}
	if account == nil {
		return nil, &AccountNotFoundError{DocsPath: "/docs/actions/wallet/sendTransactionSync"}
	}

	// Resolve timeout (mirrors viem's: parameters.timeout ?? Math.max((chain?.blockTime ?? 0) * 3, 5_000))
	timeout := resolveTimeout(params.Timeout, params.Chain, client.Chain())

	// Resolve data suffix: param > client
	dataSuffix := params.DataSuffix
	if dataSuffix == "" && len(client.DataSuffix()) > 0 {
		dataSuffix = encoding.BytesToHex(client.DataSuffix())
	}

	// Apply data suffix if data is present
	txData := params.Data
	if txData != "" && dataSuffix != "" {
		txData = data.ConcatHex(txData, dataSuffix)
	}

	// Validate the request
	if err := transaction.AssertRequest(transaction.AssertRequestParams{
		Account:              account.Address().Hex(),
		To:                   params.To,
		MaxFeePerGas:         params.MaxFeePerGas,
		MaxPriorityFeePerGas: params.MaxPriorityFeePerGas,
	}); err != nil {
		return nil, err
	}

	// Check if this is a local account that can sign transactions
	signable, isLocal := account.(TransactionSignableAccount)

	if !isLocal {
		// ---- JSON-RPC Account path ----
		return sendTransactionSyncViaRPC(ctx, client, account, params, txData, timeout)
	}

	// ---- Local Account path ----
	return sendTransactionSyncViaLocalSign(ctx, client, account, signable, params, txData, timeout)
}

// sendTransactionSyncViaRPC handles the JSON-RPC account path:
// eth_sendTransaction -> waitForTransactionReceipt.
func sendTransactionSyncViaRPC(
	ctx context.Context,
	client Client,
	account Account,
	params SendTransactionSyncParameters,
	txData string,
	timeout time.Duration,
) (*formatters.TransactionReceipt, error) {
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

	// Get chain ID and assert chain matches
	var chainID *uint64
	if ch != nil {
		cid, err := public.GetChainID(ctx, client)
		if err != nil {
			return nil, fmt.Errorf("failed to get chain ID: %w", err)
		}
		chainID = &cid
		if assertChainID {
			if chainErr := viemchain.AssertCurrentChain(ch, int64(cid)); chainErr != nil {
				return nil, chainErr
			}
		}
	}

	// Format the transaction request
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
		return nil, fmt.Errorf("eth_sendTransaction failed: %w", err)
	}

	var hash string
	if unmarshalErr := json.Unmarshal(resp.Result, &hash); unmarshalErr != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction hash: %w", unmarshalErr)
	}

	// Wait for the transaction receipt
	// This mirrors viem's: waitForTransactionReceipt({ checkReplacement: false, hash, pollingInterval, timeout })
	pollingInterval := params.PollingInterval
	if pollingInterval == 0 {
		pollingInterval = client.PollingInterval()
	}

	receiptCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	receipt, waitErr := public.WaitForTransactionReceipt(receiptCtx, client, public.WaitForTransactionReceiptParameters{
		Hash:             common.HexToHash(hash),
		CheckReplacement: boolPtr(false),
		PollingInterval:  pollingInterval,
		Timeout:          timeout,
	})
	if waitErr != nil {
		return nil, fmt.Errorf("failed to wait for transaction receipt: %w", waitErr)
	}

	// Check if receipt is reverted
	throwOnRevert := true
	if params.ThrowOnReceiptRevert != nil {
		throwOnRevert = *params.ThrowOnReceiptRevert
	}
	if throwOnRevert && receipt != nil && receipt.Status == 0 {
		fmtReceipt := rpcReceiptToFormattersReceipt(receipt)
		return nil, &TransactionReceiptRevertedError{Receipt: fmtReceipt}
	}

	return rpcReceiptToFormattersReceipt(receipt), nil
}

// sendTransactionSyncViaLocalSign handles the local account path:
// prepare + sign + sendRawTransactionSync.
func sendTransactionSyncViaLocalSign(
	ctx context.Context,
	client Client,
	account Account,
	signable TransactionSignableAccount,
	params SendTransactionSyncParameters,
	txData string,
	timeout time.Duration,
) (*formatters.TransactionReceipt, error) {
	// Resolve chain
	ch := params.Chain
	if ch == nil {
		ch = client.Chain()
	}

	// Prepare the transaction request
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
		return nil, fmt.Errorf("failed to prepare transaction request: %w", err)
	}

	// Convert and sign
	tx := preparedParamsToTransaction(prepared)
	serializedTx, signErr := signable.SignTransaction(tx)
	if signErr != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", signErr)
	}

	// Send via sendRawTransactionSync
	// This mirrors viem's: sendRawTransactionSync({ serializedTransaction, throwOnReceiptRevert, timeout })
	timeoutMs := timeout.Milliseconds()
	return SendRawTransactionSync(ctx, client, SendRawTransactionSyncParameters{
		SerializedTransaction: serializedTx,
		ThrowOnReceiptRevert:  params.ThrowOnReceiptRevert,
		Timeout:               &timeoutMs,
	})
}

// resolveTimeout resolves the timeout for sendTransactionSync.
// Mirrors viem's: parameters.timeout ?? Math.max((chain?.blockTime ?? 0) * 3, 5_000)
func resolveTimeout(paramTimeout *time.Duration, paramChain, clientChain *viemchain.Chain) time.Duration {
	if paramTimeout != nil {
		return *paramTimeout
	}

	ch := paramChain
	if ch == nil {
		ch = clientChain
	}

	var blockTimeMs int64
	if ch != nil && ch.BlockTime != nil {
		blockTimeMs = *ch.BlockTime
	}

	timeoutMs := blockTimeMs * 3
	if timeoutMs < 5000 {
		timeoutMs = 5000
	}

	return time.Duration(timeoutMs) * time.Millisecond
}

// boolPtr returns a pointer to a bool value.
func boolPtr(b bool) *bool {
	return &b
}

// rpcReceiptToFormattersReceipt converts a types.Receipt to a formatters.TransactionReceipt.
// This is a best-effort conversion for status checking.
func rpcReceiptToFormattersReceipt(receipt interface{}) *formatters.TransactionReceipt {
	if receipt == nil {
		return nil
	}

	// Marshal and unmarshal to convert between types
	jsonBytes, err := json.Marshal(receipt)
	if err != nil {
		return nil
	}

	var rpcReceipt formatters.RpcTransactionReceipt
	if err := json.Unmarshal(jsonBytes, &rpcReceipt); err != nil {
		return nil
	}

	formatted := formatters.FormatTransactionReceipt(rpcReceipt)
	return &formatted
}
