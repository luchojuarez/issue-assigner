package services

import (
	"log"
	"testing"

	"github.com/luchojuarez/issue-assigner/models"

	env "github.com/luchojuarez/issue-assigner/environment"

	"github.com/stretchr/testify/assert"
)

func TestSuccesRun(t *testing.T) {
	env.GetEnv().CleanAll()
	echoFile("../out/out.log")
	mockConfigSuccessCase()
	mockPRWhit2Reviwers()
	assignmentService, _ := NewAssignmentService(jsonResourcesPath + "config_test.json")
	prNeedsReviwes := make([]*models.PR, 0)
	for _, repo := range assignmentService.config.Repos {
		for _, pr := range repo.PullRequests {
			if len(pr.AssignedUsers) < assignmentService.config.ReviewersPerPR {
				prNeedsReviwes = append(prNeedsReviwes, pr)
			}
		}
	}

	log.Printf("total pr needs revieds: %d \n================================", len(prNeedsReviwes))
	assignmentService.Run()

	for _, pr := range prNeedsReviwes {
		assert.Equal(t, assignmentService.config.ReviewersPerPR, len(pr.AssignedUsers))
	}
	PrintAndClear("success_run")
}

func TestInvalidJsonInput(t *testing.T) {
	service, err := NewAssignmentService(jsonResourcesPath + "invalid.json")

	assertNil(t, service)
	assert.Equal(t, "invalid character 'n' looking for beginning of object key string", err.Error())
}

func mockPRWhit2Reviwers() {
	simpleStringResponderForPrSearch("luchojuarez/crypto", `[{"number": 1},{"number": 2},{"number": 3},{"number": 4}]`, 200, 0)
	simpleStringResponderForGetPR(3, "luchojuarez/crypto", `{"number": 3,"title":"Title 3 no more reviewers","body":"description 1","assignees":[{"login":"luchojuarez"},{"login":"luchojuarez2"},{"login":"user3"}],"user":{"login":"luchojuarez"},"commits": 2,"additions": 353,  "deletions": 18}`, 200, 0)
	simpleStringResponderForGetPR(4, "luchojuarez/crypto", `{"number": 4,"title":"Title 4 no more reviewers","body":"description 1","assignees":[{"login":"luchojuarez2"},{"login":"user3"}],"user":{"login":"luchojuarez"},"commits": 2,"additions": 353,  "deletions": 18}`, 200, 0)
}
