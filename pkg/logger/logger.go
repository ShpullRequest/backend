package logger

import (
	"github.com/ShpullRequest/backend/internal/config"
	"go.uber.org/zap"
)

func New(cfg config.NodeConfig) (*zap.Logger, error) {
	if cfg.ProdFlag {
		prodLogger, err := zap.NewProduction()
		if err != nil {
			return nil, err
		}

		level := prodLogger.Level()
		if err = level.Set("error"); err != nil {
			return nil, err
		}

		return prodLogger, err
	}

	return zap.NewDevelopment()
}
