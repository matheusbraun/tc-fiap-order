package postgres

import (
	"fmt"
	"log"
	"os"

	orderEntities "github.com/viniciuscluna/tc-fiap-50/internal/order/domain/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB() *gorm.DB {
	dsn := newDBConfig()
	db := newDB(dsn)
	migrate(db)
	return db
}

func newDBConfig() string {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")

	if host == "" || user == "" || password == "" || dbname == "" || port == "" || sslmode == "" {
		log.Fatalf("Database environment variables are not properly set")
	}

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, dbname, port, sslmode)
}

func newDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return db
}

func migrate(db *gorm.DB) {
	// Only migrate order-related entities (isolated microservice database)
	if err := db.AutoMigrate(
		&orderEntities.OrderEntity{},
		&orderEntities.OrderProductEntity{},
		&orderEntities.OrderStatusEntity{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
