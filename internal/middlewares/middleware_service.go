package middlewares

import (
	"github.com/ShpullRequest/backend/internal/api"
	"go.uber.org/zap"
)

type middlewareService struct {
	logger *zap.Logger
}

func ConfigureService(apiService api.Service) {
	ms := &middlewareService{
		logger: apiService.GetLogger(),
	}

	apiService.GetRouter().Use(ms.Logger)
	apiService.GetRouter().Use(ms.Compress)
}
