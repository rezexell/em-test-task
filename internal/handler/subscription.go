package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rezexell/em-test-task/internal/model"
	"log/slog"
	"net/http"
	"time"
)

func (h *Handler) CreateSub(c *gin.Context) {
	const fn = "handler.CreateSub"
	h.logger.Info("context", slog.String("fn", fn))

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

func (h *Handler) UpdateSub(c *gin.Context) {
	const fn = "handler.UpdateSub"
	h.logger.Info("context", slog.String("fn", fn))

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

func (h *Handler) DeleteSub(c *gin.Context) {
	const fn = "handler.DeleteSub"
	h.logger.Info("context", slog.String("fn", fn))

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

func (h *Handler) GetAllSubs(c *gin.Context) {
	const fn = "handler.GetAllSubs"
	h.logger.Info("context", slog.String("fn", fn))

	subs, err := h.service.ListAllSubscriptions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(subs) == 0 || subs == nil {
		subs = []*model.Subscription{}
	}
	c.JSON(http.StatusOK, subs)
	return
}

func (h *Handler) GetSubByID(c *gin.Context) {
	const fn = "handler.GetSubByID"
	h.logger.Info("context", slog.String("fn", fn))

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

func (h *Handler) GetFilteredSubs(c *gin.Context) {
	const fn = "handler.GetFilteredSubs"
	h.logger.Info("context", slog.String("fn", fn))

	userIDStr := c.Query("user_id")
	serviceName := c.Query("service_name")

	var userIDPtr *uuid.UUID
	if userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format"})
			return
		}
		userIDPtr = &userID
	}

	var serviceNamePtr *string
	if serviceName != "" {
		serviceNamePtr = &serviceName
	}

	subs, err := h.service.ListSubscriptionsWithFilters(
		c.Request.Context(),
		userIDPtr,
		serviceNamePtr,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if subs == nil {
		subs = []*model.Subscription{}
	}

	c.JSON(http.StatusOK, subs)
}

func (h *Handler) GetTotalCost(c *gin.Context) {
	const fn = "handler.GetTotalCost"
	h.logger.Info("context", slog.String("fn", fn))

	userIDStr := c.Query("user_id")
	serviceName := c.Query("service_name")
	startPeriodStr := c.Query("start_period")
	endPeriodStr := c.Query("end_period")

	if startPeriodStr == "" || endPeriodStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_period and end_period are required"})
		return
	}

	startPeriod, err := time.Parse("01/2006", startPeriodStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_period format, use MM/YYYY"})
		return
	}

	endPeriod, err := time.Parse("01/2006", endPeriodStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_period format, use MM/YYYY"})
		return
	}

	startPeriod = time.Date(startPeriod.Year(), startPeriod.Month(), 1, 0, 0, 0, 0, time.UTC)
	endPeriod = time.Date(endPeriod.Year(), endPeriod.Month()+1, 0, 0, 0, 0, 0, time.UTC)

	var userIDPtr *uuid.UUID
	if userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format"})
			return
		}
		userIDPtr = &userID
	}

	var serviceNamePtr *string
	if serviceName != "" {
		serviceNamePtr = &serviceName
	}

	total, err := h.service.TotalSubscriptionCost(
		c.Request.Context(),
		userIDPtr,
		serviceNamePtr,
		startPeriod,
		endPeriod,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total_cost": total})
}
