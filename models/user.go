package models

import "time"

type User struct {
	NickName    string `json:"login"`
	FetchedAt   time.Time
	RequestTime int64
}
