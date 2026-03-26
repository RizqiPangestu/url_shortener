package adapter

import (
	"time"

	"github.com/RizqiPangestu/url_shortener/internal/adapter/internal/postgres"
	"github.com/RizqiPangestu/url_shortener/internal/core"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type urlMongoAdapter struct {
	db *pg.DB
}

func NewURLMongoAdapter() core.URLPort {
	db := pg.Connect(&pg.Options{
		Addr:     "localhost:5432",
		Database: "localpostgre",
		User:     "rizqipangestu",
	})

	if err := db.Model(&postgres.URL{}).CreateTable(&orm.CreateTableOptions{
		IfNotExists: true,
	}); err != nil {
		panic(err)
	}

	return &urlMongoAdapter{
		db: db,
	}
}

func (a *urlMongoAdapter) SavePath(shortPath string, originURL string) error {
	now := time.Now()
	_, err := a.db.Model(&postgres.URL{
		OriginalURL: originURL,
		ShortPath:   shortPath,
		TTL:         0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}).Insert()
	if err != nil {
		return err
	}

	return nil
}

func (a *urlMongoAdapter) GetOriginURL(shortPath string) (string, error) {
	var url postgres.URL
	err := a.db.Model(&url).Where("short_path = ?", shortPath).Select()
	if err != nil {
		return "", err
	}
	return url.OriginalURL, nil
}
