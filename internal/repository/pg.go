package repository

import (
	"errors"
	"fmt"
	"github.com/ShpullRequest/backend/internal/config"
	"github.com/ShpullRequest/backend/pkg/storage/postgres"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

type Pg struct {
	db     postgres.PgSQL
	logger *zap.Logger
}

func NewPG(cfg config.NodeConfig, logger *zap.Logger) (*Pg, error) {
	db, err := postgres.New(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("postgres.New: %w", err)
	}

	return &Pg{
		db:     db,
		logger: logger,
	}, nil
}

type errorFunc func(code string) bool

func (p *Pg) IsError(f errorFunc, err error) bool {
	var pgError *pgconn.PgError
	if !errors.As(err, &pgError) {
		return false
	}

	p.logger.Info("pgerror", zap.Any("pgerror", pgError))

	return f(pgError.Code)
}
