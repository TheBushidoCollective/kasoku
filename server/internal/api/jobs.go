package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/thebushidocollective/kasoku/server/internal/jobs"
)

type JobsHandler struct {
	coordinator *jobs.Coordinator
}

func NewJobsHandler(coordinator *jobs.Coordinator) *JobsHandler {
	return &JobsHandler{
		coordinator: coordinator,
	}
}

// RegisterJob handles POST /jobs/register
// Registers a new job and returns whether this job should execute or wait
func (h *JobsHandler) RegisterJob(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var teamID *uuid.UUID
	if tid, ok := r.Context().Value("team_id").(uuid.UUID); ok {
		teamID = &tid
	}

	hash := chi.URLParam(r, "hash")
	if hash == "" {
		respondError(w, http.StatusBadRequest, "hash is required")
		return
	}

	isPrimary, jobID := h.coordinator.Register(r.Context(), hash, userID, teamID)

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"job_id":     jobID,
		"is_primary": isPrimary,
		"hash":       hash,
		"message": func() string {
			if isPrimary {
				return "you are the primary job, proceed with execution"
			}
			return "another job is already running, wait for completion"
		}(),
	})
}

// CompleteJob handles POST /jobs/complete/:hash
func (h *JobsHandler) CompleteJob(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	hash := chi.URLParam(r, "hash")
	if hash == "" {
		respondError(w, http.StatusBadRequest, "hash is required")
		return
	}

	if err := h.coordinator.Complete(r.Context(), hash); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "job marked as complete",
		"hash":    hash,
	})
}

// GetJobStatus handles GET /jobs/status/:hash
func (h *JobsHandler) GetJobStatus(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	hash := chi.URLParam(r, "hash")
	if hash == "" {
		respondError(w, http.StatusBadRequest, "hash is required")
		return
	}

	job, exists := h.coordinator.GetStatus(r.Context(), hash)
	if !exists {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"status": "not_found",
			"hash":   hash,
		})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"job_id":     job.ID,
		"hash":       hash,
		"status":     job.Status,
		"started_at": job.StartedAt,
		"ended_at":   job.EndedAt,
		"error":      job.Error,
	})
}

// WaitForJob handles GET /jobs/wait/:hash
// Long-polls waiting for a job to complete
func (h *JobsHandler) WaitForJob(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	hash := chi.URLParam(r, "hash")
	if hash == "" {
		respondError(w, http.StatusBadRequest, "hash is required")
		return
	}

	// Wait up to 5 minutes for the job to complete
	status, err := h.coordinator.WaitForCompletion(r.Context(), hash, 5*time.Minute)
	if err != nil {
		respondError(w, http.StatusRequestTimeout, "timeout waiting for job completion")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status": status,
		"hash":   hash,
	})
}
