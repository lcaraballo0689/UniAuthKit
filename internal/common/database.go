// internal/common/database.go
//  - Conexión GORM multi‑driver según YAML
package common

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func InitDatabase(cfg *Config) *gorm.DB {
	var db *gorm.DB
	var err error

	switch cfg.Database.Driver {
	case "postgres":
		db, err = gorm.Open(postgres.Open(cfg.Database.DSN), &gorm.Config{})
	case "mysql":
		db, err = gorm.Open(mysql.Open(cfg.Database.DSN), &gorm.Config{})
	case "mssql":
		db, err = gorm.Open(sqlserver.Open(cfg.Database.DSN), &gorm.Config{})
	default: // sqlite por defecto
		db, err = gorm.Open(sqlite.Open(cfg.Database.DSN), &gorm.Config{})
	}

	if err != nil {
		log.Fatal("DB connection error:", err)
	}
	return db
}
