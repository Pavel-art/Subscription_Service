package api

import (
	"SubscriptionService/internal/api/dto"
	"SubscriptionService/internal/application/app_interfaces"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Handler struct {
	route        *gin.Engine
	service      app_interfaces.ISubService
	customLogger *zerolog.Logger
}

func NewHandler(r *gin.Engine, s app_interfaces.ISubService, l *zerolog.Logger) *Handler {
	handler := &Handler{
		route:        r,
		service:      s,
		customLogger: l,
	}
	handler.registerRoutes()
	return handler
}

func (h *Handler) registerRoutes() {
	h.route.GET("/health", h.Health)

	api := h.route.Group("/api/v1")
	{
		subs := api.Group("/subscriptions")
		{
			subs.POST("", h.Create)
			subs.GET("/:id", h.GetById)
			subs.GET("", h.GetAll)
			subs.PUT("/:id", h.Update)
			subs.DELETE("/:id", h.Delete)
			subs.GET("/cost", h.CalculateCost)
		}
	}
}

func (h *Handler) Health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (h *Handler) Create(ctx *gin.Context) {
	h.customLogger.Debug().Msg("Create subscription: started")

	var request dto.CreateSubscriptionRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.customLogger.
			Error().
			Err(err).
			Msg("Create subscription: invalid request")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.service.Create(ctx, request)
	if err != nil {
		h.customLogger.Error().Err(err).Msg("Create subscription: service error")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create subscription"})
		return
	}

	h.customLogger.Info().
		Str("service", created.ServiceName).
		Str("createdId", created.Id.String()).
		Msg("Create subscription: created")

	ctx.JSON(http.StatusCreated, created)
}

func (h *Handler) GetById(ctx *gin.Context) {
	h.customLogger.Debug().Msg("Get subscription by id: started")

	idStr := ctx.Param("id")
	if idStr == "" {
		h.customLogger.
			Warn().
			Msg("Get subscription by id: empty id parameter")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		h.customLogger.
			Warn().Err(err).
			Str("id", idStr).
			Msg("Get subscription by id: invalid id format")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	sub, err := h.service.GetById(ctx, id)
	if err != nil {
		h.customLogger.
			Error().
			Err(err).
			Str("id", id.String()).
			Msg("Get subscription by id: service error")
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	h.customLogger.
		Info().
		Str("id", id.String()).
		Msg("Get subscription by id: success")

	ctx.JSON(http.StatusOK, sub)
}

func (h *Handler) GetAll(ctx *gin.Context) {
	h.customLogger.
		Debug().
		Msg("Get all subscriptions: started")

	page, err := strconv.ParseInt(ctx.DefaultQuery("page", "1"), 10, 64)
	if err != nil || page < 1 {
		h.customLogger.
			Warn().
			Int64("providedPage", page).
			Msg("Get all subscriptions: invalid page, using default")

		page = 1
	}

	pageSize, err := strconv.ParseInt(ctx.DefaultQuery("page_size", "20"), 10, 64)
	if err != nil || pageSize < 1 || pageSize > 100 {
		h.customLogger.
			Warn().
			Int64("providedPageSize", pageSize).
			Msg("Get all subscriptions: invalid page size, using default")

		pageSize = 20
	}

	h.customLogger.
		Debug().Int64("page", page).
		Int64("pageSize", pageSize).
		Msg("Get all subscriptions: fetching")

	res, err := h.service.GetAll(ctx, page, pageSize)
	if err != nil {
		h.customLogger.
			Error().
			Err(err).
			Int64("page", page).
			Int64("pageSize", pageSize).
			Msg("Get all subscriptions: service error")

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get all subscriptions"})
		return
	}

	h.customLogger.
		Info().
		Int64("page", page).
		Int64("pageSize", pageSize).
		Msg("Get all subscriptions: success")
	ctx.JSON(http.StatusOK, res)
}

func (h *Handler) Update(ctx *gin.Context) {
	h.customLogger.
		Debug().
		Msg("Update subscription: started")

	idS := ctx.Param("id")
	if idS == "" {
		h.customLogger.
			Warn().
			Msg("Update subscription: empty id parameter")

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	id, err := uuid.Parse(idS)
	if err != nil {
		h.customLogger.
			Warn().Err(err).
			Str("id", idS).
			Msg("Update subscription: invalid id format")

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	var request dto.UpdateSubscriptionRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		h.customLogger.
			Warn().
			Err(err).
			Str("id", id.String()).Msg("Update subscription: invalid JSON")

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request body",
			"details": err.Error(),
		})
		return
	}

	updated, err := h.service.Update(ctx, id, request)
	if err != nil {
		h.customLogger.Error().
			Err(err).Str("id", id.String()).
			Msg("Update subscription: service error")

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update subscription"})
		return
	}

	h.customLogger.
		Info().
		Str("id", id.String()).
		Msg("Update subscription: success")
	ctx.JSON(http.StatusOK, updated)
}

func (h *Handler) Delete(ctx *gin.Context) {
	h.customLogger.
		Debug().
		Msg("Delete subscription: started")

	idStr := ctx.Param("id")
	if idStr == "" {
		h.customLogger.
			Warn().
			Msg("Delete subscription: empty id parameter")

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		h.customLogger.
			Warn().
			Err(err).
			Str("id", idStr).
			Msg("Delete subscription: invalid id format")

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	err = h.service.Delete(ctx, id)
	if err != nil {
		h.customLogger.
			Error().Err(err).
			Str("id", id.String()).
			Msg("Delete subscription: service error")

		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	h.customLogger.
		Info().
		Str("id", id.String()).
		Msg("Delete subscription: success")

	ctx.JSON(http.StatusNoContent, nil)
}

func (h *Handler) CalculateCost(ctx *gin.Context) {
	h.customLogger.
		Debug().
		Msg("Calculate cost: started")

	var request dto.CostCalculationQueryRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {

		h.customLogger.
			Warn().Err(err).
			Msg("Calculate cost: invalid query parameters")

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	if !request.From.IsZero() && !request.To.IsZero() && request.From.After(request.To) {
		h.customLogger.
			Warn().Time("from", request.From).
			Time("to", request.To).
			Msg("Calculate cost: invalid date range")

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid date range: 'from' cannot be after 'to'",
		})
		return
	}

	h.customLogger.
		Debug().
		Interface("filters", request).
		Msg("Calculate cost: processing")

	totalCost, err := h.service.CalculateTotalCost(ctx, request)
	if err != nil {
		h.customLogger.
			Error().
			Err(err).
			Msg("Calculate cost: service error")

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to calculate total cost",
		})
		return
	}

	h.customLogger.
		Info().
		Int64("totalCost", totalCost).
		Msg("Calculate cost: success")

	ctx.JSON(http.StatusOK, gin.H{
		"total_cost": totalCost,
		"currency":   "RUB",
		"filters":    request,
	})
}
