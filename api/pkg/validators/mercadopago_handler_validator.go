package validators

import (
	"context"
	"net/url"
	"os"

	"github.com/NdoleStudio/httpsms/pkg/telemetry"
	"github.com/mercadopago/sdk-go/pkg/webhook"
)

// MercadopagoHandlerValidator is used to validate requests on MercadopagoHandler
type MercadopagoHandlerValidator struct {
	validator
	logger telemetry.Logger
	tracer telemetry.Tracer
}

// NewMercadopagoHandlerValidator creates a new instance of MercadopagoHandlerValidator
func NewMercadopagoHandlerValidator(logger telemetry.Logger, tracer telemetry.Tracer) *MercadopagoHandlerValidator {
	return &MercadopagoHandlerValidator{
		logger: logger,
		tracer: tracer,
	}
}

// ValidateEvent validates incoming webhook from Mercado Pago
func (validator *MercadopagoHandlerValidator) ValidateEvent(ctx context.Context, xSignature, xRequestID, dataID string) url.Values {
	_, span, _ := validator.tracer.StartWithLogger(ctx, validator.logger)
	defer span.End()

	secret := os.Getenv("MERCADOPAGO_WEBHOOK_SECRET")
	if secret == "" {
		// If no secret is configured, bypass validation
		return nil
	}

	err := webhook.ValidateSignature(xSignature, xRequestID, dataID, secret)
	if err != nil {
		errors := url.Values{}
		errors.Add("signature", "invalid signature")
		return errors
	}

	return nil
}
