package database

import (
	"database/sql"

	log "github.com/SoujiThenria/endsinger/logging"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db *sql.DB
)

func New(path string) (err error) {
	if path == "" {
		log.Warn("The database will be in memory. Nothing will be saved as soon as the program is stopped.")
	}
	db, err = sql.Open("sqlite3", path)
	if err != nil {
		return
	}

	tables := [...]string{
		tableCreateGuild,
		tableCreateChannel,
	}

	for _, t := range tables {
		statement, err := db.Prepare(t)
		if err != nil {
			return err
		}
		if _, err = statement.Exec(); err != nil {
			return err
		}
	}

	return
}

func Close() (err error) {
	err = db.Close()
	return
}
