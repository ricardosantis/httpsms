package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/NdoleStudio/httpsms/pkg/entities"
	"github.com/NdoleStudio/httpsms/pkg/responses"
	"github.com/NdoleStudio/httpsms/pkg/services"
	"github.com/NdoleStudio/httpsms/pkg/telemetry"
	"github.com/gofiber/fiber/v3"
	"github.com/NdoleStudio/stacktrace"
)

// AdminHandler handles HTTP requests for admin operations
type AdminHandler struct {
	handler
	logger          telemetry.Logger
	tracer          telemetry.Tracer
	userService     *services.UserService
	eventsQueueUser entities.UserID
}

// NewAdminHandler creates a new instance of AdminHandler
func NewAdminHandler(logger telemetry.Logger, tracer telemetry.Tracer, userService *services.UserService) *AdminHandler {
	return &AdminHandler{
		logger:          logger.WithService("AdminHandler"),
		tracer:          tracer,
		userService:     userService,
		eventsQueueUser: entities.UserID(os.Getenv("EVENTS_QUEUE_USER_ID")),
	}
}

func (h *AdminHandler) isSystemUser(userID entities.UserID) bool {
	return h.eventsQueueUser != "" && userID == h.eventsQueueUser
}

// RegisterRoutes registers the routes for the AdminHandler
func (h *AdminHandler) RegisterRoutes(router fiber.Router, middlewares ...fiber.Handler) {
	h.register(router, fiber.MethodGet, "/v1/admin/users", middlewares, h.IndexUsers)
	h.register(router, fiber.MethodPost, "/v1/admin/users/:userID/block", middlewares, h.BlockUser)
	h.register(router, fiber.MethodPost, "/v1/admin/users/:userID/unblock", middlewares, h.UnblockUser)
	h.register(router, fiber.MethodDelete, "/v1/admin/users/:userID", middlewares, h.DeleteUser)
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

// BlockUser blocks a user from accessing the platform
// @Summary Block a user (Admin only)
// @Description Blocks a user by setting active=false.
// @Tags Admin
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {object} responses.OkString
// @Failure 400 {object} responses.BadRequest
// @Failure 401 {object} responses.Unauthorized
// @Failure 403 {object} responses.Unauthorized
// @Router /v1/admin/users/{userID}/block [post]
// @Security BearerAuth
func (h *AdminHandler) BlockUser(c fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	userID := entities.UserID(c.Params("userID"))
	if userID == "" {
		return h.responseBadRequest(c, stacktrace.NewError("missing userID parameter"))
	}

	if h.isSystemUser(userID) {
		return h.responseBadRequest(c, stacktrace.NewError("cannot block the system user"))
	}

	if err := h.userService.Block(ctx, userID); err != nil {
		ctxLogger.Error(h.tracer.WrapErrorSpan(span, err))
		return h.responseInternalServerError(c)
	}

	return c.Status(http.StatusOK).JSON(responses.OkString{Status: "success", Message: "User blocked successfully"})
}

// UnblockUser restores access for a user
// @Summary Unblock a user (Admin only)
// @Description Unblocks a user by setting active=true.
// @Tags Admin
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {object} responses.OkString
// @Failure 400 {object} responses.BadRequest
// @Failure 401 {object} responses.Unauthorized
// @Failure 403 {object} responses.Unauthorized
// @Router /v1/admin/users/{userID}/unblock [post]
// @Security BearerAuth
func (h *AdminHandler) UnblockUser(c fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	userID := entities.UserID(c.Params("userID"))
	if userID == "" {
		return h.responseBadRequest(c, stacktrace.NewError("missing userID parameter"))
	}

	if h.isSystemUser(userID) {
		return h.responseBadRequest(c, stacktrace.NewError("cannot modify the system user"))
	}

	if err := h.userService.Unblock(ctx, userID); err != nil {
		ctxLogger.Error(h.tracer.WrapErrorSpan(span, err))
		return h.responseInternalServerError(c)
	}

	return c.Status(http.StatusOK).JSON(responses.OkString{Status: "success", Message: "User unblocked successfully"})
}

// DeleteUser deletes a user (admin-initiated)
// @Summary Delete a user (Admin only)
// @Description Deletes a user and all associated data.
// @Tags Admin
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {object} responses.OkString
// @Failure 400 {object} responses.BadRequest
// @Failure 401 {object} responses.Unauthorized
// @Failure 403 {object} responses.Unauthorized
// @Router /v1/admin/users/{userID} [delete]
// @Security BearerAuth
func (h *AdminHandler) DeleteUser(c fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	userID := entities.UserID(c.Params("userID"))
	if userID == "" {
		return h.responseBadRequest(c, stacktrace.NewError("missing userID parameter"))
	}

	if h.isSystemUser(userID) {
		return h.responseBadRequest(c, stacktrace.NewError("cannot delete the system user"))
	}

	if err := h.userService.AdminDelete(ctx, userID); err != nil {
		ctxLogger.Error(h.tracer.WrapErrorSpan(span, err))
		return h.responseInternalServerError(c)
	}

	return c.Status(http.StatusOK).JSON(responses.OkString{Status: "success", Message: "User deleted successfully"})
}
