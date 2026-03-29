package domain

import "time"

type Click struct {
	ID        int64
	LinkID    int64
	IP        string
	UserAgent string
	CreatedAt time.Time
}

type LinkStats struct {
	Link        Link
	TotalClicks int64
	Clicks       []Click
}

