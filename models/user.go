package models

import "time"

type User struct {
	NickName        string `json:"login"`
	FetchedAt       time.Time
	RequestTime     int64
	AssignedPRLines int
	AssignedPR      []*PR
}

func (this *User) AssingPR(pullRequest *PR) int {
	this.AssignedPR = append(this.AssignedPR, pullRequest)
	this.AssignedPRLines += pullRequest.Deletions + pullRequest.Additions
	return pullRequest.Deletions + pullRequest.Additions
}
