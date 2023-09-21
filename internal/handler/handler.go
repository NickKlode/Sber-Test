package handler

import (
	"sber-test/internal/service"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "sber-test/docs"

	swaggerFiles "github.com/swaggo/files"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	lists := router.Group("/lists")
	{
		lists.POST("/", h.createList)
		lists.GET("/", h.getAll)
		lists.DELETE("/:id", h.deleteList)
		lists.PUT("/:id", h.updateList)
		lists.POST("/find", h.getByDate)
	}

	return router
}
