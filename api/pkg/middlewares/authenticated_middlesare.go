package middlewares

import (
	"fmt"
	"os"

	"github.com/NdoleStudio/httpsms/pkg/entities"
	"github.com/NdoleStudio/httpsms/pkg/repositories"
	"github.com/NdoleStudio/httpsms/pkg/telemetry"
	"github.com/gofiber/fiber/v3"
)

const (
	authHeaderBearer = "Authorization"
	authHeaderAPIKey = "x-api-key"
	bearerScheme     = "Bearer"
)

const (
	// ContextKeyAuthUserID is the context key used to store the ID of an authenticated user
	ContextKeyAuthUserID = "auth.user.id"
)

// Authenticated checks if the request is authenticated and that the user is not blocked
func Authenticated(tracer telemetry.Tracer, logger telemetry.Logger, userRepo repositories.UserRepository) fiber.Handler {
	logger = logger.WithService("middlewares.Authenticated")

	return func(c fiber.Ctx) error {
		_, span := tracer.StartFromFiberCtx(c, "middlewares.Authenticated")
		defer span.End()

		if tokenUser, ok := c.Locals(ContextKeyAuthUserID).(entities.AuthContext); !ok || tokenUser.IsNoop() {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "You are not authorized to carry out this request.",
				"data":    "Make sure your API key is set in the [x-api-key] header in the request",
			})
		}

		authCtx, _ := c.Locals(ContextKeyAuthUserID).(entities.AuthContext)
		user, err := userRepo.Load(c.Context(), authCtx.ID)
		if err != nil {
			logger.Error(fmt.Errorf("failed to load user [%s] for auth check: %w", authCtx.ID, err))
			return c.Next()
		}

		if !user.Active && !user.IsAdmin && string(user.ID) != os.Getenv("EVENTS_QUEUE_USER_ID") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "error",
				"message": "Account blocked",
				"data":    "Your account has been blocked. Contact the administrator.",
			})
		}

		return c.Next()
	}
}
