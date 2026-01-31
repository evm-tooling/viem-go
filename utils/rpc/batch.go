package rpc

import (
	"context"
	"sync"
	"time"
)

// BatchScheduler batches multiple RPC requests together.
type BatchScheduler struct {
	client    *HTTPClient
	batchSize int
	wait      time.Duration
	mu        sync.Mutex
	pending   []pendingRequest
	timer     *time.Timer
	ctx       context.Context
	cancel    context.CancelFunc
}

// pendingRequest represents a request waiting to be batched.
type pendingRequest struct {
	body   RPCRequest
	respCh chan batchResult
}

// batchResult contains the result of a batched request.
type batchResult struct {
	resp *RPCResponse
	err  error
}

// BatchSchedulerOptions contains options for the batch scheduler.
type BatchSchedulerOptions struct {
	// BatchSize is the maximum number of requests per batch.
	BatchSize int
	// Wait is the maximum time to wait before sending a batch.
	Wait time.Duration
}

// DefaultBatchSchedulerOptions returns default options.
func DefaultBatchSchedulerOptions() BatchSchedulerOptions {
	return BatchSchedulerOptions{
		BatchSize: 1000,
		Wait:      0,
	}
}

// NewBatchScheduler creates a new batch scheduler.
func NewBatchScheduler(client *HTTPClient, opts ...BatchSchedulerOptions) *BatchScheduler {
	opt := DefaultBatchSchedulerOptions()
	if len(opts) > 0 {
		opt = opts[0]
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &BatchScheduler{
		client:    client,
		batchSize: opt.BatchSize,
		wait:      opt.Wait,
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Schedule adds a request to the batch.
func (s *BatchScheduler) Schedule(ctx context.Context, body RPCRequest) (*RPCResponse, error) {
	respCh := make(chan batchResult, 1)

	s.mu.Lock()

	// Add to pending requests
	s.pending = append(s.pending, pendingRequest{
		body:   body,
		respCh: respCh,
	})

	// Check if we should flush immediately
	shouldFlush := len(s.pending) >= s.batchSize

	// Start timer if this is the first request and we have a wait time
	if len(s.pending) == 1 && s.wait > 0 && !shouldFlush {
		s.timer = time.AfterFunc(s.wait, func() {
			s.flush()
		})
	}

	s.mu.Unlock()

	// Flush if batch is full
	if shouldFlush {
		s.flush()
	}

	// Wait for result
	select {
	case result := <-respCh:
		return result.resp, result.err
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-s.ctx.Done():
		return nil, s.ctx.Err()
	}
}

// flush sends all pending requests.
func (s *BatchScheduler) flush() {
	s.mu.Lock()

	// Stop timer if running
	if s.timer != nil {
		s.timer.Stop()
		s.timer = nil
	}

	// Get pending requests
	pending := s.pending
	s.pending = nil

	s.mu.Unlock()

	if len(pending) == 0 {
		return
	}

	// Build batch request
	bodies := make([]RPCRequest, len(pending))
	for i, p := range pending {
		bodies[i] = p.body
	}

	// Send batch request
	responses, err := s.client.BatchRequest(s.ctx, bodies)

	// Map responses back to requests
	responseMap := make(map[any]RPCResponse)
	if err == nil {
		for _, resp := range responses {
			responseMap[resp.ID] = resp
		}
	}

	// Send results to waiting goroutines
	for _, p := range pending {
		result := batchResult{}
		if err != nil {
			result.err = err
		} else if resp, ok := responseMap[p.body.ID]; ok {
			result.resp = &resp
		} else {
			result.err = NewHTTPRequestError(s.client.URL(), 0, "", p.body, nil)
		}

		select {
		case p.respCh <- result:
		default:
		}
	}
}

// Close stops the batch scheduler.
func (s *BatchScheduler) Close() {
	s.cancel()

	s.mu.Lock()
	if s.timer != nil {
		s.timer.Stop()
	}
	s.mu.Unlock()
}

// Flush forces a flush of pending requests.
func (s *BatchScheduler) Flush() {
	s.flush()
}

// PendingCount returns the number of pending requests.
func (s *BatchScheduler) PendingCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.pending)
}
