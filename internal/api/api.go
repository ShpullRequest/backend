package api

import (
	"github.com/ShpullRequest/backend/internal/config"
	"github.com/ShpullRequest/backend/internal/repository"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Service interface {
	GetRouter() *gin.Engine
	GetPg() *repository.Pg
	GetLogger() *zap.Logger
}

type API struct {
	router *gin.Engine
	pg     *repository.Pg
	logger *zap.Logger
}

func New(cfg config.NodeConfig, pg *repository.Pg, logger *zap.Logger) *API {
	if cfg.ProdFlag {
		gin.SetMode(gin.ReleaseMode)
	}

	return &API{
		router: gin.New(),
		pg:     pg,
		logger: logger,
	}
}

func (a *API) GetRouter() *gin.Engine {
	return a.router
}

func (a *API) GetPg() *repository.Pg {
	return a.pg
}

func (a *API) GetLogger() *zap.Logger {
	return a.logger
}
