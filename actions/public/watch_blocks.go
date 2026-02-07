package public

import (
	"context"
	"fmt"
	"time"

	json "github.com/goccy/go-json"

	"github.com/ChefBingbong/viem-go/client/transport"
	"github.com/ChefBingbong/viem-go/types"
	"github.com/ChefBingbong/viem-go/utils/observe"
	"github.com/ChefBingbong/viem-go/utils/poll"
)

// WatchBlocksParameters contains the parameters for the WatchBlocks action.
// This mirrors viem's WatchBlocksParameters type.
type WatchBlocksParameters struct {
	// BlockTag is the block tag to watch. Defaults to "latest".
	BlockTag BlockTag

	// EmitMissed determines whether to emit missed blocks to the callback.
	// When true, if the watcher detects a gap between blocks, it will fetch
	// and emit all the missed blocks in sequence.
	// Default: false
	EmitMissed bool

	// EmitOnBegin determines whether to emit the current block immediately
	// when the watcher starts.
	// Default: false
	EmitOnBegin bool

	// IncludeTransactions determines whether to include full transaction objects
	// in the block data.
	// Default: false
	IncludeTransactions bool

	// Poll forces polling mode even when WebSocket transport is available.
	// If nil, automatically detects based on transport type.
	Poll *bool

	// PollingInterval is the interval between polls when using polling mode.
	// If zero, uses the client's default polling interval.
	PollingInterval time.Duration
}

// WatchBlocksEvent represents an event from WatchBlocks.
type WatchBlocksEvent struct {
	// Block is the current block.
	Block *types.Block

	// PrevBlock is the previous block (nil for first event).
	PrevBlock *types.Block

	// Error is any error that occurred while fetching the block.
	Error error
}

// blocksObserver is the global observer for block subscriptions.
var blocksObserver = observe.New[WatchBlocksEvent]()

// WatchBlocks watches and returns information for incoming blocks.
//
// This is equivalent to viem's `watchBlocks` action with full Go optimization:
//   - Channels instead of callbacks for native Go concurrency
//   - context.Context for cancellation
//   - Observer pattern for deduplication (multiple watchers share one source)
//   - Automatic transport detection (polling vs subscription)
//
// JSON-RPC Methods:
//   - When polling: calls eth_getBlockByNumber on a polling interval
//   - When subscribing: uses eth_subscribe with "newHeads" event, then fetches full block
//
// Example:
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	events := public.WatchBlocks(ctx, client, public.WatchBlocksParameters{
//	    EmitOnBegin:         true,
//	    EmitMissed:          true,
//	    IncludeTransactions: true,
//	})
//
//	for event := range events {
//	    if event.Error != nil {
//	        log.Printf("error: %v", event.Error)
//	        continue
//	    }
//	    fmt.Printf("Block %d: %d transactions\n",
//	        event.Block.Number,
//	        len(event.Block.Transactions))
//	}
func WatchBlocks(
	ctx context.Context,
	client WatchClient,
	params WatchBlocksParameters,
) <-chan WatchBlocksEvent {
	// Set default block tag
	blockTag := params.BlockTag
	if blockTag == "" {
		blockTag = client.ExperimentalBlockTag()
		if blockTag == "" {
			blockTag = BlockTagLatest
		}
	}

	// Determine if we should poll or subscribe
	enablePolling := ShouldPoll(client, params.Poll)

	// Get polling interval
	pollingInterval := GetPollingInterval(client, params.PollingInterval)

	// Create output channel
	ch := make(chan WatchBlocksEvent, 10)

	go func() {
		defer close(ch)

		if enablePolling {
			pollBlocks(ctx, client, params, blockTag, pollingInterval, ch)
		} else {
			subscribeBlocks(ctx, client, params, blockTag, ch)
		}
	}()

	return ch
}

// pollBlocks implements block watching using polling.
func pollBlocks(
	ctx context.Context,
	client WatchClient,
	params WatchBlocksParameters,
	blockTag BlockTag,
	interval time.Duration,
	ch chan<- WatchBlocksEvent,
) {
	var prevBlock *types.Block

	// Create observer ID for deduplication
	observerID := fmt.Sprintf("watchBlocks.%s.%s.%v.%v.%v.%v",
		client.UID(),
		blockTag,
		params.EmitMissed,
		params.EmitOnBegin,
		params.IncludeTransactions,
		interval,
	)

	// Use observer for deduplication
	eventCh := blocksObserver.Subscribe(observerID, func() (<-chan WatchBlocksEvent, func()) {
		sourceCh := make(chan WatchBlocksEvent, 10)

		// Start polling
		pollResults := poll.Poll(ctx, func(ctx context.Context) (*types.Block, error) {
			return GetBlock(ctx, client, GetBlockParameters{
				BlockTag:            blockTag,
				IncludeTransactions: params.IncludeTransactions,
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
					case sourceCh <- WatchBlocksEvent{Error: result.Error}:
					case <-ctx.Done():
						return
					}
					continue
				}

				block := result.Value

				// Skip if same as previous
				if prevBlock != nil {
					if block.Number == prevBlock.Number {
						continue
					}

					// Emit missed blocks if enabled
					if params.EmitMissed && block.Number-prevBlock.Number > 1 {
						for i := prevBlock.Number + 1; i < block.Number; i++ {
							missedBlockNum := i
							missedBlock, err := GetBlock(ctx, client, GetBlockParameters{
								BlockNumber:         &missedBlockNum,
								IncludeTransactions: params.IncludeTransactions,
							})
							if err != nil {
								// Skip errors fetching missed blocks, continue with current
								continue
							}

							select {
							case sourceCh <- WatchBlocksEvent{
								Block:     missedBlock,
								PrevBlock: prevBlock,
							}:
								prevBlock = missedBlock
							case <-ctx.Done():
								return
							}
						}
					}
				}

				// Emit current block if it's newer
				shouldEmit := prevBlock == nil ||
					(blockTag == BlockTagPending && block.Number == 0) ||
					block.Number > prevBlock.Number

				if shouldEmit {
					select {
					case sourceCh <- WatchBlocksEvent{
						Block:     block,
						PrevBlock: prevBlock,
					}:
						prevBlock = block
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

// subscribeBlocks implements block watching using WebSocket subscription.
func subscribeBlocks(
	ctx context.Context,
	client WatchClient,
	params WatchBlocksParameters,
	blockTag BlockTag,
	ch chan<- WatchBlocksEvent,
) {
	var prevBlock *types.Block
	emitFetched := true

	// Emit on begin if requested
	if params.EmitOnBegin {
		block, err := GetBlock(ctx, client, GetBlockParameters{
			BlockTag:            blockTag,
			IncludeTransactions: params.IncludeTransactions,
		})
		if err != nil {
			select {
			case ch <- WatchBlocksEvent{Error: err}:
			case <-ctx.Done():
				return
			}
		} else {
			select {
			case ch <- WatchBlocksEvent{
				Block:     block,
				PrevBlock: nil,
			}:
				prevBlock = block
				emitFetched = false
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
				case ch <- WatchBlocksEvent{Error: fmt.Errorf("failed to parse block header: %w", err)}:
				case <-ctx.Done():
				}
				return
			}

			// Parse hex block number
			blockNumber, err := parseHexUint64(header.Number)
			if err != nil {
				select {
				case ch <- WatchBlocksEvent{Error: fmt.Errorf("failed to parse block number: %w", err)}:
				case <-ctx.Done():
				}
				return
			}

			// Fetch full block
			block, err := GetBlock(ctx, client, GetBlockParameters{
				BlockNumber:         &blockNumber,
				IncludeTransactions: params.IncludeTransactions,
			})
			if err != nil {
				// Ignore errors fetching block, continue waiting
				return
			}

			// Skip if we already emitted this from emitOnBegin
			if emitFetched {
				emitFetched = false
				return
			}

			// Emit missed blocks if enabled
			if params.EmitMissed && prevBlock != nil {
				if block.Number-prevBlock.Number > 1 {
					for i := prevBlock.Number + 1; i < block.Number; i++ {
						missedBlockNum := i
						missedBlock, err := GetBlock(ctx, client, GetBlockParameters{
							BlockNumber:         &missedBlockNum,
							IncludeTransactions: params.IncludeTransactions,
						})
						if err != nil {
							continue
						}

						select {
						case ch <- WatchBlocksEvent{
							Block:     missedBlock,
							PrevBlock: prevBlock,
						}:
							prevBlock = missedBlock
						case <-ctx.Done():
							return
						}
					}
				}
			}

			// Emit current block
			select {
			case ch <- WatchBlocksEvent{
				Block:     block,
				PrevBlock: prevBlock,
			}:
				emitFetched = false
				prevBlock = block
			case <-ctx.Done():
			}
		},
		func(err error) {
			select {
			case ch <- WatchBlocksEvent{Error: err}:
			case <-ctx.Done():
			}
		},
	)

	if err != nil {
		select {
		case ch <- WatchBlocksEvent{Error: fmt.Errorf("failed to subscribe: %w", err)}:
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
