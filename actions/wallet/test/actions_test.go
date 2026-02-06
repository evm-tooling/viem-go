package wallet_test

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ChefBingbong/viem-go/actions/wallet"
	"github.com/ChefBingbong/viem-go/chain"
	"github.com/ChefBingbong/viem-go/client/transport"
	"github.com/ChefBingbong/viem-go/types"
	"github.com/ChefBingbong/viem-go/utils/formatters"
	"github.com/ChefBingbong/viem-go/utils/signature"
	utiltx "github.com/ChefBingbong/viem-go/utils/transaction"
)

// ============================================================================
// Mock Client & Account Types
// ============================================================================

// mockAccount implements wallet.Account for testing JSON-RPC account paths.
type mockAccount struct {
	address common.Address
}

func (a *mockAccount) Address() common.Address { return a.address }

// mockLocalAccount implements wallet.Account + wallet.LocalAccount for testing
// the local account path where GetAddresses returns the account directly.
type mockLocalAccount struct {
	address common.Address
}

func (a *mockLocalAccount) Address() common.Address { return a.address }
func (a *mockLocalAccount) IsLocal()                {}

// mockSignableAccount implements wallet.SignableAccount for local message signing.
type mockSignableAccount struct {
	address   common.Address
	signFn    func(msg signature.SignableMessage) (string, error)
}

func (a *mockSignableAccount) Address() common.Address { return a.address }
func (a *mockSignableAccount) IsLocal()                {}
func (a *mockSignableAccount) SignMessage(msg signature.SignableMessage) (string, error) {
	return a.signFn(msg)
}

// mockTypedDataSignableAccount implements wallet.TypedDataSignableAccount.
type mockTypedDataSignableAccount struct {
	address common.Address
	signFn  func(data signature.TypedDataDefinition) (string, error)
}

func (a *mockTypedDataSignableAccount) Address() common.Address { return a.address }
func (a *mockTypedDataSignableAccount) IsLocal()                {}
func (a *mockTypedDataSignableAccount) SignTypedData(data signature.TypedDataDefinition) (string, error) {
	return a.signFn(data)
}

// mockTransactionSignableAccount implements wallet.TransactionSignableAccount.
type mockTransactionSignableAccount struct {
	address common.Address
	signFn  func(tx *utiltx.Transaction) (string, error)
}

func (a *mockTransactionSignableAccount) Address() common.Address { return a.address }
func (a *mockTransactionSignableAccount) IsLocal()                {}
func (a *mockTransactionSignableAccount) SignTransaction(tx *utiltx.Transaction) (string, error) {
	return a.signFn(tx)
}

// mockAuthorizationSignableAccount implements wallet.AuthorizationSignableAccount.
type mockAuthorizationSignableAccount struct {
	address common.Address
	signFn  func(auth types.AuthorizationRequest) (*types.SignedAuthorization, error)
}

func (a *mockAuthorizationSignableAccount) Address() common.Address { return a.address }
func (a *mockAuthorizationSignableAccount) IsLocal()                {}
func (a *mockAuthorizationSignableAccount) SignAuthorization(auth types.AuthorizationRequest) (*types.SignedAuthorization, error) {
	return a.signFn(auth)
}

// mockClient implements wallet.Client for testing.
type mockClient struct {
	transport       transport.Transport
	chain           *chain.Chain
	cacheTime       time.Duration
	blockTag        types.BlockTag
	batch           *types.BatchOptions
	ccipRead        *types.CCIPReadOptions
	uid             string
	dataSuffix      []byte
	pollingInterval time.Duration
	account         wallet.Account
	requestRecorder func(method string, params []any)
}

func (c *mockClient) Request(ctx context.Context, method string, params ...any) (*transport.RPCResponse, error) {
	if c.requestRecorder != nil {
		c.requestRecorder(method, params)
	}
	return c.transport.Request(ctx, transport.RPCRequest{Method: method, Params: params})
}

func (c *mockClient) Chain() *chain.Chain {
	return c.chain
}

func (c *mockClient) CacheTime() time.Duration {
	if c.cacheTime == 0 {
		return 4 * time.Second
	}
	return c.cacheTime
}

func (c *mockClient) ExperimentalBlockTag() types.BlockTag {
	return c.blockTag
}

func (c *mockClient) Batch() *types.BatchOptions {
	return c.batch
}

func (c *mockClient) CCIPRead() *types.CCIPReadOptions {
	return c.ccipRead
}

func (c *mockClient) UID() string {
	if c.uid == "" {
		return "test-wallet-mock-client"
	}
	return c.uid
}

func (c *mockClient) DataSuffix() []byte {
	return c.dataSuffix
}

func (c *mockClient) PollingInterval() time.Duration {
	if c.pollingInterval == 0 {
		return 100 * time.Millisecond
	}
	return c.pollingInterval
}

func (c *mockClient) Account() wallet.Account {
	return c.account
}

// ============================================================================
// Test Helpers
// ============================================================================

// createTestServer creates a test HTTP server that responds to JSON-RPC requests.
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

// createErrorTestServer creates a test server that responds with a JSON-RPC error.
func createErrorTestServer(t *testing.T, code int, message string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			JSONRPC string `json:"jsonrpc"`
			ID      any    `json:"id"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)

		resp := map[string]any{
			"jsonrpc": "2.0",
			"id":      req.ID,
			"error": map[string]any{
				"code":    code,
				"message": message,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
}

// createMockClient creates a mock client for testing.
func createMockClient(t *testing.T, serverURL string) *mockClient {
	tr, err := transport.HTTP(serverURL)(transport.TransportParams{})
	require.NoError(t, err)

	return &mockClient{
		transport: tr,
	}
}

// testChain creates a minimal chain config for testing with the given chain ID.
func testChain(id int64) *chain.Chain {
	return &chain.Chain{
		ID:   id,
		Name: "Test Chain",
		NativeCurrency: chain.ChainNativeCurrency{
			Name:     "Ether",
			Symbol:   "ETH",
			Decimals: 18,
		},
		RpcUrls: map[string]chain.ChainRpcUrls{
			"default": {HTTP: []string{"http://localhost:8545"}},
		},
	}
}

var (
	sourceAddr = common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
	targetAddr = common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8")
)

// ============================================================================
// SendTransaction Tests
// ============================================================================

func TestSendTransaction_Default(t *testing.T) {
	var capturedMethod string
	server := createTestServer(t, func(method string, params []any) any {
		capturedMethod = method
		switch method {
		case "eth_chainId":
			return "0x1"
		case "eth_sendTransaction":
			return "0xabc123def456abc123def456abc123def456abc123def456abc123def456abc1"
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	client.chain = testChain(1)
	ctx := context.Background()

	hash, err := wallet.SendTransaction(ctx, client, wallet.SendTransactionParameters{
		Account: &mockAccount{address: sourceAddr},
		To:      targetAddr.Hex(),
		Value:   big.NewInt(1000000000000000000), // 1 ETH
	})

	require.NoError(t, err)
	assert.Equal(t, "eth_sendTransaction", capturedMethod)
	assert.Equal(t, "0xabc123def456abc123def456abc123def456abc123def456abc123def456abc1", hash)
}

func TestSendTransaction_InferredAccount(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		switch method {
		case "eth_chainId":
			return "0x1"
		case "eth_sendTransaction":
			return "0xabc123def456abc123def456abc123def456abc123def456abc123def456abc1"
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	client.chain = testChain(1)
	client.account = &mockAccount{address: sourceAddr}
	ctx := context.Background()

	hash, err := wallet.SendTransaction(ctx, client, wallet.SendTransactionParameters{
		To:    targetAddr.Hex(),
		Value: big.NewInt(1000000000000000000),
	})

	require.NoError(t, err)
	assert.NotEmpty(t, hash)
}

func TestSendTransaction_NoAccount(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	_, err := wallet.SendTransaction(ctx, client, wallet.SendTransactionParameters{
		To:    targetAddr.Hex(),
		Value: big.NewInt(1000000000000000000),
	})

	require.Error(t, err)
	_, ok := err.(*wallet.AccountNotFoundError)
	assert.True(t, ok, "expected AccountNotFoundError, got %T: %v", err, err)
}

func TestSendTransaction_ChainMismatch(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_chainId" {
			return "0x1" // mainnet
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	client.chain = testChain(1)
	ctx := context.Background()

	// Provide a chain with a different ID (optimism = 10)
	optimismChain := testChain(10)
	optimismChain.Name = "OP Mainnet"

	_, err := wallet.SendTransaction(ctx, client, wallet.SendTransactionParameters{
		Account: &mockAccount{address: sourceAddr},
		Chain:   optimismChain,
		To:      targetAddr.Hex(),
		Value:   big.NewInt(1000000000000000000),
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "chain")
}

func TestSendTransaction_WithGasPrice(t *testing.T) {
	var capturedParams []any
	server := createTestServer(t, func(method string, params []any) any {
		switch method {
		case "eth_chainId":
			return "0x1"
		case "eth_sendTransaction":
			capturedParams = params
			return "0xabc123def456abc123def456abc123def456abc123def456abc123def456abc1"
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	client.chain = testChain(1)
	ctx := context.Background()

	hash, err := wallet.SendTransaction(ctx, client, wallet.SendTransactionParameters{
		Account:  &mockAccount{address: sourceAddr},
		To:       targetAddr.Hex(),
		Value:    big.NewInt(1000000000000000000),
		GasPrice: big.NewInt(20000000000), // 20 gwei
	})

	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	require.NotEmpty(t, capturedParams)
}

func TestSendTransaction_WithMaxFeePerGas(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		switch method {
		case "eth_chainId":
			return "0x1"
		case "eth_sendTransaction":
			return "0xabc123def456abc123def456abc123def456abc123def456abc123def456abc1"
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	client.chain = testChain(1)
	ctx := context.Background()

	hash, err := wallet.SendTransaction(ctx, client, wallet.SendTransactionParameters{
		Account:      &mockAccount{address: sourceAddr},
		To:           targetAddr.Hex(),
		Value:        big.NewInt(1000000000000000000),
		MaxFeePerGas: big.NewInt(50000000000), // 50 gwei
	})

	require.NoError(t, err)
	assert.NotEmpty(t, hash)
}

func TestSendTransaction_WithMaxFeeAndPriority(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		switch method {
		case "eth_chainId":
			return "0x1"
		case "eth_sendTransaction":
			return "0xabc123def456abc123def456abc123def456abc123def456abc123def456abc1"
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	client.chain = testChain(1)
	ctx := context.Background()

	hash, err := wallet.SendTransaction(ctx, client, wallet.SendTransactionParameters{
		Account:              &mockAccount{address: sourceAddr},
		To:                   targetAddr.Hex(),
		Value:                big.NewInt(1000000000000000000),
		MaxFeePerGas:         big.NewInt(50000000000),
		MaxPriorityFeePerGas: big.NewInt(2000000000),
	})

	require.NoError(t, err)
	assert.NotEmpty(t, hash)
}

func TestSendTransaction_WithNonce(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		switch method {
		case "eth_chainId":
			return "0x1"
		case "eth_sendTransaction":
			return "0xabc123def456abc123def456abc123def456abc123def456abc123def456abc1"
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	client.chain = testChain(1)
	ctx := context.Background()

	nonce := 42
	hash, err := wallet.SendTransaction(ctx, client, wallet.SendTransactionParameters{
		Account: &mockAccount{address: sourceAddr},
		To:      targetAddr.Hex(),
		Value:   big.NewInt(1000000000000000000),
		Nonce:   &nonce,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, hash)
}

func TestSendTransaction_LocalAccount(t *testing.T) {
	var capturedMethod string
	server := createTestServer(t, func(method string, params []any) any {
		capturedMethod = method
		switch method {
		case "eth_chainId":
			return "0x1"
		case "eth_getTransactionCount":
			return "0x0"
		case "eth_getBlockByNumber":
			return map[string]any{
				"number":      "0x10",
				"baseFeePerGas": "0x3b9aca00",
				"gasLimit":    "0x1c9c380",
				"gasUsed":     "0x0",
				"timestamp":   "0x60000000",
				"hash":        "0x1234567890123456789012345678901234567890123456789012345678901234",
				"parentHash":  "0x0000000000000000000000000000000000000000000000000000000000000000",
				"transactions": []string{},
			}
		case "eth_maxPriorityFeePerGas":
			return "0x3b9aca00" // 1 gwei
		case "eth_estimateGas":
			return "0x5208" // 21000
		case "eth_sendRawTransaction":
			return "0xlocalhash123456789012345678901234567890123456789012345678901234"
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	client.chain = testChain(1)
	ctx := context.Background()

	localAccount := &mockTransactionSignableAccount{
		address: sourceAddr,
		signFn: func(tx *utiltx.Transaction) (string, error) {
			return "0x02f850018203118080825208808080c080a04012522854168b27e5dc3d5839bab5e6b39e1a0ffd343901ce1622e3d64b48f1a04e00902ae0502c4728cbf12156290df99c3ed7de85b1dbfe20b5c36931733a33", nil
		},
	}

	hash, err := wallet.SendTransaction(ctx, client, wallet.SendTransactionParameters{
		Account: localAccount,
		To:      targetAddr.Hex(),
		Value:   big.NewInt(1000000000000000000),
	})

	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	// Local account path should end with eth_sendRawTransaction
	assert.Equal(t, "eth_sendRawTransaction", capturedMethod)
}

func TestSendTransaction_DataSuffix(t *testing.T) {
	var capturedParams []any
	server := createTestServer(t, func(method string, params []any) any {
		switch method {
		case "eth_chainId":
			return "0x1"
		case "eth_sendTransaction":
			capturedParams = params
			return "0xabc123def456abc123def456abc123def456abc123def456abc123def456abc1"
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	client.chain = testChain(1)
	client.dataSuffix = []byte{0x12, 0x34, 0x56, 0x78}
	ctx := context.Background()

	hash, err := wallet.SendTransaction(ctx, client, wallet.SendTransactionParameters{
		Account: &mockAccount{address: sourceAddr},
		To:      targetAddr.Hex(),
		Data:    "0xdeadbeef",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	require.NotEmpty(t, capturedParams)
}

func TestSendTransaction_ParamDataSuffixOverridesClient(t *testing.T) {
	var capturedParams []any
	server := createTestServer(t, func(method string, params []any) any {
		switch method {
		case "eth_chainId":
			return "0x1"
		case "eth_sendTransaction":
			capturedParams = params
			return "0xabc123def456abc123def456abc123def456abc123def456abc123def456abc1"
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	client.chain = testChain(1)
	client.dataSuffix = []byte{0xAA, 0xBB}
	ctx := context.Background()

	hash, err := wallet.SendTransaction(ctx, client, wallet.SendTransactionParameters{
		Account:    &mockAccount{address: sourceAddr},
		To:         targetAddr.Hex(),
		Data:       "0xdeadbeef",
		DataSuffix: "0xccdd",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	require.NotEmpty(t, capturedParams)
}

// ============================================================================
// SendRawTransaction Tests
// ============================================================================

func TestSendRawTransaction_Default(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_sendRawTransaction" {
			return "0xhash0000000000000000000000000000000000000000000000000000000001"
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	hash, err := wallet.SendRawTransaction(ctx, client, wallet.SendRawTransactionParameters{
		SerializedTransaction: "0x02f850018203118080825208808080c080a04012522854168b27e5dc3d5839bab5e6b39e1a0ffd343901ce1622e3d64b48f1a04e00902ae0502c4728cbf12156290df99c3ed7de85b1dbfe20b5c36931733a33",
	})

	require.NoError(t, err)
	assert.Equal(t, "0xhash0000000000000000000000000000000000000000000000000000000001", hash)
}

func TestSendRawTransaction_RPCError(t *testing.T) {
	server := createErrorTestServer(t, -32000, "nonce too low")
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	_, err := wallet.SendRawTransaction(ctx, client, wallet.SendRawTransactionParameters{
		SerializedTransaction: "0x02f850018203118080825208808080c080a04012",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "eth_sendRawTransaction failed")
}

// ============================================================================
// SignMessage Tests
// ============================================================================

func TestSignMessage_JSONRPC(t *testing.T) {
	expectedSig := "0xa461f509887bd19e312c0c58467ce8ff8e300d3c1a90b608a760c5b80318eaf15fe57c96f9175d6cd4daad4663763baa7e78836e067d0163e9a2ccf2ff753f5b1b"
	server := createTestServer(t, func(method string, params []any) any {
		if method == "personal_sign" {
			return expectedSig
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	sig, err := wallet.SignMessage(ctx, client, wallet.SignMessageParameters{
		Account: &mockAccount{address: sourceAddr},
		Message: signature.NewSignableMessage("hello world"),
	})

	require.NoError(t, err)
	assert.Equal(t, expectedSig, sig)
}

func TestSignMessage_Raw(t *testing.T) {
	expectedSig := "0xa461f509887bd19e312c0c58467ce8ff8e300d3c1a90b608a760c5b80318eaf15fe57c96f9175d6cd4daad4663763baa7e78836e067d0163e9a2ccf2ff753f5b1b"
	server := createTestServer(t, func(method string, params []any) any {
		if method == "personal_sign" {
			return expectedSig
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	sig, err := wallet.SignMessage(ctx, client, wallet.SignMessageParameters{
		Account: &mockAccount{address: sourceAddr},
		Message: signature.NewSignableMessageRawHex("0x68656c6c6f20776f726c64"),
	})

	require.NoError(t, err)
	assert.Equal(t, expectedSig, sig)
}

func TestSignMessage_RawBytes(t *testing.T) {
	expectedSig := "0xa461f509887bd19e312c0c58467ce8ff8e300d3c1a90b608a760c5b80318eaf15fe57c96f9175d6cd4daad4663763baa7e78836e067d0163e9a2ccf2ff753f5b1b"
	server := createTestServer(t, func(method string, params []any) any {
		if method == "personal_sign" {
			return expectedSig
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	sig, err := wallet.SignMessage(ctx, client, wallet.SignMessageParameters{
		Account: &mockAccount{address: sourceAddr},
		Message: signature.NewSignableMessageRaw([]byte{104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100}),
	})

	require.NoError(t, err)
	assert.Equal(t, expectedSig, sig)
}

func TestSignMessage_InferredAccount(t *testing.T) {
	expectedSig := "0xa461f509887bd19e312c0c58467ce8ff8e300d3c1a90b608a760c5b80318eaf15fe57c96f9175d6cd4daad4663763baa7e78836e067d0163e9a2ccf2ff753f5b1b"
	server := createTestServer(t, func(method string, params []any) any {
		if method == "personal_sign" {
			return expectedSig
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	client.account = &mockAccount{address: sourceAddr}
	ctx := context.Background()

	sig, err := wallet.SignMessage(ctx, client, wallet.SignMessageParameters{
		Message: signature.NewSignableMessage("hello world"),
	})

	require.NoError(t, err)
	assert.Equal(t, expectedSig, sig)
}

func TestSignMessage_LocalAccount(t *testing.T) {
	expectedSig := "0xlocalsig123"
	localAccount := &mockSignableAccount{
		address: sourceAddr,
		signFn: func(msg signature.SignableMessage) (string, error) {
			return expectedSig, nil
		},
	}

	// Server should NOT be called for local accounts
	server := createTestServer(t, func(method string, params []any) any {
		t.Fatal("RPC should not be called for local account signing")
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	sig, err := wallet.SignMessage(ctx, client, wallet.SignMessageParameters{
		Account: localAccount,
		Message: signature.NewSignableMessage("hello world"),
	})

	require.NoError(t, err)
	assert.Equal(t, expectedSig, sig)
}

func TestSignMessage_NoAccount(t *testing.T) {
	client := &mockClient{}
	ctx := context.Background()

	_, err := wallet.SignMessage(ctx, client, wallet.SignMessageParameters{
		Message: signature.NewSignableMessage("hello world"),
	})

	require.Error(t, err)
	_, ok := err.(*wallet.AccountNotFoundError)
	assert.True(t, ok, "expected AccountNotFoundError, got %T: %v", err, err)
}

// ============================================================================
// SignTransaction Tests
// ============================================================================

func TestSignTransaction_JSONRPC(t *testing.T) {
	expectedSig := "0x02f854018203118084607d7d8a825208808080c080a02591128fce3fce3e2c4feaafb1cadfcafe81fa66f00b0eec2ca5bb9bf05ebeb9a019edec10144ec5e05de3f5fff2b792cbe6e7a946f659a2020f8fee4d4689df6a"
	server := createTestServer(t, func(method string, params []any) any {
		switch method {
		case "eth_chainId":
			return "0x1"
		case "eth_signTransaction":
			return expectedSig
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	client.chain = testChain(1)
	ctx := context.Background()

	nonce := 785
	sig, err := wallet.SignTransaction(ctx, client, wallet.SignTransactionParameters{
		Account: &mockAccount{address: sourceAddr},
		Gas:     big.NewInt(21000),
		Nonce:   &nonce,
	})

	require.NoError(t, err)
	assert.Equal(t, expectedSig, sig)
}

func TestSignTransaction_LocalAccount(t *testing.T) {
	expectedSig := "0x02f850018203118080825208808080c080a04012522854168b27e5dc3d5839bab5e6b39e1a0ffd343901ce1622e3d64b48f1a04e00902ae0502c4728cbf12156290df99c3ed7de85b1dbfe20b5c36931733a33"
	localAccount := &mockTransactionSignableAccount{
		address: sourceAddr,
		signFn: func(tx *utiltx.Transaction) (string, error) {
			assert.Equal(t, 1, tx.ChainId)
			return expectedSig, nil
		},
	}

	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_chainId" {
			return "0x1"
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	client.chain = testChain(1)
	ctx := context.Background()

	nonce := 785
	sig, err := wallet.SignTransaction(ctx, client, wallet.SignTransactionParameters{
		Account: localAccount,
		Gas:     big.NewInt(21000),
		Nonce:   &nonce,
	})

	require.NoError(t, err)
	assert.Equal(t, expectedSig, sig)
}

func TestSignTransaction_WithMaxFeeAndPriority(t *testing.T) {
	expectedSig := "0x02f859018203118477359400850....."
	server := createTestServer(t, func(method string, params []any) any {
		switch method {
		case "eth_chainId":
			return "0x1"
		case "eth_signTransaction":
			return expectedSig
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	client.chain = testChain(1)
	ctx := context.Background()

	nonce := 785
	sig, err := wallet.SignTransaction(ctx, client, wallet.SignTransactionParameters{
		Account:              &mockAccount{address: sourceAddr},
		Gas:                  big.NewInt(21000),
		Nonce:                &nonce,
		MaxFeePerGas:         big.NewInt(20000000000),
		MaxPriorityFeePerGas: big.NewInt(2000000000),
	})

	require.NoError(t, err)
	assert.Equal(t, expectedSig, sig)
}

func TestSignTransaction_WithData(t *testing.T) {
	expectedSig := "0x02f852018203118080825208808082123..."
	server := createTestServer(t, func(method string, params []any) any {
		switch method {
		case "eth_chainId":
			return "0x1"
		case "eth_signTransaction":
			return expectedSig
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	client.chain = testChain(1)
	ctx := context.Background()

	nonce := 785
	sig, err := wallet.SignTransaction(ctx, client, wallet.SignTransactionParameters{
		Account: &mockAccount{address: sourceAddr},
		Gas:     big.NewInt(21000),
		Nonce:   &nonce,
		Data:    "0x1234",
	})

	require.NoError(t, err)
	assert.Equal(t, expectedSig, sig)
}

func TestSignTransaction_ChainMismatch(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_chainId" {
			return "0x1" // mainnet
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	client.chain = testChain(10) // optimism
	ctx := context.Background()

	_, err := wallet.SignTransaction(ctx, client, wallet.SignTransactionParameters{
		Account: &mockAccount{address: sourceAddr},
		Gas:     big.NewInt(21000),
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "chain")
}

func TestSignTransaction_NoAccount(t *testing.T) {
	client := &mockClient{}
	ctx := context.Background()

	_, err := wallet.SignTransaction(ctx, client, wallet.SignTransactionParameters{
		Gas: big.NewInt(21000),
	})

	require.Error(t, err)
	_, ok := err.(*wallet.AccountNotFoundError)
	assert.True(t, ok, "expected AccountNotFoundError, got %T: %v", err, err)
}

// ============================================================================
// SignTypedData Tests
// ============================================================================

func TestSignTypedData_JSONRPC(t *testing.T) {
	expectedSig := "0xsig_typed_data_123"
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_signTypedData_v4" {
			return expectedSig
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	sig, err := wallet.SignTypedData(ctx, client, wallet.SignTypedDataParameters{
		Account: &mockAccount{address: sourceAddr},
		Domain: signature.TypedDataDomain{
			Name:              "Ether Mail",
			Version:           "1",
			ChainId:           big.NewInt(1),
			VerifyingContract: "0xCcCCccccCCCCcCCCCCCcCcCccCcCCCcCcccccccC",
		},
		Types: map[string][]signature.TypedDataField{
			"Person": {
				{Name: "name", Type: "string"},
				{Name: "wallet", Type: "address"},
			},
			"Mail": {
				{Name: "from", Type: "Person"},
				{Name: "to", Type: "Person"},
				{Name: "contents", Type: "string"},
			},
		},
		PrimaryType: "Mail",
		Message: map[string]any{
			"from":     map[string]any{"name": "Cow", "wallet": "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826"},
			"to":       map[string]any{"name": "Bob", "wallet": "0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB"},
			"contents": "Hello, Bob!",
		},
	})

	require.NoError(t, err)
	assert.Equal(t, expectedSig, sig)
}

func TestSignTypedData_LocalAccount(t *testing.T) {
	expectedSig := "0xlocal_typed_data_sig_123"
	localAccount := &mockTypedDataSignableAccount{
		address: sourceAddr,
		signFn: func(data signature.TypedDataDefinition) (string, error) {
			assert.Equal(t, "Mail", data.PrimaryType)
			return expectedSig, nil
		},
	}

	server := createTestServer(t, func(method string, params []any) any {
		t.Fatal("RPC should not be called for local typed data signing")
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	sig, err := wallet.SignTypedData(ctx, client, wallet.SignTypedDataParameters{
		Account: localAccount,
		Domain: signature.TypedDataDomain{
			Name:              "Ether Mail",
			Version:           "1",
			ChainId:           big.NewInt(1),
			VerifyingContract: "0xCcCCccccCCCCcCCCCCCcCcCccCcCCCcCcccccccC",
		},
		Types: map[string][]signature.TypedDataField{
			"Person": {
				{Name: "name", Type: "string"},
				{Name: "wallet", Type: "address"},
			},
			"Mail": {
				{Name: "from", Type: "Person"},
				{Name: "to", Type: "Person"},
				{Name: "contents", Type: "string"},
			},
		},
		PrimaryType: "Mail",
		Message: map[string]any{
			"from":     map[string]any{"name": "Cow", "wallet": "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826"},
			"to":       map[string]any{"name": "Bob", "wallet": "0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB"},
			"contents": "Hello, Bob!",
		},
	})

	require.NoError(t, err)
	assert.Equal(t, expectedSig, sig)
}

func TestSignTypedData_NoAccount(t *testing.T) {
	client := &mockClient{}
	ctx := context.Background()

	_, err := wallet.SignTypedData(ctx, client, wallet.SignTypedDataParameters{
		Domain:      signature.TypedDataDomain{Name: "Test"},
		Types:       map[string][]signature.TypedDataField{"Test": {{Name: "value", Type: "uint256"}}},
		PrimaryType: "Test",
		Message:     map[string]any{"value": big.NewInt(1)},
	})

	require.Error(t, err)
	_, ok := err.(*wallet.AccountNotFoundError)
	assert.True(t, ok, "expected AccountNotFoundError")
}

func TestSignTypedData_InvalidVerifyingContract(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	_, err := wallet.SignTypedData(ctx, client, wallet.SignTypedDataParameters{
		Account: &mockAccount{address: sourceAddr},
		Domain: signature.TypedDataDomain{
			Name:              "Test",
			VerifyingContract: "not-an-address",
		},
		Types:       map[string][]signature.TypedDataField{"Test": {{Name: "value", Type: "uint256"}}},
		PrimaryType: "Test",
		Message:     map[string]any{"value": big.NewInt(1)},
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid verifying contract")
}

func TestSignTypedData_InvalidPrimaryType(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	_, err := wallet.SignTypedData(ctx, client, wallet.SignTypedDataParameters{
		Account: &mockAccount{address: sourceAddr},
		Domain: signature.TypedDataDomain{
			Name: "Test",
		},
		Types:       map[string][]signature.TypedDataField{"Test": {{Name: "value", Type: "uint256"}}},
		PrimaryType: "NonExistent",
		Message:     map[string]any{"value": big.NewInt(1)},
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid primary type")
}

// ============================================================================
// WriteContract Tests
// ============================================================================

func TestWriteContract_Default(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		switch method {
		case "eth_chainId":
			return "0x1"
		case "eth_sendTransaction":
			return "0xwritehash000000000000000000000000000000000000000000000000000001"
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	client.chain = testChain(1)
	ctx := context.Background()

	abiJSON := `[{"inputs":[{"name":"tokenId","type":"uint32"}],"name":"mint","outputs":[],"stateMutability":"nonpayable","type":"function"}]`

	hash, err := wallet.WriteContract(ctx, client, wallet.WriteContractParameters{
		Account:      &mockAccount{address: sourceAddr},
		Address:      "0xFBA3912Ca04dd458c843e2EE08967fC04f3579c2",
		ABI:          abiJSON,
		FunctionName: "mint",
		Args:         []any{uint32(69420)},
	})

	require.NoError(t, err)
	assert.Equal(t, "0xwritehash000000000000000000000000000000000000000000000000000001", hash)
}

func TestWriteContract_NoAccount(t *testing.T) {
	client := &mockClient{}
	ctx := context.Background()

	abiJSON := `[{"inputs":[],"name":"mint","outputs":[],"stateMutability":"nonpayable","type":"function"}]`

	_, err := wallet.WriteContract(ctx, client, wallet.WriteContractParameters{
		Address:      "0xFBA3912Ca04dd458c843e2EE08967fC04f3579c2",
		ABI:          abiJSON,
		FunctionName: "mint",
	})

	require.Error(t, err)
	_, ok := err.(*wallet.AccountNotFoundError)
	assert.True(t, ok, "expected AccountNotFoundError")
}

func TestWriteContract_InvalidABI(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	_, err := wallet.WriteContract(ctx, client, wallet.WriteContractParameters{
		Account:      &mockAccount{address: sourceAddr},
		Address:      "0xFBA3912Ca04dd458c843e2EE08967fC04f3579c2",
		ABI:          "invalid-abi",
		FunctionName: "mint",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "ABI")
}

func TestWriteContract_WithValue(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		switch method {
		case "eth_chainId":
			return "0x1"
		case "eth_sendTransaction":
			return "0xwritehash000000000000000000000000000000000000000000000000000002"
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	client.chain = testChain(1)
	ctx := context.Background()

	abiJSON := `[{"inputs":[],"name":"mint","outputs":[],"stateMutability":"payable","type":"function"}]`

	hash, err := wallet.WriteContract(ctx, client, wallet.WriteContractParameters{
		Account:      &mockAccount{address: sourceAddr},
		Address:      "0xFBA3912Ca04dd458c843e2EE08967fC04f3579c2",
		ABI:          abiJSON,
		FunctionName: "mint",
		Value:        big.NewInt(1000000000000000000), // 1 ETH
	})

	require.NoError(t, err)
	assert.NotEmpty(t, hash)
}

// ============================================================================
// DeployContract Tests
// ============================================================================

func TestDeployContract_Default(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		switch method {
		case "eth_chainId":
			return "0x1"
		case "eth_sendTransaction":
			return "0xdeployhash00000000000000000000000000000000000000000000000000001"
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	client.chain = testChain(1)
	client.account = &mockAccount{address: sourceAddr}
	ctx := context.Background()

	abiJSON := `[{"inputs":[],"stateMutability":"nonpayable","type":"constructor"}]`

	hash, err := wallet.DeployContract(ctx, client, wallet.DeployContractParameters{
		ABI:      abiJSON,
		Bytecode: "0x608060405234801561001057600080fd5b50",
	})

	require.NoError(t, err)
	assert.Equal(t, "0xdeployhash00000000000000000000000000000000000000000000000000001", hash)
}

func TestDeployContract_WithConstructorArgs(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		switch method {
		case "eth_chainId":
			return "0x1"
		case "eth_sendTransaction":
			return "0xdeployhash00000000000000000000000000000000000000000000000000002"
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	client.chain = testChain(1)
	client.account = &mockAccount{address: sourceAddr}
	ctx := context.Background()

	abiJSON := `[{"inputs":[{"name":"_value","type":"uint256"}],"stateMutability":"nonpayable","type":"constructor"}]`

	hash, err := wallet.DeployContract(ctx, client, wallet.DeployContractParameters{
		ABI:      abiJSON,
		Bytecode: "0x608060405234801561001057600080fd5b50",
		Args:     []any{big.NewInt(42)},
	})

	require.NoError(t, err)
	assert.NotEmpty(t, hash)
}

func TestDeployContract_InvalidBytecode(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	client.account = &mockAccount{address: sourceAddr}
	ctx := context.Background()

	abiJSON := `[{"inputs":[],"stateMutability":"nonpayable","type":"constructor"}]`

	_, err := wallet.DeployContract(ctx, client, wallet.DeployContractParameters{
		ABI:      abiJSON,
		Bytecode: "not-hex",
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "encode deploy data")
}

// ============================================================================
// GetAddresses Tests
// ============================================================================

func TestGetAddresses_JSONRPC(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_accounts" {
			return []string{
				"0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266",
				"0x70997970c51812dc3a010c7d01b50e0d17dc79c8",
			}
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	addresses, err := wallet.GetAddresses(ctx, client)

	require.NoError(t, err)
	assert.Len(t, addresses, 2)
}

func TestGetAddresses_LocalAccount(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		t.Fatal("RPC should not be called for local account")
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	client.account = &mockLocalAccount{address: sourceAddr}
	ctx := context.Background()

	addresses, err := wallet.GetAddresses(ctx, client)

	require.NoError(t, err)
	assert.Len(t, addresses, 1)
}

// ============================================================================
// RequestAddresses Tests
// ============================================================================

func TestRequestAddresses_Default(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_requestAccounts" {
			return []string{
				"0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266",
			}
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	addresses, err := wallet.RequestAddresses(ctx, client)

	require.NoError(t, err)
	assert.Len(t, addresses, 1)
}

// ============================================================================
// GetPermissions Tests
// ============================================================================

func TestGetPermissions_Default(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		if method == "wallet_getPermissions" {
			return []map[string]any{
				{
					"caveats":          []any{},
					"date":             1234567890,
					"id":               "perm-1",
					"invoker":          "https://example.com",
					"parentCapability": "eth_accounts",
				},
			}
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	permissions, err := wallet.GetPermissions(ctx, client)

	require.NoError(t, err)
	assert.Len(t, permissions, 1)
	assert.Equal(t, "eth_accounts", permissions[0].ParentCapability)
}

// ============================================================================
// RequestPermissions Tests
// ============================================================================

func TestRequestPermissions_Default(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		if method == "wallet_requestPermissions" {
			return []map[string]any{
				{
					"caveats":          []any{},
					"date":             1234567890,
					"id":               "perm-1",
					"invoker":          "https://example.com",
					"parentCapability": "eth_accounts",
				},
			}
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	permissions, err := wallet.RequestPermissions(ctx, client, wallet.RequestPermissionsParameters{
		"eth_accounts": {},
	})

	require.NoError(t, err)
	assert.Len(t, permissions, 1)
	assert.Equal(t, "eth_accounts", permissions[0].ParentCapability)
}

// ============================================================================
// AddChain Tests
// ============================================================================

func TestAddChain_Default(t *testing.T) {
	var capturedMethod string
	var capturedParams []any
	server := createTestServer(t, func(method string, params []any) any {
		capturedMethod = method
		capturedParams = params
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	optimism := &chain.Chain{
		ID:   10,
		Name: "OP Mainnet",
		NativeCurrency: chain.ChainNativeCurrency{
			Name:     "Ether",
			Symbol:   "ETH",
			Decimals: 18,
		},
		RpcUrls: map[string]chain.ChainRpcUrls{
			"default": {HTTP: []string{"https://mainnet.optimism.io"}},
		},
	}

	err := wallet.AddChain(ctx, client, wallet.AddChainParameters{
		Chain: optimism,
	})

	require.NoError(t, err)
	assert.Equal(t, "wallet_addEthereumChain", capturedMethod)
	require.NotEmpty(t, capturedParams)
}

func TestAddChain_NilChain(t *testing.T) {
	client := &mockClient{}
	ctx := context.Background()

	err := wallet.AddChain(ctx, client, wallet.AddChainParameters{
		Chain: nil,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "chain is required")
}

// ============================================================================
// SwitchChain Tests
// ============================================================================

func TestSwitchChain_Default(t *testing.T) {
	var capturedMethod string
	server := createTestServer(t, func(method string, params []any) any {
		capturedMethod = method
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	err := wallet.SwitchChain(ctx, client, wallet.SwitchChainParameters{
		ID: 10, // Optimism
	})

	require.NoError(t, err)
	assert.Equal(t, "wallet_switchEthereumChain", capturedMethod)
}

func TestSwitchChain_RPCError(t *testing.T) {
	server := createErrorTestServer(t, 4902, "Unrecognized chain ID")
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	err := wallet.SwitchChain(ctx, client, wallet.SwitchChainParameters{
		ID: 999999,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "wallet_switchEthereumChain failed")
}

// ============================================================================
// WatchAsset Tests
// ============================================================================

func TestWatchAsset_Default(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		if method == "wallet_watchAsset" {
			return true
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	success, err := wallet.WatchAsset(ctx, client, wallet.WatchAssetParameters{
		Type: "ERC20",
		Options: wallet.WatchAssetOptions{
			Address:  "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			Decimals: 18,
			Symbol:   "WETH",
		},
	})

	require.NoError(t, err)
	assert.True(t, success)
}

func TestWatchAsset_Rejected(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		if method == "wallet_watchAsset" {
			return false
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	success, err := wallet.WatchAsset(ctx, client, wallet.WatchAssetParameters{
		Type: "ERC20",
		Options: wallet.WatchAssetOptions{
			Address:  "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
			Decimals: 18,
			Symbol:   "WETH",
		},
	})

	require.NoError(t, err)
	assert.False(t, success)
}

// ============================================================================
// GetCapabilities Tests
// ============================================================================

func TestGetCapabilities_Default(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		if method == "wallet_getCapabilities" {
			return map[string]map[string]any{
				"0x1": {
					"paymasterService": map[string]any{
						"supported": true,
					},
				},
				"0xa": {
					"paymasterService": map[string]any{
						"supported": true,
					},
				},
			}
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	client.account = &mockAccount{address: sourceAddr}
	ctx := context.Background()

	caps, err := wallet.GetCapabilities(ctx, client, wallet.GetCapabilitiesParameters{})

	require.NoError(t, err)
	assert.NotNil(t, caps)
	assert.Contains(t, caps, int64(1))
	assert.Contains(t, caps, int64(10))
}

func TestGetCapabilities_NormalizesAddSubAccount(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		if method == "wallet_getCapabilities" {
			return map[string]map[string]any{
				"0x1": {
					"addSubAccount": map[string]any{
						"supported": true,
					},
				},
			}
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	client.account = &mockAccount{address: sourceAddr}
	ctx := context.Background()

	caps, err := wallet.GetCapabilities(ctx, client, wallet.GetCapabilitiesParameters{})

	require.NoError(t, err)
	assert.Contains(t, caps[1], "unstable_addSubAccount")
}

// ============================================================================
// SendCalls Tests
// ============================================================================

func TestSendCalls_Default(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		if method == "wallet_sendCalls" {
			return map[string]any{
				"id": "0xcallbatch123",
			}
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	client.chain = testChain(1)
	client.account = &mockAccount{address: sourceAddr}
	ctx := context.Background()

	result, err := wallet.SendCalls(ctx, client, wallet.SendCallsParameters{
		Calls: []wallet.Call{
			{Data: "0xdeadbeef", To: targetAddr.Hex()},
			{To: targetAddr.Hex(), Value: big.NewInt(69420)},
		},
	})

	require.NoError(t, err)
	assert.Equal(t, "0xcallbatch123", result.ID)
}

func TestSendCalls_NoChain(t *testing.T) {
	client := &mockClient{
		account: &mockAccount{address: sourceAddr},
	}
	ctx := context.Background()

	_, err := wallet.SendCalls(ctx, client, wallet.SendCallsParameters{
		Calls: []wallet.Call{
			{To: targetAddr.Hex()},
		},
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "chain is required")
}

func TestSendCalls_DefaultVersion(t *testing.T) {
	var capturedParams []any
	server := createTestServer(t, func(method string, params []any) any {
		if method == "wallet_sendCalls" {
			capturedParams = params
			return map[string]any{"id": "0x123"}
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	client.chain = testChain(1)
	client.account = &mockAccount{address: sourceAddr}
	ctx := context.Background()

	_, err := wallet.SendCalls(ctx, client, wallet.SendCallsParameters{
		Calls: []wallet.Call{
			{To: targetAddr.Hex()},
		},
	})

	require.NoError(t, err)
	require.NotEmpty(t, capturedParams)
	// The version should default to "2.0.0"
	if paramMap, ok := capturedParams[0].(map[string]any); ok {
		assert.Equal(t, "2.0.0", paramMap["version"])
	}
}

// ============================================================================
// GetCallsStatus Tests
// ============================================================================

func TestGetCallsStatus_Default(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		if method == "wallet_getCallsStatus" {
			return map[string]any{
				"atomic":     true,
				"chainId":    "0x1",
				"status":     float64(200),
				"version":    "2.0.0",
				"receipts":   []any{},
			}
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	status, err := wallet.GetCallsStatus(ctx, client, wallet.GetCallsStatusParameters{
		ID: "0xdeadbeef",
	})

	require.NoError(t, err)
	assert.Equal(t, 200, status.StatusCode)
	assert.Equal(t, "success", status.Status)
	assert.True(t, status.Atomic)
}

func TestGetCallsStatus_Pending(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		if method == "wallet_getCallsStatus" {
			return map[string]any{
				"atomic":   false,
				"status":   float64(100),
				"version":  "2.0.0",
				"receipts": []any{},
			}
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	status, err := wallet.GetCallsStatus(ctx, client, wallet.GetCallsStatusParameters{
		ID: "0xpending",
	})

	require.NoError(t, err)
	assert.Equal(t, 100, status.StatusCode)
	assert.Equal(t, "pending", status.Status)
}

// ============================================================================
// ShowCallsStatus Tests
// ============================================================================

func TestShowCallsStatus_Default(t *testing.T) {
	var capturedMethod string
	server := createTestServer(t, func(method string, params []any) any {
		capturedMethod = method
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	err := wallet.ShowCallsStatus(ctx, client, wallet.ShowCallsStatusParameters{
		ID: "0xdeadbeef",
	})

	require.NoError(t, err)
	assert.Equal(t, "wallet_showCallsStatus", capturedMethod)
}

// ============================================================================
// PrepareAuthorization Tests
// ============================================================================

func TestPrepareAuthorization_Default(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		switch method {
		case "eth_chainId":
			return "0x1"
		case "eth_getTransactionCount":
			return "0x5"
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	client.chain = testChain(1)
	ctx := context.Background()

	auth, err := wallet.PrepareAuthorization(ctx, client, wallet.PrepareAuthorizationParameters{
		Account:         &mockAccount{address: sourceAddr},
		ContractAddress: "0xA0Cf798816D4b9b9866b5330EEa46a18382f251e",
	})

	require.NoError(t, err)
	assert.Equal(t, "0xA0Cf798816D4b9b9866b5330EEa46a18382f251e", auth.Address)
	assert.Equal(t, 1, auth.ChainId)
	assert.Equal(t, 5, auth.Nonce) // From eth_getTransactionCount
}

func TestPrepareAuthorization_WithChainID(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_getTransactionCount" {
			return "0xa"
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	chainID := 42
	auth, err := wallet.PrepareAuthorization(ctx, client, wallet.PrepareAuthorizationParameters{
		Account:         &mockAccount{address: sourceAddr},
		ContractAddress: "0xA0Cf798816D4b9b9866b5330EEa46a18382f251e",
		ChainID:         &chainID,
	})

	require.NoError(t, err)
	assert.Equal(t, 42, auth.ChainId)
}

func TestPrepareAuthorization_WithNonce(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_chainId" {
			return "0x1"
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	client.chain = testChain(1)
	ctx := context.Background()

	nonce := 99
	auth, err := wallet.PrepareAuthorization(ctx, client, wallet.PrepareAuthorizationParameters{
		Account:         &mockAccount{address: sourceAddr},
		ContractAddress: "0xA0Cf798816D4b9b9866b5330EEa46a18382f251e",
		Nonce:           &nonce,
	})

	require.NoError(t, err)
	assert.Equal(t, 99, auth.Nonce)
}

func TestPrepareAuthorization_SelfExecutor(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		switch method {
		case "eth_chainId":
			return "0x1"
		case "eth_getTransactionCount":
			return "0x5"
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	client.chain = testChain(1)
	ctx := context.Background()

	auth, err := wallet.PrepareAuthorization(ctx, client, wallet.PrepareAuthorizationParameters{
		Account:         &mockAccount{address: sourceAddr},
		ContractAddress: "0xA0Cf798816D4b9b9866b5330EEa46a18382f251e",
		Executor:        "self",
	})

	require.NoError(t, err)
	// Nonce should be incremented by 1 for self-executor
	assert.Equal(t, 6, auth.Nonce) // 5 + 1
}

func TestPrepareAuthorization_NoAccount(t *testing.T) {
	client := &mockClient{}
	ctx := context.Background()

	_, err := wallet.PrepareAuthorization(ctx, client, wallet.PrepareAuthorizationParameters{
		ContractAddress: "0xA0Cf798816D4b9b9866b5330EEa46a18382f251e",
	})

	require.Error(t, err)
	_, ok := err.(*wallet.AccountNotFoundError)
	assert.True(t, ok, "expected AccountNotFoundError")
}

// ============================================================================
// SignAuthorization Tests
// ============================================================================

func TestSignAuthorization_Default(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		switch method {
		case "eth_chainId":
			return "0x1"
		case "eth_getTransactionCount":
			return "0x5"
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	client.chain = testChain(1)
	ctx := context.Background()

	localAccount := &mockAuthorizationSignableAccount{
		address: sourceAddr,
		signFn: func(auth types.AuthorizationRequest) (*types.SignedAuthorization, error) {
			return &types.SignedAuthorization{
				Address: auth.Address,
				ChainId: auth.ChainId,
				Nonce:   auth.Nonce,
				R:       "0xabc",
				S:       "0xdef",
				YParity: 0,
			}, nil
		},
	}

	signed, err := wallet.SignAuthorization(ctx, client, wallet.SignAuthorizationParameters{
		Account:         localAccount,
		ContractAddress: "0xA0Cf798816D4b9b9866b5330EEa46a18382f251e",
	})

	require.NoError(t, err)
	assert.Equal(t, "0xA0Cf798816D4b9b9866b5330EEa46a18382f251e", signed.Address)
	assert.Equal(t, 1, signed.ChainId)
	assert.Equal(t, "0xabc", signed.R)
	assert.Equal(t, "0xdef", signed.S)
}

func TestSignAuthorization_NoAccount(t *testing.T) {
	client := &mockClient{}
	ctx := context.Background()

	_, err := wallet.SignAuthorization(ctx, client, wallet.SignAuthorizationParameters{
		ContractAddress: "0xA0Cf798816D4b9b9866b5330EEa46a18382f251e",
	})

	require.Error(t, err)
	_, ok := err.(*wallet.AccountNotFoundError)
	assert.True(t, ok, "expected AccountNotFoundError")
}

func TestSignAuthorization_NonLocalAccount(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	// JSON-RPC account (not local, doesn't implement AuthorizationSignableAccount)
	_, err := wallet.SignAuthorization(ctx, client, wallet.SignAuthorizationParameters{
		Account:         &mockAccount{address: sourceAddr},
		ContractAddress: "0xA0Cf798816D4b9b9866b5330EEa46a18382f251e",
	})

	require.Error(t, err)
	_, ok := err.(*wallet.AccountTypeNotSupportedError)
	assert.True(t, ok, "expected AccountTypeNotSupportedError, got %T: %v", err, err)
}

// ============================================================================
// SendRawTransactionSync Tests
// ============================================================================

func TestSendRawTransactionSync_Default(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_sendRawTransactionSync" {
			return map[string]any{
				"transactionHash":   "0xhash00000000000000000000000000000000000000000000000000000000001",
				"transactionIndex":  "0x0",
				"blockHash":         "0x1234567890123456789012345678901234567890123456789012345678901234",
				"blockNumber":       "0x10",
				"from":              sourceAddr.Hex(),
				"to":                targetAddr.Hex(),
				"cumulativeGasUsed": "0x5208",
				"gasUsed":           "0x5208",
				"contractAddress":   nil,
				"logs":              []any{},
				"status":            "0x1",
				"logsBloom":         "0x00",
				"effectiveGasPrice": "0x3b9aca00",
				"type":              "0x2",
			}
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	receipt, err := wallet.SendRawTransactionSync(ctx, client, wallet.SendRawTransactionSyncParameters{
		SerializedTransaction: "0x02f850018203118080825208808080c080a04012522854168b27e5dc3d5839bab5e6b39e1a0ffd343901ce1622e3d64b48f1a04e00902ae0502c4728cbf12156290df99c3ed7de85b1dbfe20b5c36931733a33",
	})

	require.NoError(t, err)
	assert.NotNil(t, receipt)
}

func TestSendRawTransactionSync_RevertedThrows(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_sendRawTransactionSync" {
			return map[string]any{
				"transactionHash": "0xhash00000000000000000000000000000000000000000000000000000000001",
				"blockHash":       "0x1234567890123456789012345678901234567890123456789012345678901234",
				"blockNumber":     "0x10",
				"status":          "0x0", // reverted
				"logsBloom":       "0x00",
				"type":            "0x2",
				"logs":            []any{},
			}
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	_, err := wallet.SendRawTransactionSync(ctx, client, wallet.SendRawTransactionSyncParameters{
		SerializedTransaction: "0x02f850",
	})

	require.Error(t, err)
	_, ok := err.(*wallet.TransactionReceiptRevertedError)
	assert.True(t, ok, "expected TransactionReceiptRevertedError, got %T: %v", err, err)
}

func TestSendRawTransactionSync_RevertedNoThrow(t *testing.T) {
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_sendRawTransactionSync" {
			return map[string]any{
				"transactionHash": "0xhash00000000000000000000000000000000000000000000000000000000001",
				"blockHash":       "0x1234567890123456789012345678901234567890123456789012345678901234",
				"blockNumber":     "0x10",
				"status":          "0x0",
				"logsBloom":       "0x00",
				"type":            "0x2",
				"logs":            []any{},
			}
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	throwOnRevert := false
	receipt, err := wallet.SendRawTransactionSync(ctx, client, wallet.SendRawTransactionSyncParameters{
		SerializedTransaction: "0x02f850",
		ThrowOnReceiptRevert:  &throwOnRevert,
	})

	require.NoError(t, err)
	assert.NotNil(t, receipt)
}

func TestSendRawTransactionSync_WithTimeout(t *testing.T) {
	var capturedParams []any
	server := createTestServer(t, func(method string, params []any) any {
		if method == "eth_sendRawTransactionSync" {
			capturedParams = params
			return map[string]any{
				"transactionHash": "0xhash00000000000000000000000000000000000000000000000000000000001",
				"blockHash":       "0x1234567890123456789012345678901234567890123456789012345678901234",
				"blockNumber":     "0x10",
				"status":          "0x1",
				"logsBloom":       "0x00",
				"type":            "0x2",
				"logs":            []any{},
			}
		}
		return nil
	})
	defer server.Close()

	client := createMockClient(t, server.URL)
	ctx := context.Background()

	timeout := int64(10000)
	receipt, err := wallet.SendRawTransactionSync(ctx, client, wallet.SendRawTransactionSyncParameters{
		SerializedTransaction: "0x02f850",
		Timeout:               &timeout,
	})

	require.NoError(t, err)
	assert.NotNil(t, receipt)
	// Should have two params: serializedTx and timeout hex
	require.Len(t, capturedParams, 2)
}

// ============================================================================
// Error Type Tests
// ============================================================================

func TestAccountNotFoundError_Message(t *testing.T) {
	err := &wallet.AccountNotFoundError{DocsPath: "/docs/actions/wallet/sendTransaction"}
	assert.Contains(t, err.Error(), "could not find an Account")
	assert.Contains(t, err.Error(), "https://viem.sh/docs/actions/wallet/sendTransaction")
}

func TestAccountNotFoundError_NoDocs(t *testing.T) {
	err := &wallet.AccountNotFoundError{}
	assert.Contains(t, err.Error(), "could not find an Account")
	assert.NotContains(t, err.Error(), "https://viem.sh")
}

func TestAccountTypeNotSupportedError_Message(t *testing.T) {
	err := &wallet.AccountTypeNotSupportedError{
		DocsPath: "/docs/eip7702/signAuthorization",
		MetaMessages: []string{
			"The `signAuthorization` Action does not support JSON-RPC Accounts.",
		},
	}
	assert.Contains(t, err.Error(), "account type not supported")
	assert.Contains(t, err.Error(), "does not support JSON-RPC Accounts")
	assert.Contains(t, err.Error(), "https://viem.sh/docs/eip7702/signAuthorization")
}

func TestTransactionReceiptRevertedError_Message(t *testing.T) {
	err := &wallet.TransactionReceiptRevertedError{
		Receipt: &formatters.TransactionReceipt{
			TransactionHash: "0xdeadbeef",
		},
	}
	assert.Contains(t, err.Error(), "transaction reverted")
	assert.Contains(t, err.Error(), "0xdeadbeef")
}
