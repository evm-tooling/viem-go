package public

import (
	"context"
	"fmt"
	"time"

	json "github.com/goccy/go-json"

	"github.com/ChefBingbong/viem-go/client/transport"
	"github.com/ChefBingbong/viem-go/utils/observe"
	"github.com/ChefBingbong/viem-go/utils/poll"
)

// WatchBlockNumberParameters contains the parameters for the WatchBlockNumber action.
// This mirrors viem's WatchBlockNumberParameters type.
type WatchBlockNumberParameters struct {
	// EmitOnBegin determines whether to emit the current block number immediately
	// when the watcher starts.
	// Default: false
	EmitOnBegin bool

	// EmitMissed determines whether to emit missed block numbers to the callback.
	// When true, if the watcher detects a gap between block numbers, it will emit
	// all the missed block numbers in sequence.
	// Default: false
	EmitMissed bool

	// Poll forces polling mode even when WebSocket transport is available.
	// If nil, automatically detects based on transport type.
	Poll *bool

	// PollingInterval is the interval between polls when using polling mode.
	// If zero, uses the client's default polling interval.
	PollingInterval time.Duration
}

// WatchBlockNumberEvent represents an event from WatchBlockNumber.
type WatchBlockNumberEvent struct {
	// BlockNumber is the current block number.
	BlockNumber uint64

	// PrevBlockNumber is the previous block number (nil for first event).
	PrevBlockNumber *uint64

	// Error is any error that occurred while fetching the block number.
	Error error
}

// blockNumberObserver is the global observer for block number subscriptions.
var blockNumberObserver = observe.New[WatchBlockNumberEvent]()

// WatchBlockNumber watches and returns incoming block numbers.
//
// This is equivalent to viem's `watchBlockNumber` action with full Go optimization:
//   - Channels instead of callbacks for native Go concurrency
//   - context.Context for cancellation
//   - Observer pattern for deduplication (multiple watchers share one source)
//   - Automatic transport detection (polling vs subscription)
//
// JSON-RPC Methods:
//   - When polling: calls eth_blockNumber on a polling interval
//   - When subscribing: uses eth_subscribe with "newHeads" event
//
// Example:
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	events := public.WatchBlockNumber(ctx, client, public.WatchBlockNumberParameters{
//	    EmitOnBegin: true,
//	    EmitMissed:  true,
//	})
//
//	for event := range events {
//	    if event.Error != nil {
//	        log.Printf("error: %v", event.Error)
//	        continue
//	    }
//	    fmt.Printf("Block: %d (prev: %v)\n", event.BlockNumber, event.PrevBlockNumber)
//	}
func WatchBlockNumber(
	ctx context.Context,
	client WatchClient,
	params WatchBlockNumberParameters,
) <-chan WatchBlockNumberEvent {
	// Determine if we should poll or subscribe
	enablePolling := ShouldPoll(client, params.Poll)

	// Get polling interval
	pollingInterval := GetPollingInterval(client, params.PollingInterval)

	// Create output channel
	ch := make(chan WatchBlockNumberEvent, 10)

	go func() {
		defer close(ch)

		if enablePolling {
			pollBlockNumber(ctx, client, params, pollingInterval, ch)
		} else {
			subscribeBlockNumber(ctx, client, params, ch)
		}
	}()

	return ch
}

// pollBlockNumber implements block number watching using polling.
func pollBlockNumber(
	ctx context.Context,
	client WatchClient,
	params WatchBlockNumberParameters,
	interval time.Duration,
	ch chan<- WatchBlockNumberEvent,
) {
	var prevBlockNumber *uint64

	// Create observer ID for deduplication
	observerID := fmt.Sprintf("watchBlockNumber.%s.%v.%v.%v",
		client.UID(),
		params.EmitOnBegin,
		params.EmitMissed,
		interval,
	)

	// Use observer for deduplication
	eventCh := blockNumberObserver.Subscribe(observerID, func() (<-chan WatchBlockNumberEvent, func()) {
		sourceCh := make(chan WatchBlockNumberEvent, 10)

		// Start polling
		pollResults := poll.Poll(ctx, func(ctx context.Context) (uint64, error) {
			// Get block number with no caching
			cacheDuration := time.Duration(0)
			return GetBlockNumber(ctx, client, GetBlockNumberParameters{
				CacheTime: &cacheDuration,
			})
		}, poll.Options{
			Interval:    interval,
			EmitOnBegin: params.EmitOnBegin,
		})

		// Process poll results
		go func() {
			defer close(sourceCh)

			for result := range pollResults {
				if result.Error != nil {
					select {
					case sourceCh <- WatchBlockNumberEvent{Error: result.Error}:
					case <-ctx.Done():
						return
					}
					continue
				}

				blockNumber := result.Value

				// Skip if same as previous
				if prevBlockNumber != nil && blockNumber == *prevBlockNumber {
					continue
				}

				// Emit missed blocks if enabled
				if params.EmitMissed && prevBlockNumber != nil && blockNumber-*prevBlockNumber > 1 {
					for i := *prevBlockNumber + 1; i < blockNumber; i++ {
						prev := i - 1
						select {
						case sourceCh <- WatchBlockNumberEvent{
							BlockNumber:     i,
							PrevBlockNumber: &prev,
						}:
							prevCopy := i
							prevBlockNumber = &prevCopy
						case <-ctx.Done():
							return
						}
					}
				}

				// Emit current block number if it's newer
				if prevBlockNumber == nil || blockNumber > *prevBlockNumber {
					select {
					case sourceCh <- WatchBlockNumberEvent{
						BlockNumber:     blockNumber,
						PrevBlockNumber: prevBlockNumber,
					}:
						prevCopy := blockNumber
						prevBlockNumber = &prevCopy
					case <-ctx.Done():
						return
					}
				}
			}
		}()

		return sourceCh, func() {}
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

// subscribeBlockNumber implements block number watching using WebSocket subscription.
func subscribeBlockNumber(
	ctx context.Context,
	client WatchClient,
	params WatchBlockNumberParameters,
	ch chan<- WatchBlockNumberEvent,
) {
	var prevBlockNumber *uint64

	// Emit on begin if requested
	if params.EmitOnBegin {
		blockNumber, err := GetBlockNumber(ctx, client, GetBlockNumberParameters{})
		if err != nil {
			select {
			case ch <- WatchBlockNumberEvent{Error: err}:
			case <-ctx.Done():
				return
			}
		} else {
			select {
			case ch <- WatchBlockNumberEvent{
				BlockNumber:     blockNumber,
				PrevBlockNumber: nil,
			}:
				prevBlockNumber = &blockNumber
			case <-ctx.Done():
				return
			}
		}
	}

	// Subscribe to newHeads
	sub, err := client.Subscribe(
		transport.NewHeadsSubscribeParams(),
		func(data json.RawMessage) {
			// Parse block number from newHeads notification
			var header struct {
				Number string `json:"number"`
			}
			if err := json.Unmarshal(data, &header); err != nil {
				select {
				case ch <- WatchBlockNumberEvent{Error: fmt.Errorf("failed to parse block header: %w", err)}:
				case <-ctx.Done():
				}
				return
			}

			// Parse hex block number
			blockNumber, err := parseHexUint64(header.Number)
			if err != nil {
				select {
				case ch <- WatchBlockNumberEvent{Error: fmt.Errorf("failed to parse block number: %w", err)}:
				case <-ctx.Done():
				}
				return
			}

			// Emit missed blocks if enabled
			if params.EmitMissed && prevBlockNumber != nil && blockNumber-*prevBlockNumber > 1 {
				for i := *prevBlockNumber + 1; i < blockNumber; i++ {
					prev := i - 1
					select {
					case ch <- WatchBlockNumberEvent{
						BlockNumber:     i,
						PrevBlockNumber: &prev,
					}:
						prevCopy := i
						prevBlockNumber = &prevCopy
					case <-ctx.Done():
						return
					}
				}
			}

			// Emit current block number
			select {
			case ch <- WatchBlockNumberEvent{
				BlockNumber:     blockNumber,
				PrevBlockNumber: prevBlockNumber,
			}:
				prevCopy := blockNumber
				prevBlockNumber = &prevCopy
			case <-ctx.Done():
			}
		},
		func(err error) {
			select {
			case ch <- WatchBlockNumberEvent{Error: err}:
			case <-ctx.Done():
			}
		},
	)

	if err != nil {
		select {
		case ch <- WatchBlockNumberEvent{Error: fmt.Errorf("failed to subscribe: %w", err)}:
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
