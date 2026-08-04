package models

type TestCaseBasic struct {
	BaseModel
	Identity        string `json:"identity"`
	ProblemIdentity string `json:"problem_identity"`
	Input           string `json:"input"`
	Output          string `json:"output"`
}

func (table *TestCaseBasic) TableName() string {
	return "test_case_basic"
}
