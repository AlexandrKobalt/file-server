package sqlite

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Config struct {
	Path string
}

type Sqlite struct {
	DB *sqlx.DB
}

func New(cfg Config) (sqlite *Sqlite, err error) {
	db, err := sqlx.Open("sqlite3", cfg.Path)
	if err != nil {
		return nil, err
	}

	return &Sqlite{
		DB: db,
	}, nil
}

func (s *Sqlite) Start(_ context.Context) error {
	return s.DB.Ping()
}

func (s *Sqlite) Stop(_ context.Context) error {
	return s.DB.Close()
}
