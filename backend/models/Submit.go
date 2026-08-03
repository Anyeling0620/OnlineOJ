package models

import "gorm.io/gorm"

type Submit struct {
	gorm.Model
	Identity        string `json:"identity"`
	ProblemIdentity string `json:"problem_identity"`
	UserIdentity    string `json:"user_identity"`
	Path            string `json:"path"`
	Status          int    `gorm:"type:smallint" json:"status"`
}

func (submit *Submit) TableName() string {
	return "submit"
}
