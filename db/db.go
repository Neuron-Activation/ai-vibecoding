package db

import (
	"fmt"
	"go-app/models"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

var db *gorm.DB

// InitDB открывает соединение с указанной СУБД и выполняет миграции.
func InitDB(dialect, conn string) error {
	connDB, err := gorm.Open(dialect, conn)
	if err != nil {
		return err
	}
	db = connDB
	log.Info("[DB] Connected")
	if err := migrateSchema(); err != nil {
		return err
	}
	log.Info("[DB] Schema migrated successfully")
	return nil
}

// InitDBFromEnv инициализирует БД на основе переменных окружения.
// Поддерживает:
// - DB_DIALECT=postgres (по умолчанию) с наборами переменных db_user/db_pass/db_name/db_host/db_port
// - DB_DIALECT=sqlite3 с DB_CONN (например ":memory:" или "file:test.db?cache=shared")
func InitDBFromEnv() error {
	dialect := os.Getenv("DB_DIALECT")
	if dialect == "" {
		dialect = "postgres"
	}
	var conn string
	if dialect == "sqlite3" {
		conn = os.Getenv("DB_CONN")
		if conn == "" {
			conn = ":memory:"
		}
	} else {
		user := os.Getenv("db_user")
		pass := os.Getenv("db_pass")
		name := os.Getenv("db_name")
		host := os.Getenv("db_host")
		port := os.Getenv("db_port")
		conn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, name, pass)
	}
	return InitDB(dialect, conn)
}

func GetDB() *gorm.DB {
	return db
}

func CloseDB() error {
	if db == nil {
		return nil
	}
	return db.Close()
}

func migrateSchema() error {
	if db == nil {
		return fmt.Errorf("db not initialized")
	}
	err := db.AutoMigrate(
		models.Note{},
	).Error

	return err
}
