package public

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"

	viemabi "github.com/ChefBingbong/viem-go/abi"
	"github.com/ChefBingbong/viem-go/client/transport"
	"github.com/ChefBingbong/viem-go/utils/batch"
	"github.com/ChefBingbong/viem-go/utils/formatters"
	"github.com/ChefBingbong/viem-go/utils/observe"
	"github.com/ChefBingbong/viem-go/utils/poll"
)

// WatchContractEventParameters contains the parameters for the WatchContractEvent action.
// This mirrors viem's WatchContractEventParameters type.
type WatchContractEventParameters struct {
	// Address is the contract address(es) to filter logs from.
	// Can be a single address or a slice of addresses.
	Address any // common.Address or []common.Address

	// ABI is the contract ABI for decoding event logs.
	ABI *viemabi.ABI

	// EventName is the name of the event to filter for.
	// Must match an event in the ABI.
	EventName string

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

	// WorkerPoolSize is the number of workers for parallel log decoding.
	// Default: 4
	WorkerPoolSize int
}

// WatchContractEventEvent represents an event from WatchContractEvent.
type WatchContractEventEvent struct {
	// Logs are the decoded event logs.
	// Each log includes Args and EventName from ABI decoding.
	// When Batch is true, this may contain multiple logs.
	// When Batch is false, this will contain a single log.
	Logs []formatters.Log

	// Error is any error that occurred.
	Error error
}

// contractEventObserver is the global observer for contract event subscriptions.
var contractEventObserver = observe.New[WatchContractEventEvent]()

// WatchContractEvent watches and returns emitted contract event logs with ABI decoding.
//
// This is equivalent to viem's `watchContractEvent` action with full Go optimization:
//   - Channels instead of callbacks for native Go concurrency
//   - context.Context for cancellation
//   - Observer pattern for deduplication (multiple watchers share one source)
//   - Worker pool for parallel ABI decoding
//   - BatchCollector for efficient batching
//   - Automatic transport detection (polling vs subscription)
//   - Filter fallback when eth_newFilter is not supported
//
// Unlike WatchEvent, WatchContractEvent:
//   - Requires an ABI for event decoding
//   - Automatically encodes event topics from ABI
//   - Decodes event logs using the ABI
//   - Filters by EventName
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
//	events := public.WatchContractEvent(ctx, client, public.WatchContractEventParameters{
//	    Address:   contractAddress,
//	    ABI:       erc20ABI,
//	    EventName: "Transfer",
//	    Args:      map[string]any{"from": senderAddress},
//	    Batch:     true,
//	})
//
//	for event := range events {
//	    if event.Error != nil {
//	        log.Printf("error: %v", event.Error)
//	        continue
//	    }
//	    for _, log := range event.Logs {
//	        fmt.Printf("Transfer from %v to %v: %v\n",
//	            log.Args["from"],
//	            log.Args["to"],
//	            log.Args["value"])
//	    }
//	}
func WatchContractEvent(
	ctx context.Context,
	client WatchClient,
	params WatchContractEventParameters,
) <-chan WatchContractEventEvent {
	// Default batch to true
	batchMode := params.Batch
	if !batchMode {
		batchMode = true
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
	ch := make(chan WatchContractEventEvent, 10)

	go func() {
		defer close(ch)

		if enablePolling {
			pollContractEvent(ctx, client, params, batchMode, pollingInterval, ch)
		} else {
			subscribeContractEvent(ctx, client, params, batchMode, ch)
		}
	}()

	return ch
}

// pollContractEvent implements contract event watching using polling.
func pollContractEvent(
	ctx context.Context,
	client WatchClient,
	params WatchContractEventParameters,
	batchMode bool,
	interval time.Duration,
	ch chan<- WatchContractEventEvent,
) {
	// Build topics from ABI
	topics := buildContractEventTopics(params.ABI, params.EventName, params.Args)
	strict := params.Strict

	// Create observer ID for deduplication
	observerID := fmt.Sprintf("watchContractEvent.%s.%v.%v.%s.%v.%v.%v",
		client.UID(),
		params.Address,
		params.Args,
		params.EventName,
		batchMode,
		strict,
		interval,
	)

	// Use observer for deduplication
	eventCh := contractEventObserver.Subscribe(observerID, func() (<-chan WatchContractEventEvent, func()) {
		sourceCh := make(chan WatchContractEventEvent, 100)

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

			if previousBlockNumber != 0 && previousBlockNumber >= blockNumber {
				return nil, nil
			}

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

		// Process poll results with ABI decoding
		go func() {
			defer close(sourceCh)
			defer func() {
				if filterID != "" {
					_, _ = UninstallFilter(context.Background(), client, filterID)
				}
			}()

			for result := range pollResults {
				if result.Error != nil {
					if filterID != "" && isInvalidInputError(result.Error) {
						initialized = false
						filterID = ""
					}
					select {
					case sourceCh <- WatchContractEventEvent{Error: result.Error}:
					case <-ctx.Done():
						return
					}
					continue
				}

				logs := result.Value
				if len(logs) == 0 {
					continue
				}

				// Decode logs using ABI
				decodedLogs := decodeContractEventLogs(logs, params.ABI, params.EventName, strict)
				if len(decodedLogs) == 0 {
					continue
				}

				// Emit logs
				if batchMode {
					select {
					case sourceCh <- WatchContractEventEvent{Logs: decodedLogs}:
					case <-ctx.Done():
						return
					}
				} else {
					for _, log := range decodedLogs {
						select {
						case sourceCh <- WatchContractEventEvent{Logs: []formatters.Log{log}}:
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

// subscribeContractEvent implements contract event watching using WebSocket subscription.
func subscribeContractEvent(
	ctx context.Context,
	client WatchClient,
	params WatchContractEventParameters,
	batchMode bool,
	ch chan<- WatchContractEventEvent,
) {
	// Build topics from ABI
	topics := buildContractEventTopics(params.ABI, params.EventName, params.Args)

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
		subscribeContractEventBatched(ctx, client, addressFilter, topics, params, ch)
	} else {
		subscribeContractEventDirect(ctx, client, addressFilter, topics, params, ch)
	}
}

// subscribeContractEventBatched subscribes with batching enabled.
func subscribeContractEventBatched(
	ctx context.Context,
	client WatchClient,
	addressFilter any,
	topics []any,
	params WatchContractEventParameters,
	ch chan<- WatchContractEventEvent,
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
				case ch <- WatchContractEventEvent{Logs: batch}:
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
			log := parseContractLogFromSubscription(data, params)
			if log != nil {
				select {
				case logCh <- *log:
				case <-ctx.Done():
				}
			}
		},
		func(err error) {
			select {
			case ch <- WatchContractEventEvent{Error: err}:
			case <-ctx.Done():
			}
		},
	)

	if err != nil {
		close(logCh)
		select {
		case ch <- WatchContractEventEvent{Error: fmt.Errorf("failed to subscribe: %w", err)}:
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

// subscribeContractEventDirect subscribes without batching.
func subscribeContractEventDirect(
	ctx context.Context,
	client WatchClient,
	addressFilter any,
	topics []any,
	params WatchContractEventParameters,
	ch chan<- WatchContractEventEvent,
) {
	// Subscribe to logs
	sub, err := client.Subscribe(
		transport.LogsSubscribeParams(addressFilter, topics),
		func(data json.RawMessage) {
			log := parseContractLogFromSubscription(data, params)
			if log != nil {
				select {
				case ch <- WatchContractEventEvent{Logs: []formatters.Log{*log}}:
				case <-ctx.Done():
				}
			}
		},
		func(err error) {
			select {
			case ch <- WatchContractEventEvent{Error: err}:
			case <-ctx.Done():
			}
		},
	)

	if err != nil {
		select {
		case ch <- WatchContractEventEvent{Error: fmt.Errorf("failed to subscribe: %w", err)}:
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

// parseContractLogFromSubscription parses and decodes a log from a subscription notification.
func parseContractLogFromSubscription(data json.RawMessage, params WatchContractEventParameters) *formatters.Log {
	var rpcLog formatters.RpcLog
	if err := json.Unmarshal(data, &rpcLog); err != nil {
		return nil
	}

	log := formatters.FormatLog(rpcLog, nil)

	// Decode using ABI
	if params.ABI != nil {
		decodedLogs := decodeContractEventLogs([]formatters.Log{log}, params.ABI, params.EventName, params.Strict)
		if len(decodedLogs) > 0 {
			return &decodedLogs[0]
		}
		// In non-strict mode, return the log anyway
		if !params.Strict {
			return &log
		}
		return nil
	}

	return &log
}

// buildContractEventTopics builds topic filters from an ABI event.
func buildContractEventTopics(abi *viemabi.ABI, eventName string, args map[string]any) []any {
	if abi == nil || eventName == "" {
		return nil
	}

	// Find the event in the ABI
	var event *viemabi.Event
	for _, e := range abi.Events {
		if e.Name == eventName {
			eventCopy := e
			event = &eventCopy
			break
		}
	}

	if event == nil {
		return nil
	}

	var topics []any

	// Topic0: event signature
	topics = append(topics, event.Topic.Hex())

	// Add indexed argument topics
	if len(args) > 0 {
		for _, input := range event.Inputs {
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

// decodeContractEventLogs decodes event logs using an ABI.
func decodeContractEventLogs(logs []formatters.Log, abi *viemabi.ABI, eventName string, strict bool) []formatters.Log {
	if abi == nil {
		return logs
	}

	// Find the event in the ABI
	var event *viemabi.Event
	for _, e := range abi.Events {
		if eventName == "" || e.Name == eventName {
			eventCopy := e
			event = &eventCopy
			break
		}
	}

	if event == nil {
		if strict {
			return nil
		}
		return logs
	}

	var decodedLogs []formatters.Log
	for _, log := range logs {
		// Verify topic matches event signature
		if len(log.Topics) == 0 {
			if !strict {
				decodedLogs = append(decodedLogs, log)
			}
			continue
		}

		// Check if first topic matches event signature
		if log.Topics[0] != event.Topic.Hex() {
			if !strict {
				decodedLogs = append(decodedLogs, log)
			}
			continue
		}

		// Try to decode the event
		rawTopics := make([]common.Hash, len(log.Topics))
		for i, t := range log.Topics {
			rawTopics[i] = common.HexToHash(t)
		}

		decoded, err := abi.DecodeEventLogByName(eventName, rawTopics, common.FromHex(log.Data))
		if err != nil {
			if strict {
				continue
			}
			// Add without decoding in non-strict mode
			log.EventName = eventName
			decodedLogs = append(decodedLogs, log)
			continue
		}

		// Add decoded args to log
		log.EventName = decoded.EventName
		log.Args = decoded.Args
		decodedLogs = append(decodedLogs, log)
	}

	return decodedLogs
}
