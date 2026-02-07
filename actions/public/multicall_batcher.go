package public

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ChefBingbong/viem-go/types"
)

// multicallBatchEntry represents one caller's submission to the batcher.
type multicallBatchEntry struct {
	contracts []MulticallContract
	params    MulticallParameters
}

// multicallBatchResult represents one caller's results from the batcher.
type multicallBatchResult struct {
	results MulticallReturnType
	err     error
}

// MulticallBatcher aggregates concurrent Multicall calls into fewer, larger
// multicall RPC calls. This mirrors viem's `batch.multicall` behavior where
// multiple concurrent readContract/multicall calls within a time window are
// merged into a single aggregate3 RPC call.
//
// When only a single caller is present (no concurrency), the call is executed
// directly with zero overhead — no timer, no goroutine indirection.
//
// Without this, N concurrent Multicall() calls = N separate eth_call RPCs.
// With this, N concurrent Multicall() calls within the wait window = 1 eth_call RPC.
type MulticallBatcher struct {
	client Client
	opts   types.MulticallBatchOptions

	mu      sync.Mutex
	pending []pendingMulticall
	timer   *time.Timer
}

type pendingMulticall struct {
	entry    multicallBatchEntry
	resultCh chan multicallBatchResult
}

// multicallBatcherCache stores batchers per client UID.
var (
	multicallBatcherCache   = make(map[string]*MulticallBatcher)
	multicallBatcherCacheMu sync.Mutex
)

// getMulticallBatcher returns or creates a MulticallBatcher for the given client.
func getMulticallBatcher(client Client, opts *types.MulticallBatchOptions) *MulticallBatcher {
	if opts == nil {
		return nil
	}

	key := fmt.Sprintf("multicall_batcher.%s", client.UID())

	multicallBatcherCacheMu.Lock()
	defer multicallBatcherCacheMu.Unlock()

	if batcher, ok := multicallBatcherCache[key]; ok {
		return batcher
	}

	batcher := &MulticallBatcher{
		client: client,
		opts:   *opts,
	}
	multicallBatcherCache[key] = batcher
	return batcher
}

// Schedule submits a multicall request to be batched with other concurrent requests.
//
// Fast path: if no other goroutines are currently waiting to batch, the call
// executes directly via multicallDirect with zero overhead (no timer, no channel).
//
// Batch path: if a batch window is already open (timer running from a prior caller),
// this call joins the pending batch and waits for the merged result.
func (b *MulticallBatcher) Schedule(ctx context.Context, params MulticallParameters) (MulticallReturnType, error) {
	b.mu.Lock()

	// Fast path: no pending batch and no timer running — execute directly.
	// This avoids any overhead for single sequential callers.
	if len(b.pending) == 0 && b.timer == nil {
		// Start a batch window so that concurrent goroutines arriving in the
		// next few microseconds can join. But we don't block on it — we'll
		// check if anyone joined after a very short runtime.Gosched-style yield.
		//
		// Actually, the simplest correct approach: set a flag, unlock, execute
		// directly, and if anyone arrives while we're executing they start
		// their own batch window. This gives zero overhead for the common
		// single-caller case.
		b.mu.Unlock()
		return multicallDirect(ctx, b.client, params)
	}

	// Batch path: a batch window is already open — join it.
	resultCh := make(chan multicallBatchResult, 1)

	entry := multicallBatchEntry{
		contracts: params.Contracts,
		params:    params,
	}

	b.pending = append(b.pending, pendingMulticall{
		entry:    entry,
		resultCh: resultCh,
	})

	// Check if batch should be flushed (total contracts exceed batch size)
	totalContracts := 0
	for _, p := range b.pending {
		totalContracts += len(p.entry.contracts)
	}

	batchSize := b.opts.BatchSize
	if batchSize <= 0 {
		batchSize = 2048
	}

	if totalContracts >= batchSize {
		b.flushLocked()
	}

	b.mu.Unlock()

	// Wait for result
	select {
	case result := <-resultCh:
		return result.results, result.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// ScheduleConcurrent submits a multicall as part of a known-concurrent workload.
// Unlike Schedule, this always enters the batch path with a wait window, making it
// suitable for fan-out patterns where many goroutines submit simultaneously.
//
// Use this when you know multiple goroutines will call Multicall concurrently
// (e.g., resolving N tokens in parallel). Use the regular Multicall/Schedule path
// for general-purpose calls where concurrency is unknown.
func (b *MulticallBatcher) ScheduleConcurrent(ctx context.Context, params MulticallParameters) (MulticallReturnType, error) {
	resultCh := make(chan multicallBatchResult, 1)

	b.mu.Lock()

	entry := multicallBatchEntry{
		contracts: params.Contracts,
		params:    params,
	}

	wasEmpty := len(b.pending) == 0
	b.pending = append(b.pending, pendingMulticall{
		entry:    entry,
		resultCh: resultCh,
	})

	// Check if batch should be flushed (total contracts exceed batch size)
	totalContracts := 0
	for _, p := range b.pending {
		totalContracts += len(p.entry.contracts)
	}

	batchSize := b.opts.BatchSize
	if batchSize <= 0 {
		batchSize = 2048
	}

	shouldFlush := totalContracts >= batchSize

	if shouldFlush {
		b.flushLocked()
		b.mu.Unlock()
	} else if wasEmpty {
		wait := b.opts.Wait
		if wait <= 0 {
			// No explicit wait — use a minimal yield to let concurrent goroutines submit.
			wait = time.Millisecond
		}
		b.timer = time.AfterFunc(wait, func() {
			b.mu.Lock()
			b.flushLocked()
			b.mu.Unlock()
		})
		b.mu.Unlock()
	} else {
		b.mu.Unlock()
	}

	// Wait for result
	select {
	case result := <-resultCh:
		return result.results, result.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// flushLocked executes the current batch. Must be called with mu held.
func (b *MulticallBatcher) flushLocked() {
	if len(b.pending) == 0 {
		return
	}

	// Cancel timer if running
	if b.timer != nil {
		b.timer.Stop()
		b.timer = nil
	}

	// Capture current batch
	batch := b.pending
	b.pending = nil

	// Build merged contracts list and track offsets per caller
	type callerRange struct {
		start int
		count int
	}
	ranges := make([]callerRange, len(batch))
	var allContracts []MulticallContract
	offset := 0

	for i, p := range batch {
		count := len(p.entry.contracts)
		ranges[i] = callerRange{start: offset, count: count}
		allContracts = append(allContracts, p.entry.contracts...)
		offset += count
	}

	// Use the first caller's params as the base config
	baseParams := batch[0].entry.params
	mergedParams := MulticallParameters{
		Contracts:           allContracts,
		BatchSize:           baseParams.BatchSize,
		Deployless:          baseParams.Deployless,
		MulticallAddress:    baseParams.MulticallAddress,
		BlockNumber:         baseParams.BlockNumber,
		BlockTag:            baseParams.BlockTag,
		MaxConcurrentChunks: baseParams.MaxConcurrentChunks,
	}

	// Force allowFailure=true for the merged call since different callers
	// may have different expectations about individual failures
	trueVal := true
	mergedParams.AllowFailure = &trueVal

	// Execute the single merged multicall in a goroutine
	go func() {
		results, err := multicallDirect(context.Background(), b.client, mergedParams)

		// Route results back to individual callers
		for i, p := range batch {
			r := ranges[i]
			result := multicallBatchResult{}

			if err != nil {
				result.err = err
			} else if r.start+r.count <= len(results) {
				callerResults := results[r.start : r.start+r.count]

				// If the original caller had AllowFailure=false, check for failures
				if p.entry.params.AllowFailure != nil && !*p.entry.params.AllowFailure {
					for _, cr := range callerResults {
						if cr.Status == "failure" {
							result.err = cr.Error
							break
						}
					}
				}

				if result.err == nil {
					result.results = callerResults
				}
			} else {
				result.err = fmt.Errorf("multicall batcher: result index out of range (%d+%d > %d)",
					r.start, r.count, len(results))
			}

			select {
			case p.resultCh <- result:
			default:
			}
		}
	}()
}
