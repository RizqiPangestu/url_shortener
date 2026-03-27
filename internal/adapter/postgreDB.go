package adapter

import "github.com/go-pg/pg/v10"

func NewPostgresDB() *pg.DB {
	db := pg.Connect(&pg.Options{
		Addr:     "localhost:5432",
		Database: "localpostgre",
		User:     "rizqipangestu",
	})
	return db
}
