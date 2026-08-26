package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/NdoleStudio/httpsms/pkg/services"
	"github.com/NdoleStudio/httpsms/pkg/telemetry"
	"github.com/NdoleStudio/httpsms/pkg/validators"
	"github.com/NdoleStudio/stacktrace"
	"github.com/gofiber/fiber/v3"
)

// MercadopagoHandler handles Mercado Pago webhooks and checkout sessions
type MercadopagoHandler struct {
	handler
	logger    telemetry.Logger
	tracer    telemetry.Tracer
	service   *services.MercadopagoService
	validator *validators.MercadopagoHandlerValidator
}

// NewMercadopagoHandler creates a new instance of MercadopagoHandler
func NewMercadopagoHandler(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	service *services.MercadopagoService,
	validator *validators.MercadopagoHandlerValidator,
) *MercadopagoHandler {
	h := &MercadopagoHandler{
		logger:    logger.WithService(fmt.Sprintf("%T", &MercadopagoHandler{})),
		tracer:    tracer,
		service:   service,
		validator: validator,
	}
	return h
}

// RegisterRoutes registers endpoints for Mercado Pago
func (h *MercadopagoHandler) RegisterRoutes(app *fiber.App, middlewares ...fiber.Handler) {
	router := app.Group("mercadopago")
	h.register(router, fiber.MethodPost, "/webhook", nil, h.Webhook)
	h.register(router, fiber.MethodPost, "/checkout-session", middlewares, h.CreateCheckoutSession)

	v1Router := app.Group("v1/mercadopago")
	h.register(v1Router, fiber.MethodPost, "/webhook", nil, h.Webhook)
	h.register(v1Router, fiber.MethodPost, "/checkout-session", middlewares, h.CreateCheckoutSession)
}

// MercadopagoCheckoutRequest request payload for checkout session
type MercadopagoCheckoutRequest struct {
	PlanID     string `json:"plan_id" validate:"required"`
	PriceID    string `json:"price_id"`
	SuccessURL string `json:"success_url"`
	CancelURL  string `json:"cancel_url"`
}

// CreateCheckoutSession handles POST /v1/mercadopago/checkout-session
func (h *MercadopagoHandler) CreateCheckoutSession(c fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	userID := h.userIDFomContext(c)
	if userID == "" {
		return h.responseUnauthorized(c)
	}

	var req MercadopagoCheckoutRequest
	if err := c.Bind().Body(&req); err != nil {
		return h.responseBadRequest(c, errors.New("invalid request body"))
	}

	checkoutURL, err := h.service.CreateCheckoutSession(ctx, services.CreateCheckoutSessionParams{
		UserID:     userID,
		PlanID:     req.PlanID,
		PriceID:    req.PriceID,
		SuccessURL: req.SuccessURL,
		CancelURL:  req.CancelURL,
	})
	if err != nil {
		ctxLogger.Error(stacktrace.Propagatef(err, "cannot create mercadopago checkout session for user [%s]", userID))
		return h.responseInternalServerError(c)
	}

	return h.responseOK(c, "checkout session created successfully", fiber.Map{
		"url": checkoutURL,
	})
}

// WebhookEvent represents a generic Mercado Pago webhook payload
type WebhookEvent struct {
	Action string `json:"action"`
	Type   string `json:"type"`
	Data   struct {
		ID string `json:"id"`
	} `json:"data"`
}

// Webhook handles incoming Mercado Pago webhook events
func (h *MercadopagoHandler) Webhook(c fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	xSignature := c.Get("x-signature")
	xRequestID := c.Get("x-request-id")
	
	// Query params: MP webhooks might send data.id in query string for V1, or in body for V2
	dataID := c.Query("data.id")

	var event WebhookEvent
	// Parse the body if available
	if len(c.Body()) > 0 {
		if err := json.Unmarshal(c.Body(), &event); err != nil {
			ctxLogger.Warn(stacktrace.Propagatef(err, "cannot unmarshal mercadopago webhook payload"))
		}
	}
	
	if dataID == "" && event.Data.ID != "" {
		dataID = event.Data.ID
	}

	// Validate Event Signature
	if errorsList := h.validator.ValidateEvent(ctx, xSignature, xRequestID, dataID); len(errorsList) != 0 {
		ctxLogger.Warn(stacktrace.NewErrorf("mercadopago webhook validation error: %v", errorsList))
		return h.responseUnprocessableEntity(c, errorsList, "invalid mercadopago webhook signature")
	}

	if err := h.handleEvent(ctx, c.OriginalURL(), dataID, c.Query("type"), event.Type, event.Action); err != nil {
		ctxLogger.Error(stacktrace.Propagatef(err, "cannot handle mercadopago event [%s]", dataID))
		return h.responseInternalServerError(c)
	}

	return h.responseNoContent(c, "mercadopago event processed successfully")
}

func (h *MercadopagoHandler) handleEvent(ctx context.Context, source, dataID, queryType, bodyType, action string) error {
	if dataID == "" {
		return nil
	}
	
	eventType := queryType
	if eventType == "" {
		eventType = bodyType
	}
	
	// "subscription_preapproval" or "preapproval"
	if strings.Contains(eventType, "preapproval") {
		return h.service.HandleSubscriptionUpdated(ctx, source, dataID)
	}
	
	if eventType == "payment" || action == "payment.created" || action == "payment.updated" {
		return h.service.HandlePaymentUpdated(ctx, source, dataID)
	}
	
	return nil
}
