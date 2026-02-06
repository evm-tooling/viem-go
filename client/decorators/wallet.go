package decorators

import (
	"github.com/ChefBingbong/viem-go/client"
)

// WalletActions returns wallet action methods as a map.
// This mirrors viem's walletActions decorator for extension purposes.
//
// Example:
//
//	client := client.CreateWalletClient(config)
//	actions := decorators.WalletActions(client)
func WalletActions(c *client.WalletClient) map[string]any {
	return map[string]any{
		// Signing
		"signMessage":          c.SignMessage,
		"signTypedData":        c.SignTypedData,
		"signTransaction":      c.SignTransaction,
		"signAuthorization":    c.SignAuthorization,
		"prepareAuthorization": c.PrepareAuthorization,

		// Transactions
		"sendTransaction":           c.SendTransaction,
		"sendTransactionSync":       c.SendTransactionSync,
		"sendRawTransaction":        c.SendRawTransaction,
		"sendRawTransactionSync":    c.SendRawTransactionSync,
		"prepareTransactionRequest": c.PrepareTransactionRequest,

		// Contracts
		"writeContract":     c.WriteContract,
		"writeContractSync": c.WriteContractSync,
		"deployContract":    c.DeployContract,

		// Account Management
		"getAddresses":     c.GetAddresses,
		"requestAddresses": c.RequestAddresses,
		"addChain":         c.AddChain,
		"switchChain":      c.SwitchChain,

		// Permissions & Assets
		"getPermissions":     c.GetPermissions,
		"requestPermissions": c.RequestPermissions,
		"watchAsset":         c.WatchAsset,

		// EIP-5792 Batch Calls
		"getCapabilities":    c.GetCapabilities,
		"sendCalls":          c.SendCalls,
		"sendCallsSync":      c.SendCallsSync,
		"getCallsStatus":     c.GetCallsStatus,
		"waitForCallsStatus": c.WaitForCallsStatus,
		"showCallsStatus":    c.ShowCallsStatus,
	}
}
