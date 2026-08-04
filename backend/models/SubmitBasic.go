package models

import "gorm.io/gorm"

type SubmitBasic struct {
	BaseModel
	Identity        string `json:"identity"`
	ProblemIdentity string `json:"problem_identity"`
	UserIdentity    string `json:"user_identity"`
	Path            string `json:"path"`
	Status          int    `gorm:"type:smallint" json:"status"` // [0-待判断 1-答案正确 2-答案错误 3-超出时间限制 4-超出内存限制 5-编译错误]

	ProblemBasic *ProblemBasic `gorm:"foreignKey:ProblemIdentity;references:Identity" json:"problem_basic"`
	UserBasic    *UserBasic    `gorm:"foreignKey:UserIdentity;references:Identity" json:"user_basic"`
}

func (submit *SubmitBasic) TableName() string {
	return "submit_basic"
}

func GetSubmitList(problemIdentity, userIdentity string, status int) *gorm.DB {
	tx := DB.Model(new(SubmitBasic)).
		Preload("ProblemBasic", func(db *gorm.DB) *gorm.DB {
			return db.Omit("Content")
		}).
		Preload("UserBasic", func(db *gorm.DB) *gorm.DB {
			return db.Omit("Password")
		})

	if problemIdentity != "" {
		tx = tx.Where(
			"submit_basic.problem_identity = ?",
			problemIdentity,
		)
	}

	if userIdentity != "" {
		tx = tx.Where(
			"submit_basic.user_identity = ?",
			userIdentity,
		)
	}

	if status != -1 {
		tx = tx.Where("submit_basic.status = ?", status)
	}

	return tx
}
