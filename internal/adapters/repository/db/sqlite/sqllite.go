package sqlite

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

const (
	comicTable        = "Comic"
	keywordTable      = "Keyword"
	keywordComicTable = "comic_keyword"
	keywordComicView  = "kw_comic"
)

func NewSqliteDB() (*sqlx.DB, error) {
	const op = "op.new_sqlite_db"

	db, err := sqlx.Open("sqlite3", viper.GetString("database.dsn"))
	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %s", op, err)
	}

	return db, nil
}
