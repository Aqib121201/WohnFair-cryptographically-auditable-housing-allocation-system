package queue

import (
	"container/heap"
	"time"
)

// Ticket represents a housing request in the queue
type Ticket struct {
	ID            string
	UserID        string
	UserGroup     string
	Urgency       int
	EnqueueTime   time.Time
	PriorityScore float64
	Constraints   interface{} // Will be the protobuf request
}

// PriorityQueue implements heap.Interface for managing tickets by priority
type PriorityQueue struct {
	tickets []*Ticket
}

// Len returns the number of tickets in the queue
func (pq PriorityQueue) Len() int { return len(pq.tickets) }

// Less determines the ordering of tickets (higher priority first)
func (pq PriorityQueue) Less(i, j int) bool {
	// Higher priority score comes first
	if pq.tickets[i].PriorityScore != pq.tickets[j].PriorityScore {
		return pq.tickets[i].PriorityScore > pq.tickets[j].PriorityScore
	}
	
	// If priority scores are equal, earlier enqueue time comes first
	return pq.tickets[i].EnqueueTime.Before(pq.tickets[j].EnqueueTime)
}

// Swap exchanges tickets at positions i and j
func (pq PriorityQueue) Swap(i, j int) {
	pq.tickets[i], pq.tickets[j] = pq.tickets[j], pq.tickets[i]
}

// Push adds a ticket to the queue
func (pq *PriorityQueue) Push(x interface{}) {
	ticket := x.(*Ticket)
	pq.tickets = append(pq.tickets, ticket)
}

// Pop removes and returns the highest priority ticket
func (pq *PriorityQueue) Pop() interface{} {
	old := pq.tickets
	n := len(old)
	ticket := old[n-1]
	old[n-1] = nil // avoid memory leak
	pq.tickets = old[0 : n-1]
	return ticket
}

// Peek returns the highest priority ticket without removing it
func (pq *PriorityQueue) Peek() *Ticket {
	if pq.Len() == 0 {
		return nil
	}
	return pq.tickets[0]
}

// GetByID returns a ticket by its ID
func (pq *PriorityQueue) GetByID(id string) *Ticket {
	for _, ticket := range pq.tickets {
		if ticket.ID == id {
			return ticket
		}
	}
	return nil
}

// RemoveByID removes a ticket by its ID
func (pq *PriorityQueue) RemoveByID(id string) bool {
	for i, ticket := range pq.tickets {
		if ticket.ID == id {
			// Remove the ticket and re-heapify
			heap.Remove(pq, i)
			return true
		}
	}
	return false
}

// UpdatePriority updates a ticket's priority and re-heapifies
func (pq *PriorityQueue) UpdatePriority(id string, newPriority float64) bool {
	for i, ticket := range pq.tickets {
		if ticket.ID == id {
			ticket.PriorityScore = newPriority
			heap.Fix(pq, i)
			return true
		}
	}
	return false
}

// GetQueueStats returns basic statistics about the queue
func (pq *PriorityQueue) GetQueueStats() QueueStats {
	if pq.Len() == 0 {
		return QueueStats{}
	}

	stats := QueueStats{
		TotalTickets: pq.Len(),
		OldestTicket: pq.tickets[0].EnqueueTime,
		NewestTicket: pq.tickets[0].EnqueueTime,
	}

	// Find oldest and newest tickets
	for _, ticket := range pq.tickets {
		if ticket.EnqueueTime.Before(stats.OldestTicket) {
			stats.OldestTicket = ticket.EnqueueTime
		}
		if ticket.EnqueueTime.After(stats.NewestTicket) {
			stats.NewestTicket = ticket.EnqueueTime
		}
	}

	// Calculate average priority
	totalPriority := 0.0
	for _, ticket := range pq.tickets {
		totalPriority += ticket.PriorityScore
	}
	stats.AveragePriority = totalPriority / float64(pq.Len())

	return stats
}

// QueueStats contains queue statistics
type QueueStats struct {
	TotalTickets    int
	OldestTicket    time.Time
	NewestTicket    time.Time
	AveragePriority float64
}

// IsEmpty returns true if the queue has no tickets
func (pq *PriorityQueue) IsEmpty() bool {
	return pq.Len() == 0
}

// Size returns the number of tickets in the queue
func (pq *PriorityQueue) Size() int {
	return pq.Len()
}

// Clear removes all tickets from the queue
func (pq *PriorityQueue) Clear() {
	pq.tickets = nil
}

// GetTickets returns a copy of all tickets (for debugging/monitoring)
func (pq *PriorityQueue) GetTickets() []*Ticket {
	tickets := make([]*Ticket, len(pq.tickets))
	copy(tickets, pq.tickets)
	return tickets
}
