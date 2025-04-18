// internal/common/config.go
//  - Lee config.yaml en un struct global
package common

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

// Config refleja exactamente las claves del YAML
type Config struct {
	App struct{ Name, Environment, Version string }
	Server struct{ Port string }
	Auth struct {
		JWTSecret            string `yaml:"jwt_secret"`
		TokenDurationMinutes int    `yaml:"token_duration_minutes"`
		RefreshTokenMinutes  int    `yaml:"refresh_token_minutes"`
		MFAEnabled           bool   `yaml:"mfa_enabled"`
		AttemptsLimit        int    `yaml:"attempts_limit"`
		AttemptsWindowMin    int    `yaml:"attempts_window_minutes"`
	}
	Database struct{ Driver, DSN string }
	Multitenancy struct{ Enabled bool }
	OAuth  map[string]any `yaml:"oauth"`
	I18n   map[string]any `yaml:"i18n"`
	ResponseTemplate struct {
		StatusField, CodeField, MessageField, DataField,
		ErrorField, TimestampField string
		UserFields []string `yaml:"user_fields"`
	} `yaml:"response_template"`
	Endpoints map[string]bool
	Queries   map[string]string
}

// LoadConfig lee el YAML y devuelve Config
func LoadConfig(path string) Config {
	var c Config
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("No se pudo leer config.yaml:", err)
	}
	if err := yaml.Unmarshal(b, &c); err != nil {
		log.Fatal("YAML inválido:", err)
	}
	return c
}
