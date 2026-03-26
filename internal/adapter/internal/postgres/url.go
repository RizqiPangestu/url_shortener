package postgres

import (
	"fmt"
	"time"

	"github.com/RizqiPangestu/url_shortener/internal/core"
)

type URL struct {
	ShortPath      string        `pg:"short_path,pk,unique"`
	OriginalURL    string        `pg:"original_url"`
	TTL            time.Duration `pg:"ttl"`
	CreatedAt      time.Time     `pg:"created_at"`
	UpdatedAt      time.Time     `pg:"updated_at"`
	LastAccessedAt time.Time     `pg:"last_accessed_at"`
}

func (u *URL) Entity() core.URL {
	fmt.Printf("%+v\n", u)
	return core.URL{
		ShortPath:      u.ShortPath,
		OriginalURL:    u.OriginalURL,
		TTL:            u.TTL,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
		LastAccessedAt: u.LastAccessedAt,
	}
}
