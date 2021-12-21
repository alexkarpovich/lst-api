package db

import "github.com/jmoiron/sqlx"

type DB interface {
	Db() *sqlx.DB
	//Exec(statement string)
	//Query(statement string) Row
}

// type Row interface {
// 	Scan(dest ...interface{})
// 	Next() bool
// }
