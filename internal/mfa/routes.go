// internal/mfa/routes.go – demo de MFA basico (OTP secreto fijo)
package mfa

import (
	"crypto/rand"
	"encoding/base32"

	"login-service/internal/common"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var secretStore = make(map[uint]string)

// generateSecret crea un secreto base32
func generateSecret() string {
	b := make([]byte, 10)
	rand.Read(b)
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b)
}

func Register(app *fiber.App, cfg *common.Config, db *gorm.DB) {
	if !cfg.Auth.MFAEnabled {
		return
	}

	// Setup MFA → devuelve secreto al usuario
	app.Post("/mfa/setup", func(c *fiber.Ctx) error {
		var body struct{ UserID uint }
		c.BodyParser(&body)
		sec := generateSecret()
		secretStore[body.UserID] = sec
		return c.JSON(fiber.Map{"secret": sec})
	})

	// Verificar código (aquí usamos el propio secreto como OTP demo)
	app.Post("/mfa/verify", func(c *fiber.Ctx) error {
		var body struct {
			UserID uint
			Code   string
		}
		c.BodyParser(&body)

		if secretStore[body.UserID] == body.Code {
			return c.JSON(fiber.Map{"verified": true})
		}
		return c.Status(401).JSON(fiber.Map{"verified": false})
	})
}
