package services

import (
	"log"
	"sync"

	"github.com/luchojuarez/issue-assigner/models"
	"github.com/ztrue/tracerr"
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

func (this *AssignmentService) Run() error {
	if this == nil {
		return tracerr.New("Cant load config...")
	}
	for _, currentIssue := range this.config.IssueList {
		assignedUsers := len(currentIssue.GetAssignedUsers())
		issueAuthor := currentIssue.GetAuthor().NickName
		for assignedUsers < this.config.ReviewersPerIssue {
			usersList := this.UserServiceInstance.GetSortedUsersByAssignations(&this.config)
			iddleUser := usersList[0]
			if iddleUser.NickName == issueAuthor {
				log.Printf("owner can't revive their own issue (%s)", iddleUser.NickName)
				iddleUser = usersList[1]
			}
			this.assingn(iddleUser, &currentIssue)
			assignedUsers += 1
		}
	}
	return nil
}

func (this *AssignmentService) assingn(user *models.User, issue *models.Issue) {
	user.AssingIssue(*issue)
	(*issue).Assing(user)
	TraceInfof("NEW assing from issue:'%v', to user '%s', assigned lines %d", (*issue).ToString(), user.NickName, (*issue).Weight())
}
