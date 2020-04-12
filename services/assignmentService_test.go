package services

import (
	"log"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	env "github.com/luchojuarez/issue-assigner/environment"
	"github.com/stretchr/testify/assert"
)

func TestSuccesRun(t *testing.T) {
	httpmock.Reset()
	env.GetEnv().ClearUserStorage()
	defer PrintAndClearWhithBeginTime("../out/success_run.log", time.Now())

	mockConfigSuccessCase()
	mockPRWhit2Reviwers()
	assignmentService, err := NewAssignmentService(jsonResourcesPath + "config_test.json")
	log.Printf("cargo la config --- %v", err)

	assignmentService.Run()
	log.Printf("Corrio la config --- %v", assignmentService.config)

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
	simpleStringResponderForGetPR(7, "user1/foo", `{"number": 4,"title":"Title 4 no more reviewers","body":"description 1","user":{"login":"luchojuarez"},"commits": 2,"additions": 353,  "deletions": 18}`, 200, 0)
	simpleStringResponderForGetPR(10, "user2/bar", `{"number": 4,"title":"Title 4 no more reviewers","body":"description 1","user":{"login":"luchojuarez"},"commits": 2,"additions": 353,  "deletions": 18}`, 200, 0)
}
