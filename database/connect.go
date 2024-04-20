package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var slowLogger = logger.New(
	log.New(os.Stdout, "\r\n", log.LstdFlags),
	logger.Config{
		SlowThreshold: 1 * time.Microsecond,
		LogLevel:      logger.Silent,
		Colorful:      true,
	},
)

type Database interface {
	Connect() (*gorm.DB, error)
	Migration(data interface{}) error
	CreateTestData(c *fiber.Ctx) error
}

type Data struct {
	DB *gorm.DB
}

func (method *Data) Connect() (*gorm.DB, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("Ошибка чтения переменного окружения: %w", err)
	}
	dsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: slowLogger})
	if err != nil {
		return nil, fmt.Errorf("Ошабка подключения: %w", err)
	}
	method.DB = db
	log.Println("Подключение в бд прошло успешно.")
	return db, nil
}
