package models

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB, SqlDB = Init()

func Init() (*gorm.DB, *sql.DB) {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("load .env file failed,", err)
	}
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_IP"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("load db failed,", err)
	}

	sqlDB, _ := db.DB()
	return db, sqlDB
}
