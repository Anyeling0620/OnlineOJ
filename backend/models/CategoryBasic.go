package models

type CategoryBasic struct {
	BaseModel
	Identity          string             `json:"identity"`
	Name              string             `json:"name"`
	ParentId          int                `json:"parent_id"`
	ProblemCategories []*ProblemCategory `gorm:"foreignKey:CategoryId;references:ID" json:"problem_categories"`
}

func (category *CategoryBasic) TableName() string {
	return "category_basic"
}
