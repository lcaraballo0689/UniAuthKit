// internal/oauth/routes.go – Stub de OAuth Google
package oauth

import (
	"fmt"

	"login-service/internal/common"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Register(app *fiber.App, cfg *common.Config, db *gorm.DB) {
	if !cfg.Endpoints["oauth"] || !cfg.OAuth["enabled"].(bool) {
		return
	}
	g := cfg.OAuth["providers"].(map[string]any)["google"].(map[string]any)
	clientID := g["client_id"].(string)
	redirect := g["redirect_uri"].(string)

	// Paso 1: redirigir al login de Google
	app.Get("/oauth/google/login", func(c *fiber.Ctx) error {
		url := fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&scope=email&response_type=code&redirect_uri=%s", clientID, redirect)
		return c.Redirect(url)
	})

	// Paso 2: callback (solo stub)
	app.Get("/oauth/google/callback", func(c *fiber.Ctx) error {
		code := c.Query("code")
		return c.SendString("Recibido code=" + code + " (stub, intercambiar por token)")
	})
}
