package models

type Issue struct {
	Title string `json:"title"`
	User  *User
}
