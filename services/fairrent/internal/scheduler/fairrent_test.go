package scheduler

import (
	"context"
	"fmt"
	"math"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wohnfair/wohnfair/services/gen/wohnfair/common/v1"
	"github.com/wohnfair/wohnfair/services/gen/wohnfair/fairrent/v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestNewFairRent(t *testing.T) {
	logger := zap.NewNop()
	
	// Test with default config
	fr := NewFairRent(nil, logger)
	assert.NotNil(t, fr)
	assert.Equal(t, 2.0, fr.alpha)
	assert.NotNil(t, fr.groupWeights)
	assert.Equal(t, 1.5, fr.groupWeights["USER_GROUP_REFUGEE"])
	
	// Test with custom config
	customConfig := &Config{
		Alpha: 1.5,
		GroupWeights: map[string]float64{
			"USER_GROUP_STUDENT": 2.0,
		},
	}
	
	fr2 := NewFairRent(customConfig, logger)
	assert.Equal(t, 1.5, fr2.alpha)
	assert.Equal(t, 2.0, fr2.groupWeights["USER_GROUP_STUDENT"])
}

func TestFairRent_Enqueue(t *testing.T) {
	logger := zap.NewNop()
	fr := NewFairRent(nil, logger)
	ctx := context.Background()
	
	// Test basic enqueue
	req := &fairrentv1.EnqueueRequest{
		UserId: &commonv1.UserID{Value: "user1"},
		UserGroup: commonv1.UserGroup_USER_GROUP_STUDENT,
		Urgency: commonv1.UrgencyLevel_URGENCY_LEVEL_HIGH,
		PriorityScore: 0.5,
	}
	
	resp, err := fr.Enqueue(ctx, req)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.TicketId.Value)
	assert.Equal(t, commonv1.AllocationStatus_ALLOCATION_STATUS_QUEUED, resp.Status)
	assert.Equal(t, int32(1), resp.QueuePosition)
	
	// Verify ticket was added to map
	ticket, exists := fr.ticketMap[resp.TicketId.Value]
	assert.True(t, exists)
	assert.Equal(t, "user1", ticket.UserID)
	assert.Equal(t, "USER_GROUP_STUDENT", ticket.UserGroup)
	assert.Equal(t, 5, ticket.Urgency)
	
	// Test enqueue with different user group
	req2 := &fairrentv1.EnqueueRequest{
		UserId: &commonv1.UserID{Value: "user2"},
		UserGroup: commonv1.UserGroup_USER_GROUP_REFUGEE,
		Urgency: commonv1.UrgencyLevel_URGENCY_LEVEL_CRITICAL,
		PriorityScore: 0.8,
	}
	
	resp2, err := fr.Enqueue(ctx, req2)
	require.NoError(t, err)
	assert.NotEmpty(t, resp2.TicketId.Value)
	
	// Verify queue length
	assert.Equal(t, 2, fr.queue.Len())
}

func TestFairRent_ScheduleNext(t *testing.T) {
	logger := zap.NewNop()
	fr := NewFairRent(nil, logger)
	ctx := context.Background()
	
	// Enqueue multiple requests with different priorities
	requests := []struct {
		userID    string
		userGroup commonv1.UserGroup
		urgency   commonv1.UrgencyLevel
		priority  float64
	}{
		{"user1", commonv1.UserGroup_USER_GROUP_STUDENT, commonv1.UrgencyLevel_URGENCY_LEVEL_LOW, 0.1},
		{"user2", commonv1.UserGroup_USER_GROUP_REFUGEE, commonv1.UrgencyLevel_URGENCY_LEVEL_CRITICAL, 0.9},
		{"user3", commonv1.UserGroup_USER_GROUP_SENIOR, commonv1.UrgencyLevel_URGENCY_LEVEL_HIGH, 0.5},
	}
	
	for _, req := range requests {
		enqueueReq := &fairrentv1.EnqueueRequest{
			UserId: &commonv1.UserID{Value: req.userID},
			UserGroup: req.userGroup,
			Urgency: req.urgency,
			PriorityScore: req.priority,
		}
		_, err := fr.Enqueue(ctx, enqueueReq)
		require.NoError(t, err)
	}
	
	// Schedule next request
	scheduleReq := &fairrentv1.ScheduleNextRequest{
		Horizon: &commonv1.SchedulingHorizon{
			LookAhead: &durationpb.Duration{Seconds: 3600}, // 1 hour
		},
	}
	
	resp, err := fr.ScheduleNext(ctx, scheduleReq)
	require.NoError(t, err)
	
	// Should schedule the refugee user (highest priority due to group weight + urgency)
	assert.Equal(t, "user2", resp.UserId.Value)
	assert.Equal(t, int32(2), resp.QueuePosition)
	
	// Verify queue length decreased
	assert.Equal(t, 2, fr.queue.Len())
	
	// Verify ticket was removed from map
	_, exists := fr.ticketMap[resp.TicketId.Value]
	assert.False(t, exists)
}

func TestFairRent_PeekPosition(t *testing.T) {
	logger := zap.NewNop()
	fr := NewFairRent(nil, logger)
	ctx := context.Background()
	
	// Enqueue a request
	req := &fairrentv1.EnqueueRequest{
		UserId: &commonv1.UserID{Value: "user1"},
		UserGroup: commonv1.UserGroup_USER_GROUP_STUDENT,
		Urgency: commonv1.UrgencyLevel_URGENCY_LEVEL_MEDIUM,
		PriorityScore: 0.3,
	}
	
	resp, err := fr.Enqueue(ctx, req)
	require.NoError(t, err)
	
	// Peek position
	peekReq := &fairrentv1.PeekPositionRequest{
		TicketId: resp.TicketId,
	}
	
	peekResp, err := fr.PeekPosition(ctx, peekReq)
	require.NoError(t, err)
	
	assert.Equal(t, resp.TicketId.Value, peekResp.TicketId.Value)
	assert.Equal(t, int32(1), peekResp.CurrentPosition)
	assert.Equal(t, int32(1), peekResp.TotalInQueue)
	assert.Equal(t, commonv1.AllocationStatus_ALLOCATION_STATUS_QUEUED, peekResp.Status)
	
	// Test with non-existent ticket
	nonExistentReq := &fairrentv1.PeekPositionRequest{
		TicketId: &commonv1.TicketID{Value: "nonexistent"},
	}
	
	_, err = fr.PeekPosition(ctx, nonExistentReq)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ticket not found")
}

func TestFairRent_GetMetrics(t *testing.T) {
	logger := zap.NewNop()
	fr := NewFairRent(nil, logger)
	ctx := context.Background()
	
	// Enqueue some requests
	requests := []struct {
		userID    string
		userGroup commonv1.UserGroup
		urgency   commonv1.UrgencyLevel
	}{
		{"user1", commonv1.UserGroup_USER_GROUP_STUDENT, commonv1.UrgencyLevel_URGENCY_LEVEL_LOW},
		{"user2", commonv1.UserGroup_USER_GROUP_REFUGEE, commonv1.UrgencyLevel_URGENCY_LEVEL_HIGH},
	}
	
	for _, req := range requests {
		enqueueReq := &fairrentv1.EnqueueRequest{
			UserId: &commonv1.UserID{Value: req.userID},
			UserGroup: req.userGroup,
			Urgency: req.urgency,
		}
		_, err := fr.Enqueue(ctx, enqueueReq)
		require.NoError(t, err)
	}
	
	// Get metrics
	metrics, err := fr.GetMetrics(ctx)
	require.NoError(t, err)
	
	assert.Equal(t, 2.0, metrics.Alpha)
	assert.Equal(t, int32(2), metrics.TotalRequests)
	assert.Equal(t, int32(0), metrics.TotalAllocations)
	assert.Equal(t, int32(2), metrics.ActiveRequests)
	
	// Verify group weights are included
	assert.Equal(t, 1.5, metrics.GroupWeights["USER_GROUP_REFUGEE"])
	assert.Equal(t, 1.0, metrics.GroupWeights["USER_GROUP_STUDENT"])
	
	// Verify group metrics
	assert.Len(t, metrics.GroupMetrics, 2)
	
	// Find refugee group metrics
	var refugeeMetrics *fairrentv1.GroupFairnessMetrics
	for _, gm := range metrics.GroupMetrics {
		if gm.UserGroup == commonv1.UserGroup_USER_GROUP_REFUGEE {
			refugeeMetrics = gm
			break
		}
	}
	
	require.NotNil(t, refugeeMetrics)
	assert.Equal(t, int32(1), refugeeMetrics.RequestsCount)
	assert.Equal(t, float64(0.5), refugeeMetrics.AllocationRate) // 1 out of 2 total
}

func TestFairRent_CalculatePriorityScore(t *testing.T) {
	logger := zap.NewNop()
	fr := NewFairRent(nil, logger)
	
	// Test with different urgency levels
	testCases := []struct {
		urgency      commonv1.UrgencyLevel
		userGroup    commonv1.UserGroup
		priorityBonus float64
		expected     float64
	}{
		{
			urgency:       commonv1.UrgencyLevel_URGENCY_LEVEL_LOW,
			userGroup:     commonv1.UserGroup_USER_GROUP_STUDENT,
			priorityBonus: 0.0,
			expected:      math.Pow(0.2*1.0+0.0, 2.0), // (0.2)^2 = 0.04
		},
		{
			urgency:       commonv1.UrgencyLevel_URGENCY_LEVEL_HIGH,
			userGroup:     commonv1.UserGroup_USER_GROUP_REFUGEE,
			priorityBonus: 0.5,
			expected:      math.Pow(0.6*1.5+0.5, 2.0), // (1.4)^2 = 1.96
		},
		{
			urgency:       commonv1.UrgencyLevel_URGENCY_LEVEL_CRITICAL,
			userGroup:     commonv1.UserGroup_USER_GROUP_DISABLED,
			priorityBonus: 0.8,
			expected:      math.Pow(0.8*1.3+0.8, 2.0), // (1.84)^2 = 3.3856
		},
	}
	
	for _, tc := range testCases {
		req := &fairrentv1.EnqueueRequest{
			UserId: &commonv1.UserID{Value: "test"},
			UserGroup: tc.userGroup,
			Urgency: tc.urgency,
			PriorityScore: tc.priorityBonus,
		}
		
		score := fr.calculatePriorityScore(req)
		assert.InDelta(t, tc.expected, score, 0.01, 
			"Priority score mismatch for urgency=%v, group=%v, bonus=%v", 
			tc.urgency, tc.userGroup, tc.priorityBonus)
	}
}

func TestFairRent_EstimateWaitTime(t *testing.T) {
	logger := zap.NewNop()
	fr := NewFairRent(nil, logger)
	
	// Create a ticket
	ticket := &Ticket{
		ID:           "test",
		UserID:       "user1",
		UserGroup:    "USER_GROUP_STUDENT",
		Urgency:      3,
		EnqueueTime:  time.Now(),
		PriorityScore: 1.0,
	}
	
	// Test wait time estimation
	waitTime := fr.estimateWaitTime(ticket)
	
	// Should be positive
	assert.True(t, waitTime > 0)
	
	// Should not exceed max wait time
	assert.True(t, waitTime <= fr.config.MaxWaitTime)
}

func TestFairRent_CalculatePosition(t *testing.T) {
	logger := zap.NewNop()
	fr := NewFairRent(nil, logger)
	
	// Add some tickets to the queue
	tickets := []*Ticket{
		{ID: "1", PriorityScore: 1.0},
		{ID: "2", PriorityScore: 2.0},
		{ID: "3", PriorityScore: 3.0},
	}
	
	for _, ticket := range tickets {
		heap.Push(fr.queue, ticket)
		fr.ticketMap[ticket.ID] = ticket
	}
	
	// Test position calculation
	position := fr.calculatePosition(tickets[0]) // Lowest priority
	assert.Equal(t, 3, position) // Should be last
	
	position = fr.calculatePosition(tickets[2]) // Highest priority
	assert.Equal(t, 1, position) // Should be first
}

func TestFairRent_CalculateGroupMetrics(t *testing.T) {
	logger := zap.NewNop()
	fr := NewFairRent(nil, logger)
	
	// Add tickets for different groups
	groups := []string{"USER_GROUP_STUDENT", "USER_GROUP_REFUGEE", "USER_GROUP_SENIOR"}
	
	for i, group := range groups {
		ticket := &Ticket{
			ID:           fmt.Sprintf("ticket_%d", i),
			UserID:       fmt.Sprintf("user_%d", i),
			UserGroup:    group,
			Urgency:      3,
			EnqueueTime:  time.Now().Add(-time.Duration(i) * time.Hour),
			PriorityScore: float64(i + 1),
		}
		fr.ticketMap[ticket.ID] = ticket
	}
	
	// Calculate group metrics
	metrics := fr.calculateGroupMetrics()
	
	// Should have metrics for all groups
	assert.Len(t, metrics, 3)
	
	// Verify each group has metrics
	groupNames := make(map[string]bool)
	for _, gm := range metrics {
		groupNames[gm.UserGroup.String()] = true
	}
	
	for _, group := range groups {
		assert.True(t, groupNames[group], "Missing metrics for group: %s", group)
	}
}

func TestFairRent_Concurrency(t *testing.T) {
	logger := zap.NewNop()
	fr := NewFairRent(nil, logger)
	ctx := context.Background()
	
	// Test concurrent enqueue operations
	const numGoroutines = 10
	const requestsPerGoroutine = 10
	
	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines*requestsPerGoroutine)
	
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			
			for j := 0; j < requestsPerGoroutine; j++ {
				req := &fairrentv1.EnqueueRequest{
					UserId: &commonv1.UserID{Value: fmt.Sprintf("user_%d_%d", goroutineID, j)},
					UserGroup: commonv1.UserGroup_USER_GROUP_STUDENT,
					Urgency: commonv1.UrgencyLevel_URGENCY_LEVEL_MEDIUM,
					PriorityScore: float64(j) / float64(requestsPerGoroutine),
				}
				
				_, err := fr.Enqueue(ctx, req)
				if err != nil {
					errors <- err
				}
			}
		}(i)
	}
	
	wg.Wait()
	close(errors)
	
	// Check for errors
	for err := range errors {
		t.Errorf("Concurrent enqueue failed: %v", err)
	}
	
	// Verify all requests were enqueued
	assert.Equal(t, numGoroutines*requestsPerGoroutine, fr.queue.Len())
	assert.Equal(t, numGoroutines*requestsPerGoroutine, len(fr.ticketMap))
}

func TestFairRent_StarvationPrevention(t *testing.T) {
	logger := zap.NewNop()
	
	// Create config with very short max wait time
	config := &Config{
		Alpha:       1.0,
		MaxWaitTime: 100 * time.Millisecond,
		GroupWeights: map[string]float64{
			"USER_GROUP_STUDENT": 1.0,
			"USER_GROUP_REFUGEE": 1.5,
		},
	}
	
	fr := NewFairRent(config, logger)
	ctx := context.Background()
	
	// Enqueue a low-priority request first
	lowPriorityReq := &fairrentv1.EnqueueRequest{
		UserId: &commonv1.UserID{Value: "low_priority"},
		UserGroup: commonv1.UserGroup_USER_GROUP_STUDENT,
		Urgency: commonv1.UrgencyLevel_URGENCY_LEVEL_LOW,
		PriorityScore: 0.0,
	}
	
	lowResp, err := fr.Enqueue(ctx, lowPriorityReq)
	require.NoError(t, err)
	
	// Wait a bit
	time.Sleep(50 * time.Millisecond)
	
	// Enqueue a high-priority request
	highPriorityReq := &fairrentv1.EnqueueRequest{
		UserId: &commonv1.UserID{Value: "high_priority"},
		UserGroup: commonv1.UserGroup_USER_GROUP_REFUGEE,
		Urgency: commonv1.UrgencyLevel_URGENCY_LEVEL_CRITICAL,
		PriorityScore: 1.0,
	}
	
	highResp, err := fr.Enqueue(ctx, highPriorityReq)
	require.NoError(t, err)
	
	// Wait for starvation protection to kick in
	time.Sleep(100 * time.Millisecond)
	
	// Schedule next - should prioritize the low-priority request due to starvation protection
	scheduleReq := &fairrentv1.ScheduleNextRequest{
		Horizon: &commonv1.SchedulingHorizon{
			LookAhead: &durationpb.Duration{Seconds: 3600},
		},
	}
	
	resp, err := fr.ScheduleNext(ctx, scheduleReq)
	require.NoError(t, err)
	
	// Should schedule the low-priority request first due to starvation protection
	assert.Equal(t, "low_priority", resp.UserId.Value)
}

// Benchmark tests
func BenchmarkFairRent_Enqueue(b *testing.B) {
	logger := zap.NewNop()
	fr := NewFairRent(nil, logger)
	ctx := context.Background()
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		req := &fairrentv1.EnqueueRequest{
			UserId: &commonv1.UserID{Value: fmt.Sprintf("user_%d", i)},
			UserGroup: commonv1.UserGroup_USER_GROUP_STUDENT,
			Urgency: commonv1.UrgencyLevel_URGENCY_LEVEL_MEDIUM,
			PriorityScore: float64(i) / float64(b.N),
		}
		
		_, err := fr.Enqueue(ctx, req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFairRent_ScheduleNext(b *testing.B) {
	logger := zap.NewNop()
	fr := NewFairRent(nil, logger)
	ctx := context.Background()
	
	// Pre-populate queue
	for i := 0; i < 1000; i++ {
		req := &fairrentv1.EnqueueRequest{
			UserId: &commonv1.UserID{Value: fmt.Sprintf("user_%d", i)},
			UserGroup: commonv1.UserGroup_USER_GROUP_STUDENT,
			Urgency: commonv1.UrgencyLevel_URGENCY_LEVEL_MEDIUM,
			PriorityScore: float64(i) / 1000.0,
		}
		
		_, err := fr.Enqueue(ctx, req)
		if err != nil {
			b.Fatal(err)
		}
	}
	
	scheduleReq := &fairrentv1.ScheduleNextRequest{
		Horizon: &commonv1.SchedulingHorizon{
			LookAhead: &durationpb.Duration{Seconds: 3600},
		},
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, err := fr.ScheduleNext(ctx, scheduleReq)
		if err != nil {
			b.Fatal(err)
		}
	}
}
