package models

type ProblemCategory struct {
	BaseModel
	ProblemID  uint `json:"problem_id"`
	CategoryID uint `json:"category_id"`

	ProblemBasic  *ProblemBasic  `gorm:"foreignKey:ProblemID;references:ID" json:"problem_basic"`
	CategoryBasic *CategoryBasic `gorm:"foreignKey:CategoryID;references:ID" json:"category_basic"`
}

func (problemCategory *ProblemCategory) TableName() string {
	return "problem_category"
}
