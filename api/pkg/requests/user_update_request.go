package requests

import (
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/NdoleStudio/httpsms/pkg/services"
)

// UserUpdate is the payload for updating a phone
type UserUpdate struct {
	request
	Timezone      string `json:"timezone" example:"Europe/Helsinki"`
	Locale        string `json:"locale" example:"pt-BR"`
	ActivePhoneID string `json:"active_phone_id" example:"32343a19-da5e-4b1b-a767-3298a73703cb"`
}

// Sanitize sets defaults to MessageOutstanding
func (input *UserUpdate) Sanitize() UserUpdate {
	input.ActivePhoneID = strings.TrimSpace(input.ActivePhoneID)
	input.Timezone = strings.TrimSpace(input.Timezone)
	input.Locale = strings.TrimSpace(input.Locale)
	return *input
}

// ToUpdateParams converts UserUpdate to services.UserUpdateParams
func (input *UserUpdate) ToUpdateParams() services.UserUpdateParams {
	location, err := time.LoadLocation(input.Timezone)
	if err != nil {
		location = time.UTC
	}

	var activePhoneID *uuid.UUID
	if input.ActivePhoneID != "" {
		val := uuid.MustParse(input.ActivePhoneID)
		activePhoneID = &val
	}

	return services.UserUpdateParams{
		ActivePhoneID: activePhoneID,
		Timezone:      location,
		Locale:        input.Locale,
	}
}
