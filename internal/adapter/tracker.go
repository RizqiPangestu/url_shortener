package adapter

import (
	"time"

	"github.com/RizqiPangestu/url_shortener/internal/adapter/internal/postgres"
	"github.com/RizqiPangestu/url_shortener/internal/core"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type trackerPostgresAdapter struct {
	db *pg.DB
}

func NewTrackerPostgresAdapter(db *pg.DB) core.TrackerPort {
	if err := db.Model(&postgres.Tracker{}).CreateTable(&orm.CreateTableOptions{
		IfNotExists: true,
	}); err != nil {
		panic(err)
	}

	return &trackerPostgresAdapter{db: db}
}

func (a *trackerPostgresAdapter) Track(shortPath string) error {
	now := time.Now()
	dateOnly, err := time.Parse(time.DateOnly, now.Format(time.DateOnly))
	if err != nil {
		return err
	}

	// do upsert
	if _, err := a.db.Model(&postgres.Tracker{
		ShortPath:     shortPath,
		Date:          dateOnly,
		ClickCount:    1,
		LastClickedAt: now,
	}).
		OnConflict(`("short_path", "date") DO UPDATE`).
		Set(`click_count = tracker.click_count + 1, 
			last_clicked_at = EXCLUDED.last_clicked_at`).
		Insert(); err != nil {
		return err
	}
	return nil
}
