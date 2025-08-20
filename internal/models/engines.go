package models

type Engine struct {
	BaseModel
	Remark string `json:"remark"`
	Url    string `json:"url"`
	UserId uint   `json:"user_id"`
	Token  string `json:"token"`
}
