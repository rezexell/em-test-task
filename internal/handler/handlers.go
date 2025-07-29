package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rezexell/em-test-task/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRouter() *gin.Engine {
	router := gin.Default()

	sub := router.Group("/sub")
	{
		sub.POST("/", h.createSub)
		sub.PUT("/", h.updateSub)
		sub.DELETE("/:id", h.deleteSub)
		sub.GET("/all", h.getAllSubs) // Изменен путь
		sub.GET("/:id", h.getSubByID)
		sub.GET("/by-user", h.getSubsByUser)
	}
	return router
}
