package logger

import (
	"github.com/ShpullRequest/backend/internal/config"
	"go.uber.org/zap"
)

func New() (*zap.Logger, error) {
	if config.Config.ProdFlag {
		return zap.NewProduction()
	}

	return zap.NewDevelopment()
}
