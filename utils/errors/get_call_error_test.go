package errors

import (
	"errors"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ChefBingbong/viem-go/utils/rpc"
)

func TestCallExecutionError_Format(t *testing.T) {
	to := common.HexToAddress("0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359")
	data := hexutil.MustDecode("0x06fdde03")
	chainID := int64(137)

	// Simulate HTTP error (like polygon 401)
	httpErr := rpc.NewHTTPRequestErrorWithDetails(
		"https://polygon-rpc.com",
		401,
		"Unauthorized",
		map[string]any{"method": "eth_call", "params": []any{}},
		"message: API key disabled, reason: tenant disabled",
	)

	err := GetCallError(httpErr, CallErrorParams{
		To:      &to,
		Data:    data,
		ChainID: &chainID,
	})

	got := err.Error()

	// Verify error type prefix (mirrors viem's CallExecutionError)
	if !strings.Contains(got, "CallExecutionError:") {
		t.Errorf("Error() should contain 'CallExecutionError:' prefix, got:\n%s", got)
	}
	if !strings.Contains(got, "HTTP request failed") {
		t.Errorf("Error() should contain cause message, got:\n%s", got)
	}
	if !strings.Contains(got, "Raw Call Arguments:") {
		t.Errorf("Error() should contain Raw Call Arguments section, got:\n%s", got)
	}
	if !strings.Contains(got, "to:") {
		t.Errorf("Error() should contain 'to' in raw args, got:\n%s", got)
	}
	if !strings.Contains(got, "data:") {
		t.Errorf("Error() should contain 'data' in raw args, got:\n%s", got)
	}
	if !strings.Contains(got, "chainId: 137") {
		t.Errorf("Error() should contain chainId, got:\n%s", got)
	}
	if !strings.Contains(got, "Details:") {
		t.Errorf("Error() should contain Details from cause, got:\n%s", got)
	}
}

func TestCallExecutionError_WithFromAndValue(t *testing.T) {
	from := common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
	to := common.HexToAddress("0x1234567890123456789012345678901234567890")
	chainID := int64(137)
	value := mustBigInt("100000000000000000") // 0.1 ETH

	err := GetCallError(
		rpc.NewHTTPRequestError("https://rpc.example.com", 500, "Internal Server Error", nil, nil),
		CallErrorParams{
			From:                 &from,
			To:                   &to,
			Value:                value,
			ChainID:              &chainID,
			NativeCurrencySymbol: "POL",
		},
	)

	got := err.Error()
	if !strings.Contains(got, "from:") {
		t.Errorf("Error() should contain 'from' when Account set, got:\n%s", got)
	}
	if !strings.Contains(got, "value:") {
		t.Errorf("Error() should contain 'value' when Value set, got:\n%s", got)
	}
	if !strings.Contains(got, "POL") {
		t.Errorf("Error() should contain native currency symbol, got:\n%s", got)
	}
}

func TestCallExecutionError_Unwrap(t *testing.T) {
	cause := rpc.NewHTTPRequestError("https://x.com", 401, "", nil, nil)
	err := GetCallError(cause, CallErrorParams{})
	if !isSameError(err, cause) {
		t.Error("Unwrap() should return the cause")
	}
}

func mustBigInt(s string) *big.Int {
	v := new(big.Int)
	v.SetString(s, 10)
	return v
}

func isSameError(err, target error) bool {
	return errors.Unwrap(err) == target
}
