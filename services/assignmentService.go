package services

import (
	"sync"

	"github.com/luchojuarez/issue-assigner/models"
)

type AssignmentService struct {
	config              JsonConfig
	UserServiceInstance *UserService
	lock                *sync.Mutex
	taskList            []*models.Issue
}

func NewAssignmentService(configFilePath string) (*AssignmentService, error) {
	c, err := Load(configFilePath)
	if err != nil {
		return nil, err
	}
	service := AssignmentService{
		config:              *c,
		UserServiceInstance: NewUserService0(),
		lock:                &sync.Mutex{},
	}
	return &service, nil
}

func (this *AssignmentService) Run() {
	for _, currentIssue := range this.config.IssueList {
		assignedUsers := len(currentIssue.GetAssignedUsers())
		for assignedUsers < this.config.ReviewersPerIssue {
			iddleUser := this.UserServiceInstance.GetSortedUsersByAssignations(&this.config)[0]
			this.assingn(iddleUser, &currentIssue)
			assignedUsers += 1
		}
	}
}

func (this *AssignmentService) assingn(user *models.User, issue *models.Issue) {
	user.AssingIssue(*issue)
	(*issue).Assing(user)
	TraceInfof("NEW assing from issue:'%v', to user '%s', assigned lines %d", (*issue).ToString(), user.NickName, (*issue).Weight())
}
