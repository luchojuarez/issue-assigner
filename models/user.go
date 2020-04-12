package models

import "time"

type User struct {
	NickName          string `json:"login"`
	FetchedAt         time.Time
	RequestTime       int64
	AssignedTaskValue int
	AssignedIssues    []*Issue
}

func (this *User) AssingIssue(issue Issue) int {
	this.AssignedIssues = append(this.AssignedIssues, &issue)
	this.AssignedTaskValue += issue.Weight()
	return issue.Weight()
}
