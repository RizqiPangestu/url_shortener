package core

import "time"

type URL struct {
	ShortPath      string
	OriginalURL    string
	TTL            time.Duration
	CreatedAt      time.Time
	UpdatedAt      time.Time
	LastAccessedAt time.Time
}
