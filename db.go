package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB(dataSourceName string) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}

	_, err = db.Exec(initSQL)
	if err != nil {
		log.Panic(err)
	}
}

const initSQL = `
CREATE TABLE IF NOT EXISTS assembly
(
  id     INTEGER PRIMARY KEY,
  name   TEXT NOT NULL,
  url    TEXT NOT NULL,
  color  TEXT,
  active INTEGER DEFAULT 1
);
CREATE UNIQUE INDEX IF NOT EXISTS assembly_name_uindex
  ON assembly (name);
CREATE UNIQUE INDEX IF NOT EXISTS assembly_url_uindex
  ON assembly (url);

CREATE TABLE IF NOT EXISTS entries (
  id  INTEGER PRIMARY KEY,
  url TEXT NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS entries_url_uindex
  ON entries (url);
`
