package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "time-tracker/docs"

	"time-tracker/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		users := api.Group("/users")
		{
			users.GET("/", h.getUsers)
			users.POST("/", h.createUser)
			users.GET("/:id", h.getUser)
			users.DELETE("/:id", h.deleteUser)
			users.PUT("/:id", h.updateUser)
			users.GET("/:id/stats", h.getUsersStats)

		}

		tasks := api.Group("/tasks")
		{
			tasks.POST("/:id", h.createTask)
			tasks.POST("/:id/start", h.startTask)
			tasks.POST("/:id/stop", h.stopTask)
		}
	}

	return router
}
