package test

import (
	"fmt"
	"github.com/Anyeling0620/OnlineOJ/backend/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"testing"
)

func TestGorm(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatal("load .env file failed,", err)
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
		t.Fatal("load db failed,", err)
	}

	data := make([]models.ProblemBasic, 0)
	err = db.Find(&data).Error
	if err != nil {
		t.Fatal("test load table failed", err)
	}

	for _, v := range data {
		fmt.Printf("Proble => %+v\n", v)
	}
}
