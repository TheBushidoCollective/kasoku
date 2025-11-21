package cache

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/thebushidocollective/kasoku/server/internal/db"
	"github.com/thebushidocollective/kasoku/server/internal/models"
	"github.com/thebushidocollective/kasoku/server/internal/storage"
)

type Service struct {
	db      *db.DB
	storage storage.Storage
}

func NewService(database *db.DB, storageBackend storage.Storage) *Service {
	return &Service{
		db:      database,
		storage: storageBackend,
	}
}

// Put stores a cache artifact
func (s *Service) Put(ctx context.Context, userID uuid.UUID, teamID *uuid.UUID, hash string, command string, reader io.Reader, size int64) error {
	// Check if entry already exists
	existing, err := s.db.GetCacheEntry(userID, hash)
	if err == nil && existing != nil {
		// Already cached
		return nil
	}

	// Check storage quota before storing
	if err := s.checkAndEnforceQuota(ctx, userID, teamID, size); err != nil {
		return fmt.Errorf("quota exceeded: %w", err)
	}

	// Generate storage key
	storageKey := generateStorageKey(userID, hash)

	// Store artifact
	if err := s.storage.Put(ctx, storageKey, reader, size); err != nil {
		return fmt.Errorf("failed to store artifact: %w", err)
	}

	// Create database entry
	entry := &models.CacheEntry{
		UserID:     userID,
		TeamID:     teamID,
		Hash:       hash,
		Command:    command,
		Size:       size,
		StorageKey: storageKey,
	}

	if err := s.db.CreateCacheEntry(entry); err != nil {
		// Clean up storage on database error
		s.storage.Delete(ctx, storageKey)
		return fmt.Errorf("failed to create cache entry: %w", err)
	}

	// Record analytics
	s.recordEvent(userID, teamID, "put", hash, command, size, 0, 0)

	return nil
}

// Get retrieves a cache artifact
func (s *Service) Get(ctx context.Context, userID uuid.UUID, teamID *uuid.UUID, hash string, writer io.Writer) (*models.CacheEntry, error) {
	var entry *models.CacheEntry
	var err error

	// Try user's cache first
	entry, err = s.db.GetCacheEntry(userID, hash)
	if err != nil && teamID != nil {
		// If not found and user is in a team, try team cache
		entry, err = s.db.GetCacheEntryByTeam(*teamID, hash)
	}

	if err != nil {
		// Record cache miss
		s.recordEvent(userID, teamID, "miss", hash, "", 0, 0, 0)
		return nil, fmt.Errorf("cache miss: %w", err)
	}

	// Retrieve artifact from storage
	if err := s.storage.Get(ctx, entry.StorageKey, writer); err != nil {
		return nil, fmt.Errorf("failed to retrieve artifact: %w", err)
	}

	// Record cache hit (time_saved will be calculated by client)
	s.recordEvent(userID, teamID, "hit", hash, entry.Command, entry.Size, 0, 0)

	return entry, nil
}

// Delete removes a cache artifact
func (s *Service) Delete(ctx context.Context, userID uuid.UUID, hash string) error {
	entry, err := s.db.GetCacheEntry(userID, hash)
	if err != nil {
		return fmt.Errorf("cache entry not found: %w", err)
	}

	// Delete from storage
	if err := s.storage.Delete(ctx, entry.StorageKey); err != nil {
		return fmt.Errorf("failed to delete artifact: %w", err)
	}

	// Delete from database
	if err := s.db.DeleteCacheEntry(entry.ID); err != nil {
		return fmt.Errorf("failed to delete cache entry: %w", err)
	}

	return nil
}

// List returns paginated cache entries for a user
func (s *Service) List(ctx context.Context, userID uuid.UUID, teamID *uuid.UUID, limit, offset int) ([]models.CacheEntry, error) {
	if teamID != nil {
		return s.db.ListCacheEntriesByTeam(*teamID, limit, offset)
	}
	return s.db.ListCacheEntries(userID, limit, offset)
}

// checkAndEnforceQuota checks if the user has enough quota and evicts old entries if needed (FIFO)
func (s *Service) checkAndEnforceQuota(ctx context.Context, userID uuid.UUID, teamID *uuid.UUID, newSize int64) error {
	// Get or create storage quota
	quota, err := s.db.GetStorageQuota(userID)
	if err != nil {
		return fmt.Errorf("failed to get storage quota: %w", err)
	}

	// If no quota exists, create default based on user's plan
	if quota == nil {
		user, err := s.db.GetUserByID(userID)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}

		limitBytes := s.getQuotaForPlan(user.Plan)
		quota = &models.StorageQuota{
			UserID:     userID,
			TeamID:     teamID,
			UsedBytes:  0,
			LimitBytes: limitBytes,
		}

		if err := s.db.CreateStorageQuota(quota); err != nil {
			return fmt.Errorf("failed to create storage quota: %w", err)
		}
	}

	// Calculate total used storage
	var totalUsed int64
	if teamID != nil {
		totalUsed, err = s.db.GetTotalUsedStorageByTeam(*teamID)
	} else {
		totalUsed, err = s.db.GetTotalUsedStorage(userID)
	}

	if err != nil {
		return fmt.Errorf("failed to calculate used storage: %w", err)
	}

	// FIFO eviction: Remove oldest entries until we have space
	neededSpace := (totalUsed + newSize) - quota.LimitBytes
	if neededSpace > 0 {
		if err := s.evictOldestEntries(ctx, userID, teamID, neededSpace); err != nil {
			return fmt.Errorf("failed to evict entries: %w", err)
		}
	}

	return nil
}

// evictOldestEntries removes the oldest cache entries until enough space is freed (FIFO)
func (s *Service) evictOldestEntries(ctx context.Context, userID uuid.UUID, teamID *uuid.UUID, neededSpace int64) error {
	var freedSpace int64

	for freedSpace < neededSpace {
		var entries []models.CacheEntry
		var err error

		if teamID != nil {
			entries, err = s.db.GetOldestCacheEntriesByTeam(*teamID, 10)
		} else {
			entries, err = s.db.GetOldestCacheEntries(userID, 10)
		}

		if err != nil {
			return err
		}

		if len(entries) == 0 {
			// No more entries to evict
			break
		}

		for _, entry := range entries {
			// Delete from storage
			s.storage.Delete(ctx, entry.StorageKey)

			// Delete from database
			if err := s.db.DeleteCacheEntry(entry.ID); err != nil {
				continue
			}

			freedSpace += entry.Size

			// Record eviction event
			s.recordEvent(userID, teamID, "evict", entry.Hash, entry.Command, entry.Size, 0, 0)

			if freedSpace >= neededSpace {
				break
			}
		}
	}

	return nil
}

// getQuotaForPlan returns the storage quota in bytes for a given plan
func (s *Service) getQuotaForPlan(plan string) int64 {
	switch plan {
	case "free":
		return 0 // No remote cache for free tier
	case "individual":
		return 2 * 1024 * 1024 * 1024 // 2GB
	case "team":
		return 50 * 1024 * 1024 * 1024 // 50GB default for teams
	default:
		return 0
	}
}

// recordEvent creates an analytics event
func (s *Service) recordEvent(userID uuid.UUID, teamID *uuid.UUID, eventType, hash, command string, size, duration, timeSaved int64) {
	event := &models.AnalyticsEvent{
		UserID:    userID,
		TeamID:    teamID,
		EventType: eventType,
		Hash:      hash,
		Command:   command,
		Size:      size,
		Duration:  duration,
		TimeSaved: timeSaved,
	}

	// Best effort - don't fail the operation if analytics fails
	s.db.CreateAnalyticsEvent(event)
}

// GetAnalytics returns analytics summary for a user
func (s *Service) GetAnalytics(ctx context.Context, userID uuid.UUID, since time.Time) (map[string]interface{}, error) {
	return s.db.GetAnalyticsSummary(userID, since)
}

// generateStorageKey creates a unique storage key for an artifact
func generateStorageKey(userID uuid.UUID, hash string) string {
	return fmt.Sprintf("%s/%s.tar.gz", userID.String(), hash)
}
