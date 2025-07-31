package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rezexell/em-test-task/internal/model"
	"log/slog"
	"net/http"
	"time"
)

// CreateSub
// @Summary Создать новую подписку
// @Description Создает новую подписочную запись
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param input body model.Subscription true "Данные подписки"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string "Пример: {\"error\": \"invalid UUID format\"}"
// @Failure 500 {object} map[string]string "Пример: {\"error\": \"database connection failed\"}"
// @Router /sub [post]
func (h *Handler) CreateSub(c *gin.Context) {
	const fn = "handler.CreateSub"
	h.logger.Info("context", slog.String("fn", fn))

	var sub model.Subscription
	if err := c.ShouldBindJSON(&sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := sub.AfterBind(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if sub.ID == uuid.Nil {
		sub.ID = uuid.New()
	}

	if err := h.service.CreateSubscription(c.Request.Context(), &sub); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "subscription created", "id": sub.ID})
	return
}

// UpdateSub
// @Summary Обновить существующую подписку
// @Description Обновляет данные подписки по ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param input body model.Subscription true "Обновленные данные подписки"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string "Пример: {\"error\": \"start_date: required field\"}"
// @Failure 500 {object} map[string]string "Пример: {\"error\": \"subscription not found\"}"
// @Router /sub [put]
func (h *Handler) UpdateSub(c *gin.Context) {
	const fn = "handler.UpdateSub"
	h.logger.Info("context", slog.String("fn", fn))

	var sub model.Subscription
	if err := c.ShouldBindJSON(&sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := sub.AfterBind(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if sub.ID == uuid.Nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "id is required"})
		return
	}

	if err := h.service.UpdateSubscription(c.Request.Context(), &sub); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "subscription updated"})
	return
}

// DeleteSub
// @Summary Удалить подписку
// @Description Удаляет подписку по ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "ID подписки (UUID)"
// @Success 204
// @Failure 400 {object} map[string]string "Пример: {\"error\": \"invalid id format\"}"
// @Failure 500 {object} map[string]string "Пример: {\"error\": \"delete operation failed\"}"
// @Router /sub/{id} [delete]
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

// GetAllSubs
// @Summary Получить все подписки
// @Description Возвращает список всех подписок
// @Tags subscriptions
// @Produce json
// @Success 200 {array} model.Subscription
// @Failure 500 {object} map[string]string "Пример: {\"error\": \"failed to fetch subscriptions\"}"
// @Router /sub [get]
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

	response := make([]gin.H, 0, len(subs))
	for _, sub := range subs {
		response = append(response, sub.ToResponse())
	}
	c.JSON(http.StatusOK, response)
	return
}

// GetSubByID
// @Summary Получить подписку по ID
// @Description Возвращает подписку по её идентификатору
// @Tags subscriptions
// @Produce json
// @Param id path string true "ID подписки (UUID)"
// @Success 200 {object} model.Subscription
// @Failure 400 {object} map[string]string "Пример: {\"error\": \"invalid UUID format\"}"
// @Failure 404 {object} map[string]string "Пример: {\"error\": \"subscription not found\"}"
// @Failure 500 {object} map[string]string "Пример: {\"error\": \"database query failed\"}"
// @Router /sub/{id} [get]
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

	c.JSON(http.StatusOK, sub.ToResponse())
	return
}

// GetFilteredSubs
// @Summary Фильтрация подписок
// @Description Возвращает подписки по фильтрам
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "ID пользователя (UUID)"
// @Param service_name query string false "Название сервиса"
// @Success 200 {array} model.Subscription
// @Failure 400 {object} map[string]string "Пример: {\"error\": \"invalid user_id format\"}"
// @Failure 500 {object} map[string]string "Пример: {\"error\": \"filtering failed\"}"
// @Router /sub/filter [get]
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
	return
}

// GetTotalCost
// @Summary Расчет общей стоимости
// @Description Рассчитывает общую стоимость подписок за период
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "ID пользователя (UUID)"
// @Param service_name query string false "Название сервиса"
// @Param start_period query string true "Начало периода (MM/YYYY)" Example(01/2023)
// @Param end_period query string true "Конец периода (MM/YYYY)" Example(12/2023)
// @Success 200 {object} map[string]int "Пример: {\"total_cost\": 150}"
// @Failure 400 {object} map[string]string "Пример: {\"error\": \"invalid end_period format, use MM/YYYY\"}"
// @Failure 500 {object} map[string]string "Пример: {\"error\": \"cost calculation failed\"}"
// @Router /sub/total-cost [get]
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
	return
}
