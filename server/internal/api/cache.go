package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/thebushidocollective/kasoku/server/internal/cache"
)

type CacheHandler struct {
	cache *cache.Service
}

func NewCacheHandler(cacheService *cache.Service) *CacheHandler {
	return &CacheHandler{
		cache: cacheService,
	}
}

// PutCache handles PUT /cache/:hash
func (h *CacheHandler) PutCache(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var teamID *uuid.UUID
	if tid, ok := r.Context().Value("team_id").(uuid.UUID); ok {
		teamID = &tid
	}

	// Get hash from URL
	hash := chi.URLParam(r, "hash")
	if hash == "" {
		respondError(w, http.StatusBadRequest, "hash is required")
		return
	}

	// Get command from query params
	command := r.URL.Query().Get("command")

	// Get content length
	size := r.ContentLength
	if size <= 0 {
		respondError(w, http.StatusBadRequest, "content-length is required")
		return
	}

	// Store artifact
	if err := h.cache.Put(r.Context(), userID, teamID, hash, command, r.Body, size); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "cache entry created",
		"hash":    hash,
	})
}

// GetCache handles GET /cache/:hash
func (h *CacheHandler) GetCache(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var teamID *uuid.UUID
	if tid, ok := r.Context().Value("team_id").(uuid.UUID); ok {
		teamID = &tid
	}

	// Get hash from URL
	hash := chi.URLParam(r, "hash")
	if hash == "" {
		respondError(w, http.StatusBadRequest, "hash is required")
		return
	}

	// Stream artifact directly to response
	entry, err := h.cache.Get(r.Context(), userID, teamID, hash, w)
	if err != nil {
		respondError(w, http.StatusNotFound, "cache miss")
		return
	}

	// Set headers
	w.Header().Set("Content-Type", "application/gzip")
	w.Header().Set("X-Cache-Command", entry.Command)
	w.Header().Set("X-Cache-Created-At", entry.CreatedAt.Format(time.RFC3339))
	w.Header().Set("Content-Length", strconv.FormatInt(entry.Size, 10))
}

// DeleteCache handles DELETE /cache/:hash
func (h *CacheHandler) DeleteCache(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Get hash from URL
	hash := chi.URLParam(r, "hash")
	if hash == "" {
		respondError(w, http.StatusBadRequest, "hash is required")
		return
	}

	// Delete cache entry
	if err := h.cache.Delete(r.Context(), userID, hash); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "cache entry deleted",
		"hash":    hash,
	})
}

// ListCache handles GET /cache
func (h *CacheHandler) ListCache(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var teamID *uuid.UUID
	if tid, ok := r.Context().Value("team_id").(uuid.UUID); ok {
		teamID = &tid
	}

	// Parse pagination params
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	// Get cache entries
	entries, err := h.cache.List(r.Context(), userID, teamID, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"entries": entries,
		"limit":   limit,
		"offset":  offset,
	})
}

// GetAnalytics handles GET /analytics
func (h *CacheHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user from context
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse since parameter (default to last 30 days)
	since := time.Now().AddDate(0, 0, -30)
	if sinceParam := r.URL.Query().Get("since"); sinceParam != "" {
		if t, err := time.Parse(time.RFC3339, sinceParam); err == nil {
			since = t
		}
	}

	// Get analytics summary
	summary, err := h.cache.GetAnalytics(r.Context(), userID, since)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, summary)
}

// Helper functions

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{
		"error": message,
	})
}

// respondStream writes data directly to response without JSON encoding
func respondStream(w http.ResponseWriter, status int, contentType string, reader io.Reader) error {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)
	_, err := io.Copy(w, reader)
	return err
}
