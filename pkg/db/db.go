package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

var db *sql.DB

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
id INTEGER PRIMARY KEY AUTOINCREMENT,
date CHAR(8) NOT NULL DEFAULT '',
title VARCHAR(255) NOT NULL,
comment TEXT,
repeat VARCHAR(128)
);

CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

func Init(dbfile string) error {
	d, err := sql.Open("sqlite", dbfile)
	if err != nil {
		return fmt.Errorf("failed to open DB: %w", err)
	}

	if err := d.Ping(); err != nil {
		d.Close()
		return fmt.Errorf("failed to ping DB: %w", err)
	}

	if _, err := d.Exec(schema); err != nil {
		d.Close()
		return fmt.Errorf("failed to create schema: %w", err)
	}

	db = d
	return nil
}

func GetDB() *sql.DB {
	return db
}
