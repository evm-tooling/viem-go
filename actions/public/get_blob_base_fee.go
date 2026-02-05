package public

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
)

// GetBlobBaseFeeReturnType is the return type for the GetTransactionCount action.
// It represents the tx count in wei.
type GetBlobBaseFeeReturnType = *big.Int


func GetBlobBaseFee(ctx context.Context, client Client) (GetBlobBaseFeeReturnType, error) {
	// Execute the request
	resp, err := client.Request(ctx, "eth_blobBaseFee")
	if err != nil {
		return nil, fmt.Errorf("eth_blobBaseFee failed: %w", err)
	}

	var blobFeeHex string
	if unmarshalErr := json.Unmarshal(resp.Result, &blobFeeHex); unmarshalErr != nil {
		return nil, fmt.Errorf("failed to unmarshal blob base fee: %w", unmarshalErr)
	}

	// Parse the tx count
	blobFee, err := parseHexBigInt(blobFeeHex)
	if err != nil {
		return nil, fmt.Errorf("failed to parseblob base fee: %w", err)
	}

	return blobFee, nil
}
