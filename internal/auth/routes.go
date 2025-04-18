// internal/auth/routes.go
//  - Endpoints /login y /register
//  - Fallback a "superadmin" + full‑access si no hay permisos configurados
package auth

import (
	"time"

	"login-service/internal/common"
	"login-service/internal/tokens"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"uniqueIndex"`
	Password  string
	Email     string
	Role      string
	LastLogin time.Time
}

// Register registra las rutas de autenticación
func Register(app *fiber.App, cfg *common.Config, db *gorm.DB) {
	// Auto-migrar tabla de usuarios (demo)
	db.AutoMigrate(&User{})

	// ---------- /register ----------
	if cfg.Endpoints["register"] {
		app.Post("/register", func(c *fiber.Ctx) error {
			var body struct {
				Username string
				Password string
				Email    string
				Role     string
			}
			if err := c.BodyParser(&body); err != nil {
				return c.Status(400).JSON(fiber.Map{"error": "bad_request"})
			}
			hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 14)
			u := User{Username: body.Username, Password: string(hash), Email: body.Email, Role: body.Role}
			if err := db.Create(&u).Error; err != nil {
				return c.Status(409).JSON(fiber.Map{"error": "user_exists"})
			}
			return c.JSON(fiber.Map{"message": "registered"})
		})
	}

	// ---------- /login ----------
	if cfg.Endpoints["login"] {
		app.Post("/login", func(c *fiber.Ctx) error {
			var body struct{ Username, Password string }
			if err := c.BodyParser(&body); err != nil {
				return c.Status(400).JSON(fiber.Map{"error": "bad_request"})
			}

			var u User
			db.Where("username = ?", body.Username).First(&u)
			if u.ID == 0 || bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(body.Password)) != nil {
				return c.Status(401).JSON(fiber.Map{"error": "invalid_credentials"})
			}

			// ─── Obtener permisos vía query parametrizada ───
			perms := []string{}
			if q, ok := cfg.Queries["permisos_usuario"]; ok && q != "" {
				rows, _ := db.Raw(q, u.ID).Rows()
				defer rows.Close()
				for rows.Next() {
					var p string
					rows.Scan(&p)
					perms = append(perms, p)
				}
			}

			// ─── Fallback ► superadmin + full‑access ───
			if len(perms) == 0 {
				u.Role = "superadmin"
				perms = []string{"*"} // comodín full‑access
			}

			// Actualizar last_login y generar tokens
			u.LastLogin = time.Now()
			db.Save(&u)

			access, refresh := tokens.GeneratePair(cfg, u.ID, u.Role, perms)
			return c.JSON(fiber.Map{
				"access_token":  access,
				"refresh_token": refresh,
				"expires_in":    cfg.Auth.TokenDurationMinutes * 60,
				"user": fiber.Map{
					"id":          u.ID,
					"username":    u.Username,
					"email":       u.Email,
					"role":        u.Role,
					"permissions": perms,
				},
			})
		})
	}
}
