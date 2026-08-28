package middlewares

import (
	"fmt"
	"github.com/NdoleStudio/httpsms/pkg/telemetry"
	"github.com/gofiber/fiber/v3"
	"github.com/NdoleStudio/httpsms/pkg/entities"
)

// AdminOnly creates a middleware that restricts access to users with IsAdmin == true
func AdminOnly(logger telemetry.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		userInter := c.Locals("user")
		if userInter == nil {
			logger.Error(fmt.Errorf("AdminOnly middleware called without user in context"))
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		
		user, ok := userInter.(*entities.User)
		if !ok || !user.IsAdmin {
			logger.Error(fmt.Errorf("Access denied for non-admin user"))
			return c.SendStatus(fiber.StatusForbidden)
		}

		return c.Next()
	}
}
