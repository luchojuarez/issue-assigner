package models

import (
	"fmt"
	"log"

	"github.com/luchojuarez/issue-assigner/utils"
)

type ExeptionRules struct {
	Type    string   `json:"type"`
	ApplyTo []string `json:"apply_to,omitempty"`
	ExeptTo []string `json:"exept_to,omitempty"`
	Content []string `json:"content"`
}
type TaskSource struct {
	ResourceType  string           `json:"resource_type"`
	Resources     []string         `json:"resources"`
	ExeptionRules []*ExeptionRules `json:"exeption_rules,omitempty"`
}

func (this *TaskSource) Evaluate(pr *PR) *string {
	if this == nil {
		return nil
	}
	for _, e := range this.ExeptionRules {
		if messagge := e.exclude(pr); messagge != nil {
			return messagge
		}
	}
	return nil
}

func (this *ExeptionRules) exclude(pr *PR) *string {
	switch this.Type {
	case "by_label":
		return this.excludedByLable(pr)
	case "needed_lebel":
		return this.excludedByLableNeeded(pr)
	}

	return nil
}

// exlude current ISSUE iff repo contains some especific lable
func (this *ExeptionRules) excludedByLable(pr *PR) *string {
	for _, currentL := range this.Content {
		for _, prLabels := range pr.Labels {
			if currentL != prLabels {
				continue
			}
			//chek exeptirions rules
			if !utils.ContainsAny(this.ApplyTo, pr.Labels) {
				continue
			}
			if utils.ContainsAny(this.ExeptTo, pr.Labels) {
				continue
			}
			errorMessagge := fmt.Sprintf("%s excluded reason '%s' by '%s'", pr.ToString(), currentL, this.Type)
			log.Print(errorMessagge)
			return &errorMessagge
		}
	}
	return nil
}

//Exclude current issue if missing lables attached in "content"
func (this *ExeptionRules) excludedByLableNeeded(pr *PR) *string {
	if utils.Contains(this.ApplyTo, pr.Repo.FullName) {
		return nil
	}
	found := false

	for _, currentL := range this.Content {
		for _, prLabels := range pr.Labels {
			if currentL == prLabels {
				found = true
			}
		}
	}
	if !found {
		errorMessagge := fmt.Sprintf("%s excluded reason '%s' by '%v'", pr.ToString(), this.Content, this.Type)
		log.Print(errorMessagge)
		return &errorMessagge
	}

	return nil
}
