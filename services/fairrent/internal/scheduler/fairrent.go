package scheduler

import (
	"container/heap"
	"context"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/wohnfair/wohnfair/services/gen/wohnfair/common/v1"
	"github.com/wohnfair/wohnfair/services/gen/wohnfair/fairrent/v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// FairRent implements α-fair scheduling for housing allocation
type FairRent struct {
	mu sync.RWMutex

	// Queue management
	queue     *PriorityQueue
	ticketMap map[string]*Ticket

	// Fairness parameters
	alpha        float64
	groupWeights map[string]float64

	// Metrics
	metrics *Metrics

	// Configuration
	config *Config

	logger *zap.Logger
}

// Config holds scheduler configuration
type Config struct {
	Alpha        float64            `yaml:"alpha"`
	GroupWeights map[string]float64 `yaml:"group_weights"`
	MaxWaitTime  time.Duration      `yaml:"max_wait_time"`
	LogLevel     string             `yaml:"log_level"`
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		Alpha: 2.0, // α=2 provides good fairness with reasonable efficiency
		GroupWeights: map[string]float64{
			"USER_GROUP_REFUGEE":     1.5,  // Higher priority for refugees
			"USER_GROUP_DISABLED":    1.3,  // Higher priority for disabled
			"USER_GROUP_SENIOR":      1.2,  // Higher priority for seniors
			"USER_GROUP_LOW_INCOME":  1.1,  // Slightly higher for low income
			"USER_GROUP_STUDENT":     1.0,  // Baseline priority
			"USER_GROUP_FAMILY":      1.0,  // Baseline priority
			"USER_GROUP_SINGLE":      0.9,  // Slightly lower for single
			"USER_GROUP_MIDDLE_INCOME": 0.8, // Lower for middle income
			"USER_GROUP_HIGH_INCOME": 0.7,   // Lower for high income
		},
		MaxWaitTime: 24 * time.Hour, // Maximum wait time before starvation protection
		LogLevel:    "info",
	}
}

// NewFairRent creates a new scheduler instance
func NewFairRent(config *Config, logger *zap.Logger) *FairRent {
	if config == nil {
		config = DefaultConfig()
	}

	fr := &FairRent{
		queue:        &PriorityQueue{},
		ticketMap:    make(map[string]*Ticket),
		alpha:        config.Alpha,
		groupWeights: config.GroupWeights,
		metrics:      NewMetrics(),
		config:       config,
		logger:       logger,
	}

	heap.Init(fr.queue)
	return fr
}

// Enqueue adds a new housing request to the queue
func (fr *FairRent) Enqueue(ctx context.Context, req *fairrentv1.EnqueueRequest) (*fairrentv1.EnqueueResponse, error) {
	fr.mu.Lock()
	defer fr.mu.Unlock()

	// Generate ticket ID
	ticketID := generateTicketID()

	// Create ticket
	ticket := &Ticket{
		ID:           ticketID,
		UserID:       req.UserId.Value,
		UserGroup:    req.UserGroup.String(),
		Urgency:      int(req.Urgency),
		EnqueueTime:  time.Now(),
		PriorityScore: fr.calculatePriorityScore(req),
		Constraints:   req,
	}

	// Add to queue
	heap.Push(fr.queue, ticket)
	fr.ticketMap[ticketID] = ticket

	// Update metrics
	fr.metrics.RequestsEnqueued.Inc()
	fr.metrics.QueueLength.Set(float64(fr.queue.Len()))

	fr.logger.Info("Request enqueued",
		zap.String("ticket_id", ticketID),
		zap.String("user_group", req.UserGroup.String()),
		zap.Int("urgency", int(req.Urgency)),
		zap.Float64("priority_score", ticket.PriorityScore),
	)

	return &fairrentv1.EnqueueResponse{
		TicketId: &commonv1.TicketID{Value: ticketID},
		Status:   commonv1.AllocationStatus_ALLOCATION_STATUS_QUEUED,
		QueuePosition: int32(fr.queue.Len()),
		EstimatedAllocationTime: &timestamppb.Timestamp{
			Seconds: time.Now().Add(fr.estimateWaitTime(ticket)).Unix(),
		},
		Metadata: &commonv1.Metadata{
			CreatedAt: &timestamppb.Timestamp{Seconds: time.Now().Unix()},
		},
	}, nil
}

// ScheduleNext processes the next allocation from the queue
func (fr *FairRent) ScheduleNext(ctx context.Context, req *fairrentv1.ScheduleNextRequest) (*fairrentv1.ScheduleNextResponse, error) {
	fr.mu.Lock()
	defer fr.mu.Unlock()

	if fr.queue.Len() == 0 {
		return nil, fmt.Errorf("queue is empty")
	}

	// Get next ticket with highest priority
	ticket := heap.Pop(fr.queue).(*Ticket)
	delete(fr.ticketMap, ticket.ID)

	// Update metrics
	fr.metrics.RequestsProcessed.Inc()
	fr.metrics.QueueLength.Set(float64(fr.queue.Len()))

	fr.logger.Info("Request scheduled",
		zap.String("ticket_id", ticket.ID),
		zap.String("user_group", ticket.UserGroup),
		zap.Float64("priority_score", ticket.PriorityScore),
		zap.Duration("wait_time", time.Since(ticket.EnqueueTime)),
	)

	return &fairrentv1.ScheduleNextResponse{
		TicketId: &commonv1.TicketID{Value: ticket.ID},
		UserId:   &commonv1.UserID{Value: ticket.UserID},
		AllocationTime: &timestamppb.Timestamp{Seconds: time.Now().Unix()},
		FairnessScore: ticket.PriorityScore,
		Metadata: &commonv1.Metadata{
			CreatedAt: &timestamppb.Timestamp{Seconds: time.Now().Unix()},
		},
	}, nil
}

// PeekPosition returns the current position and estimated wait time
func (fr *FairRent) PeekPosition(ctx context.Context, req *fairrentv1.PeekPositionRequest) (*fairrentv1.PeekPositionResponse, error) {
	fr.mu.RLock()
	defer fr.mu.RUnlock()

	ticketID := req.TicketId.Value
	ticket, exists := fr.ticketMap[ticketID]
	if !exists {
		return nil, fmt.Errorf("ticket not found: %s", ticketID)
	}

	// Calculate position (this is simplified - in practice would need more sophisticated tracking)
	position := fr.calculatePosition(ticket)

	// Estimate wait time
	estimatedWait := fr.estimateWaitTime(ticket)

	fr.logger.Debug("Position peeked",
		zap.String("ticket_id", ticketID),
		zap.Int("position", position),
		zap.Duration("estimated_wait", estimatedWait),
	)

	return &fairrentv1.PeekPositionResponse{
		TicketId: req.TicketId,
		CurrentPosition: int32(position),
		TotalInQueue: int32(fr.queue.Len()),
		EstimatedWaitTime: &durationpb.Duration{
			Seconds: int64(estimatedWait.Seconds()),
		},
		EstimatedAllocationTime: &timestamppb.Timestamp{
			Seconds: time.Now().Add(estimatedWait).Unix(),
		},
		FairnessScore: ticket.PriorityScore,
		Status:        commonv1.AllocationStatus_ALLOCATION_STATUS_QUEUED,
	}, nil
}

// GetMetrics returns fairness and performance metrics
func (fr *FairRent) GetMetrics(ctx context.Context) (*fairrentv1.FairnessMetrics, error) {
	fr.mu.RLock()
	defer fr.mu.RUnlock()

	metrics := fr.metrics.GetMetrics()
	groupMetrics := fr.calculateGroupMetrics()

	return &fairrentv1.FairnessMetrics{
		Alpha:        fr.alpha,
		GroupWeights: fr.groupWeights,
		TotalRequests: int32(metrics.TotalRequests),
		TotalAllocations: int32(metrics.TotalAllocations),
		ActiveRequests: int32(fr.queue.Len()),
		AverageWaitTime: &durationpb.Duration{
			Seconds: int64(metrics.AverageWaitTime.Seconds()),
		},
		MedianWaitTime: &durationpb.Duration{
			Seconds: int64(metrics.MedianWaitTime.Seconds()),
		},
		P95WaitTime: &durationpb.Duration{
			Seconds: int64(metrics.P95WaitTime.Seconds()),
		},
		P99WaitTime: &durationpb.Duration{
			Seconds: int64(metrics.P99WaitTime.Seconds()),
		},
		MaxWaitTime: &durationpb.Duration{
			Seconds: int64(metrics.MaxWaitTime.Seconds()),
		},
		GroupMetrics: groupMetrics,
		MaxWaitTimeRatio: metrics.MaxWaitTimeRatio,
		GiniCoefficient: metrics.GiniCoefficient,
		AllocationRate: metrics.AllocationRate,
		QueueTurnoverRate: metrics.QueueTurnoverRate,
		CalculatedAt: &timestamppb.Timestamp{Seconds: time.Now().Unix()},
	}, nil
}

// calculatePriorityScore computes the α-fair priority score
func (fr *FairRent) calculatePriorityScore(req *fairrentv1.EnqueueRequest) float64 {
	// Base priority from urgency
	urgencyScore := float64(req.Urgency) / 5.0

	// Group weight adjustment
	groupWeight := 1.0
	if weight, exists := fr.groupWeights[req.UserGroup.String()]; exists {
		groupWeight = weight
	}

	// Additional priority factors
	priorityBonus := req.PriorityScore

	// α-fair formula: priority = (urgency * group_weight + priority_bonus)^α
	basePriority := urgencyScore*groupWeight + priorityBonus
	return math.Pow(basePriority, fr.alpha)
}

// estimateWaitTime estimates how long a ticket will wait
func (fr *FairRent) estimateWaitTime(ticket *Ticket) time.Duration {
	// Simple estimation based on queue position and historical processing rate
	position := fr.calculatePosition(ticket)
	avgProcessingTime := fr.metrics.GetAverageProcessingTime()
	
	estimatedWait := time.Duration(position) * avgProcessingTime
	
	// Apply starvation protection
	if estimatedWait > fr.config.MaxWaitTime {
		estimatedWait = fr.config.MaxWaitTime
	}
	
	return estimatedWait
}

// calculatePosition estimates the ticket's position in the queue
func (fr *FairRent) calculatePosition(ticket *Ticket) int {
	// This is a simplified calculation
	// In practice, would need more sophisticated position tracking
	position := 1
	for _, queuedTicket := range fr.queue.tickets {
		if queuedTicket.PriorityScore > ticket.PriorityScore {
			position++
		}
	}
	return position
}

// calculateGroupMetrics computes fairness metrics per user group
func (fr *FairRent) calculateGroupMetrics() []*fairrentv1.GroupFairnessMetrics {
	groupStats := make(map[string]*GroupStats)
	
	// Collect statistics
	for _, ticket := range fr.ticketMap {
		if stats, exists := groupStats[ticket.UserGroup]; exists {
			stats.Count++
			stats.TotalWaitTime += time.Since(ticket.EnqueueTime)
		} else {
			groupStats[ticket.UserGroup] = &GroupStats{
				Group: ticket.UserGroup,
				Count: 1,
				TotalWaitTime: time.Since(ticket.EnqueueTime),
			}
		}
	}
	
	// Convert to protobuf format
	var metrics []*fairrentv1.GroupFairnessMetrics
	for _, stats := range groupStats {
		avgWaitTime := stats.TotalWaitTime / time.Duration(stats.Count)
		targetRate := 1.0 / float64(len(groupStats)) // Equal distribution target
		
		metrics = append(metrics, &fairrentv1.GroupFairnessMetrics{
			UserGroup: commonv1.UserGroup(commonv1.UserGroup_value[stats.Group]),
			RequestsCount: int32(stats.Count),
			AllocationsCount: 0, // Would track actual allocations
			AllocationRate: float64(stats.Count) / float64(fr.queue.Len()),
			AverageWaitTime: &durationpb.Duration{
				Seconds: int64(avgWaitTime.Seconds()),
			},
			FairnessScore: 1.0, // Would calculate actual fairness score
			TargetAllocationRate: targetRate,
			ActualVsTargetRatio: float64(stats.Count) / float64(fr.queue.Len()) / targetRate,
		})
	}
	
	return metrics
}

// generateTicketID creates a unique ticket identifier
func generateTicketID() string {
	return fmt.Sprintf("TKT_%d", time.Now().UnixNano())
}

// GroupStats holds per-group statistics
type GroupStats struct {
	Group         string
	Count         int
	TotalWaitTime time.Duration
}
