package models

type TaskSoruce struct {
	ResourceType string   `json:"resource_type"`
	Resources    []string `json:"resources"`
}
