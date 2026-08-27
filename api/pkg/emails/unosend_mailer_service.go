package emails

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/NdoleStudio/httpsms/pkg/telemetry"
	"github.com/NdoleStudio/stacktrace"
)

// UnosendConfig is the config for setting up the unosendMailer
type UnosendConfig struct {
	FromName  string
	FromEmail string
	APIKey    string
	BaseURL   string
}

type unosendMailer struct {
	from       string
	apiKey     string
	baseURL    string
	tracer     telemetry.Tracer
	httpClient *http.Client
}

// NewUnosendEmailService creates a new instance of the unosendMailer
func NewUnosendEmailService(tracer telemetry.Tracer, httpClient *http.Client, config UnosendConfig) Mailer {
	return &unosendMailer{
		tracer:     tracer,
		httpClient: httpClient,
		apiKey:     config.APIKey,
		baseURL:    config.BaseURL,
		from:       fmt.Sprintf("%s <%s>", config.FromName, config.FromEmail),
	}
}

// unosendEmailRequest is the payload sent to the Unosend API
type unosendEmailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html,omitempty"`
	Text    string   `json:"text,omitempty"`
}

// Send a new email
func (mailer *unosendMailer) Send(ctx context.Context, email *Email) (err error) {
	ctx, span := mailer.tracer.Start(ctx)
	defer span.End()

	payload := &unosendEmailRequest{
		From:    mailer.from,
		To:      []string{email.toAddress()},
		Subject: email.Subject,
		HTML:    email.HTML,
		Text:    email.Text,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return stacktrace.Propagatef(err, "cannot marshal Unosend email request")
	}

	url := fmt.Sprintf("%s/emails", mailer.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return stacktrace.Propagatef(err, "cannot create HTTP request for Unosend")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", mailer.apiKey))

	resp, err := mailer.httpClient.Do(req)
	if err != nil {
		return stacktrace.Propagatef(err, "cannot send email to Unosend at [%s]", url)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return stacktrace.Propagatef(err, "cannot read Unosend response body")
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return stacktrace.NewErrorf("Unosend returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
