package main

import (
	"database/sql"
)

const helloQuery = "SELECT 'Hello, #' || id FROM names WHERE name = ?"
const addNameQuery = "INSERT INTO names (name) VALUES (?)"
const createTableQuery = `
	CREATE TABLE IF NOT EXISTS names (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR NOT NULL
	)
`


func initDB(db *sql.DB) error {
	_, err := db.Exec(createTableQuery)
	return err
}


func helloSql_NotCoolAtAll(db *sql.DB, name string) (string, error) {
	// too much bad code to write here
	return "", nil
}


func helloSql_Cool(db *sql.DB, name string) (string, error) {
	var result string
	err := RunTransaction(db).Use(func(tx *sql.Tx) error {

		_, err := tx.Exec(addNameQuery, name)
		if err != nil {
			return err
		}

		err = QueryRows(tx, helloQuery, name).Use(func(rows *sql.Rows) error {
			_ = rows.Next()
			return rows.Scan(&result)
		})

		return err
	})
	return result, err
}


////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////



type DBResource struct {
	Use func (func (db *sql.DB) error) error
}

func NewDBResource(driverName, datasourceName string) DBResource {
	return DBResource{
		Use: func(callback func(db *sql.DB) error) error {
			db, err := sql.Open(driverName, datasourceName)
			if err != nil {
				return err
			}
			err = callback(db)
			if err != nil {
				_ = db.Close()
				return err
			} else {
				return db.Close()
			}
		},
	}
}

type TxResource struct {
	Use func (func (tx *sql.Tx) error) error
}

func RunTransaction(db *sql.DB) TxResource {
	return TxResource{
		Use: func(callback func(tx *sql.Tx) error) error {
			tx, err := db.Begin()
			if err != nil {
				return err
			}
			err = callback(tx)
			if err != nil {
				_ = tx.Rollback()
				return err
			} else {
				return tx.Commit()
			}
		},
	}
}

type RowsResource struct {
	Use func (func (rows *sql.Rows) error) error
}

func QueryRows(tx *sql.Tx, query string, args ...interface{}) RowsResource {
	return RowsResource{
		Use: func(callback func(rows *sql.Rows) error) error {
			rows, err := tx.Query(query, args...)
			if err != nil {
				return err
			}
			err = callback(rows)
			if err != nil {
				_ = rows.Close()
				return err
			} else {
				return rows.Close()
			}
		},
	}
}

