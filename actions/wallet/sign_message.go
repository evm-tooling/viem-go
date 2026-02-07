package wallet

import (
	"context"
	"fmt"

	json "github.com/goccy/go-json"

	"github.com/ChefBingbong/viem-go/utils/encoding"
	"github.com/ChefBingbong/viem-go/utils/signature"
)

// SignMessageParameters contains the parameters for the SignMessage action.
// This mirrors viem's SignMessageParameters type.
type SignMessageParameters struct {
	// Account is the account to sign with. If nil, uses the client's account.
	Account Account

	// Message is the message to sign.
	Message signature.SignableMessage
}

// SignMessageReturnType is the return type for the SignMessage action (hex string).
type SignMessageReturnType = string

// SignMessage calculates an Ethereum-specific signature in EIP-191 format:
// keccak256("\x19Ethereum Signed Message:\n" + len(message) + message))
//
// - For local accounts (implementing SignableAccount), signs locally without an RPC call.
// - For JSON-RPC accounts, delegates to the `personal_sign` RPC method.
//
// This is equivalent to viem's `signMessage` action.
//
// JSON-RPC Method: personal_sign
//
// Example:
//
//	sig, err := wallet.SignMessage(ctx, client, wallet.SignMessageParameters{
//	    Message: signature.NewSignableMessage("hello world"),
//	})
//
// Example with account override:
//
//	sig, err := wallet.SignMessage(ctx, client, wallet.SignMessageParameters{
//	    Account: myAccount,
//	    Message: signature.NewSignableMessage("hello world"),
//	})
func SignMessage(ctx context.Context, client Client, params SignMessageParameters) (SignMessageReturnType, error) {
	// Resolve account: param > client
	account := params.Account
	if account == nil {
		account = client.Account()
	}
	if account == nil {
		return "", &AccountNotFoundError{DocsPath: "/docs/actions/wallet/signMessage"}
	}

	// If the account can sign locally, use it directly (mirrors viem's account.signMessage check)
	if signable, ok := account.(SignableAccount); ok {
		return signable.SignMessage(params.Message)
	}

	// Otherwise, encode the message and send via personal_sign RPC
	message := encodeSignableMessage(params.Message)

	var hexResult string
	resp, err := client.Request(ctx, "personal_sign", message, account.Address().Hex())
	if err != nil {
		return "", fmt.Errorf("personal_sign failed: %w", err)
	}

	if unmarshalErr := json.Unmarshal(resp.Result, &hexResult); unmarshalErr != nil {
		return "", fmt.Errorf("failed to unmarshal signature: %w", unmarshalErr)
	}

	return hexResult, nil
}

// encodeSignableMessage converts a SignableMessage to a hex string for RPC.
// This mirrors viem's message encoding logic:
//   - string message -> stringToHex
//   - raw []byte -> toHex
//   - raw hex string -> pass through
func encodeSignableMessage(msg signature.SignableMessage) string {
	if msg.Raw != nil {
		switch v := msg.Raw.(type) {
		case []byte:
			return encoding.BytesToHex(v)
		case string:
			// Already a hex string
			return v
		}
	}
	// String message -> hex encode
	return encoding.StringToHex(msg.Message)
}
