package services

import (
	"testing"
)

func TestSuccesRun(t *testing.T) {
	mockConfigSuccessCase()
	assignmentService, _ := NewAssignmentService(jsonResourcesPath + "config_test.json")

	assignmentService.run()
}
