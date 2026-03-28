package adapter

import (
	"errors"
	"time"

	"github.com/RizqiPangestu/url_shortener/internal/adapter/internal/postgres"
	"github.com/RizqiPangestu/url_shortener/internal/core"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type urlPostgresAdapter struct {
	db *pg.DB
}

func NewURLPostgresAdapter(db *pg.DB) core.URLPort {
	if err := db.Model(&postgres.URL{}).CreateTable(&orm.CreateTableOptions{
		IfNotExists: true,
	}); err != nil {
		panic(err)
	}

	return &urlPostgresAdapter{
		db: db,
	}
}

func (a *urlPostgresAdapter) SavePath(shortPath string, originURL string, ttl time.Duration) error {
	now := time.Now()
	_, err := a.db.Model(&postgres.URL{
		OriginalURL:    originURL,
		ShortPath:      shortPath,
		TTL:            ttl,
		CreatedAt:      now,
		UpdatedAt:      now,
		LastAccessedAt: now,
	}).Insert()
	if err != nil {
		if pgErr, ok := err.(pg.Error); ok && pgErr.IntegrityViolation() {
			// 23505 is the code for unique violation
			// ref: postgresql.org/docs/10/static/errcodes-appendix.html
			if pgErr.Field('C') == "23505" {
				return core.ErrURLAlreadyExists
			}

		}
		return err
	}

	return nil
}

func (a *urlPostgresAdapter) GetByOriginalURL(originalURL string) (core.URL, error) {
	var url postgres.URL
	// TODO: create an index on original_url
	err := a.db.Model(&url).Where("original_url = ?", originalURL).Select()
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return core.URL{}, core.ErrURLNotFound
		}
		return core.URL{}, err
	}
	return url.Entity(), nil
}

func (a *urlPostgresAdapter) GetByShortPath(shortPath string) (core.URL, error) {
	var url postgres.URL
	// TODO: handle SQL injection
	err := a.db.Model(&url).Where("short_path = ?", shortPath).Select()
	if err != nil {
		return core.URL{}, err
	}

	return url.Entity(), nil
}

func (a *urlPostgresAdapter) UpdateLastAccessedAt(shortPath string) error {
	now := time.Now()
	_, err := a.db.Model(&postgres.URL{
		ShortPath:      shortPath,
		UpdatedAt:      now,
		LastAccessedAt: now,
	}).WherePK().UpdateNotZero()
	if err != nil {
		return err
	}
	return nil
}

func (a *urlPostgresAdapter) DeleteByShortPath(shortPath string) error {
	_, err := a.db.Model(&postgres.URL{
		ShortPath: shortPath,
	}).WherePK().Delete()
	if err != nil {
		return err
	}
	return nil
}
