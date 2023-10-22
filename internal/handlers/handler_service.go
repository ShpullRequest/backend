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
	apiService.GetRouter().GET("/users/:vkId/", hs.GetUserByVkID)
	apiService.GetRouter().PATCH("/users/", hs.EditUser)

	apiService.GetRouter().GET("/companies/", hs.GetAllCompanies)
	apiService.GetRouter().GET("/companies/:companyId/", hs.GetCompany)
	apiService.GetRouter().GET("/companies/my/", hs.GetMyCompanies)
	apiService.GetRouter().POST("/companies/", hs.NewCompany)
	apiService.GetRouter().POST("/companies/:companyId/accept/", hs.AcceptCompany)

	apiService.GetRouter().GET("/events/", hs.GetAllEvents)
	apiService.GetRouter().GET("/events/company/:companyId", hs.GetCompanyEvents)
	apiService.GetRouter().GET("/events/search/:query/", hs.SearchEvents)
	apiService.GetRouter().GET("/events/:eventId/", hs.GetEvent)
	apiService.GetRouter().GET("/events/:eventId/reviews/", hs.GetReviewsEvent)
	apiService.GetRouter().POST("/events/:eventId/reviews/", hs.NewReviewEvent)
	apiService.GetRouter().POST("/events/", hs.NewEvent)
	apiService.GetRouter().PATCH("/events/:eventId/reviews/", hs.EditReviewsEvent)
	apiService.GetRouter().PATCH("/events/:eventId/", hs.EditEvent)

	apiService.GetRouter().NoRoute(hs.NoRoute)
	apiService.GetRouter().NoMethod(hs.NoRoute)
}
