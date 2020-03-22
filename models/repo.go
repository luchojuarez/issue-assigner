package models

type Repo struct {
	Name          string `json:"name"`
	FullName      string `json:"full_name"`
	Description   string `json:"description"`
	Contributors  []*User
	Collaborators []*User
}
