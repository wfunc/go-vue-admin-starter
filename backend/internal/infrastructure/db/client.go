package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/internal/config"
	_ "github.com/lib/pq"
)

type Store struct {
	SQL *sql.DB
	Ent *ent.Client
}

func New(cfg config.DatabaseConfig) (*Store, error) {
	db, err := sql.Open(cfg.DriverName(), cfg.ConnectionString())
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	if cfg.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}

	ctx, cancel := context.WithTimeout(context.Background(), Timeout())
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("connect database %s:%d/%s: %w", cfg.Host, cfg.Port, cfg.Name, err)
	}

	driver := entsql.OpenDB(dialect.Postgres, db)
	client := ent.NewClient(ent.Driver(driver))
	return &Store{SQL: db, Ent: client}, nil
}

func (s *Store) Close() error {
	if s == nil || s.SQL == nil {
		return nil
	}
	s.Ent.Close()
	return s.SQL.Close()
}

func Ping(ctx context.Context, db *sql.DB) error {
	if db == nil {
		return nil
	}
	return db.PingContext(ctx)
}

func Timeout() time.Duration {
	return 5 * time.Second
}
