package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rezexell/em-test-task/internal/service"
	sloggin "github.com/samber/slog-gin"
	"log/slog"
)

type Handler struct {
	service *service.Service
	logger  *slog.Logger
}

func NewHandler(service *service.Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) InitRouter() *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(sloggin.New(h.logger))

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
