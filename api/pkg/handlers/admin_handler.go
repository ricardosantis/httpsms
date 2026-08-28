package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/NdoleStudio/httpsms/pkg/responses"
	"github.com/NdoleStudio/httpsms/pkg/services"
	"github.com/NdoleStudio/httpsms/pkg/telemetry"
	"github.com/gofiber/fiber/v3"
)

// AdminHandler handles HTTP requests for admin operations
type AdminHandler struct {
	handler
	logger      telemetry.Logger
	tracer      telemetry.Tracer
	userService *services.UserService
}

// NewAdminHandler creates a new instance of AdminHandler
func NewAdminHandler(logger telemetry.Logger, tracer telemetry.Tracer, userService *services.UserService) *AdminHandler {
	return &AdminHandler{
		logger:      logger.WithService("AdminHandler"),
		tracer:      tracer,
		userService: userService,
	}
}

// RegisterRoutes registers the routes for the AdminHandler
func (h *AdminHandler) RegisterRoutes(router fiber.Router, middlewares ...fiber.Handler) {
	h.register(router, fiber.MethodGet, "/v1/admin/users", middlewares, h.IndexUsers)
}

// IndexUsers fetches all users
// @Summary Fetch all users (Admin only)
// @Description Fetches a paginated list of all users.
// @Tags Admin
// @Accept json
// @Produce json
// @Param skip query int false "Skip" default(0)
// @Param limit query int false "Limit" default(20)
// @Param query query string false "Search by email"
// @Success 200 {object} responses.UserListResponse
// @Failure 401 {object} responses.Unauthorized
// @Failure 403 {object} responses.Unauthorized
// @Router /v1/admin/users [get]
// @Security BearerAuth
func (h *AdminHandler) IndexUsers(c fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	skip, _ := strconv.Atoi(c.Query("skip", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	query := c.Query("query", "")

	users, totalCount, err := h.userService.IndexAll(ctx, skip, limit, query)
	if err != nil {
		ctxLogger.Error(h.tracer.WrapErrorSpan(span, err))
		return h.responseInternalServerError(c)
	}

	response := &responses.UserListResponse{}
	response.Data.Items = users
	response.Data.TotalCount = totalCount
	response.Message = "Users loaded successfully"
	response.Status = "success"

	ctxLogger.Info(fmt.Sprintf("loaded [%d] users", len(users)))
	return c.Status(http.StatusOK).JSON(response)
}
