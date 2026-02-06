package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/accounts"
	"github.com/ChefBingbong/viem-go/actions/public"
	"github.com/ChefBingbong/viem-go/chain/definitions"
	"github.com/ChefBingbong/viem-go/client"
	"github.com/ChefBingbong/viem-go/client/transport"
	"github.com/ChefBingbong/viem-go/utils/unit"
)

// ---------------------------------------------------------------------------
// Constants
// ---------------------------------------------------------------------------

const polygonRPC = "https://rough-purple-market.matic.quiknode.pro/c1a568726a34041d3c5d58603f5981951e6a8503"

// defaultPrivateKey is the Anvil account #0 key — used only for offline signing demos.
// NOT funded on Polygon mainnet.
const defaultPrivateKey = "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

// privateKeyHex reads from the PRIVATE_KEY env var, falling back to the default.
// Set PRIVATE_KEY in a .env file or export it in your shell.
var privateKeyHex = getEnvOrDefault("PRIVATE_KEY", defaultPrivateKey)

func getEnvOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

var (
	SourceAddr = common.HexToAddress("0x83d83708D7E01977454cE1859c1E188060274297")
	TargetAddr = common.HexToAddress("0x5d9339C29f1582e08F2b69bfa5740D11E0898D1F")
	// USDC on Polygon
	USDCAddr = common.HexToAddress("0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359")
)

// ---------------------------------------------------------------------------
// Clients (lazily initialised by Setup)
// ---------------------------------------------------------------------------

var (
	PublicClient *client.PublicClient
	WalletCl     *client.WalletClient
	LocalAccount *accounts.PrivateKeyAccount
)

// Setup initialises clients. Call once from main.
func Setup() error {
	var err error
	PublicClient, err = client.CreatePublicClient(client.PublicClientConfig{
		Chain:     &definitions.Polygon,
		Transport: transport.HTTP(polygonRPC),
	})
	if err != nil {
		return fmt.Errorf("public client: %w", err)
	}
	LocalAccount, err = accounts.PrivateKeyToAccount(privateKeyHex)
	if err != nil {
		return fmt.Errorf("local account: %w", err)
	}

	// Create wallet client with the local PrivateKeyAccount.
	// Because PrivateKeyAccount now satisfies client.Account and all wallet signing
	// interfaces, SendTransaction will take the local path (prepare + sign + sendRaw)
	// instead of asking the RPC node to sign.
	WalletCl, err = client.CreateWalletClient(client.WalletClientConfig{
		Account:   LocalAccount,
		Chain:     &definitions.Polygon,
		Transport: transport.HTTP(polygonRPC),
	})
	if err != nil {
		return fmt.Errorf("wallet client: %w", err)
	}

	return nil
}

// Teardown cleans up clients.
func Teardown() {
	if PublicClient != nil {
		PublicClient.Close()
	}
	if WalletCl != nil {
		WalletCl.Close()
	}
}

// PrintAccountInfo prints the account address and balance on Polygon.
func PrintAccountInfo() {
	fmt.Printf("Account:  %s\n", SourceAddr.Hex())
	fmt.Printf("Chain:    Polygon (%d)\n", definitions.Polygon.ID)
	ctx := context.Background()
	balance, err := public.GetBalance(ctx, PublicClient, public.GetBalanceParameters{Address: SourceAddr})
	if err == nil {
		fmt.Printf("Balance:  %s POL\n", unit.FormatEther(balance))
	} else {
		fmt.Println("Balance:  (could not fetch)")
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func MustParseGwei(s string) *big.Int {
	v, err := unit.ParseGwei(s)
	if err != nil {
		panic(fmt.Sprintf("invalid gwei %q: %v", s, err))
	}
	return v
}

func MustParseEther(s string) *big.Int {
	v, err := unit.ParseEther(s)
	if err != nil {
		panic(fmt.Sprintf("invalid ether %q: %v", s, err))
	}
	return v
}

func PrintHeader(title string) {
	fmt.Println()
	fmt.Println(strings.Repeat("=", 70))
	fmt.Printf("  %s\n", title)
	fmt.Println(strings.Repeat("=", 70))
}

func PrintSection(title string) {
	fmt.Printf("\n--- %s ---\n", title)
}

func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
