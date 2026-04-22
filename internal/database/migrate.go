package database

import (
	"database/sql"
	"fmt"
	"log/slog"
)

func Migrate(db *DB) error {
	driver := db.Driver()
	if err := createAuthorsTable(db, driver); err != nil {
		return fmt.Errorf("authors table: %w", err)
	}
	if err := createLoansTable(db, driver); err != nil {
		return fmt.Errorf("loans table: %w", err)
	}
	if err := migrateBooksTable(db, driver); err != nil {
		return fmt.Errorf("books table: %w", err)
	}
	return nil
}

func createAuthorsTable(db *DB, driver string) error {
	var query string
	if driver == "postgres" {
		query = `CREATE TABLE IF NOT EXISTS authors (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			bio TEXT,
			birth_year INTEGER,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`
	} else {
		query = `CREATE TABLE IF NOT EXISTS authors (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			bio TEXT,
			birth_year INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`
	}
	_, err := db.Exec(query)
	return err
}

func createLoansTable(db *DB, driver string) error {
	var query string
	if driver == "postgres" {
		query = `CREATE TABLE IF NOT EXISTS loans (
			id SERIAL PRIMARY KEY,
			book_id INTEGER NOT NULL,
			borrower_name TEXT NOT NULL,
			loan_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			due_date TIMESTAMP,
			return_date TIMESTAMP,
			status TEXT DEFAULT 'active'
		);`
	} else {
		query = `CREATE TABLE IF NOT EXISTS loans (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			book_id INTEGER NOT NULL,
			borrower_name TEXT NOT NULL,
			loan_date DATETIME DEFAULT CURRENT_TIMESTAMP,
			due_date DATETIME,
			return_date DATETIME,
			status TEXT DEFAULT 'active'
		);`
	}
	_, err := db.Exec(query)
	return err
}

func migrateBooksTable(db *DB, driver string) error {
	if driver == "postgres" {
		query := `CREATE TABLE IF NOT EXISTS books (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			author TEXT NOT NULL,
			author_id INTEGER NOT NULL DEFAULT 0,
			genre TEXT NOT NULL,
			pages INTEGER NOT NULL,
			publicationyear INTEGER NOT NULL
		);`
		_, err := db.Exec(query)
		return err
	}

	var count int
	err := db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='books'").Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		query := `CREATE TABLE books (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			author TEXT NOT NULL,
			author_id INTEGER NOT NULL DEFAULT 0,
			genre TEXT NOT NULL,
			pages INTEGER NOT NULL,
			publicationyear INTEGER NOT NULL
		);`
		_, err := db.Exec(query)
		return err
	}

	rows, err := db.Query("PRAGMA table_info(books)")
	if err != nil {
		return err
	}
	defer rows.Close()

	hasAuthorID := false
	for rows.Next() {
		var cid int
		var name string
		var typeName string
		var notNull int
		var dfltValue sql.NullString
		var pk int
		if err := rows.Scan(&cid, &name, &typeName, &notNull, &dfltValue, &pk); err != nil {
			return err
		}
		if name == "author_id" {
			hasAuthorID = true
			break
		}
	}

	if !hasAuthorID {
		slog.Info("migrating books table: adding author_id column")
		_, err = db.Exec("ALTER TABLE books ADD COLUMN author_id INTEGER NOT NULL DEFAULT 0")
		if err != nil {
			return fmt.Errorf("failed to add author_id column: %w", err)
		}
	}

	return nil
}
