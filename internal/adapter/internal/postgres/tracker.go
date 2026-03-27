package postgres

import "time"

type Tracker struct {
	ShortPath     string    `pg:"short_path,pk,unique"`
	ClickCount    int       `pg:"click_count"`
	LastClickedAt time.Time `pg:"last_clicked_at"`
	Timestamp     time.Time `pg:"date"`
}
