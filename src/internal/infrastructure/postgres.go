package infrastructure

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	Conn *sqlx.DB
}

func NewPostgresDB(dataSourceName string) (*PostgresDB, error) {
	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	return &PostgresDB{Conn: db}, nil
}

func (h *PostgresDB) Db() *sqlx.DB {
	return h.Conn
}

func (h *PostgresDB) Close() error {
	return h.Conn.Close()
}
