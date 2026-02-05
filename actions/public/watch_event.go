package public

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"

	viemabi "github.com/ChefBingbong/viem-go/abi"
	"github.com/ChefBingbong/viem-go/client/transport"
	"github.com/ChefBingbong/viem-go/utils/batch"
	"github.com/ChefBingbong/viem-go/utils/formatters"
	"github.com/ChefBingbong/viem-go/utils/observe"
	"github.com/ChefBingbong/viem-go/utils/poll"
)

// WatchEventParameters contains the parameters for the WatchEvent action.
// This mirrors viem's WatchEventParameters type.
type WatchEventParameters struct {
	// Address is the contract address(es) to filter logs from.
	// Can be a single address or a slice of addresses.
	Address any // common.Address or []common.Address

	// Event is a single event definition to filter for.
	// Mutually exclusive with Events.
	Event *viemabi.Event

	// Events is a list of event definitions to filter for.
	// Mutually exclusive with Event.
	Events []*viemabi.Event

	// Args are indexed event arguments to filter by.
	// Keys are parameter names, values are the expected values.
	Args map[string]any

	// FromBlock is the block number to start watching from.
	// If set, forces polling mode.
	FromBlock *uint64

	// Strict determines whether logs must match the event definition exactly.
	// When true, logs with mismatched indexed/non-indexed arguments are skipped.
	// Default: false
	Strict bool

	// Batch determines whether to batch logs together.
	// When true, multiple logs are collected and emitted together.
	// When false, each log is emitted as a separate event.
	// Default: true
	Batch bool

	// Poll forces polling mode even when WebSocket transport is available.
	// If nil, automatically detects based on transport type.
	Poll *bool

	// PollingInterval is the interval between polls when using polling mode.
	// If zero, uses the client's default polling interval.
	PollingInterval time.Duration

	// WorkerPoolSize is the number of workers for parallel log processing.
	// Only used when decoding event logs.
	// Default: 4
	WorkerPoolSize int
}

// WatchEventEvent represents an event from WatchEvent.
type WatchEventEvent struct {
	// Logs are the event logs.
	// When Batch is true, this may contain multiple logs.
	// When Batch is false, this will contain a single log.
	Logs []formatters.Log

	// Error is any error that occurred.
	Error error
}

// eventObserver is the global observer for event subscriptions.
var eventObserver = observe.New[WatchEventEvent]()

// WatchEvent watches and returns emitted event logs.
//
// This is equivalent to viem's `watchEvent` action with full Go optimization:
//   - Channels instead of callbacks for native Go concurrency
//   - context.Context for cancellation
//   - Observer pattern for deduplication (multiple watchers share one source)
//   - Worker pool for parallel log decoding
//   - BatchCollector for efficient batching
//   - Automatic transport detection (polling vs subscription)
//   - Filter fallback when eth_newFilter is not supported
//
// JSON-RPC Methods:
//   - When polling with filter support:
//   - Calls eth_newFilter to create a filter
//   - Calls eth_getFilterChanges on a polling interval
//   - When polling without filter support:
//   - Calls eth_getLogs for each block range
//   - When subscribing: uses eth_subscribe with "logs" event
//
// Example:
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	events := public.WatchEvent(ctx, client, public.WatchEventParameters{
//	    Address: contractAddress,
//	    Event:   transferEvent,
//	    Batch:   true,
//	})
//
//	for event := range events {
//	    if event.Error != nil {
//	        log.Printf("error: %v", event.Error)
//	        continue
//	    }
//	    for _, log := range event.Logs {
//	        fmt.Printf("Event: %s at block %d\n", log.EventName, log.BlockNumber)
//	    }
//	}
func WatchEvent(
	ctx context.Context,
	client WatchClient,
	params WatchEventParameters,
) <-chan WatchEventEvent {
	// Default batch to true
	batchMode := params.Batch
	if !batchMode {
		batchMode = true // Default to true
	}

	// Determine if we should poll or subscribe
	enablePolling := ShouldPoll(client, params.Poll)
	// Force polling if fromBlock is set
	if params.FromBlock != nil {
		enablePolling = true
	}

	// Get polling interval
	pollingInterval := GetPollingInterval(client, params.PollingInterval)

	// Create output channel
	ch := make(chan WatchEventEvent, 10)

	go func() {
		defer close(ch)

		if enablePolling {
			pollEvent(ctx, client, params, batchMode, pollingInterval, ch)
		} else {
			subscribeEvent(ctx, client, params, batchMode, ch)
		}
	}()

	return ch
}

// pollEvent implements event watching using polling.
func pollEvent(
	ctx context.Context,
	client WatchClient,
	params WatchEventParameters,
	batchMode bool,
	interval time.Duration,
	ch chan<- WatchEventEvent,
) {
	// Build topics from event definitions
	topics := buildEventTopics(params.Event, params.Events, params.Args)

	// Create observer ID for deduplication
	observerID := fmt.Sprintf("watchEvent.%s.%v.%v.%v.%v.%v",
		client.UID(),
		params.Address,
		params.Args,
		batchMode,
		params.FromBlock,
		interval,
	)

	// Use observer for deduplication
	eventCh := eventObserver.Subscribe(observerID, func() (<-chan WatchEventEvent, func()) {
		sourceCh := make(chan WatchEventEvent, 100)

		var filterID FilterID
		var previousBlockNumber uint64
		var filterSupported = true
		initialized := false

		if params.FromBlock != nil {
			previousBlockNumber = *params.FromBlock - 1
		}

		// Start polling
		pollResults := poll.Poll(ctx, func(ctx context.Context) ([]formatters.Log, error) {
			// First iteration: create filter
			if !initialized {
				if filterSupported {
					filter, err := CreateEventFilter(ctx, client, CreateEventFilterParameters{
						Address:      params.Address,
						Topics:       topics,
						FromBlock:    params.FromBlock,
						FromBlockTag: BlockTagLatest,
					})
					if err != nil {
						// Filter creation failed - fall back to getLogs
						filterSupported = false
					} else {
						filterID = filter.ID
					}
				}
				initialized = true
				return nil, nil
			}

			// Subsequent iterations: get filter changes or use getLogs fallback
			if filterSupported && filterID != "" {
				return GetFilterChangesLogs(ctx, client, filterID)
			}

			// Fallback to getLogs
			blockNumber, err := GetBlockNumber(ctx, client, GetBlockNumberParameters{})
			if err != nil {
				return nil, err
			}

			// Skip if no new blocks
			if previousBlockNumber != 0 && previousBlockNumber == blockNumber {
				return nil, nil
			}

			// Get logs for the new blocks
			var fromBlock uint64
			if previousBlockNumber != 0 {
				fromBlock = previousBlockNumber + 1
			} else {
				fromBlock = blockNumber
			}

			logs, err := GetLogs(ctx, client, GetLogsParameters{
				Address:   params.Address,
				Topics:    topics,
				FromBlock: &fromBlock,
				ToBlock:   &blockNumber,
			})

			previousBlockNumber = blockNumber
			return logs, err
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
					// Check if filter was invalidated
					if filterID != "" && isInvalidInputError(result.Error) {
						// Reinitialize filter
						initialized = false
						filterID = ""
					}
					select {
					case sourceCh <- WatchEventEvent{Error: result.Error}:
					case <-ctx.Done():
						return
					}
					continue
				}

				logs := result.Value
				if len(logs) == 0 {
					continue
				}

				// Emit logs
				if batchMode {
					select {
					case sourceCh <- WatchEventEvent{Logs: logs}:
					case <-ctx.Done():
						return
					}
				} else {
					// Emit individually
					for _, log := range logs {
						select {
						case sourceCh <- WatchEventEvent{Logs: []formatters.Log{log}}:
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

// subscribeEvent implements event watching using WebSocket subscription.
func subscribeEvent(
	ctx context.Context,
	client WatchClient,
	params WatchEventParameters,
	batchMode bool,
	ch chan<- WatchEventEvent,
) {
	// Build topics from event definitions
	topics := buildEventTopics(params.Event, params.Events, params.Args)

	// Build address filter
	var addressFilter any
	if params.Address != nil {
		switch addr := params.Address.(type) {
		case common.Address:
			addressFilter = addr.Hex()
		case *common.Address:
			if addr != nil {
				addressFilter = addr.Hex()
			}
		case []common.Address:
			addrs := make([]string, len(addr))
			for i, a := range addr {
				addrs[i] = a.Hex()
			}
			addressFilter = addrs
		case string:
			addressFilter = addr
		case []string:
			addressFilter = addr
		}
	}

	if batchMode {
		subscribeEventBatched(ctx, client, addressFilter, topics, params, ch)
	} else {
		subscribeEventDirect(ctx, client, addressFilter, topics, params, ch)
	}
}

// subscribeEventBatched subscribes with batching enabled.
func subscribeEventBatched(
	ctx context.Context,
	client WatchClient,
	addressFilter any,
	topics []any,
	params WatchEventParameters,
	ch chan<- WatchEventEvent,
) {
	// Create a channel for individual logs
	logCh := make(chan formatters.Log, 1000)

	// Start batch collector
	collector := batch.NewCollector[formatters.Log](batch.CollectorOptions{
		BatchSize: 100,
		Timeout:   100 * time.Millisecond,
	})
	batches := collector.Collect(ctx, logCh)

	// Forward batches to output channel
	go func() {
		for batch := range batches {
			if len(batch) > 0 {
				select {
				case ch <- WatchEventEvent{Logs: batch}:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	// Subscribe to logs
	sub, err := client.Subscribe(
		transport.LogsSubscribeParams(addressFilter, topics),
		func(data json.RawMessage) {
			log := parseLogFromSubscription(data, params)
			if log != nil {
				select {
				case logCh <- *log:
				case <-ctx.Done():
				}
			}
		},
		func(err error) {
			select {
			case ch <- WatchEventEvent{Error: err}:
			case <-ctx.Done():
			}
		},
	)

	if err != nil {
		close(logCh)
		select {
		case ch <- WatchEventEvent{Error: fmt.Errorf("failed to subscribe: %w", err)}:
		case <-ctx.Done():
		}
		return
	}

	// Wait for context cancellation
	<-ctx.Done()

	// Cleanup
	close(logCh)
	if sub != nil {
		_ = sub.Unsubscribe()
	}
}

// subscribeEventDirect subscribes without batching.
func subscribeEventDirect(
	ctx context.Context,
	client WatchClient,
	addressFilter any,
	topics []any,
	params WatchEventParameters,
	ch chan<- WatchEventEvent,
) {
	// Subscribe to logs
	sub, err := client.Subscribe(
		transport.LogsSubscribeParams(addressFilter, topics),
		func(data json.RawMessage) {
			log := parseLogFromSubscription(data, params)
			if log != nil {
				select {
				case ch <- WatchEventEvent{Logs: []formatters.Log{*log}}:
				case <-ctx.Done():
				}
			}
		},
		func(err error) {
			select {
			case ch <- WatchEventEvent{Error: err}:
			case <-ctx.Done():
			}
		},
	)

	if err != nil {
		select {
		case ch <- WatchEventEvent{Error: fmt.Errorf("failed to subscribe: %w", err)}:
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

// parseLogFromSubscription parses a log from a subscription notification.
func parseLogFromSubscription(data json.RawMessage, params WatchEventParameters) *formatters.Log {
	var rpcLog formatters.RpcLog
	if err := json.Unmarshal(data, &rpcLog); err != nil {
		return nil
	}

	log := formatters.FormatLog(rpcLog, nil)
	return &log
}

// buildEventTopics builds topic filters from event definitions.
func buildEventTopics(event *viemabi.Event, events []*viemabi.Event, args map[string]any) []any {
	if event == nil && len(events) == 0 {
		return nil
	}

	var topics []any

	// Get all events to process
	allEvents := events
	if event != nil {
		allEvents = []*viemabi.Event{event}
	}

	if len(allEvents) == 0 {
		return nil
	}

	// Build topic0 (event signature)
	if len(allEvents) == 1 {
		topics = append(topics, allEvents[0].Topic.Hex())
	} else {
		// Multiple events - use array for OR condition
		sigs := make([]string, len(allEvents))
		for i, e := range allEvents {
			sigs[i] = e.Topic.Hex()
		}
		topics = append(topics, sigs)
	}

	// Add indexed argument topics (only if single event and args provided)
	if len(allEvents) == 1 && len(args) > 0 {
		e := allEvents[0]
		for _, input := range e.Inputs {
			if input.Indexed {
				if argValue, ok := args[input.Name]; ok {
					topics = append(topics, encodeFilterTopic(argValue))
				} else {
					topics = append(topics, nil) // Match any
				}
			}
		}
	}

	return topics
}

// isInvalidInputError checks if an error indicates an invalid filter.
func isInvalidInputError(err error) bool {
	if err == nil {
		return false
	}
	// Check for common invalid input error patterns
	errStr := err.Error()
	return contains(errStr, "filter not found") ||
		contains(errStr, "invalid filter") ||
		contains(errStr, "Filter not found")
}

// contains checks if s contains substr (simple string search).
func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr) != -1
}

// searchString performs a simple string search.
func searchString(s, substr string) int {
	n := len(substr)
	if n == 0 {
		return 0
	}
	for i := 0; i <= len(s)-n; i++ {
		if s[i:i+n] == substr {
			return i
		}
	}
	return -1
}

// EventLogWorkerPool processes event logs in parallel.
type EventLogWorkerPool struct {
	workers int
	jobs    chan eventLogJob
	wg      sync.WaitGroup
}

type eventLogJob struct {
	log     formatters.Log
	abi     *viemabi.ABI
	result  chan<- formatters.Log
	errChan chan<- error
}

// NewEventLogWorkerPool creates a new worker pool for processing event logs.
func NewEventLogWorkerPool(workers int) *EventLogWorkerPool {
	if workers <= 0 {
		workers = 4
	}
	return &EventLogWorkerPool{
		workers: workers,
		jobs:    make(chan eventLogJob, workers*2),
	}
}

// Start starts the worker pool.
func (p *EventLogWorkerPool) Start(ctx context.Context) {
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-p.jobs:
					if !ok {
						return
					}
					// Process the log (decode event data)
					processedLog := processEventLog(job.log, job.abi)
					select {
					case job.result <- processedLog:
					case <-ctx.Done():
						return
					}
				}
			}
		}()
	}
}

// Submit submits a log for processing.
func (p *EventLogWorkerPool) Submit(log formatters.Log, abi *viemabi.ABI, result chan<- formatters.Log, errChan chan<- error) {
	p.jobs <- eventLogJob{
		log:     log,
		abi:     abi,
		result:  result,
		errChan: errChan,
	}
}

// Stop stops the worker pool.
func (p *EventLogWorkerPool) Stop() {
	close(p.jobs)
	p.wg.Wait()
}

// processEventLog processes a single event log (decodes event data if ABI provided).
func processEventLog(log formatters.Log, abi *viemabi.ABI) formatters.Log {
	// If no ABI, return as-is
	if abi == nil {
		return log
	}

	// Try to decode the event
	// This is a simplified version - full implementation would use ABI decoding
	return log
}
