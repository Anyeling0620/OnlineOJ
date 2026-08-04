package models

import "gorm.io/gorm"

type SubmitBasic struct {
	BaseModel
	Identity  string `json:"identity"`
	ProblemId int    `json:"problem_id"`
	UserId    int    `json:"user_id"`
	Path      string `json:"path"`
	Status    int    `gorm:"type:smallint" json:"status"` // [0-待判断 1-答案正确 2-答案错误 3-超出时间限制 4-超出内存限制 5-编译错误]

	ProblemBasic *ProblemBasic `gorm:"foreignKey:ProblemId;references:ID" json:"problem_basic"`
	UserBasic    *UserBasic    `gorm:"foreignKey:UserId;references:ID" json:"user_basic"`
}

func (submit *SubmitBasic) TableName() string {
	return "submit_basic"
}

func GetSubmitList(problemIdentity, userIdentity string, status int) *gorm.DB {
	tx := DB.Model(new(SubmitBasic)).Preload("ProblemBasic", func(db *gorm.DB) *gorm.DB {
		return db.Omit("Content")
	}).Preload("UserBasic", func(db *gorm.DB) *gorm.DB {
		return db.Omit("Password")
	})
	if problemIdentity != "" {
		tx = tx.Joins("JOIN problem_basic AS problem_filter ON problem_filter.id = submit_basic.problem_id").
			Where("problem_filter.identity = ?", problemIdentity)
	}
	if userIdentity != "" {
		tx = tx.Joins("JOIN user_basic AS user_filter ON user_filter.id = submit_basic.user_id").
			Where("user_filter.identity = ?", userIdentity)
	}
	if status != -1 {
		tx = tx.Where("submit_basic.status = ?", status)
	}
	return tx
}
