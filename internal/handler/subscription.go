package handler

import (
	"github.com/gin-gonic/gin"
)

func (h *Handler) getAllSubs(c *gin.Context) {
	c.JSON(200, gin.H{"message": h.service.GetAllSubs()})
}
