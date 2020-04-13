package models

type Issue interface {
	Weight() int
	GetAssignedUsers() []*User
	GetAuthor() *User
	Assing(u *User)
	ToString() string
}
