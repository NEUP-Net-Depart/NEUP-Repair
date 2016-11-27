package model

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Storager interface {
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Unsafe() *sqlx.DB
	Beginx() (*sqlx.Tx, error)
}
