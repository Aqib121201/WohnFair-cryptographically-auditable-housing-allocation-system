package scheduler

import (
	"math"
	"sort"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics collects and exposes scheduler metrics
type Metrics struct {
	// Prometheus metrics
	RequestsEnqueued   prometheus.Counter
	RequestsProcessed  prometheus.Counter
	QueueLength        prometheus.Gauge
	ProcessingDuration prometheus.Histogram
	PriorityScores     prometheus.Histogram

	// Internal metrics
	mu sync.RWMutex

	// Wait time tracking
	waitTimes []time.Duration
	maxWaitTime time.Duration
	minWaitTime time.Duration

	// Processing time tracking
	processingTimes []time.Duration
	lastProcessTime time.Time

	// Request counts
	totalRequests   int64
	totalAllocations int64

	// Fairness metrics
	groupAllocations map[string]int64
	groupWaitTimes   map[string][]time.Duration
}

// NewMetrics creates a new metrics instance
func NewMetrics() *Metrics {
	m := &Metrics{
		RequestsEnqueued: promauto.NewCounter(prometheus.CounterOpts{
			Name: "fairrent_requests_enqueued_total",
			Help: "Total number of requests enqueued",
		}),
		RequestsProcessed: promauto.NewCounter(prometheus.CounterOpts{
			Name: "fairrent_requests_processed_total",
			Help: "Total number of requests processed",
		}),
		QueueLength: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "fairrent_queue_length",
			Help: "Current number of requests in queue",
		}),
		ProcessingDuration: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "fairrent_processing_duration_seconds",
			Help:    "Time taken to process requests",
			Buckets: prometheus.DefBuckets,
		}),
		PriorityScores: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "fairrent_priority_scores",
			Help:    "Distribution of priority scores",
			Buckets: prometheus.LinearBuckets(0, 1, 20),
		}),
		waitTimes:        make([]time.Duration, 0),
		processingTimes:  make([]time.Duration, 0),
		groupAllocations: make(map[string]int64),
		groupWaitTimes:   make(map[string][]time.Duration),
	}

	return m
}

// RecordRequestEnqueued records a new request being enqueued
func (m *Metrics) RecordRequestEnqueued(userGroup string) {
	m.RequestsEnqueued.Inc()
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.totalRequests++
	if _, exists := m.groupAllocations[userGroup]; !exists {
		m.groupAllocations[userGroup] = 0
		m.groupWaitTimes[userGroup] = make([]time.Duration, 0)
	}
}

// RecordRequestProcessed records a request being processed
func (m *Metrics) RecordRequestProcessed(userGroup string, waitTime time.Duration, priorityScore float64) {
	m.RequestsProcessed.Inc()
	
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.totalAllocations++
	m.groupAllocations[userGroup]++
	
	// Record wait time
	m.waitTimes = append(m.waitTimes, waitTime)
	if len(m.waitTimes) > 1000 { // Keep only last 1000 for memory efficiency
		m.waitTimes = m.waitTimes[1:]
	}
	
	// Update min/max wait times
	if waitTime > m.maxWaitTime {
		m.maxWaitTime = waitTime
	}
	if m.minWaitTime == 0 || waitTime < m.minWaitTime {
		m.minWaitTime = waitTime
	}
	
	// Record group-specific wait time
	if groupTimes, exists := m.groupWaitTimes[userGroup]; exists {
		m.groupWaitTimes[userGroup] = append(groupTimes, waitTime)
		if len(m.groupWaitTimes[userGroup]) > 100 { // Keep only last 100 per group
			m.groupWaitTimes[userGroup] = m.groupWaitTimes[userGroup][1:]
		}
	}
	
	// Record priority score
	m.PriorityScores.Observe(priorityScore)
	
	// Record processing duration
	now := time.Now()
	if !m.lastProcessTime.IsZero() {
		processingTime := now.Sub(m.lastProcessTime)
		m.processingTimes = append(m.processingTimes, processingTime)
		if len(m.processingTimes) > 1000 {
			m.processingTimes = m.processingTimes[1:]
		}
		m.ProcessingDuration.Observe(processingTime.Seconds())
	}
	m.lastProcessTime = now
}

// GetMetrics returns computed metrics
func (m *Metrics) GetMetrics() *SchedulerMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	metrics := &SchedulerMetrics{
		TotalRequests:   m.totalRequests,
		TotalAllocations: m.totalAllocations,
		QueueLength:     int64(m.QueueLength.(prometheus.Gauge).(prometheus.Gauge)),
	}
	
	// Calculate wait time statistics
	if len(m.waitTimes) > 0 {
		metrics.AverageWaitTime = m.calculateAverageWaitTime()
		metrics.MedianWaitTime = m.calculateMedianWaitTime()
		metrics.P95WaitTime = m.calculatePercentileWaitTime(95)
		metrics.P99WaitTime = m.calculatePercentileWaitTime(99)
		metrics.MaxWaitTime = m.maxWaitTime
		metrics.MinWaitTime = m.minWaitTime
		
		// Calculate fairness metrics
		metrics.MaxWaitTimeRatio = float64(m.maxWaitTime) / float64(m.minWaitTime)
		metrics.GiniCoefficient = m.calculateGiniCoefficient()
	}
	
	// Calculate processing statistics
	if len(m.processingTimes) > 0 {
		metrics.AverageProcessingTime = m.calculateAverageProcessingTime()
		metrics.AllocationRate = m.calculateAllocationRate()
		metrics.QueueTurnoverRate = m.calculateQueueTurnoverRate()
	}
	
	return metrics
}

// calculateAverageWaitTime computes the average wait time
func (m *Metrics) calculateAverageWaitTime() time.Duration {
	if len(m.waitTimes) == 0 {
		return 0
	}
	
	total := time.Duration(0)
	for _, waitTime := range m.waitTimes {
		total += waitTime
	}
	return total / time.Duration(len(m.waitTimes))
}

// calculateMedianWaitTime computes the median wait time
func (m *Metrics) calculateMedianWaitTime() time.Duration {
	if len(m.waitTimes) == 0 {
		return 0
	}
	
	// Create a copy to avoid modifying the original slice
	times := make([]time.Duration, len(m.waitTimes))
	copy(times, m.waitTimes)
	
	sort.Slice(times, func(i, j int) bool {
		return times[i] < times[j]
	})
	
	mid := len(times) / 2
	if len(times)%2 == 0 {
		return (times[mid-1] + times[mid]) / 2
	}
	return times[mid]
}

// calculatePercentileWaitTime computes the nth percentile wait time
func (m *Metrics) calculatePercentileWaitTime(percentile int) time.Duration {
	if len(m.waitTimes) == 0 {
		return 0
	}
	
	// Create a copy to avoid modifying the original slice
	times := make([]time.Duration, len(m.waitTimes))
	copy(times, m.waitTimes)
	
	sort.Slice(times, func(i, j int) bool {
		return times[i] < times[j]
	})
	
	index := int(float64(percentile) / 100.0 * float64(len(times)-1))
	return times[index]
}

// calculateAverageProcessingTime computes the average processing time
func (m *Metrics) calculateAverageProcessingTime() time.Duration {
	if len(m.processingTimes) == 0 {
		return 0
	}
	
	total := time.Duration(0)
	for _, processingTime := range m.processingTimes {
		total += processingTime
	}
	return total / time.Duration(len(m.processingTimes))
}

// calculateAllocationRate computes allocations per hour
func (m *Metrics) calculateAllocationRate() float64 {
	if len(m.processingTimes) == 0 {
		return 0
	}
	
	// Calculate based on recent processing times
	recentCount := min(100, len(m.processingTimes))
	recentTimes := m.processingTimes[len(m.processingTimes)-recentCount:]
	
	totalTime := time.Duration(0)
	for _, t := range recentTimes {
		totalTime += t
	}
	
	if totalTime == 0 {
		return 0
	}
	
	// Convert to allocations per hour
	return float64(recentCount) / totalTime.Hours()
}

// calculateQueueTurnoverRate computes requests processed per hour
func (m *Metrics) calculateQueueTurnoverRate() float64 {
	if len(m.processingTimes) == 0 {
		return 0
	}
	
	// Similar to allocation rate but for all requests
	recentCount := min(100, len(m.processingTimes))
	recentTimes := m.processingTimes[len(m.processingTimes)-recentCount:]
	
	totalTime := time.Duration(0)
	for _, t := range recentTimes {
		totalTime += t
	}
	
	if totalTime == 0 {
		return 0
	}
	
	return float64(recentCount) / totalTime.Hours()
}

// calculateGiniCoefficient computes wait time inequality
func (m *Metrics) calculateGiniCoefficient() float64 {
	if len(m.waitTimes) < 2 {
		return 0
	}
	
	// Create a copy and sort
	times := make([]time.Duration, len(m.waitTimes))
	copy(times, m.waitTimes)
	sort.Slice(times, func(i, j int) bool {
		return times[i] < times[j]
	})
	
	// Convert to float64 for calculations
	values := make([]float64, len(times))
	for i, t := range times {
		values[i] = float64(t.Milliseconds())
	}
	
	// Calculate Gini coefficient
	n := float64(len(values))
	sum := 0.0
	for i, value := range values {
		sum += (2*float64(i+1) - n - 1) * value
	}
	
	totalSum := 0.0
	for _, value := range values {
		totalSum += value
	}
	
	if totalSum == 0 {
		return 0
	}
	
	return sum / (n * totalSum)
}

// GetAverageProcessingTime returns the average processing time
func (m *Metrics) GetAverageProcessingTime() time.Duration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.calculateAverageProcessingTime()
}

// SchedulerMetrics contains computed scheduler metrics
type SchedulerMetrics struct {
	TotalRequests        int64
	TotalAllocations     int64
	QueueLength          int64
	AverageWaitTime      time.Duration
	MedianWaitTime       time.Duration
	P95WaitTime          time.Duration
	P99WaitTime          time.Duration
	MaxWaitTime          time.Duration
	MinWaitTime          time.Duration
	MaxWaitTimeRatio     float64
	GiniCoefficient      float64
	AverageProcessingTime time.Duration
	AllocationRate       float64
	QueueTurnoverRate    float64
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
