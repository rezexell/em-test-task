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
		sub.GET("/", h.getAllSubs)
	}
	return router
}
