package services

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/NdoleStudio/httpsms/pkg/entities"
	"github.com/NdoleStudio/httpsms/pkg/events"
	"github.com/NdoleStudio/httpsms/pkg/repositories"
	"github.com/NdoleStudio/httpsms/pkg/telemetry"
	"github.com/NdoleStudio/stacktrace"

	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/payment"
	"github.com/mercadopago/sdk-go/pkg/preapproval"
	"github.com/mercadopago/sdk-go/pkg/preference"
)

// MercadopagoService is responsible for handling Mercado Pago payments and webhooks
type MercadopagoService struct {
	service
	logger          telemetry.Logger
	tracer          telemetry.Tracer
	eventDispatcher *EventDispatcher
	userRepository  repositories.UserRepository
	mpClient        preapproval.Client
	prefClient      preference.Client
	paymentClient   payment.Client
}

// NewMercadopagoService creates a new instance of MercadopagoService
func NewMercadopagoService(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	userRepository repositories.UserRepository,
	eventDispatcher *EventDispatcher,
) *MercadopagoService {
	accessToken := os.Getenv("MERCADOPAGO_ACCESS_TOKEN")
	var mpClient preapproval.Client
	var prefClient preference.Client
	var paymentClient payment.Client
	if accessToken != "" {
		cfg, err := config.New(accessToken)
		if err == nil {
			mpClient = preapproval.NewClient(cfg)
			prefClient = preference.NewClient(cfg)
			paymentClient = payment.NewClient(cfg)
		}
	}

	return &MercadopagoService{
		logger:          logger.WithService(fmt.Sprintf("%T", &MercadopagoService{})),
		tracer:          tracer,
		userRepository:  userRepository,
		eventDispatcher: eventDispatcher,
		mpClient:        mpClient,
		prefClient:      prefClient,
		paymentClient:   paymentClient,
	}
}

// CreateCheckoutSession creates a new Mercado Pago Preapproval (Subscription) for a user
func (service *MercadopagoService) CreateCheckoutSession(ctx context.Context, params CreateCheckoutSessionParams) (string, error) {
	ctx, span, _ := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	if service.mpClient == nil {
		return "", stacktrace.NewError("mercadopago client is not initialized")
	}

	user, err := service.userRepository.Load(ctx, params.UserID)
	if err != nil {
		return "", stacktrace.Propagatef(err, "cannot load user [%s]", params.UserID)
	}

	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		appURL = "https://smsandroid.com.br"
	}

	successURL := params.SuccessURL
	if successURL == "" {
		successURL = appURL + "/billing?status=success"
	}

	// In Mercado Pago, we link the subscription to a Plan ID created in the dashboard.
	// If the frontend passes PriceID, we use it as the PreapprovalPlanID.
	
	// For plans that exceed MP subscription limits (> 4000 BRL), we use Checkout Pro Preference (one-off payment).
	isHighTierYearly := params.PlanID == "50k-yearly" || params.PlanID == "100k-yearly" || params.PlanID == "200k-yearly"
	
	if isHighTierYearly {
		priceMap := map[string]float64{
			"50k-yearly": 4990,
			"100k-yearly": 9900,
			"200k-yearly": 19900,
		}
		
		req := preference.Request{
			Items: []preference.ItemRequest{
				{
					Title:       "Assinatura " + params.PlanID,
					Quantity:    1,
					UnitPrice:   priceMap[params.PlanID],
					CurrencyID:  "BRL",
				},
			},
			Payer: &preference.PayerRequest{
				Email: user.Email,
			},
			BackURLs: &preference.BackURLsRequest{
				Success: successURL,
				Pending: successURL,
				Failure: appURL + "/billing?status=canceled",
			},
			AutoReturn:        "approved",
			ExternalReference: string(user.ID) + "|" + params.PlanID,
		}
		
		res, err := service.prefClient.Create(ctx, req)
		if err != nil {
			return "", stacktrace.Propagatef(err, "cannot create mercadopago preference for user [%s]", params.UserID)
		}
		return res.InitPoint, nil
	}

	planID := params.PriceID
	if planID == "" {
		planID = os.Getenv("MERCADOPAGO_PLAN_" + strings.ToUpper(strings.ReplaceAll(params.PlanID, "-", "_")))
	}

	if planID == "" {
		return "", stacktrace.NewError("no mercadopago plan ID configured for plan [%s]", params.PlanID)
	}

	// Create the Preapproval request
	req := preapproval.Request{
		PreapprovalPlanID: planID,
		PayerEmail:        user.Email,
		BackURL:           successURL,
		Reason:            "Assinatura " + params.PlanID,
		ExternalReference: string(user.ID) + "|" + params.PlanID,
		Status:            "pending",
	}

	res, err := service.mpClient.Create(ctx, req)
	if err != nil {
		return "", stacktrace.Propagatef(err, "cannot create mercadopago preapproval for user [%s]", params.UserID)
	}

	return res.InitPoint, nil
}

// HandleSubscriptionUpdated processes subscription (preapproval) created/updated webhooks
func (service *MercadopagoService) HandleSubscriptionUpdated(ctx context.Context, source string, preapprovalID string) error {
	ctx, span, ctxLogger := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	if service.mpClient == nil {
		return stacktrace.NewError("mercadopago client is not initialized")
	}

	sub, err := service.mpClient.Get(ctx, preapprovalID)
	if err != nil {
		return stacktrace.Propagatef(err, "cannot fetch mercadopago preapproval [%s]", preapprovalID)
	}

	// ExternalReference contains userID|planID
	parts := strings.Split(sub.ExternalReference, "|")
	var userIDStr, planID string
	if len(parts) >= 1 {
		userIDStr = parts[0]
	}
	if len(parts) >= 2 {
		planID = parts[1]
	}

	if userIDStr == "" {
		// Try to fallback by email
		user, err := service.userRepository.LoadByEmail(ctx, sub.PayerEmail)
		if err != nil {
			return stacktrace.Propagatef(err, "cannot find user for customer email [%s]", sub.PayerEmail)
		}
		userIDStr = string(user.ID)
	}

	user, err := service.userRepository.Load(ctx, entities.UserID(userIDStr))
	if err != nil {
		return stacktrace.Propagatef(err, "cannot load user [%s]", userIDStr)
	}

	status := sub.Status
	renewsAt := sub.NextPaymentDate

	// Based on status, we dispatch the corresponding event
	if status == "authorized" {
		payload := &events.UserSubscriptionUpdatedPayload{
			UserID:                user.ID,
			SubscriptionUpdatedAt: time.Now().UTC(),
			SubscriptionID:        sub.ID,
			SubscriptionName:      service.mapPlanToSubscriptionName(planID),
			SubscriptionRenewsAt:  renewsAt,
			SubscriptionStatus:    status,
		}

		// Also check if this is the first time the subscription is active, maybe trigger Created event.
		// For simplicity, UserSubscriptionUpdated covers limits refresh as well.
		event, err := service.createEvent(events.UserSubscriptionUpdated, source, payload)
		if err != nil {
			return stacktrace.Propagatef(err, "cannot create [%s] event for user [%s]", events.UserSubscriptionUpdated, user.ID)
		}
		if err = service.eventDispatcher.Dispatch(ctx, event); err != nil {
			return stacktrace.Propagatef(err, "cannot dispatch [%s] event for user [%s]", events.UserSubscriptionUpdated, user.ID)
		}
		ctxLogger.Info(fmt.Sprintf("mercadopago subscription [%s] authorized/updated for user [%s]", sub.ID, user.ID))

	} else if status == "cancelled" {
		payload := &events.UserSubscriptionCancelledPayload{
			UserID:                  user.ID,
			SubscriptionCancelledAt: time.Now().UTC(),
			SubscriptionID:          sub.ID,
			SubscriptionName:        entities.SubscriptionNameFree,
			SubscriptionEndsAt:      time.Now().UTC(),
			SubscriptionStatus:      status,
		}
		event, err := service.createEvent(events.UserSubscriptionCancelled, source, payload)
		if err != nil {
			return stacktrace.Propagatef(err, "cannot create [%s] event for user [%s]", events.UserSubscriptionCancelled, user.ID)
		}
		if err = service.eventDispatcher.Dispatch(ctx, event); err != nil {
			return stacktrace.Propagatef(err, "cannot dispatch [%s] event for user [%s]", events.UserSubscriptionCancelled, user.ID)
		}
		ctxLogger.Info(fmt.Sprintf("mercadopago subscription [%s] cancelled for user [%s]", sub.ID, user.ID))
	}

	return nil
}

func (service *MercadopagoService) mapPlanToSubscriptionName(plan string) entities.SubscriptionName {
	planLower := strings.ToLower(plan)
	if strings.Contains(planLower, "pro") {
		if strings.Contains(planLower, "yearly") {
			return entities.SubscriptionNameProYearly
		}
		if strings.Contains(planLower, "lifetime") {
			return entities.SubscriptionNameProLifetime
		}
		return entities.SubscriptionNameProMonthly
	}
	if strings.Contains(planLower, "ultra") {
		if strings.Contains(planLower, "yearly") {
			return entities.SubscriptionNameUltraYearly
		}
		return entities.SubscriptionNameUltraMonthly
	}
	if strings.Contains(planLower, "20k") {
		if strings.Contains(planLower, "yearly") {
			return entities.SubscriptionName20KYearly
		}
		return entities.SubscriptionName20KMonthly
	}
	if strings.Contains(planLower, "50k") {
		if strings.Contains(planLower, "yearly") {
			return entities.SubscriptionName50KYearly
		}
		return entities.SubscriptionName50KMonthly
	}
	if strings.Contains(planLower, "100k") {
		if strings.Contains(planLower, "yearly") {
			return entities.SubscriptionName100KYearly
		}
		return entities.SubscriptionName100KMonthly
	}
	if strings.Contains(planLower, "200k") {
		if strings.Contains(planLower, "yearly") {
			return entities.SubscriptionName200KYearly
		}
		return entities.SubscriptionName200KMonthly
	}
	return entities.SubscriptionNameFree
}

// HandlePaymentUpdated processes one-off payment webhooks (Checkout Pro)
func (service *MercadopagoService) HandlePaymentUpdated(ctx context.Context, source string, paymentIDStr string) error {
	ctx, span, ctxLogger := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	if service.paymentClient == nil {
		return stacktrace.NewError("mercadopago payment client is not initialized")
	}

	paymentID, err := strconv.ParseInt(paymentIDStr, 10, 64)
	if err != nil {
		return stacktrace.Propagatef(err, "invalid mercadopago payment id [%s]", paymentIDStr)
	}

	pay, err := service.paymentClient.Get(ctx, int(paymentID))
	if err != nil {
		return stacktrace.Propagatef(err, "cannot fetch mercadopago payment [%d]", paymentID)
	}

	// Only process approved payments for this flow (or refunds if we want to cancel)
	if pay.Status != "approved" {
		ctxLogger.Info(fmt.Sprintf("ignoring mercadopago payment [%d] with status [%s]", paymentID, pay.Status))
		return nil
	}

	// ExternalReference contains userID|planID
	parts := strings.Split(pay.ExternalReference, "|")
	var userIDStr, planID string
	if len(parts) >= 1 {
		userIDStr = parts[0]
	}
	if len(parts) >= 2 {
		planID = parts[1]
	}

	if userIDStr == "" {
		if pay.Payer.Email != "" {
			user, err := service.userRepository.LoadByEmail(ctx, pay.Payer.Email)
			if err != nil {
				return stacktrace.Propagatef(err, "cannot find user for customer email [%s]", pay.Payer.Email)
			}
			userIDStr = string(user.ID)
		} else {
			ctxLogger.Warn(stacktrace.NewError("cannot resolve user from payment [%d]", paymentID))
			return nil
		}
	}

	user, err := service.userRepository.Load(ctx, entities.UserID(userIDStr))
	if err != nil {
		return stacktrace.Propagatef(err, "cannot load user [%s]", userIDStr)
	}

	renewsAt := time.Now().UTC().AddDate(1, 0, 0) // 1 year from now

	payload := &events.UserSubscriptionUpdatedPayload{
		UserID:                user.ID,
		SubscriptionUpdatedAt: time.Now().UTC(),
		SubscriptionID:        strconv.FormatInt(paymentID, 10),
		SubscriptionName:      service.mapPlanToSubscriptionName(planID),
		SubscriptionRenewsAt:  renewsAt,
		SubscriptionStatus:    "authorized",
	}

	event, err := service.createEvent(events.UserSubscriptionUpdated, source, payload)
	if err != nil {
		return stacktrace.Propagatef(err, "cannot create [%s] event for user [%s]", events.UserSubscriptionUpdated, user.ID)
	}
	if err = service.eventDispatcher.Dispatch(ctx, event); err != nil {
		return stacktrace.Propagatef(err, "cannot dispatch [%s] event for user [%s]", events.UserSubscriptionUpdated, user.ID)
	}
	ctxLogger.Info(fmt.Sprintf("mercadopago one-off payment [%d] approved for user [%s]", paymentID, user.ID))

	return nil
}
