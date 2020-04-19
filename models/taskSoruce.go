package models

import (
	"log"
)

type ExeptionRules struct {
	Type    string   `json:"type"`
	Content []string `json:"content"`
}
type TaskSource struct {
	ResourceType  string           `json:"resource_type"`
	Resources     []string         `json:"resources"`
	ExeptionRules []*ExeptionRules `json:"exeption_rules,omitempty"`
}

func (this *TaskSource) Evaluate(i Issue) *string {
	if this == nil {
		return nil
	}
	for _, e := range this.ExeptionRules {
		if messagge := e.exclude(i); messagge != nil {
			return messagge
		}
	}
	return nil
}

func (this *ExeptionRules) exclude(i Issue) *string {
	switch this.Type {
	case "by_label":
		return this.excludedByLable(i)
	case "needed_lebel":
		return this.excludedByLableNeeded(i)
	}

	return nil
}

// exlude current ISSUE iff repo contains some especific lable
func (this *ExeptionRules) excludedByLable(i Issue) *string {
	log.Printf("esto llega al excludedByLable '%v', '%s'", this, i)
	return nil
}

//Exclude current issue if missing lables attached in "content"
func (this *ExeptionRules) excludedByLableNeeded(i Issue) *string {
	log.Printf("esto llega al excludedByLableNeeded '%v', '%s'", this, i)

	return nil
}
