package middlewares

import (
	"fmt"

	"github.com/NdoleStudio/httpsms/pkg/entities"
	"github.com/NdoleStudio/httpsms/pkg/repositories"
	"github.com/NdoleStudio/httpsms/pkg/telemetry"
	"github.com/gofiber/fiber/v3"
)

// AdminOnly creates a middleware that restricts access to users with IsAdmin == true
func AdminOnly(logger telemetry.Logger, userRepo repositories.UserRepository) fiber.Handler {
	return func(c fiber.Ctx) error {
		authCtx, ok := c.Locals(ContextKeyAuthUserID).(entities.AuthContext)
		if !ok || authCtx.IsNoop() {
			logger.Error(fmt.Errorf("AdminOnly middleware called without valid auth context"))
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		user, err := userRepo.Load(c.Context(), authCtx.ID)
		if err != nil {
			logger.Error(fmt.Errorf("AdminOnly failed to load user [%s]: %w", authCtx.ID, err))
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if !user.IsAdmin {
			logger.Error(fmt.Errorf("Access denied for non-admin user [%s]", authCtx.ID))
			return c.SendStatus(fiber.StatusForbidden)
		}

		// store the loaded user in context for downstream handlers
		c.Locals("user", user)

		return c.Next()
	}
}
