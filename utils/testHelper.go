package utils

import (
	"encoding/json"
	"fmt"

	"github.com/ztrue/tracerr"
)

const (
	jsonExamplePrFromSearch = `{
        "number": %d,
        "state": "open",
        "locked": false,
        "title": "%s",
        "user": %s,
        "body": "%s",
        "assignee": null,
        "assignees": [],
        "requested_reviewers": [],
        "requested_teams": [],
        "labels": [],
        "milestone": null,
        "draft": false,
        "author_association": "NONE"
    }`

	jsonExamplePrFromGET = `{
        "number": %d,
        "state": "open",
        "locked": false,
        "title": "%s",
        "user": %s,
        "body": "%s",
        "assignee": null,
        "assignees": [],
        "requested_reviewers": [],
        "requested_teams": [],
        "labels": [],
        "milestone": null,
        "draft": false,
        "author_association": "NONE",
        "additions":%d,
        "deletions":%d
    }`

	jsonExamplePrFromGETWithAssignees = `{
          "number": %d,
          "state": "open",
          "locked": false,
          "title": "%s",
          "user": %s,
          "body": "%s",
          "assignee": null,
          "assignees": %s,
          "requested_reviewers": [],
          "requested_teams": [],
          "labels": [],
          "milestone": null,
          "draft": false,
          "author_association": "NONE",
          "additions":%d,
          "deletions":%d
      }`
)

func GetJsonPrFromSearch(number int, title, user, body string) (string, error) {
	var mapp map[string]interface{}
	if err := json.Unmarshal([]byte(user), &mapp); err != nil {
		return "", tracerr.New("user must be a valid JSON string")
	}
	return fmt.Sprintf(jsonExamplePrFromSearch, number, title, user, body), nil
}

func GetJsonPrFromGET(number int, title, user, body string, additions, deletions int) (string, error) {
	var mapp map[string]interface{}
	if err := json.Unmarshal([]byte(user), &mapp); err != nil {
		return "", tracerr.New("user must be a valid JSON string")
	}
	return fmt.Sprintf(jsonExamplePrFromGET, number, title, user, body, additions, deletions), nil
}

func GetJsonPrFromGETWithAssignees(number int, title, user, body string, additions, deletions int, assignees string) (string, error) {
	var mapp []interface{}
	if err := json.Unmarshal([]byte(assignees), &mapp); err != nil {
		return "", tracerr.New("assignees must be a valid JSON string")
	}
	return fmt.Sprintf(jsonExamplePrFromGETWithAssignees, number, title, user, body, assignees, additions, deletions), nil
}
