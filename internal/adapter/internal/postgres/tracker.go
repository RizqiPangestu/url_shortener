package postgres

import "time"

type Tracker struct {
	ShortPath     string    `pg:"short_path,unique:short_path_date"`
	Date          time.Time `pg:"date,unique:short_path_date"`
	ClickCount    int       `pg:"click_count"`
	LastClickedAt time.Time `pg:"last_clicked_at"`
}
