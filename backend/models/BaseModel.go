package models

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        uint `gorm:"primaryKey;autoIncrement" json:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
