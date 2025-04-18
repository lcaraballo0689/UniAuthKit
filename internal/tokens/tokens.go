package tokens

import (
	"time"

	"login-service/internal/common"

	"github.com/golang-jwt/jwt/v4"
)

var blacklist = map[string]bool{}

// GeneratePair → ahora recibe userID y role como valores simples.
func GeneratePair(cfg *common.Config, userID uint, role string, perms []string) (string, string) {
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":     userID,
		"role":        role,
		"permissions": perms,
		"exp":         time.Now().Add(time.Minute * time.Duration(cfg.Auth.TokenDurationMinutes)).Unix(),
	})
	accessStr, _ := access.SignedString([]byte(cfg.Auth.JWTSecret))

	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"type":    "refresh",
		"exp":     time.Now().Add(time.Minute * time.Duration(cfg.Auth.RefreshTokenMinutes)).Unix(),
	})
	refreshStr, _ := refresh.SignedString([]byte(cfg.Auth.JWTSecret))

	return accessStr, refreshStr
}

func Blacklist(t string)            { blacklist[t] = true }
func IsBlacklisted(t string) bool   { return blacklist[t] }