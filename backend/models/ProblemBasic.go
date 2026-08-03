package models

import "gorm.io/gorm"

type ProblemBasic struct {
	gorm.Model
	Identity   string `json:"identity"`
	CategoryId string `json:"category_id"`
	Title      string `json:"title"`
	PassNum    int    `json:"pass_num"`
	SubmitNum  int    `json:"submit_num"`
	MaxRuntime int    `json:"max_runtime"`
	MaxMem     int    `json:"max_mem"`
}

func (problemBasic *ProblemBasic) TableName() string {
	return "problem_basic"
}
