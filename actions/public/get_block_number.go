package public

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"
)

// GetBlockNumberParameters contains the parameters for the GetBlockNumber action.
// This mirrors viem's GetBlockNumberParameters type.
type GetBlockNumberParameters struct {
	// CacheTime is the time (in duration) that cached block number will remain in memory.
	// If nil, uses the client's cache time.
	CacheTime *time.Duration
}

// GetBlockNumberReturnType is the return type for the GetBlockNumber action.
type GetBlockNumberReturnType = uint64

// blockNumberCache is a simple cache for block numbers.
var (
	blockNumberCacheMu   sync.RWMutex
	blockNumberCacheData = make(map[string]cachedBlockNumber)
)

type cachedBlockNumber struct {
	blockNumber uint64
	expiresAt   time.Time
}

// GetBlockNumber returns the number of the most recent block seen.
//
// This is equivalent to viem's `getBlockNumber` action.
//
// JSON-RPC Method: eth_blockNumber
//
// Example:
//
//	blockNumber, err := public.GetBlockNumber(ctx, client, public.GetBlockNumberParameters{})
//	// blockNumber is the latest block number
func GetBlockNumber(ctx context.Context, client Client, params GetBlockNumberParameters) (GetBlockNumberReturnType, error) {
	// Determine cache time
	cacheTime := client.CacheTime()
	if params.CacheTime != nil {
		cacheTime = *params.CacheTime
	}

	// Check cache
	cacheKey := fmt.Sprintf("blockNumber.%s", client.UID())
	if cacheTime > 0 {
		blockNumberCacheMu.RLock()
		if cached, ok := blockNumberCacheData[cacheKey]; ok && time.Now().Before(cached.expiresAt) {
			blockNumberCacheMu.RUnlock()
			return cached.blockNumber, nil
		}
		blockNumberCacheMu.RUnlock()
	}

	// Execute the request
	resp, err := client.Request(ctx, "eth_blockNumber")
	if err != nil {
		return 0, fmt.Errorf("eth_blockNumber failed: %w", err)
	}

	var hexBlockNumber string
	if unmarshalErr := json.Unmarshal(resp.Result, &hexBlockNumber); unmarshalErr != nil {
		return 0, fmt.Errorf("failed to unmarshal block number: %w", unmarshalErr)
	}

	// Parse the block number
	blockNumber, err := parseHexUint64(hexBlockNumber)
	if err != nil {
		return 0, fmt.Errorf("failed to parse block number: %w", err)
	}

	// Cache the result
	if cacheTime > 0 {
		blockNumberCacheMu.Lock()
		blockNumberCacheData[cacheKey] = cachedBlockNumber{
			blockNumber: blockNumber,
			expiresAt:   time.Now().Add(cacheTime),
		}
		blockNumberCacheMu.Unlock()
	}

	return blockNumber, nil
}

// parseHexUint64 parses a hex string to uint64.
func parseHexUint64(hexStr string) (uint64, error) {
	hexStr = strings.TrimPrefix(hexStr, "0x")
	if hexStr == "" {
		return 0, nil
	}

	n := new(big.Int)
	_, ok := n.SetString(hexStr, 16)
	if !ok {
		return 0, fmt.Errorf("invalid hex string: %s", hexStr)
	}
	return n.Uint64(), nil
}
