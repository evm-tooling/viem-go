package transport

import (
	"context"
	"errors"
	"sync"
	"time"
)

// FallbackTransportConfig contains configuration for the fallback transport.
type FallbackTransportConfig struct {
	// Key is the transport key.
	Key string
	// Name is the transport name.
	Name string
	// Methods specifies which RPC methods to allow/block.
	Methods *MethodFilter
	// Rank enables ranking transports by latency/success rate.
	Rank *RankConfig
	// RetryCount is the maximum number of retry attempts per transport.
	RetryCount int
	// RetryDelay is the base delay between retries.
	RetryDelay time.Duration
	// Timeout is the request timeout.
	Timeout time.Duration
}

// RankConfig contains configuration for transport ranking.
type RankConfig struct {
	// Enabled enables transport ranking.
	Enabled bool
	// Interval is the interval to re-evaluate transport rankings.
	Interval time.Duration
	// Weights contains weights for ranking metrics.
	Weights *RankWeights
}

// RankWeights contains weights for ranking metrics.
type RankWeights struct {
	// Latency weight (0-1).
	Latency float64
	// Stability weight (0-1).
	Stability float64
}

// DefaultFallbackTransportConfig returns default fallback transport configuration.
func DefaultFallbackTransportConfig() FallbackTransportConfig {
	return FallbackTransportConfig{
		Key:        "fallback",
		Name:       "Fallback JSON-RPC",
		RetryCount: 3,
		RetryDelay: 150 * time.Millisecond,
		Timeout:    10 * time.Second,
		Rank: &RankConfig{
			Enabled:  true,
			Interval: 10 * time.Second,
			Weights: &RankWeights{
				Latency:   0.3,
				Stability: 0.7,
			},
		},
	}
}

// transportStats tracks statistics for a transport.
type transportStats struct {
	latency   time.Duration
	successes int
	failures  int
	mu        sync.RWMutex
}

// FallbackTransport implements a fallback transport that tries multiple transports.
type FallbackTransport struct {
	config     FallbackTransportConfig
	transports []Transport
	stats      []*transportStats
	order      []int
	orderMu    sync.RWMutex
}

// Fallback creates a new fallback transport factory.
func Fallback(factories ...TransportFactory) TransportFactory {
	return FallbackWithConfig(factories, DefaultFallbackTransportConfig())
}

// FallbackWithConfig creates a new fallback transport factory with config.
func FallbackWithConfig(factories []TransportFactory, config FallbackTransportConfig) TransportFactory {
	return func(params TransportParams) (Transport, error) {
		if len(factories) == 0 {
			return nil, errors.New("at least one transport factory is required")
		}

		// Create transports
		transports := make([]Transport, 0, len(factories))
		for _, factory := range factories {
			t, err := factory(params)
			if err != nil {
				// Continue to next transport on error
				continue
			}
			transports = append(transports, t)
		}

		if len(transports) == 0 {
			return nil, errors.New("failed to create any transports")
		}

		// Apply parameter overrides
		if params.RetryCount != nil {
			config.RetryCount = *params.RetryCount
		}
		if params.Timeout != nil {
			config.Timeout = *params.Timeout
		}

		return NewFallbackTransport(transports, config)
	}
}

// NewFallbackTransport creates a new fallback transport.
func NewFallbackTransport(transports []Transport, config FallbackTransportConfig) (*FallbackTransport, error) {
	if len(transports) == 0 {
		return nil, errors.New("at least one transport is required")
	}

	// Initialize stats
	stats := make([]*transportStats, len(transports))
	order := make([]int, len(transports))
	for i := range transports {
		stats[i] = &transportStats{}
		order[i] = i
	}

	ft := &FallbackTransport{
		config:     config,
		transports: transports,
		stats:      stats,
		order:      order,
	}

	// Start ranking if enabled
	if config.Rank != nil && config.Rank.Enabled {
		go ft.rankingLoop()
	}

	return ft, nil
}

// Config returns the transport configuration.
func (t *FallbackTransport) Config() TransportConfig {
	return TransportConfig{
		Name:       t.config.Name,
		Key:        t.config.Key,
		Type:       "fallback",
		Methods:    t.config.Methods,
		RetryCount: t.config.RetryCount,
		RetryDelay: t.config.RetryDelay,
		Timeout:    t.config.Timeout,
	}
}

// Request sends a JSON-RPC request, trying transports in order.
func (t *FallbackTransport) Request(ctx context.Context, req RPCRequest) (*RPCResponse, error) {
	// Check method filter
	if t.config.Methods != nil && !t.config.Methods.IsAllowed(req.Method) {
		return nil, ErrMethodNotSupported
	}

	// Ensure request has required fields
	if req.ID == nil {
		req.ID = NextID()
	}
	if req.JSONRPC == "" {
		req.JSONRPC = "2.0"
	}

	// Get transport order
	t.orderMu.RLock()
	order := make([]int, len(t.order))
	copy(order, t.order)
	t.orderMu.RUnlock()

	var lastErr error

	// Try each transport in order
	for _, idx := range order {
		transport := t.transports[idx]
		stats := t.stats[idx]

		// Send request
		start := time.Now()
		resp, err := transport.Request(ctx, req)
		latency := time.Since(start)

		if err == nil {
			// Update stats
			stats.mu.Lock()
			stats.latency = (stats.latency + latency) / 2
			stats.successes++
			stats.mu.Unlock()

			return resp, nil
		}

		// Update stats
		stats.mu.Lock()
		stats.failures++
		stats.mu.Unlock()

		lastErr = err

		// Check context
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
	}

	return nil, lastErr
}

// Value returns transport-specific attributes.
func (t *FallbackTransport) Value() *TransportValue {
	// Return value from first transport
	if len(t.transports) > 0 {
		return t.transports[0].Value()
	}
	return &TransportValue{}
}

// Close closes all transports.
func (t *FallbackTransport) Close() error {
	var errs []error
	for _, transport := range t.transports {
		if err := transport.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// rankingLoop periodically re-ranks transports.
func (t *FallbackTransport) rankingLoop() {
	if t.config.Rank == nil || t.config.Rank.Interval <= 0 {
		return
	}

	ticker := time.NewTicker(t.config.Rank.Interval)
	defer ticker.Stop()

	for range ticker.C {
		t.updateRanking()
	}
}

// updateRanking updates the transport order based on stats.
func (t *FallbackTransport) updateRanking() {
	weights := &RankWeights{Latency: 0.3, Stability: 0.7}
	if t.config.Rank != nil && t.config.Rank.Weights != nil {
		weights = t.config.Rank.Weights
	}

	// Calculate scores
	scores := make([]float64, len(t.transports))

	// Find max latency for normalization
	var maxLatency time.Duration
	for _, stats := range t.stats {
		stats.mu.RLock()
		if stats.latency > maxLatency {
			maxLatency = stats.latency
		}
		stats.mu.RUnlock()
	}

	if maxLatency == 0 {
		maxLatency = time.Second // Default
	}

	for i, stats := range t.stats {
		stats.mu.RLock()
		total := stats.successes + stats.failures
		stability := 0.0
		if total > 0 {
			stability = float64(stats.successes) / float64(total)
		}
		latencyScore := 1.0 - (float64(stats.latency) / float64(maxLatency))
		if latencyScore < 0 {
			latencyScore = 0
		}
		stats.mu.RUnlock()

		scores[i] = weights.Latency*latencyScore + weights.Stability*stability
	}

	// Sort by score (descending)
	order := make([]int, len(t.transports))
	for i := range order {
		order[i] = i
	}

	// Simple bubble sort (small array)
	for i := 0; i < len(order)-1; i++ {
		for j := 0; j < len(order)-i-1; j++ {
			if scores[order[j]] < scores[order[j+1]] {
				order[j], order[j+1] = order[j+1], order[j]
			}
		}
	}

	t.orderMu.Lock()
	t.order = order
	t.orderMu.Unlock()
}

// Transports returns the underlying transports.
func (t *FallbackTransport) Transports() []Transport {
	return t.transports
}
