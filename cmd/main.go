// main.go – arranque del microservicio
package main

import (
	"log"

	"login-service/internal/common"
	"login-service/internal/auth"
	"login-service/internal/mfa"
	"login-service/internal/oauth"
	"login-service/internal/tokens"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// 1) Cargar configuración desde YAML
	cfg := common.LoadConfig("config/config.yaml")

	// 2) Conectar a la base de datos elegida
	db := common.InitDatabase(&cfg)

	// 3) Crear instancia Fiber (HTTP)
	app := fiber.New()
	app.Use(logger.New())

	// 4) Registrar módulos
	auth.Register(app, &cfg, db)
	mfa.Register(app, &cfg, db)
	oauth.Register(app, &cfg, db)
	tokens.Register(app, &cfg, db)

	// 5) Lanzar servidor
	log.Fatal(app.Listen(":" + cfg.Server.Port))
}
