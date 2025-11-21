package jobs

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

// JobStatus represents the status of a running job
type JobStatus string

const (
	JobStatusRunning  JobStatus = "running"
	JobStatusComplete JobStatus = "complete"
	JobStatusFailed   JobStatus = "failed"
)

// Job represents a cache build job
type Job struct {
	ID        uuid.UUID
	Hash      string
	UserID    uuid.UUID
	TeamID    *uuid.UUID
	Status    JobStatus
	StartedAt time.Time
	EndedAt   *time.Time
	Error     string
}

// Coordinator manages distributed job coordination
type Coordinator struct {
	jobs      map[string]*Job // key is hash
	mu        sync.RWMutex
	completed chan string // channel to notify job completion
}

// NewCoordinator creates a new job coordinator
func NewCoordinator() *Coordinator {
	return &Coordinator{
		jobs:      make(map[string]*Job),
		completed: make(chan string, 100),
	}
}

// Register registers a new job and returns whether this is the primary job
// Returns (isPrimary, jobID)
func (c *Coordinator) Register(ctx context.Context, hash string, userID uuid.UUID, teamID *uuid.UUID) (bool, uuid.UUID) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if job is already running
	if existing, ok := c.jobs[hash]; ok {
		if existing.Status == JobStatusRunning {
			// Another job is already running, this is a secondary job
			return false, existing.ID
		}
	}

	// This is the primary job
	job := &Job{
		ID:        uuid.New(),
		Hash:      hash,
		UserID:    userID,
		TeamID:    teamID,
		Status:    JobStatusRunning,
		StartedAt: time.Now(),
	}

	c.jobs[hash] = job
	return true, job.ID
}

// Complete marks a job as complete
func (c *Coordinator) Complete(ctx context.Context, hash string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	job, ok := c.jobs[hash]
	if !ok {
		return nil // Job doesn't exist, that's okay
	}

	now := time.Now()
	job.Status = JobStatusComplete
	job.EndedAt = &now

	// Notify waiting jobs
	select {
	case c.completed <- hash:
	default:
	}

	// Clean up after a short delay to allow waiting jobs to poll
	go func() {
		time.Sleep(5 * time.Second)
		c.mu.Lock()
		delete(c.jobs, hash)
		c.mu.Unlock()
	}()

	return nil
}

// Fail marks a job as failed
func (c *Coordinator) Fail(ctx context.Context, hash string, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	job, ok := c.jobs[hash]
	if !ok {
		return
	}

	now := time.Now()
	job.Status = JobStatusFailed
	job.EndedAt = &now
	if err != nil {
		job.Error = err.Error()
	}

	// Clean up failed jobs quickly
	go func() {
		time.Sleep(1 * time.Second)
		c.mu.Lock()
		delete(c.jobs, hash)
		c.mu.Unlock()
	}()
}

// GetStatus returns the current status of a job by hash
func (c *Coordinator) GetStatus(ctx context.Context, hash string) (*Job, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	job, ok := c.jobs[hash]
	if !ok {
		return nil, false
	}

	// Return a copy to prevent race conditions
	jobCopy := *job
	return &jobCopy, true
}

// WaitForCompletion waits for a job to complete with a timeout
// Returns true if the job completed successfully, false if timeout or failed
func (c *Coordinator) WaitForCompletion(ctx context.Context, hash string, timeout time.Duration) (JobStatus, error) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	deadline := time.Now().Add(timeout)

	for {
		select {
		case <-ctx.Done():
			return JobStatusFailed, ctx.Err()
		case <-ticker.C:
			if time.Now().After(deadline) {
				return JobStatusFailed, context.DeadlineExceeded
			}

			job, ok := c.GetStatus(ctx, hash)
			if !ok {
				// Job disappeared, might have completed
				return JobStatusComplete, nil
			}

			if job.Status == JobStatusComplete {
				return JobStatusComplete, nil
			}

			if job.Status == JobStatusFailed {
				return JobStatusFailed, nil
			}
		}
	}
}
