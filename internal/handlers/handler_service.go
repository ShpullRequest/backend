package handlers

import (
	_ "github.com/ShpullRequest/backend/docs"
	"github.com/ShpullRequest/backend/internal/api"
	"github.com/ShpullRequest/backend/internal/repository"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
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

	apiService.GetRouter().GET("/achievements/", hs.GetAllAchievements)
	apiService.GetRouter().GET("/achievements/:achievementId", hs.GetAchievement)
	apiService.GetRouter().POST("/achievements/", hs.NewAchievement)
	apiService.GetRouter().PATCH("/achievements/:achievementId", hs.EditAchievement)

	apiService.GetRouter().GET("/users/", hs.GetMe)
	apiService.GetRouter().GET("/users/:vkId/", hs.GetUserByVkID)
	apiService.GetRouter().PATCH("/users/", hs.EditUser)

	apiService.GetRouter().GET("/companies/", hs.GetAllCompanies)
	apiService.GetRouter().GET("/companies/:companyId/", hs.GetCompany)
	apiService.GetRouter().GET("/companies/my/", hs.GetMyCompanies)
	apiService.GetRouter().POST("/companies/", hs.NewCompany)
	apiService.GetRouter().POST("/companies/:companyId/accept/", hs.AcceptCompany)

	apiService.GetRouter().GET("/places/", hs.GetAllPlaces)
	apiService.GetRouter().GET("/places/search/:query/", hs.SearchPlaces)
	apiService.GetRouter().GET("/places/:placeId", hs.GetPlace)
	apiService.GetRouter().GET("/places/:placeId/reviews", hs.GetReviewsPlace)
	apiService.GetRouter().POST("/places/", hs.NewPlace)
	apiService.GetRouter().POST("/places/:placeId/reviews/", hs.NewReviewPlace)
	apiService.GetRouter().PATCH("/places/:placeId", hs.EditPlace)
	apiService.GetRouter().PATCH("/places/:placeId/reviews/", hs.EditReviewPlace)

	apiService.GetRouter().GET("/events/", hs.GetAllEvents)
	apiService.GetRouter().GET("/events/company/:companyId", hs.GetCompanyEvents)
	apiService.GetRouter().GET("/events/search/:query/", hs.SearchEvents)
	apiService.GetRouter().GET("/events/:eventId/", hs.GetEvent)
	apiService.GetRouter().GET("/events/:eventId/reviews/", hs.GetReviewsEvent)
	apiService.GetRouter().POST("/events/", hs.NewEvent)
	apiService.GetRouter().POST("/events/:eventId/reviews/", hs.NewReviewEvent)
	apiService.GetRouter().PATCH("/events/:eventId/", hs.EditEvent)
	apiService.GetRouter().PATCH("/events/:eventId/reviews/", hs.EditReviewsEvent)

	apiService.GetRouter().GET("/routes/", hs.GetAllRoutes)
	apiService.GetRouter().GET("/routes/company/:companyId/", hs.GetCompanyRoutes)
	apiService.GetRouter().GET("/routes/search/:query/", hs.SearchRoutes)
	apiService.GetRouter().GET("/routes/:routeId/", hs.GetRoute)
	apiService.GetRouter().GET("/routes/:routeId/reviews/", hs.GetReviewsRoutes)
	apiService.GetRouter().POST("/routes/", hs.NewRoute)
	apiService.GetRouter().POST("/routes/:routeId/reviews/", hs.NewReviewRoute)
	apiService.GetRouter().PATCH("/routes/:routeId/", hs.EditRoute)
	apiService.GetRouter().PATCH("/routes/:routeId/reviews/", hs.EditReviewRoute)

	apiService.GetRouter().GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiService.GetRouter().NoRoute(hs.NoRoute)
	apiService.GetRouter().NoMethod(hs.NoRoute)
}
