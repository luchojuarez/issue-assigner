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
