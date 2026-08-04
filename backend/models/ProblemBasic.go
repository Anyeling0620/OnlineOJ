package models

import (
	"gorm.io/gorm"
)

type ProblemBasic struct {
	BaseModel
	Identity   string `json:"identity"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	PassNum    int    `json:"pass_num"`
	SubmitNum  int    `json:"submit_num"`
	MaxRuntime int    `json:"max_runtime"`
	MaxMem     int    `json:"max_mem"`

	ProblemCategories []*ProblemCategory `gorm:"foreignKey:ProblemId;references:ID" json:"problem_categories"`
}

func (problemBasic *ProblemBasic) TableName() string {
	return "problem_basic"
}

func GetProblemList(keyword, categoryIdentity string) *gorm.DB {
	tx := DB.Model(new(ProblemBasic)).Preload("ProblemCategories").Preload("ProblemCategories.CategoryBasic").Where("title like ? OR content like ? ", "%"+keyword+"%", "%"+keyword+"%")
	if categoryIdentity != "" {
		tx.Joins("right join problem_category pc on pc.problem_id = problem_basic.id").Where("pc.category_id = (select cb.id from category_basic cb where cb.identity = ?)", categoryIdentity)
	}
	return tx
}
