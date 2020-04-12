package models

import (
	"fmt"
	"time"
)

type PR struct {
	Commits       int           `json:"commits"`
	Additions     int           `json:"additions"`
	Deletions     int           `json:"deletions"`
	Number        int           `json:"number"`
	Title         string        `json:"title"`
	Body          string        `json:"body"`
	Assignees     []interface{} `json:"assignees"`
	Repo          *Repo         `json:"repo"`
	AssignedUsers []*User
	AssigneesSize int
	User          *User `json:"user"`
	FetchedAt     time.Time
	RequestTime   int64
}

func (this *PR) Weight() int {
	return this.Additions + this.Deletions
}

func (this *PR) GetAssignedUsers() []*User {
	return this.AssignedUsers
}

func (this *PR) GetAuthor() *User {
	return this.User
}

func (this *PR) Assing(u *User) {
	this.AssignedUsers = append(this.AssignedUsers, u)
}
func (this *PR) ToString() string {
	return fmt.Sprintf("PR %s(%d) by '%s' weight: %d", "this.Repo.Name", this.Number, this.User.NickName, this.Weight())
}

func (this *PR) SetEndTime(initTime time.Time) {
	this.FetchedAt = initTime
	endMillis := time.Now().UnixNano() / int64(time.Millisecond)
	this.RequestTime = endMillis - (initTime.UnixNano() / int64(time.Millisecond))
}
