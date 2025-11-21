package db

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/thebushidocollective/kasoku/server/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	*gorm.DB
}

// Connect establishes a database connection
func Connect(driver, dsn string) (*DB, error) {
	var dialector gorm.Dialector

	switch driver {
	case "postgres":
		dialector = postgres.Open(dsn)
	case "sqlite":
		dialector = sqlite.Open(dsn)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate models
	if err := db.AutoMigrate(
		&models.User{},
		&models.Team{},
		&models.Project{},
		&models.Token{},
		&models.CacheEntry{},
		&models.AnalyticsEvent{},
		&models.StorageQuota{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &DB{db}, nil
}

// Cache operations

func (db *DB) GetCacheEntry(userID uuid.UUID, hash string) (*models.CacheEntry, error) {
	var entry models.CacheEntry
	result := db.Where("user_id = ? AND hash = ?", userID, hash).First(&entry)
	if result.Error != nil {
		return nil, result.Error
	}

	// Update accessed_at timestamp
	db.Model(&entry).Update("accessed_at", time.Now())

	return &entry, nil
}

func (db *DB) GetCacheEntryByTeam(teamID uuid.UUID, hash string) (*models.CacheEntry, error) {
	var entry models.CacheEntry
	result := db.Where("team_id = ? AND hash = ?", teamID, hash).First(&entry)
	if result.Error != nil {
		return nil, result.Error
	}

	// Update accessed_at timestamp
	db.Model(&entry).Update("accessed_at", time.Now())

	return &entry, nil
}

func (db *DB) CreateCacheEntry(entry *models.CacheEntry) error {
	return db.Create(entry).Error
}

func (db *DB) DeleteCacheEntry(id uuid.UUID) error {
	return db.Delete(&models.CacheEntry{}, id).Error
}

func (db *DB) ListCacheEntries(userID uuid.UUID, limit, offset int) ([]models.CacheEntry, error) {
	var entries []models.CacheEntry
	result := db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&entries)

	return entries, result.Error
}

func (db *DB) ListCacheEntriesByTeam(teamID uuid.UUID, limit, offset int) ([]models.CacheEntry, error) {
	var entries []models.CacheEntry
	result := db.Where("team_id = ?", teamID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&entries)

	return entries, result.Error
}

// GetOldestCacheEntries returns the oldest cache entries for FIFO eviction
func (db *DB) GetOldestCacheEntries(userID uuid.UUID, limit int) ([]models.CacheEntry, error) {
	var entries []models.CacheEntry
	result := db.Where("user_id = ?", userID).
		Order("created_at ASC").
		Limit(limit).
		Find(&entries)

	return entries, result.Error
}

func (db *DB) GetOldestCacheEntriesByTeam(teamID uuid.UUID, limit int) ([]models.CacheEntry, error) {
	var entries []models.CacheEntry
	result := db.Where("team_id = ?", teamID).
		Order("created_at ASC").
		Limit(limit).
		Find(&entries)

	return entries, result.Error
}

// Storage quota operations

func (db *DB) GetStorageQuota(userID uuid.UUID) (*models.StorageQuota, error) {
	var quota models.StorageQuota
	result := db.Where("user_id = ?", userID).First(&quota)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &quota, nil
}

func (db *DB) CreateStorageQuota(quota *models.StorageQuota) error {
	return db.Create(quota).Error
}

func (db *DB) UpdateStorageQuota(quota *models.StorageQuota) error {
	return db.Save(quota).Error
}

func (db *DB) GetTotalUsedStorage(userID uuid.UUID) (int64, error) {
	var total int64
	result := db.Model(&models.CacheEntry{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(size), 0)").
		Scan(&total)

	return total, result.Error
}

func (db *DB) GetTotalUsedStorageByTeam(teamID uuid.UUID) (int64, error) {
	var total int64
	result := db.Model(&models.CacheEntry{}).
		Where("team_id = ?", teamID).
		Select("COALESCE(SUM(size), 0)").
		Scan(&total)

	return total, result.Error
}

// User operations

func (db *DB) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	result := db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (db *DB) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (db *DB) GetUserByStripeSubscriptionID(subscriptionID string) (*models.User, error) {
	var user models.User
	result := db.Where("stripe_subscription_id = ?", subscriptionID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (db *DB) CreateUser(user *models.User) error {
	return db.Create(user).Error
}

func (db *DB) UpdateUser(user *models.User) error {
	return db.Save(user).Error
}

// Token operations

func (db *DB) GetTokenByValue(tokenValue string) (*models.Token, error) {
	var token models.Token
	result := db.Where("token = ?", tokenValue).First(&token)
	if result.Error != nil {
		return nil, result.Error
	}

	// Update last_used timestamp
	now := time.Now()
	db.Model(&token).Update("last_used", &now)

	return &token, nil
}

func (db *DB) CreateToken(token *models.Token) error {
	return db.Create(token).Error
}

func (db *DB) DeleteToken(id uuid.UUID) error {
	return db.Delete(&models.Token{}, id).Error
}

func (db *DB) ListTokens(userID uuid.UUID) ([]models.Token, error) {
	var tokens []models.Token
	result := db.Where("user_id = ?", userID).Order("created_at DESC").Find(&tokens)
	return tokens, result.Error
}

// Analytics operations

func (db *DB) CreateAnalyticsEvent(event *models.AnalyticsEvent) error {
	return db.Create(event).Error
}

func (db *DB) GetAnalyticsSummary(userID uuid.UUID, since time.Time) (map[string]interface{}, error) {
	var results []struct {
		EventType string
		Count     int64
		TotalSize int64
	}

	err := db.Model(&models.AnalyticsEvent{}).
		Select("event_type, COUNT(*) as count, SUM(size) as total_size").
		Where("user_id = ? AND created_at >= ?", userID, since).
		Group("event_type").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	summary := make(map[string]interface{})
	for _, r := range results {
		summary[r.EventType] = map[string]int64{
			"count":      r.Count,
			"total_size": r.TotalSize,
		}
	}

	// Get total time saved
	var timeSaved int64
	db.Model(&models.AnalyticsEvent{}).
		Select("COALESCE(SUM(time_saved), 0)").
		Where("user_id = ? AND created_at >= ? AND event_type = ?", userID, since, "hit").
		Scan(&timeSaved)

	summary["time_saved_ms"] = timeSaved

	return summary, nil
}

// Team operations

func (db *DB) GetTeamByID(id uuid.UUID) (*models.Team, error) {
	var team models.Team
	result := db.Preload("Members").First(&team, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &team, nil
}

func (db *DB) CreateTeam(team *models.Team) error {
	return db.Create(team).Error
}

func (db *DB) UpdateTeam(team *models.Team) error {
	return db.Save(team).Error
}

// Project operations

func (db *DB) GetProjectByRemoteAndPath(userID uuid.UUID, gitRemote, repoPath string) (*models.Project, error) {
	var project models.Project
	result := db.Where("user_id = ? AND git_remote = ? AND repo_path = ?", userID, gitRemote, repoPath).First(&project)
	if result.Error != nil {
		return nil, result.Error
	}
	return &project, nil
}

func (db *DB) GetProjectByRemoteAndPathForTeam(teamID uuid.UUID, gitRemote, repoPath string) (*models.Project, error) {
	var project models.Project
	result := db.Where("team_id = ? AND git_remote = ? AND repo_path = ?", teamID, gitRemote, repoPath).First(&project)
	if result.Error != nil {
		return nil, result.Error
	}
	return &project, nil
}

func (db *DB) CreateProject(project *models.Project) error {
	return db.Create(project).Error
}

func (db *DB) UpdateProjectAccess(id uuid.UUID) error {
	return db.Model(&models.Project{}).Where("id = ?", id).Update("last_accessed", time.Now()).Error
}

func (db *DB) ListProjects(userID uuid.UUID) ([]models.Project, error) {
	var projects []models.Project
	result := db.Where("user_id = ?", userID).Order("last_accessed DESC").Find(&projects)
	return projects, result.Error
}

func (db *DB) ListProjectsByTeam(teamID uuid.UUID) ([]models.Project, error) {
	var projects []models.Project
	result := db.Where("team_id = ?", teamID).Order("last_accessed DESC").Find(&projects)
	return projects, result.Error
}
