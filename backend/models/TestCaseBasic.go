package models

type TestCase struct {
	BaseModel
	Identity        string `json:"identity"`
	ProblemIdentity string `json:"problem-identity"`
	Input           string `json:"input"`
	Output          string `json:"output"`
}

func (table *TestCase) TableName() string {
	return "test_case"
}
