package client_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ChefBingbong/viem-go/client"
	"github.com/ChefBingbong/viem-go/client/transport"
)

func createTestServer(t *testing.T, handler func(method string, params []any) any) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			JSONRPC string `json:"jsonrpc"`
			ID      any    `json:"id"`
			Method  string `json:"method"`
			Params  []any  `json:"params"`
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)

		result := handler(req.Method, req.Params)

		resp := map[string]any{
			"jsonrpc": "2.0",
			"id":      req.ID,
			"result":  result,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
}

func TestCreateClient(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		switch method {
		case "eth_chainId":
			return "0x1"
		case "eth_blockNumber":
			return "0x10"
		default:
			return "0x0"
		}
	})
	defer server.Close()

	c, err := client.CreateClient(client.ClientConfig{
		Transport: transport.HTTP(server.URL),
		Key:       "test",
		Name:      "Test Client",
	})
	require.NoError(t, err)
	defer c.Close()

	assert.Equal(t, "test", c.Key())
	assert.Equal(t, "Test Client", c.Name())
	assert.Equal(t, "base", c.Type())
	assert.NotEmpty(t, c.UID())
}

func TestCreateClientWithChain(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		return "0x1"
	})
	defer server.Close()

	chain := &client.Chain{
		ID:        1,
		Name:      "Ethereum",
		BlockTime: 12000,
		NativeCurrency: &client.NativeCurrency{
			Name:     "Ether",
			Symbol:   "ETH",
			Decimals: 18,
		},
	}

	c, err := client.CreateClient(client.ClientConfig{
		Transport: transport.HTTP(server.URL),
		Chain:     chain,
	})
	require.NoError(t, err)
	defer c.Close()

	assert.Equal(t, chain, c.Chain())
	assert.Equal(t, 1, c.Chain().ID)
	assert.Equal(t, "Ethereum", c.Chain().Name)
}

func TestCreateClientWithAccount(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		return "0x1"
	})
	defer server.Close()

	addr := common.HexToAddress("0x1234567890123456789012345678901234567890")
	account := client.NewAddressAccount(addr)

	c, err := client.CreateClient(client.ClientConfig{
		Transport: transport.HTTP(server.URL),
		Account:   account,
	})
	require.NoError(t, err)
	defer c.Close()

	assert.NotNil(t, c.Account())
	assert.Equal(t, addr, c.Account().Address())
}

func TestPublicClient_GetBlockNumber(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_blockNumber" {
			return "0x10"
		}
		return "0x0"
	})
	defer server.Close()

	c, err := client.CreatePublicClient(client.PublicClientConfig{
		Transport: transport.HTTP(server.URL),
	})
	require.NoError(t, err)
	defer c.Close()

	ctx := context.Background()
	blockNumber, err := c.GetBlockNumber(ctx)
	require.NoError(t, err)
	assert.Equal(t, uint64(16), blockNumber)
}

func TestPublicClient_GetChainID(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_chainId" {
			return "0x1"
		}
		return "0x0"
	})
	defer server.Close()

	c, err := client.CreatePublicClient(client.PublicClientConfig{
		Transport: transport.HTTP(server.URL),
	})
	require.NoError(t, err)
	defer c.Close()

	ctx := context.Background()
	chainID, err := c.GetChainID(ctx)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), chainID)
}

func TestPublicClient_GetBalance(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_getBalance" {
			return "0xde0b6b3a7640000" // 1 ETH in wei
		}
		return "0x0"
	})
	defer server.Close()

	c, err := client.CreatePublicClient(client.PublicClientConfig{
		Transport: transport.HTTP(server.URL),
	})
	require.NoError(t, err)
	defer c.Close()

	ctx := context.Background()
	addr := common.HexToAddress("0x1234567890123456789012345678901234567890")
	balance, err := c.GetBalance(ctx, addr)
	require.NoError(t, err)
	assert.NotNil(t, balance)
}

func TestPublicClient_Call(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_call" {
			return "0x000000000000000000000000000000000000000000000000000000000000002a" // 42
		}
		return "0x0"
	})
	defer server.Close()

	c, err := client.CreatePublicClient(client.PublicClientConfig{
		Transport: transport.HTTP(server.URL),
	})
	require.NoError(t, err)
	defer c.Close()

	ctx := context.Background()
	to := common.HexToAddress("0x1234567890123456789012345678901234567890")
	result, err := c.Call(ctx, client.CallRequest{
		To: to,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestBaseClient_Extend(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		return "0x1"
	})
	defer server.Close()

	c, err := client.CreateClient(client.ClientConfig{
		Transport: transport.HTTP(server.URL),
	})
	require.NoError(t, err)
	defer c.Close()

	// Test extend
	c.Extend("customAction", func() string { return "custom" })

	val, ok := c.GetExtension("customAction")
	assert.True(t, ok)
	assert.NotNil(t, val)

	// Test extensions map
	exts := c.Extensions()
	assert.Contains(t, exts, "customAction")
}

func TestCreatePublicClient(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		return "0x1"
	})
	defer server.Close()

	c, err := client.CreatePublicClient(client.PublicClientConfig{
		Transport: transport.HTTP(server.URL),
	})
	require.NoError(t, err)
	defer c.Close()

	assert.Equal(t, "publicClient", c.Type())
	assert.Equal(t, "Public Client", c.Name())
}

func TestCreateWalletClient(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		return "0x1"
	})
	defer server.Close()

	c, err := client.CreateWalletClient(client.WalletClientConfig{
		Transport: transport.HTTP(server.URL),
	})
	require.NoError(t, err)
	defer c.Close()

	assert.Equal(t, "walletClient", c.Type())
	assert.Equal(t, "Wallet Client", c.Name())
}

func TestClientConfig_Defaults(t *testing.T) {
	config := client.DefaultClientConfig()

	assert.Equal(t, "base", config.Key)
	assert.Equal(t, "Base Client", config.Name)
	assert.Equal(t, "base", config.Type)
	assert.Equal(t, 4000*time.Millisecond, config.CacheTime)
	assert.Equal(t, 4000*time.Millisecond, config.PollingInterval)
}

func TestPublicActions(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		switch method {
		case "eth_blockNumber":
			return "0x100"
		case "eth_chainId":
			return "0x1"
		default:
			return "0x0"
		}
	})
	defer server.Close()

	c, err := client.CreatePublicClient(client.PublicClientConfig{
		Transport: transport.HTTP(server.URL),
	})
	require.NoError(t, err)
	defer c.Close()

	ctx := context.Background()

	blockNumber, err := c.GetBlockNumber(ctx)
	require.NoError(t, err)
	assert.Equal(t, uint64(256), blockNumber)

	chainID, err := c.GetChainID(ctx)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), chainID)
}
