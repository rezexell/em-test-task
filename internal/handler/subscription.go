package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rezexell/em-test-task/internal/model"
	"net/http"
)

func (h *Handler) createSub(c *gin.Context) {
	var req model.SubReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateSubscription(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "subscription created"})
}

func (h *Handler) updateSub(c *gin.Context) {
	var req model.SubReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.service.UpdateSubscription(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "subscription updated"})
}

func (h *Handler) deleteSub(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.service.DeleteSubscription(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "subscription deleted"})
	return
}

func (h *Handler) getAllSubs(c *gin.Context) {
	subs, err := h.service.ListAllSubscriptions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(subs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no subs found"})
		return
	}
	c.JSON(http.StatusOK, subs)
	return
}

func (h *Handler) getSubByID(c *gin.Context) {
	//test
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	sub, err := h.service.GetSubscription(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if sub == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, sub)
	return
}

func (h *Handler) getSubsByUser(c *gin.Context) {
	userIDStr := c.Query("userid")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required query parameter: userid"})
		return
	}

	// Валидируем UUID
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
		return
	}

	subs, err := h.service.ListUserSubscriptions(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	if len(subs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no subs found"})
	}
	c.JSON(http.StatusOK, subs)
}
