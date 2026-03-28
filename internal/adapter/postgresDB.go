package adapter

import "github.com/go-pg/pg/v10"

func NewPostgresDB(address string, database string, user string, password string) *pg.DB {
	db := pg.Connect(&pg.Options{
		Addr:     address,
		Database: database,
		User:     user,
		Password: password,
	})
	return db
}
