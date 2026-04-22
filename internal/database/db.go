package database

import (
	"database/sql"
	"strconv"
	"strings"
)

type DB struct {
	*sql.DB
	driver string
}

func NewDB(db *sql.DB, driver string) *DB {
	return &DB{DB: db, driver: driver}
}

func (db *DB) Driver() string {
	return db.driver
}

func (db *DB) Rebind(query string) string {
	if db.driver != "postgres" {
		return query
	}

	n := 1
	var result strings.Builder

	for _, c := range query {
		if c == '?' {
			result.WriteString("$")
			result.WriteString(strconv.Itoa(n))
			n++
		} else {
			result.WriteRune(c)
		}
	}
	return result.String()
}

func (db *DB) Insert(query string, args ...interface{}) (int64, error) {
	if db.driver == "postgres" {
		returningQuery := query + " RETURNING id"
		var id int64
		err := db.QueryRow(db.Rebind(returningQuery), args...).Scan(&id)
		return id, err
	}

	result, err := db.Exec(db.Rebind(query), args...)

	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
