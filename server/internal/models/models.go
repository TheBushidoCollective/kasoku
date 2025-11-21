package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a kasoku user/account
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Name      string    `json:"name"`
	Password  string    `json:"-"`                          // bcrypt hash (optional for OAuth users)
	Plan      string    `gorm:"default:'free'" json:"plan"` // free, individual, team
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// OAuth fields
	Provider   string  `gorm:"index" json:"provider,omitempty"`    // github, gitlab, email
	ProviderID string  `gorm:"index" json:"provider_id,omitempty"` // OAuth provider user ID
	AvatarURL  *string `json:"avatar_url,omitempty"`               // Profile picture URL

	// Stripe billing
	StripeCustomerID     *string `gorm:"index" json:"-"`
	StripeSubscriptionID *string `gorm:"index" json:"-"`

	// For team plan
	TeamID *uuid.UUID `gorm:"type:uuid;index" json:"team_id,omitempty"`
	Team   *Team      `gorm:"foreignKey:TeamID" json:"team,omitempty"`

	// API tokens
	Tokens []Token `gorm:"foreignKey:UserID" json:"-"`
}

// Team represents a team subscription
type Team struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	OwnerID   uuid.UUID `gorm:"type:uuid;not null" json:"owner_id"`
	Plan      string    `gorm:"default:'team'" json:"plan"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Members
	Members []User `gorm:"foreignKey:TeamID" json:"members,omitempty"`
}

// Project represents a git repository project
type Project struct {
	ID           uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	TeamID       *uuid.UUID `gorm:"type:uuid;index" json:"team_id,omitempty"`
	Name         string     `json:"name"`                             // Display name (optional)
	GitRemote    string     `gorm:"index;not null" json:"git_remote"` // e.g., "github.com/thebushidocollective/brisk"
	RepoPath     string     `gorm:"default:'/'" json:"repo_path"`     // Path in repo where kasoku.yaml is (e.g., "/", "/web", "/server")
	LastAccessed time.Time  `json:"last_accessed"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// Token represents an API token
type Token struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	Name      string     `gorm:"not null" json:"name"`
	Token     string     `gorm:"uniqueIndex;not null" json:"-"` // SHA256 hash
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	LastUsed  *time.Time `json:"last_used,omitempty"`
}

// CacheEntry represents a cached artifact
type CacheEntry struct {
	ID         uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	TeamID     *uuid.UUID `gorm:"type:uuid;index" json:"team_id,omitempty"`
	ProjectID  *uuid.UUID `gorm:"type:uuid;index" json:"project_id,omitempty"` // Associated project
	Hash       string     `gorm:"index:idx_user_hash,unique;not null" json:"hash"`
	Command    string     `json:"command"`
	Size       int64      `json:"size"`              // bytes
	StorageKey string     `gorm:"not null" json:"-"` // S3/GCS key
	CreatedAt  time.Time  `gorm:"index" json:"created_at"`
	AccessedAt time.Time  `gorm:"index" json:"accessed_at"`
}

// AnalyticsEvent represents a cache operation for analytics
type AnalyticsEvent struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	TeamID    *uuid.UUID `gorm:"type:uuid;index" json:"team_id,omitempty"`
	ProjectID *uuid.UUID `gorm:"type:uuid;index" json:"project_id,omitempty"` // Associated project
	EventType string     `gorm:"not null;index" json:"event_type"`            // hit, miss, put
	Hash      string     `gorm:"index" json:"hash"`
	Command   string     `json:"command"`
	Size      int64      `json:"size,omitempty"`
	Duration  int64      `json:"duration,omitempty"`   // milliseconds
	TimeSaved int64      `json:"time_saved,omitempty"` // milliseconds
	CreatedAt time.Time  `gorm:"index" json:"created_at"`
}

// StorageQuota tracks user/team storage usage
type StorageQuota struct {
	UserID     uuid.UUID  `gorm:"type:uuid;primary_key" json:"user_id"`
	TeamID     *uuid.UUID `gorm:"type:uuid;unique" json:"team_id,omitempty"`
	UsedBytes  int64      `gorm:"default:0" json:"used_bytes"`
	LimitBytes int64      `gorm:"not null" json:"limit_bytes"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// BeforeCreate hooks for UUID generation
func (u *User) BeforeCreate() error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (t *Team) BeforeCreate() error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

func (t *Token) BeforeCreate() error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

func (c *CacheEntry) BeforeCreate() error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	now := time.Now()
	c.AccessedAt = now
	return nil
}

func (a *AnalyticsEvent) BeforeCreate() error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

func (p *Project) BeforeCreate() error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	now := time.Now()
	p.LastAccessed = now
	return nil
}
