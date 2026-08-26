package emails

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-hermes/hermes/v2"
)

// HermesGeneratorConfig contains details for the generator
type HermesGeneratorConfig struct {
	AppURL     string
	AppName    string
	AppLogoURL string
}

// Generator creates hermes.Hermes from HermesGeneratorConfig
func (config *HermesGeneratorConfig) Generator() hermes.Hermes {
	return hermes.Hermes{
		Theme: newHermesTheme(),
		Product: hermes.Product{
			// Appears in header & footer of e-mails
			Name: fmt.Sprintf("Equipe %s", config.AppName),
			Link: config.AppURL,
			// Optional product logo
			Copyright:   fmt.Sprintf("© %s %s. Todos os direitos reservados.", strconv.Itoa(time.Now().Year()), config.AppName),
			Logo:        config.AppLogoURL,
			TroubleText: "Se você estiver tendo problemas para clicar no botão '{ACTION}', copie e cole a URL abaixo no seu navegador:",
		},
	}
}
