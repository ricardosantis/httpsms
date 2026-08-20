package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/NdoleStudio/httpsms/pkg/services"
	"github.com/NdoleStudio/httpsms/pkg/telemetry"
	"github.com/NdoleStudio/httpsms/pkg/validators"
	"github.com/NdoleStudio/stacktrace"
	"github.com/gofiber/fiber/v3"
	"github.com/stripe/stripe-go/v78"
)

// StripeHandler handles Stripe webhooks and checkout sessions
type StripeHandler struct {
	handler
	logger    telemetry.Logger
	tracer    telemetry.Tracer
	service   *services.StripeService
	validator *validators.StripeHandlerValidator
}

// NewStripeHandler creates a new instance of StripeHandler
func NewStripeHandler(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	service *services.StripeService,
	validator *validators.StripeHandlerValidator,
) *StripeHandler {
	h := &StripeHandler{
		logger:    logger.WithService(fmt.Sprintf("%T", &StripeHandler{})),
		tracer:    tracer,
		service:   service,
		validator: validator,
	}
	return h
}

// RegisterRoutes registers endpoints for Stripe
func (h *StripeHandler) RegisterRoutes(app *fiber.App, middlewares ...fiber.Handler) {
	router := app.Group("stripe")
	h.register(router, fiber.MethodPost, "/webhook", nil, h.Webhook)
	h.register(router, fiber.MethodPost, "/checkout-session", middlewares, h.CreateCheckoutSession)
	h.register(router, fiber.MethodPost, "/customer-portal", middlewares, h.CreateCustomerPortalSession)
}

// StripeCheckoutRequest request payload for checkout session
type StripeCheckoutRequest struct {
	PlanID     string `json:"plan_id" validate:"required"`
	PriceID    string `json:"price_id"`
	SuccessURL string `json:"success_url"`
	CancelURL  string `json:"cancel_url"`
}

// CreateCheckoutSession handles POST /v1/stripe/checkout-session
func (h *StripeHandler) CreateCheckoutSession(c fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	userID := h.userIDFomContext(c)
	if userID == "" {
		return h.responseUnauthorized(c)
	}

	var req StripeCheckoutRequest
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
		ctxLogger.Error(stacktrace.Propagatef(err, "cannot create stripe checkout session for user [%s]", userID))
		return h.responseInternalServerError(c)
	}

	return h.responseOK(c, "checkout session created successfully", fiber.Map{
		"url": checkoutURL,
	})
}

// CreateCustomerPortalSession handles POST /v1/stripe/customer-portal
func (h *StripeHandler) CreateCustomerPortalSession(c fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	userID := h.userIDFomContext(c)
	if userID == "" {
		return h.responseUnauthorized(c)
	}

	portalURL, err := h.service.CreateCustomerPortalSession(ctx, userID, "")
	if err != nil {
		ctxLogger.Error(stacktrace.Propagatef(err, "cannot create stripe customer portal for user [%s]", userID))
		return h.responseInternalServerError(c)
	}

	return h.responseOK(c, "customer portal session created successfully", fiber.Map{
		"url": portalURL,
	})
}

// Webhook handles incoming Stripe webhook events
func (h *StripeHandler) Webhook(c fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	signature := c.Get("Stripe-Signature")
	if errors := h.validator.ValidateEvent(ctx, signature, c.Body()); len(errors) != 0 {
		ctxLogger.Warn(stacktrace.NewErrorf("stripe webhook validation error: %v", errors))
		return h.responseUnprocessableEntity(c, errors, "invalid stripe webhook signature")
	}

	var event stripe.Event
	if err := json.Unmarshal(c.Body(), &event); err != nil {
		ctxLogger.Error(stacktrace.Propagatef(err, "cannot unmarshal stripe webhook payload"))
		return h.responseBadRequest(c, errors.New("invalid stripe event payload"))
	}

	if err := h.handleEvent(ctx, c.OriginalURL(), &event); err != nil {
		ctxLogger.Error(stacktrace.Propagatef(err, "cannot handle stripe event [%s]", event.Type))
		return h.responseInternalServerError(c)
	}

	return h.responseNoContent(c, "stripe event processed successfully")
}

func (h *StripeHandler) handleEvent(ctx context.Context, source string, event *stripe.Event) error {
	switch event.Type {
	case "checkout.session.completed":
		var sess stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &sess); err != nil {
			return stacktrace.Propagatef(err, "cannot unmarshal checkout session event data")
		}
		return h.service.HandleCheckoutSessionCompleted(ctx, source, &sess)

	case "customer.subscription.created", "customer.subscription.updated":
		var sub stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
			return stacktrace.Propagatef(err, "cannot unmarshal subscription event data")
		}
		return h.service.HandleSubscriptionUpdated(ctx, source, &sub)

	case "customer.subscription.deleted":
		var sub stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
			return stacktrace.Propagatef(err, "cannot unmarshal subscription deleted event data")
		}
		return h.service.HandleSubscriptionDeleted(ctx, source, &sub)

	default:
		return nil
	}
}
