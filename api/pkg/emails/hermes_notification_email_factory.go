package emails

import (
	"fmt"
	"strings"
	"time"

	"github.com/NdoleStudio/httpsms/pkg/events"

	"github.com/NdoleStudio/httpsms/pkg/entities"
	"github.com/NdoleStudio/stacktrace"
	"github.com/go-hermes/hermes/v2"
)

type hermesNotificationEmailFactory struct {
	factory
	config    *HermesGeneratorConfig
	generator hermes.Hermes
}

const webhookSendFailedEventPayloadPlaceholder = "__HTTPSMS_WEBHOOK_SEND_FAILED_EVENT_PAYLOAD__"

// NewHermesNotificationEmailFactory creates a new instance of the UserEmailFactory
func NewHermesNotificationEmailFactory(config *HermesGeneratorConfig) NotificationEmailFactory {
	return &hermesNotificationEmailFactory{
		config:    config,
		generator: config.Generator(),
	}
}

func (factory *hermesNotificationEmailFactory) DiscordSendFailed(user *entities.User, payload *events.DiscordSendFailedPayload) (*Email, error) {
	email := hermes.Email{
		Body: hermes.Body{
			Title: "Olá",
			Intros: []string{
				fmt.Sprintf("Encontramos um erro ao encaminhar um SMS recebido para o seu servidor Discord em %s", user.UserTimeString(time.Now())),
			},
			Dictionary: []hermes.Entry{
				{Key: "ID do Canal Discord", Value: payload.DiscordChannelID},
				{Key: "Nome do Evento", Value: payload.EventType},
				{Key: "Número de Telefone", Value: factory.formatPhoneNumber(payload.Owner)},
				{Key: "Código de Resposta HTTP", Value: factory.formatHTTPResponseCode(payload.HTTPResponseStatusCode)},
				{Key: "Mensagem de Erro / Resposta HTTP", Value: payload.ErrorMessage},
			},
			Actions: []hermes.Action{
				{
					Instructions: "Geralmente este erro ocorre porque as permissões do aplicativo httpSMS foram revogadas no seu canal do Discord. Você pode conceder permissão para o httpSMS na página de configurações.",
					Button: hermes.Button{
						Color:     "#329ef4",
						TextColor: "#FFFFFF",
						Text:      "CONFIGURAÇÕES DO DISCORD",
						Link:      fmt.Sprintf("%s/settings/#discord-settings", strings.TrimSuffix(factory.config.AppURL, "/")),
					},
				},
			},
			Signature: "Atenciosamente",
			Outros: []string{
				fmt.Sprintf("Se tiver dúvidas, responda a este e-mail. Você pode desativar esta notificação por e-mail em %s/settings/#email-notifications", strings.TrimSuffix(factory.config.AppURL, "/")),
			},
		},
	}

	html, err := factory.generator.GenerateHTML(email)
	if err != nil {
		return nil, stacktrace.Propagatef(err, "cannot generate html email")
	}

	text, err := factory.generator.GeneratePlainText(email)
	if err != nil {
		return nil, stacktrace.Propagatef(err, "cannot generate text email")
	}

	return &Email{
		ToEmail: user.Email,
		Subject: "📢 Não foi possível encaminhar a mensagem recebida para o seu servidor Discord",
		HTML:    html,
		Text:    text,
	}, nil
}

func (factory *hermesNotificationEmailFactory) WebhookSendFailed(user *entities.User, payload *events.WebhookSendFailedPayload) (*Email, error) {
	formattedPayload, formattedPayloadHTML := formatEventPayload(payload.EventPayload)

	email := hermes.Email{
		Body: hermes.Body{
			Title: "Olá",
			Intros: []string{
				fmt.Sprintf("Encontramos um erro ao encaminhar um evento de webhook do httpSMS para o seu servidor em %s", user.UserTimeString(time.Now())),
			},
			Dictionary: []hermes.Entry{
				{Key: "URL do Servidor", Value: payload.WebhookURL},
				{Key: "Nome do Evento", Value: payload.EventType},
				{Key: "ID do Evento", Value: payload.EventID},
				{Key: "Número de Telefone", Value: factory.formatPhoneNumber(payload.Owner)},
				{Key: "Código de Resposta HTTP", Value: factory.formatHTTPResponseCode(payload.HTTPResponseStatusCode)},
				{Key: "Mensagem de Erro / Resposta HTTP", Value: payload.ErrorMessage},
				{
					Key:         "Payload do Evento",
					Value:       webhookSendFailedEventPayloadPlaceholder,
					UnsafeValue: formattedPayloadHTML,
				},
			},
			Actions: []hermes.Action{
				{
					Instructions: "Geralmente este erro ocorre porque o seu servidor está offline ou inacessível. Você pode reconfigurar o endpoint do webhook na página de configurações do httpSMS.",
					Button: hermes.Button{
						Color:     "#329ef4",
						TextColor: "#FFFFFF",
						Text:      "CONFIGURAÇÕES DO WEBHOOK",
						Link:      fmt.Sprintf("%s/settings/#webhook-settings", strings.TrimSuffix(factory.config.AppURL, "/")),
					},
				},
			},
			Signature: "Atenciosamente",
			Outros: []string{
				fmt.Sprintf("Se tiver dúvidas, responda a este e-mail. Você pode desativar esta notificação por e-mail em %s/settings/#email-notifications", strings.TrimSuffix(factory.config.AppURL, "/")),
			},
		},
	}

	html, err := factory.generator.GenerateHTML(email)
	if err != nil {
		return nil, stacktrace.Propagatef(err, "cannot generate html email")
	}

	eventPayloadEntryIndex := len(email.Body.Dictionary) - 1
	text, err := generateWebhookSendFailedPlainText(func(value string) (string, error) {
		email.Body.Dictionary[eventPayloadEntryIndex].Value = value
		return factory.generator.GeneratePlainText(email)
	}, formattedPayload)
	if err != nil {
		return nil, stacktrace.Propagatef(err, "cannot generate text email")
	}

	return &Email{
		ToEmail: user.Email,
		Subject: "📢 Não foi possível encaminhar um evento de webhook para o seu servidor",
		HTML:    html,
		Text:    text,
	}, nil
}

func generateWebhookSendFailedPlainText(generate func(string) (string, error), formattedPayload string) (string, error) {
	text, err := generate(webhookSendFailedEventPayloadPlaceholder)
	if err != nil {
		return "", err
	}

	if strings.Contains(text, webhookSendFailedEventPayloadPlaceholder) {
		// Hermes/html2text collapses dictionary whitespace, so restore the payload after plain-text generation.
		return replaceWebhookSendFailedEventPayloadPlaceholder(text, formattedPayload), nil
	}

	// Preserve the real payload if a Hermes change stops carrying the placeholder through unchanged.
	return generate(formattedPayload)
}

func replaceWebhookSendFailedEventPayloadPlaceholder(text string, formattedPayload string) string {
	before, after, found := strings.Cut(text, webhookSendFailedEventPayloadPlaceholder)
	if !found {
		return text
	}

	return before + formattedPayload + after
}

func (factory *hermesNotificationEmailFactory) MessageExpired(user *entities.User, payload *events.MessageSendExpiredPayload) (*Email, error) {
	email := hermes.Email{
		Body: hermes.Body{
			Title: "Olá",
			Intros: []string{
				fmt.Sprintf("A mensagem SMS que você enviou para %s expirou em %s e você precisará reenviá-la.", factory.formatPhoneNumber(payload.Contact), user.UserTimeString(time.Now())),
			},
			Dictionary: []hermes.Entry{
				{Key: "ID", Value: payload.MessageID.String()},
				{Key: "De", Value: factory.formatPhoneNumber(payload.Owner)},
				{Key: "Para", Value: factory.formatPhoneNumber(payload.Contact)},
				{Key: "Mensagem", Value: payload.Content},
				{Key: "Criptografado", Value: factory.formatBool(payload.Encrypted)},
			},
			Actions: []hermes.Action{
				{
					Instructions: "As mensagens normalmente expiram porque não conseguimos conectar com o seu celular Android para enviar o SMS. Verifique se o celular está conectado à internet e mantenha-o conectado ao carregador, pois o Android pode encerrar o aplicativo httpSMS se estiver ativo por muito tempo para economizar bateria.",
					Button: hermes.Button{
						Color:     "#329ef4",
						TextColor: "#FFFFFF",
						Text:      "VER MENSAGENS",
						Link:      fmt.Sprintf("%s/threads", strings.TrimSuffix(factory.config.AppURL, "/")),
					},
				},
			},
			Signature: "Atenciosamente",
			Outros: []string{
				fmt.Sprintf("Se tiver dúvidas, responda a este e-mail. Você pode desativar esta notificação por e-mail em %s/settings/#email-notifications", strings.TrimSuffix(factory.config.AppURL, "/")),
			},
		},
	}

	html, err := factory.generator.GenerateHTML(email)
	if err != nil {
		return nil, stacktrace.Propagatef(err, "cannot generate html email")
	}

	text, err := factory.generator.GeneratePlainText(email)
	if err != nil {
		return nil, stacktrace.Propagatef(err, "cannot generate text email")
	}

	return &Email{
		ToEmail: user.Email,
		Subject: "📢 Sua mensagem SMS expirou no httpSMS",
		HTML:    html,
		Text:    text,
	}, nil
}

func (factory *hermesNotificationEmailFactory) MessageFailed(user *entities.User, payload *events.MessageSendFailedPayload) (*Email, error) {
	email := hermes.Email{
		Body: hermes.Body{
			Title: "Olá",
			Intros: []string{
				fmt.Sprintf("A mensagem SMS que você enviou para %s falhou em %s e você precisará reenviá-la.", factory.formatPhoneNumber(payload.Contact), user.UserTimeString(time.Now())),
			},
			Dictionary: []hermes.Entry{
				{Key: "ID", Value: payload.ID.String()},
				{Key: "De", Value: factory.formatPhoneNumber(payload.Owner)},
				{Key: "Para", Value: factory.formatPhoneNumber(payload.Contact)},
				{Key: "Mensagem", Value: payload.Content},
				{Key: "Criptografado", Value: factory.formatBool(payload.Encrypted)},
				{Key: "Motivo da Falha", Value: payload.ErrorMessage},
			},
			Actions: []hermes.Action{
				{
					Instructions: "Verifique o aplicativo padrão de mensagens SMS do seu celular para entender o motivo exato da falha. Geralmente o envio falha se o app httpSMS foi desinstalado ou está inativo. Faça logout e login novamente no aplicativo Android e tente reenviar o SMS.",
					Button: hermes.Button{
						Color:     "#329ef4",
						TextColor: "#FFFFFF",
						Text:      "VER MENSAGENS",
						Link:      fmt.Sprintf("%s/threads", strings.TrimSuffix(factory.config.AppURL, "/")),
					},
				},
			},
			Signature: "Atenciosamente",
			Outros: []string{
				fmt.Sprintf("Se tiver dúvidas, responda a este e-mail. Você pode desativar esta notificação por e-mail em %s/settings/#email-notifications", strings.TrimSuffix(factory.config.AppURL, "/")),
			},
		},
	}

	html, err := factory.generator.GenerateHTML(email)
	if err != nil {
		return nil, stacktrace.Propagatef(err, "cannot generate html email")
	}

	text, err := factory.generator.GeneratePlainText(email)
	if err != nil {
		return nil, stacktrace.Propagatef(err, "cannot generate text email")
	}

	return &Email{
		ToEmail: user.Email,
		Subject: "📢 Seu envio de mensagem SMS falhou no httpSMS",
		HTML:    html,
		Text:    text,
	}, nil
}
