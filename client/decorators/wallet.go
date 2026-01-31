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
		"sendRawTransaction": c.SendRawTransaction,
		"sendTransaction":    c.SendTransaction,
		"signMessage":        c.SignMessage,
		"signTypedData":      c.SignTypedData,
		"getAccounts":        c.GetAccounts,
		"requestAccounts":    c.RequestAccounts,
		"switchChain":        c.SwitchChain,
		"addChain":           c.AddChain,
	}
}
