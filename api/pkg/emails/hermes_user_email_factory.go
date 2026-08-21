package emails

import (
	"fmt"
	"strings"
	"time"

	"github.com/NdoleStudio/httpsms/pkg/entities"
	"github.com/NdoleStudio/stacktrace"
	"github.com/go-hermes/hermes/v2"
)

type hermesUserEmailFactory struct {
	factory
	config    *HermesGeneratorConfig
	generator hermes.Hermes
}

// formatBillingDate renders a date like "19 June 2026" in the user's timezone.
func formatBillingDate(t time.Time, location *time.Location) string {
	return t.In(location).Format("2 January 2006")
}

func (factory *hermesUserEmailFactory) APIKeyRotated(emailAddress string, timestamp time.Time, timezone string) (*Email, error) {
	location, err := time.LoadLocation(timezone)
	if err != nil {
		location = time.UTC
	}

	email := hermes.Email{
		Body: hermes.Body{
			Intros: []string{
				fmt.Sprintf("Este é um e-mail de confirmação informando que sua Chave de API do httpSMS foi rotacionada com sucesso em %s.", timestamp.In(location).Format(time.RFC1123)),
			},
			Actions: []hermes.Action{
				{
					Instructions: "Você pode visualizar sua nova chave de API na página de configurações do httpSMS.",
					Button: hermes.Button{
						Color:     "#329ef4",
						TextColor: "#FFFFFF",
						Text:      "Configurações do httpSMS",
						Link:      fmt.Sprintf("%s/settings/", strings.TrimSuffix(factory.config.AppURL, "/")),
					},
				},
			},
			Title:     "Olá,",
			Signature: "Atenciosamente",
			Outros: []string{
				"Se você não solicitou a rotação desta chave de API, entre em contato conosco imediatamente respondendo a este e-mail.",
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
		ToEmail: emailAddress,
		Subject: "Sua Chave de API do httpSMS foi rotacionada com sucesso",
		HTML:    html,
		Text:    text,
	}, nil
}

// UsageLimitExceeded is the email sent when the plan limit is reached
func (factory *hermesUserEmailFactory) UsageLimitExceeded(user *entities.User, usage *entities.BillingUsage) (*Email, error) {
	email := hermes.Email{
		Body: hermes.Body{
			Intros: []string{
				fmt.Sprintf("Você atingiu o seu limite de %s mensagens no plano %s. Por isso, novas mensagens não serão processadas até que seu uso seja reiniciado.", factory.formatQuantity(user.SubscriptionName.Limit()), user.SubscriptionName),
				fmt.Sprintf("Entre %s e %s você enviou %s mensagens e recebeu %s, totalizando %s.", formatBillingDate(usage.StartTimestamp, user.Location()), formatBillingDate(usage.EndTimestamp, user.Location()), factory.formatQuantity(usage.SentMessages), factory.formatQuantity(usage.ReceivedMessages), factory.formatQuantity(usage.TotalMessages())),
			},
			Actions: []hermes.Action{
				{
					Instructions: "Clique no botão abaixo para atualizar seu plano e continuar enviando mensagens.",
					Button: hermes.Button{
						Color:     "#329ef4",
						TextColor: "#FFFFFF",
						Text:      "FAZER UPGRADE DO PLANO",
						Link:      fmt.Sprintf("%s/billing", strings.TrimSuffix(factory.config.AppURL, "/")),
					},
				},
			},
			Title:     "Olá,",
			Signature: "Atenciosamente",
			Outros: []string{
				"Se tiver qualquer dúvida, responda a este e-mail.",
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
		Subject: "⚠️ Você excedeu o limite do seu plano",
		HTML:    html,
		Text:    text,
	}, nil
}

// UsageLimitAlert is the email sent when the plan limit is reached
func (factory *hermesUserEmailFactory) UsageLimitAlert(user *entities.User, usage *entities.BillingUsage) (*Email, error) {
	percent := (usage.TotalMessages() * 100) / user.SubscriptionName.Limit()
	email := hermes.Email{
		Body: hermes.Body{
			Intros: []string{
				fmt.Sprintf("Aviso importante: você já utilizou %d%% do seu limite mensal de SMS no plano %s.", percent, user.SubscriptionName),
				fmt.Sprintf("Entre %s e %s você enviou %s mensagens e recebeu %s, totalizando %s do limite de %s mensagens.", formatBillingDate(usage.StartTimestamp, user.Location()), formatBillingDate(usage.EndTimestamp, user.Location()), factory.formatQuantity(usage.SentMessages), factory.formatQuantity(usage.ReceivedMessages), factory.formatQuantity(usage.TotalMessages()), factory.formatQuantity(user.SubscriptionName.Limit())),
			},
			Actions: []hermes.Action{
				{
					Instructions: "Clique no botão abaixo para atualizar seu plano e evitar interrupções no envio de mensagens.",
					Button: hermes.Button{
						Color:     "#329ef4",
						TextColor: "#FFFFFF",
						Text:      "FAZER UPGRADE DO PLANO",
						Link:      fmt.Sprintf("%s/billing", strings.TrimSuffix(factory.config.AppURL, "/")),
					},
				},
			},
			Title:     "Olá,",
			Signature: "Atenciosamente",
			Outros: []string{
				"Se tiver qualquer dúvida, responda a este e-mail.",
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
		Subject: fmt.Sprintf("⚠️ Alerta de Uso: %d%% do Limite Atingido", percent),
		HTML:    html,
		Text:    text,
	}, nil
}

// NewHermesUserEmailFactory creates a new instance of the UserEmailFactory
func NewHermesUserEmailFactory(config *HermesGeneratorConfig) UserEmailFactory {
	return &hermesUserEmailFactory{
		config:    config,
		generator: config.Generator(),
	}
}

// PhoneDead is the email sent to a user when their phone is dead
func (factory *hermesUserEmailFactory) PhoneDead(user *entities.User, lastHeartbeatTimestamp time.Time, owner string) (*Email, error) {
	location, err := time.LoadLocation(user.Timezone)
	if err != nil {
		location = time.UTC
	}

	email := hermes.Email{
		Body: hermes.Body{
			Intros: []string{
				fmt.Sprintf("Não recebemos nenhum sinal (heartbeat) do celular Android %s desde %s.", factory.formatPhoneNumber(owner), lastHeartbeatTimestamp.In(location).Format(time.RFC1123)),
				"Verifique se o aparelho celular está ligado e com uma conexão estável à internet.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Verifique os eventos de heartbeat no httpSMS",
					Button: hermes.Button{
						Color:     "#329ef4",
						TextColor: "#FFFFFF",
						Text:      "HEARTBEATS",
						Link:      fmt.Sprintf("%s/heartbeats/%s", strings.TrimSuffix(factory.config.AppURL, "/"), owner),
					},
				},
			},
			Title:     "Olá,",
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
		Subject: fmt.Sprintf("⚠️ Sem sinal do celular Android [%s]", factory.formatPhoneNumber(owner)),
		HTML:    html,
		Text:    text,
	}, nil
}
