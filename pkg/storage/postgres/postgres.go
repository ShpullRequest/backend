package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ShpullRequest/backend/internal/config"
	"github.com/ShpullRequest/backend/pkg/migrations"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type PgSQL interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	ExecContextWithReturnID(ctx context.Context, query string, args ...interface{}) (uuid.UUID, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	PrepareContextMaster(ctx context.Context, query string) (*sqlx.Stmt, error)
	PrepareContextReplica(ctx context.Context, query string) (*sqlx.Stmt, error)
	BeginMaster(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
	BeginReplica(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
	Close() error
	GetMaster() *sqlx.DB
	GetReplica() *sqlx.DB
}

type Storage struct {
	master  *sqlx.DB
	replica *sqlx.DB
}

func New(cfg config.NodeConfig, logger *zap.Logger) (*Storage, error) {
	masterDB, err := sqlx.Open("pgx/v5", cfg.MasterDSN)
	if err != nil {
		return nil, err
	}
	masterDB.SetMaxOpenConns(cfg.MasterMaxOpen)

	replicaDB, err := sqlx.Open("pgx/v5", cfg.ReplicaDSN)
	if err != nil {
		return nil, err
	}
	replicaDB.SetMaxOpenConns(cfg.ReplicaMaxOpen)

	s := &Storage{
		master:  masterDB,
		replica: replicaDB,
	}

	if cfg.MigrationsFlag {
		if err = migrations.ApplyMigrations(masterDB, logger); err != nil {
			return nil, fmt.Errorf("migrations.ApplyMigrations: %w", err)
		}
	}

	return s, nil
}

func (s *Storage) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return s.replica.GetContext(ctx, dest, query, args...)
}

func (s *Storage) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return s.replica.SelectContext(ctx, dest, query, args...)
}

func (s *Storage) NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error) {
	return s.replica.NamedQueryContext(ctx, query, arg)
}

func (s *Storage) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.master.ExecContext(ctx, query, args...)
}

func (s *Storage) ExecContextWithReturnID(ctx context.Context, query string, args ...interface{}) (uuid.UUID, error) {
	query = fmt.Sprintf("%s RETURNING id", query)

	var id uuid.UUID
	row := s.master.QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		return uuid.UUID{}, row.Err()
	}

	err := row.Scan(&id)
	return id, err
}

func (s *Storage) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	return s.master.NamedExecContext(ctx, query, arg)
}

func (s *Storage) PrepareContextMaster(ctx context.Context, query string) (*sqlx.Stmt, error) {
	return s.master.PreparexContext(ctx, query)
}

func (s *Storage) PrepareContextReplica(ctx context.Context, query string) (*sqlx.Stmt, error) {
	return s.replica.PreparexContext(ctx, query)
}

func (s *Storage) BeginMaster(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	return s.master.BeginTxx(ctx, opts)
}

func (s *Storage) BeginReplica(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	return s.replica.BeginTxx(ctx, opts)
}

func (s *Storage) Close() error {
	if err := s.master.Close(); err != nil {
		return err
	}

	return s.replica.Close()
}

func (s *Storage) GetMaster() *sqlx.DB {
	return s.master
}

func (s *Storage) GetReplica() *sqlx.DB {
	return s.replica
}
