package database

import (
	"database/sql"
	"log"
)

var db *sql.DB

func GetDB() *sql.DB {
	return db
}
func SetupDB() error {
	var err error
	db, err = sql.Open("sqlite3", "./queue.db")
	if err != nil {
		return err
	}

	// Enable WAL mode
	var journalMode string
	err = db.QueryRow("PRAGMA journal_mode=WAL;").Scan(&journalMode)
	if err != nil {
		log.Fatal(err)
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS queue (
		queue_name TEXT NOT NULL,
		id TEXT PRIMARY KEY,
		message TEXT NOT NULL,
		available_at TEXT NOT NULL,
		is_available BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return err
	}

	log.Println("Table 'queue' created successfully or already exists.")
	return nil
}
