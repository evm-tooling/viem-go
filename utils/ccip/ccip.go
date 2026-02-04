// Package ccip provides CCIP-Read (EIP-3668) offchain lookup support.
package ccip

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ChefBingbong/viem-go/abi"
)

// OffchainLookupSignature is the error selector for the OffchainLookup error.
// This is the 4-byte selector for:
// error OffchainLookup(address sender, string[] urls, bytes callData, bytes4 callbackFunction, bytes extraData)
const OffchainLookupSignature = "0x556f1830"

// OffchainLookupAbi is the ABI definition for the OffchainLookup error.
var OffchainLookupAbi = []abi.AbiParam{
	{Name: "sender", Type: "address"},
	{Name: "urls", Type: "string[]"},
	{Name: "callData", Type: "bytes"},
	{Name: "callbackFunction", Type: "bytes4"},
	{Name: "extraData", Type: "bytes"},
}

// OffchainLookupError represents a decoded OffchainLookup error.
type OffchainLookupError struct {
	Sender           common.Address
	URLs             []string
	CallData         []byte
	CallbackFunction [4]byte
	ExtraData        []byte
}

// CCIPRequestParams contains parameters for a CCIP gateway request.
type CCIPRequestParams struct {
	Data   []byte
	Sender common.Address
	URLs   []string
}

// ErrOffchainLookup is returned when an offchain lookup fails.
type ErrOffchainLookup struct {
	CallbackSelector [4]byte
	Data             []byte
	ExtraData        []byte
	Sender           common.Address
	URLs             []string
	Cause            error
}

func (e *ErrOffchainLookup) Error() string {
	return fmt.Sprintf("offchain lookup failed: %v", e.Cause)
}

func (e *ErrOffchainLookup) Unwrap() error {
	return e.Cause
}

// ErrOffchainLookupSenderMismatch is returned when the sender doesn't match.
type ErrOffchainLookupSenderMismatch struct {
	Sender common.Address
	To     common.Address
}

func (e *ErrOffchainLookupSenderMismatch) Error() string {
	return fmt.Sprintf("offchain lookup sender mismatch: expected %s, got %s", e.To.Hex(), e.Sender.Hex())
}

// ErrOffchainLookupResponseMalformed is returned when the gateway response is malformed.
type ErrOffchainLookupResponseMalformed struct {
	Result string
	URL    string
}

func (e *ErrOffchainLookupResponseMalformed) Error() string {
	return fmt.Sprintf("offchain lookup response malformed from %s: %s", e.URL, e.Result)
}

// DecodeOffchainLookupError decodes the OffchainLookup error from revert data.
func DecodeOffchainLookupError(data []byte) (*OffchainLookupError, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("data too short")
	}

	// Check selector
	selector := hexutil.Encode(data[:4])
	if selector != OffchainLookupSignature {
		return nil, fmt.Errorf("not an OffchainLookup error: got selector %s", selector)
	}

	// Decode parameters (skip selector)
	decoded, err := abi.DecodeAbiParameters(OffchainLookupAbi, data[4:])
	if err != nil {
		return nil, fmt.Errorf("failed to decode OffchainLookup: %w", err)
	}

	if len(decoded) < 5 {
		return nil, fmt.Errorf("decoded too few parameters")
	}

	// Extract values
	sender, ok := decoded[0].(common.Address)
	if !ok {
		return nil, fmt.Errorf("invalid sender type")
	}

	urlsAny, ok := decoded[1].([]any)
	if !ok {
		return nil, fmt.Errorf("invalid urls type")
	}
	urls := make([]string, len(urlsAny))
	for i, u := range urlsAny {
		urls[i], _ = u.(string)
	}

	callData, ok := decoded[2].([]byte)
	if !ok {
		return nil, fmt.Errorf("invalid callData type")
	}

	callbackBytes, ok := decoded[3].([]byte)
	if !ok {
		return nil, fmt.Errorf("invalid callbackFunction type")
	}
	var callbackFunction [4]byte
	copy(callbackFunction[:], callbackBytes)

	extraData, ok := decoded[4].([]byte)
	if !ok {
		return nil, fmt.Errorf("invalid extraData type")
	}

	return &OffchainLookupError{
		Sender:           sender,
		URLs:             urls,
		CallData:         callData,
		CallbackFunction: callbackFunction,
		ExtraData:        extraData,
	}, nil
}

// CCIPRequest makes a request to CCIP gateway URLs.
// It tries each URL in order until one succeeds.
func CCIPRequest(ctx context.Context, params CCIPRequestParams) ([]byte, error) {
	var lastErr error

	for _, url := range params.URLs {
		// Determine method based on URL format
		method := "POST"
		if strings.Contains(url, "{data}") {
			method = "GET"
		}

		// Build URL
		requestURL := url
		requestURL = strings.Replace(requestURL, "{sender}", strings.ToLower(params.Sender.Hex()), -1)
		requestURL = strings.Replace(requestURL, "{data}", hexutil.Encode(params.Data), -1)

		// Make request
		var req *http.Request
		if method == "GET" {
			var err error
			req, err = http.NewRequestWithContext(ctx, "GET", requestURL, nil)
			if err != nil {
				lastErr = err
				continue
			}
		} else {
			body := map[string]string{
				"data":   hexutil.Encode(params.Data),
				"sender": params.Sender.Hex(),
			}
			bodyBytes, _ := json.Marshal(body)
			var err error
			req, err = http.NewRequestWithContext(ctx, "POST", requestURL, bytes.NewReader(bodyBytes))
			if err != nil {
				lastErr = err
				continue
			}
			req.Header.Set("Content-Type", "application/json")
		}

		// Execute request
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		// Read response
		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = err
			continue
		}

		// Parse response
		var result string
		contentType := resp.Header.Get("Content-Type")
		if strings.HasPrefix(contentType, "application/json") {
			var jsonResp struct {
				Data  string `json:"data"`
				Error struct {
					Message string `json:"message"`
				} `json:"error"`
			}
			if err := json.Unmarshal(respBody, &jsonResp); err != nil {
				lastErr = err
				continue
			}
			if jsonResp.Error.Message != "" {
				lastErr = fmt.Errorf("gateway error: %s", jsonResp.Error.Message)
				continue
			}
			result = jsonResp.Data
		} else {
			result = string(respBody)
		}

		// Check for HTTP errors
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, result)
			continue
		}

		// Validate result is hex
		if !strings.HasPrefix(result, "0x") {
			lastErr = &ErrOffchainLookupResponseMalformed{Result: result, URL: url}
			continue
		}

		return common.FromHex(result), nil
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return nil, fmt.Errorf("no URLs provided")
}

// BuildCallbackData builds the callback calldata for CCIP-Read.
func BuildCallbackData(callbackSelector [4]byte, result, extraData []byte) ([]byte, error) {
	// Encode (bytes result, bytes extraData)
	encoded, err := abi.EncodeAbiParameters(
		[]abi.AbiParam{{Type: "bytes"}, {Type: "bytes"}},
		[]any{result, extraData},
	)
	if err != nil {
		return nil, err
	}

	// Concatenate selector + encoded args
	data := make([]byte, 4+len(encoded))
	copy(data[:4], callbackSelector[:])
	copy(data[4:], encoded)

	return data, nil
}
