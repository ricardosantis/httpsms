package validators

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/NdoleStudio/httpsms/pkg/telemetry"
	"github.com/stripe/stripe-go/v78/webhook"
)

// StripeHandlerValidator validates Stripe webhook requests
type StripeHandlerValidator struct {
	logger telemetry.Logger
	tracer telemetry.Tracer
}

// NewStripeHandlerValidator creates a new StripeHandlerValidator
func NewStripeHandlerValidator(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
) *StripeHandlerValidator {
	return &StripeHandlerValidator{
		logger: logger.WithService(fmt.Sprintf("%T", &StripeHandlerValidator{})),
		tracer: tracer,
	}
}

// ValidateEvent checks that an event payload is signed by Stripe
func (validator *StripeHandlerValidator) ValidateEvent(ctx context.Context, signature string, payload []byte) url.Values {
	_, span := validator.tracer.Start(ctx)
	defer span.End()

	secret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if secret == "" {
		return url.Values{}
	}

	_, err := webhook.ConstructEvent(payload, signature, secret)
	if err != nil {
		return url.Values{
			"signature": []string{
				fmt.Sprintf("Invalid Stripe signature: %v", err),
			},
		}
	}

	return url.Values{}
}
