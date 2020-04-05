package services

import (
	"github.com/luchojuarez/issue-assigner/models"
)

type AssignmentService struct {
	config              JsonConfig
	UserServiceInstance *UserService
}

func NewAssignmentService(configFilePath string) (*AssignmentService, error) {
	c, err := Load(configFilePath)
	if err != nil {
		return nil, err
	}
	service := AssignmentService{
		config:              *c,
		UserServiceInstance: NewUserService0(),
	}
	return &service, nil
}

func (this *AssignmentService) Run() {
	for _, currentRepo := range this.config.Repos {
		for _, currentPR := range currentRepo.PullRequests {
			assignedUsers := len(currentPR.AssignedUsers)
			for assignedUsers < this.config.ReviewersPerPR {
				iddleUser := this.UserServiceInstance.GetSortedUsersByAssignations(&this.config)[0]
				this.assingn(iddleUser, currentPR, currentRepo)
				assignedUsers += 1
			}
		}
	}
}

func (this *AssignmentService) assingn(user *models.User, pull *models.PR, repo *models.Repo) {
	linesAdded := user.AssingPR(pull)
	pull.AssignedUsers = append(pull.AssignedUsers, user)
	TraceInfof("NEW assing from repo:'%s', PR(%d) '%s' to user '%s', assigned lines %d", repo.FullName, pull.Number, pull.Body, user.NickName, linesAdded)
}
