package models

import "gorm.io/gorm"

type CategoryBasic struct {
	gorm.Model
	Name     string `json:"name"`
	ParentId int    `json:"parent_id"`
}

func (category *CategoryBasic) TableName() string {
	return "category_basic"
}
