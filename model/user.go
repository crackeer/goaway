package model

type User struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	PasswordMD5 string `json:"password_md5"`
}

func (User) TableName() string {
	return "user"
}
