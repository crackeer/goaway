package model

// User
type User struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	PasswordMD5 string `json:"password_md5"`
	UserType    string `json:"user_type"`
	CreateAt    int64  `json:"create_at"`
	ModifyAt    int64  `json:"modify_at"`
}

func (User) TableName() string {
	return "user"
}

// Map
//  @receiver user
//  @return map
func (user User) Map() map[string]interface{} {
	return map[string]interface{}{
		"id":           user.ID,
		"username":     user.Username,
		"password_md5": user.PasswordMD5,
		"user_type":    user.UserType,
	}
}

func init() {
	registerNewModelFunc("user", func() (interface{}, interface{}) {
		return &User{}, []User{}
	})
}
