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
		sub.POST("/", h.CreateSub)
		sub.PUT("/", h.UpdateSub)
		sub.DELETE("/:id", h.DeleteSub)
		sub.GET("/", h.GetAllSubs)
		sub.GET("/:id", h.GetSubByID)
		sub.GET("/filter/", h.GetFilteredSubs)
		sub.GET("/total-cost/", h.GetTotalCost)
	}
	return router
}
