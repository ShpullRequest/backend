package handlers

import (
	"github.com/ShpullRequest/backend/internal/api"
	"github.com/ShpullRequest/backend/internal/repository"
	"go.uber.org/zap"
)

type handlerService struct {
	pg     *repository.Pg
	logger *zap.Logger
}

func ConfigureService(apiService api.Service) {
	hs := &handlerService{
		pg:     apiService.GetPg(),
		logger: apiService.GetLogger(),
	}

	apiService.GetRouter().GET("/users/", hs.GetMe)
	apiService.GetRouter().GET("/users/:vkId", hs.GetUserByVkID)
	apiService.GetRouter().PATCH("/users/", hs.EditUser)

	apiService.GetRouter().NoRoute(hs.NoRoute)
	apiService.GetRouter().NoMethod(hs.NoRoute)
}
