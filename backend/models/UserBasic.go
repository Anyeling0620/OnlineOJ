package models

type UserBasic struct {
	BaseModel
	Identity  string `json:"identity"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	Phone     string `json:"phone"`
	Mail      string `json:"mail"`
	PassNum   int    `json:"pass_num"`
	SubmitNum int    `json:"submit_num"`
	IsAdmin   int    `json:"is_admin"` // 0-非管理员 1-管理员
}

func (userBasic *UserBasic) TableName() string {
	return "user_basic"
}
