package services

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/NdoleStudio/httpsms/pkg/entities"
	"github.com/NdoleStudio/httpsms/pkg/events"
	"github.com/NdoleStudio/httpsms/pkg/repositories"
	"github.com/NdoleStudio/httpsms/pkg/telemetry"
	"github.com/NdoleStudio/stacktrace"
	"github.com/stripe/stripe-go/v78"
	billingportalSession "github.com/stripe/stripe-go/v78/billingportal/session"
	checkoutSession "github.com/stripe/stripe-go/v78/checkout/session"
	"github.com/stripe/stripe-go/v78/customer"
	"github.com/stripe/stripe-go/v78/subscription"
)

// StripeService is responsible for handling Stripe payments and webhooks
type StripeService struct {
	service
	logger          telemetry.Logger
	tracer          telemetry.Tracer
	eventDispatcher *EventDispatcher
	userRepository  repositories.UserRepository
}

// NewStripeService creates a new instance of StripeService
func NewStripeService(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	userRepository repositories.UserRepository,
	eventDispatcher *EventDispatcher,
) *StripeService {
	apiKey := os.Getenv("STRIPE_SECRET_KEY")
	if apiKey != "" {
		stripe.Key = apiKey
	}
	return &StripeService{
		logger:          logger.WithService(fmt.Sprintf("%T", &StripeService{})),
		tracer:          tracer,
		userRepository:  userRepository,
		eventDispatcher: eventDispatcher,
	}
}

// CreateCheckoutSessionParams holds params for creating a Stripe checkout session
type CreateCheckoutSessionParams struct {
	UserID   entities.UserID
	PlanID   string
	PriceID  string
	SuccessURL string
	CancelURL  string
}

// CreateCheckoutSession creates a new Stripe Checkout Session for a user
func (service *StripeService) CreateCheckoutSession(ctx context.Context, params CreateCheckoutSessionParams) (string, error) {
	ctx, span, _ := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

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
	cancelURL := params.CancelURL
	if cancelURL == "" {
		cancelURL = appURL + "/billing?status=cancelled"
	}

	// Resolve Price ID from plan if not passed directly
	priceID := params.PriceID
	if priceID == "" {
		priceID = os.Getenv("STRIPE_PRICE_" + strings.ToUpper(strings.ReplaceAll(params.PlanID, "-", "_")))
	}

	sessionParams := &stripe.CheckoutSessionParams{
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		CustomerEmail: stripe.String(user.Email),
		ClientReferenceID: stripe.String(user.ID.String()),
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
		Metadata: map[string]string{
			"user_id": user.ID.String(),
			"plan_id": params.PlanID,
		},
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			Metadata: map[string]string{
				"user_id": user.ID.String(),
				"plan_id": params.PlanID,
			},
		},
	}

	if priceID == "" {
		return "", stacktrace.NewError("no stripe price ID configured for plan [%s]", params.PlanID)
	}

	sessionParams.LineItems = []*stripe.CheckoutSessionLineItemParams{
		{
			Price:    stripe.String(priceID),
			Quantity: stripe.Int64(1),
		},
	}

	sess, err := checkoutSession.New(sessionParams)
	if err != nil {
		return "", stacktrace.Propagatef(err, "cannot create stripe checkout session for user [%s]", params.UserID)
	}

	return sess.URL, nil
}

// CreateCustomerPortalSession creates a Stripe Customer Portal session URL for managing subscriptions
func (service *StripeService) CreateCustomerPortalSession(ctx context.Context, userID entities.UserID, returnURL string) (string, error) {
	ctx, span, _ := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	user, err := service.userRepository.Load(ctx, userID)
	if err != nil {
		return "", stacktrace.Propagatef(err, "cannot load user [%s]", userID)
	}

	var stripeCustomerID string
	if user.SubscriptionID != nil && strings.HasPrefix(*user.SubscriptionID, "sub_") {
		sub, err := subscription.Get(*user.SubscriptionID, nil)
		if err == nil && sub.Customer != nil {
			stripeCustomerID = sub.Customer.ID
		}
	}

	if stripeCustomerID == "" {
		// Search Stripe customer by email
		iter := customer.List(&stripe.CustomerListParams{
			Email: stripe.String(user.Email),
		})
		for iter.Next() {
			c := iter.Customer()
			if c != nil {
				stripeCustomerID = c.ID
				break
			}
		}
	}

	if stripeCustomerID == "" {
		return "", stacktrace.NewErrorf("stripe customer not found for user [%s]", userID)
	}

	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		appURL = "https://smsandroid.com.br"
	}
	if returnURL == "" {
		returnURL = appURL + "/billing"
	}

	portalParams := &stripe.BillingPortalSessionParams{
		Customer:  stripe.String(stripeCustomerID),
		ReturnURL: stripe.String(returnURL),
	}

	sess, err := billingportalSession.New(portalParams)
	if err != nil {
		return "", stacktrace.Propagatef(err, "cannot create billing portal session for customer [%s]", stripeCustomerID)
	}

	return sess.URL, nil
}

// HandleCheckoutSessionCompleted processes checkout.session.completed Stripe webhooks
func (service *StripeService) HandleCheckoutSessionCompleted(ctx context.Context, source string, sess *stripe.CheckoutSession) error {
	ctx, span, ctxLogger := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	userIDStr := sess.ClientReferenceID
	if userIDStr == "" && sess.Metadata != nil {
		userIDStr = sess.Metadata["user_id"]
	}

	if userIDStr == "" {
		user, err := service.userRepository.LoadByEmail(ctx, sess.CustomerEmail)
		if err != nil {
			return stacktrace.Propagatef(err, "cannot find user for customer email [%s]", sess.CustomerEmail)
		}
		userIDStr = user.ID.String()
	}

	userID := entities.UserID(userIDStr)
	subID := ""
	if sess.Subscription != nil {
		subID = sess.Subscription.ID
	}

	planID := ""
	if sess.Metadata != nil {
		planID = sess.Metadata["plan_id"]
	}
	if planID == "" {
		planID = "pro-monthly"
	}

	renewsAt := time.Now().AddDate(0, 1, 0)
	payload := &events.UserSubscriptionCreatedPayload{
		UserID:                userID,
		SubscriptionCreatedAt: time.Now().UTC(),
		SubscriptionID:        subID,
		SubscriptionName:      service.mapPlanToSubscriptionName(planID),
		SubscriptionRenewsAt:  renewsAt,
		SubscriptionStatus:    "active",
	}

	event, err := service.createEvent(events.UserSubscriptionCreated, source, payload)
	if err != nil {
		return stacktrace.Propagatef(err, "cannot create [%s] event for user [%s]", events.UserSubscriptionCreated, userID)
	}

	if err = service.eventDispatcher.Dispatch(ctx, event); err != nil {
		return stacktrace.Propagatef(err, "cannot dispatch [%s] event for user [%s]", events.UserSubscriptionCreated, userID)
	}

	ctxLogger.Info(fmt.Sprintf("stripe checkout session completed: created subscription [%s] for user [%s]", subID, userID))
	return nil
}

// HandleSubscriptionUpdated processes customer.subscription.updated Stripe webhooks
func (service *StripeService) HandleSubscriptionUpdated(ctx context.Context, source string, sub *stripe.Subscription) error {
	ctx, span, ctxLogger := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	user, err := service.userRepository.LoadBySubscriptionID(ctx, sub.ID)
	if err != nil {
		// Fallback: search by customer email if metadata contains user_id
		userIDStr := ""
		if sub.Metadata != nil {
			userIDStr = sub.Metadata["user_id"]
		}
		if userIDStr != "" {
			user, err = service.userRepository.Load(ctx, entities.UserID(userIDStr))
		}
	}
	if err != nil {
		return stacktrace.Propagatef(err, "cannot load user for subscription [%s]", sub.ID)
	}

	renewsAt := time.Unix(sub.CurrentPeriodEnd, 0)
	status := string(sub.Status)
	planID := ""
	if sub.Metadata != nil {
		planID = sub.Metadata["plan_id"]
	}

	payload := &events.UserSubscriptionUpdatedPayload{
		UserID:                user.ID,
		SubscriptionUpdatedAt: time.Now().UTC(),
		SubscriptionID:        sub.ID,
		SubscriptionName:      service.mapPlanToSubscriptionName(planID),
		SubscriptionRenewsAt:  renewsAt,
		SubscriptionStatus:    status,
	}

	event, err := service.createEvent(events.UserSubscriptionUpdated, source, payload)
	if err != nil {
		return stacktrace.Propagatef(err, "cannot create [%s] event for user [%s]", events.UserSubscriptionUpdated, user.ID)
	}

	if err = service.eventDispatcher.Dispatch(ctx, event); err != nil {
		return stacktrace.Propagatef(err, "cannot dispatch [%s] event for user [%s]", events.UserSubscriptionUpdated, user.ID)
	}

	ctxLogger.Info(fmt.Sprintf("stripe subscription [%s] updated for user [%s]", sub.ID, user.ID))
	return nil
}

// HandleSubscriptionDeleted processes customer.subscription.deleted Stripe webhooks
func (service *StripeService) HandleSubscriptionDeleted(ctx context.Context, source string, sub *stripe.Subscription) error {
	ctx, span, ctxLogger := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	user, err := service.userRepository.LoadBySubscriptionID(ctx, sub.ID)
	if err != nil {
		return stacktrace.Propagatef(err, "cannot load user for cancelled subscription [%s]", sub.ID)
	}

	endsAt := time.Unix(sub.CanceledAt, 0)
	payload := &events.UserSubscriptionCancelledPayload{
		UserID:                  user.ID,
		SubscriptionCancelledAt: time.Now().UTC(),
		SubscriptionID:          sub.ID,
		SubscriptionName:        entities.SubscriptionNameFree,
		SubscriptionEndsAt:      endsAt,
		SubscriptionStatus:      "canceled",
	}

	event, err := service.createEvent(events.UserSubscriptionCancelled, source, payload)
	if err != nil {
		return stacktrace.Propagatef(err, "cannot create [%s] event for user [%s]", events.UserSubscriptionCancelled, user.ID)
	}

	if err = service.eventDispatcher.Dispatch(ctx, event); err != nil {
		return stacktrace.Propagatef(err, "cannot dispatch [%s] event for user [%s]", events.UserSubscriptionCancelled, user.ID)
	}

	ctxLogger.Info(fmt.Sprintf("stripe subscription [%s] deleted for user [%s]", sub.ID, user.ID))
	return nil
}

func (service *StripeService) mapPlanToSubscriptionName(plan string) entities.SubscriptionName {
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
		return entities.SubscriptionName20KMonthly
	}
	if strings.Contains(planLower, "50k") {
		return entities.SubscriptionName50KMonthly
	}
	if strings.Contains(planLower, "100k") {
		return entities.SubscriptionName100KMonthly
	}
	if strings.Contains(planLower, "200k") {
		return entities.SubscriptionName200KMonthly
	}
	return entities.SubscriptionNameFree
}
