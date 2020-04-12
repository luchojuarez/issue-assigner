package models

type Issue interface {
	Weight() int
	GetAssignedUsers() []*User
	Assing(u *User)
	ToString() string
}
