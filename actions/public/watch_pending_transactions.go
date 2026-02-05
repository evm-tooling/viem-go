package public

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/client/transport"
	"github.com/ChefBingbong/viem-go/utils/batch"
	"github.com/ChefBingbong/viem-go/utils/observe"
	"github.com/ChefBingbong/viem-go/utils/poll"
)

// WatchPendingTransactionsParameters contains the parameters for the WatchPendingTransactions action.
// This mirrors viem's WatchPendingTransactionsParameters type.
type WatchPendingTransactionsParameters struct {
	// Batch determines whether to batch pending transaction hashes together.
	// When true, multiple hashes are collected and emitted together.
	// When false, each hash is emitted individually.
	// Default: true
	Batch bool

	// Poll forces polling mode even when WebSocket transport is available.
	// If nil, automatically detects based on transport type.
	Poll *bool

	// PollingInterval is the interval between polls when using polling mode.
	// If zero, uses the client's default polling interval.
	PollingInterval time.Duration
}

// WatchPendingTransactionsEvent represents an event from WatchPendingTransactions.
type WatchPendingTransactionsEvent struct {
	// Hashes are the pending transaction hashes.
	// When Batch is true, this may contain multiple hashes.
	// When Batch is false, this will contain a single hash.
	Hashes []common.Hash

	// Error is any error that occurred.
	Error error
}

// pendingTxObserver is the global observer for pending transaction subscriptions.
var pendingTxObserver = observe.New[WatchPendingTransactionsEvent]()

// WatchPendingTransactions watches and returns pending transaction hashes.
//
// This is equivalent to viem's `watchPendingTransactions` action with full Go optimization:
//   - Channels instead of callbacks for native Go concurrency
//   - context.Context for cancellation
//   - Observer pattern for deduplication (multiple watchers share one source)
//   - BatchCollector for efficient batching
//   - Automatic transport detection (polling vs subscription)
//
// JSON-RPC Methods:
//   - When polling:
//   - Calls eth_newPendingTransactionFilter to initialize the filter
//   - Calls eth_getFilterChanges on a polling interval
//   - When subscribing: uses eth_subscribe with "newPendingTransactions" event
//
// Example:
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	events := public.WatchPendingTransactions(ctx, client, public.WatchPendingTransactionsParameters{
//	    Batch: true,
//	})
//
//	for event := range events {
//	    if event.Error != nil {
//	        log.Printf("error: %v", event.Error)
//	        continue
//	    }
//	    fmt.Printf("Pending transactions: %d\n", len(event.Hashes))
//	    for _, hash := range event.Hashes {
//	        fmt.Printf("  - %s\n", hash.Hex())
//	    }
//	}
func WatchPendingTransactions(
	ctx context.Context,
	client WatchClient,
	params WatchPendingTransactionsParameters,
) <-chan WatchPendingTransactionsEvent {
	// Default batch to true
	batchMode := params.Batch

	// Determine if we should poll or subscribe
	enablePolling := ShouldPoll(client, params.Poll)

	// Get polling interval
	pollingInterval := GetPollingInterval(client, params.PollingInterval)

	// Create output channel
	ch := make(chan WatchPendingTransactionsEvent, 10)

	go func() {
		defer close(ch)

		if enablePolling {
			pollPendingTransactions(ctx, client, batchMode, pollingInterval, ch)
		} else {
			subscribePendingTransactions(ctx, client, batchMode, ch)
		}
	}()

	return ch
}

// pollPendingTransactions implements pending transaction watching using polling.
func pollPendingTransactions(
	ctx context.Context,
	client WatchClient,
	batchMode bool,
	interval time.Duration,
	ch chan<- WatchPendingTransactionsEvent,
) {
	// Create observer ID for deduplication
	observerID := fmt.Sprintf("watchPendingTransactions.%s.%v.%v",
		client.UID(),
		batchMode,
		interval,
	)

	// Use observer for deduplication
	eventCh := pendingTxObserver.Subscribe(observerID, func() (<-chan WatchPendingTransactionsEvent, func()) {
		sourceCh := make(chan WatchPendingTransactionsEvent, 100)

		// Create filter
		var filterID FilterID

		// Start polling
		pollResults := poll.Poll(ctx, func(ctx context.Context) ([]common.Hash, error) {
			// Create filter on first poll
			if filterID == "" {
				filter, err := CreatePendingTransactionFilter(ctx, client)
				if err != nil {
					return nil, err
				}
				filterID = filter.ID
				return nil, nil // First poll just creates filter
			}

			// Get filter changes
			return GetFilterChangesTransactions(ctx, client, filterID)
		}, poll.Options{
			Interval:    interval,
			EmitOnBegin: true,
		})

		// Process poll results
		go func() {
			defer close(sourceCh)
			defer func() {
				// Cleanup: uninstall filter
				if filterID != "" {
					_, _ = UninstallFilter(context.Background(), client, filterID)
				}
			}()

			for result := range pollResults {
				if result.Error != nil {
					select {
					case sourceCh <- WatchPendingTransactionsEvent{Error: result.Error}:
					case <-ctx.Done():
						return
					}
					continue
				}

				hashes := result.Value
				if len(hashes) == 0 {
					continue
				}

				// Emit hashes
				if batchMode {
					select {
					case sourceCh <- WatchPendingTransactionsEvent{Hashes: hashes}:
					case <-ctx.Done():
						return
					}
				} else {
					// Emit individually
					for _, hash := range hashes {
						select {
						case sourceCh <- WatchPendingTransactionsEvent{Hashes: []common.Hash{hash}}:
						case <-ctx.Done():
							return
						}
					}
				}
			}
		}()

		return sourceCh, func() {
			if filterID != "" {
				_, _ = UninstallFilter(context.Background(), client, filterID)
			}
		}
	})

	// Forward events to output channel
	for event := range eventCh {
		select {
		case ch <- event:
		case <-ctx.Done():
			return
		}
	}
}

// subscribePendingTransactions implements pending transaction watching using WebSocket subscription.
func subscribePendingTransactions(
	ctx context.Context,
	client WatchClient,
	batchMode bool,
	ch chan<- WatchPendingTransactionsEvent,
) {
	if batchMode {
		// Use batch collector for batching
		subscribePendingTransactionsBatched(ctx, client, ch)
	} else {
		// Emit individually
		subscribePendingTransactionsDirect(ctx, client, ch)
	}
}

// subscribePendingTransactionsBatched subscribes with batching enabled.
func subscribePendingTransactionsBatched(
	ctx context.Context,
	client WatchClient,
	ch chan<- WatchPendingTransactionsEvent,
) {
	// Create a channel for individual hashes
	hashCh := make(chan common.Hash, 1000)

	// Start batch collector
	collector := batch.NewCollector[common.Hash](batch.CollectorOptions{
		BatchSize: 100,
		Timeout:   100 * time.Millisecond, // Short timeout for responsiveness
	})
	batches := collector.Collect(ctx, hashCh)

	// Forward batches to output channel
	go func() {
		for batch := range batches {
			if len(batch) > 0 {
				select {
				case ch <- WatchPendingTransactionsEvent{Hashes: batch}:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	// Subscribe to newPendingTransactions
	sub, err := client.Subscribe(
		transport.NewPendingTransactionsSubscribeParams(),
		func(data json.RawMessage) {
			// Parse transaction hash
			var txHash string
			if err := json.Unmarshal(data, &txHash); err != nil {
				// Try parsing as object with hash field
				var txData struct {
					Hash string `json:"hash"`
				}
				if err := json.Unmarshal(data, &txData); err != nil {
					return
				}
				txHash = txData.Hash
			}

			if txHash != "" {
				hash := common.HexToHash(txHash)
				select {
				case hashCh <- hash:
				case <-ctx.Done():
				}
			}
		},
		func(err error) {
			select {
			case ch <- WatchPendingTransactionsEvent{Error: err}:
			case <-ctx.Done():
			}
		},
	)

	if err != nil {
		close(hashCh)
		select {
		case ch <- WatchPendingTransactionsEvent{Error: fmt.Errorf("failed to subscribe: %w", err)}:
		case <-ctx.Done():
		}
		return
	}

	// Wait for context cancellation
	<-ctx.Done()

	// Cleanup
	close(hashCh)
	if sub != nil {
		_ = sub.Unsubscribe()
	}
}

// subscribePendingTransactionsDirect subscribes without batching.
func subscribePendingTransactionsDirect(
	ctx context.Context,
	client WatchClient,
	ch chan<- WatchPendingTransactionsEvent,
) {
	// Subscribe to newPendingTransactions
	sub, err := client.Subscribe(
		transport.NewPendingTransactionsSubscribeParams(),
		func(data json.RawMessage) {
			// Parse transaction hash
			var txHash string
			if err := json.Unmarshal(data, &txHash); err != nil {
				// Try parsing as object with hash field
				var txData struct {
					Hash string `json:"hash"`
				}
				if err := json.Unmarshal(data, &txData); err != nil {
					return
				}
				txHash = txData.Hash
			}

			if txHash != "" {
				hash := common.HexToHash(txHash)
				select {
				case ch <- WatchPendingTransactionsEvent{Hashes: []common.Hash{hash}}:
				case <-ctx.Done():
				}
			}
		},
		func(err error) {
			select {
			case ch <- WatchPendingTransactionsEvent{Error: err}:
			case <-ctx.Done():
			}
		},
	)

	if err != nil {
		select {
		case ch <- WatchPendingTransactionsEvent{Error: fmt.Errorf("failed to subscribe: %w", err)}:
		case <-ctx.Done():
		}
		return
	}

	// Wait for context cancellation
	<-ctx.Done()

	// Unsubscribe
	if sub != nil {
		_ = sub.Unsubscribe()
	}
}
