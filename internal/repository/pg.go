package repository

import (
	"fmt"
	"github.com/ShpullRequest/backend/internal/config"
	"github.com/ShpullRequest/backend/pkg/storage/postgres"
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
