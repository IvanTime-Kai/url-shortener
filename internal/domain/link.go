package domain

import "time"

type Link struct {
	ID          int64
	OriginalURL string
	Code        string
	CreatedAt   time.Time
	ExpiresAt   *time.Time `json:"expires_at,omitempty"` // nil = không expire
}

func (l *Link) IsExpired() bool {
	if l.ExpiresAt == nil {
		return false
	}

	return time.Now().After(*l.ExpiresAt)
}