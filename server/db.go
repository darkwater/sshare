package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDatabase() error {
	var err error
	db, err = sql.Open("sqlite3", "./sshare.db")
	if err != nil {
		return err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		userid      INTEGER PRIMARY KEY ASC,
		url         STRING  UNIQUE,
		invited_by  INTEGER
	);
	CREATE TABLE IF NOT EXISTS public_keys (
		hash        STRING  PRIMARY KEY,
		userid      INTEGER,
		data        BLOB
	);
	`)
	if err != nil {
		return err
	}

	return nil
}
