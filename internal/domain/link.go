package domain

import "time"

type Link struct {
	ID          int64
	OriginalURL string
	Code        string
	CreatedAt   time.Time
	ExpiresAt   *time.Time // nil = không expire
}
