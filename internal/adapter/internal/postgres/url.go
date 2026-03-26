package postgres

import "time"

type URL struct {
	ShortPath   string    `pg:"short_path"`
	OriginalURL string    `pg:"original_url"`
	TTL         int       `pg:"ttl"`
	CreatedAt   time.Time `pg:"created_at"`
	UpdatedAt   time.Time `pg:"updated_at"`
}
