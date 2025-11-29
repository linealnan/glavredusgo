package db

import "database/sql"

func NewDbConnection() *sql.DB {
	db, err := sql.Open("sqlite3", "glavredus.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	return db
}
