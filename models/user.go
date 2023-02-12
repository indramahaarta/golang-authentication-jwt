package models

type User struct {
	ID       int    `json:"id"`
	UserName string `json:"user_name"`
	FullName string `json:"full_name"`
	Password string `json:"password"`
}
