package models

import "gorm.io/gorm"

type UserBasic struct {
	gorm.Model
	Identity  string `json:"identity"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	Phone     string `json:"phone"`
	Mail      string `json:"mail"`
	PassNum   int    `json:"pass_num"`
	SubmitNum int    `json:"submit_num"`
}

func (userBasic *UserBasic) TableName() string {
	return "user_basic"
}
