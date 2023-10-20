package main

import (
	"github.com/ShpullRequest/backend/internal/api"
	"github.com/ShpullRequest/backend/internal/config"
	"github.com/ShpullRequest/backend/internal/handlers"
	"github.com/ShpullRequest/backend/internal/middlewares"
	"github.com/ShpullRequest/backend/internal/repository"
	"github.com/ShpullRequest/backend/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	log, err := logger.New()
	if err != nil {
		panic(err)
	}

	config.Load()
	if err = config.Parse(); err != nil {
		log.Panic("Failed loading config server", zap.Error(err))
	}
	log.Debug("Success loaded and initialized config", zap.Any("Config", config.Config))

	pg, err := repository.NewPG(config.Config, log)
	if err != nil {
		log.Panic("Error initialization connect to postgres", zap.Error(err))
	}
	log.Debug("Success connection to database")

	apiService := api.New(config.Config, pg, log)
	middlewares.ConfigureService(apiService)
	handlers.ConfigureService(apiService)
	log.Debug("Services: API, middleware, handlers have been successfully configured and sent to launch")

	if err = apiService.GetRouter().Run(config.Config.Address); err != nil {
		log.Panic("Error run api server", zap.Error(err))
	}
}
