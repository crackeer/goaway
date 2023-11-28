package model

import "encoding/json"

// User
type User struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	PasswordMD5 string `json:"password_md5"`
	UserType    string `json:"user_type"`
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

// Map2User
//  @param tmp
//  @return User
func Map2User(tmp map[string]interface{}) User {
	bytes, _ := json.Marshal(tmp)
	user := User{}
	json.Unmarshal(bytes, &user)
	return user
}
