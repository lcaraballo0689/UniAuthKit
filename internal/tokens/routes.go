// internal/tokens/routes.go – endpoint para revocar refresh tokens
package tokens

import (
	"login-service/internal/common"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Register(app *fiber.App, cfg *common.Config, db *gorm.DB) {
	app.Post("/token/revoke", func(c *fiber.Ctx) error {
		var body struct{ Token string }
		if err := c.BodyParser(&body); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "bad_request"})
		}
		Blacklist(body.Token)
		return c.JSON(fiber.Map{"status": "revoked"})
	})
}
