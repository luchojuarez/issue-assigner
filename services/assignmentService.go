package services

import (
	"log"
)

type AssignmentService struct {
	config              JsonConfig
	UserServiceInstance *UserService
}

func NewAssignmentService(configFilePath string) (AssignmentService, error) {
	c, err := Load(configFilePath)
	if err != nil {
		return AssignmentService{}, err
	}
	return AssignmentService{
		config:              *c,
		UserServiceInstance: NewUserService0(),
	}, nil
}

func (this *AssignmentService) run() {
	for _, currentRepo := range this.config.Repos {
		for _, currentPR := range currentRepo.PullRequests {
			if len(currentPR.AssignedUsers) < this.config.ReviewersPerPR {
				iddleUser := this.UserServiceInstance.GetSortedUsersByAssignations()[0]
				iddleUser.AssingPR(currentPR)
				log.Printf("Assing from repo:'%s', PR(%d) '%s' to user '%s'", currentRepo.FullName, currentPR.Number, currentPR.Body, iddleUser.NickName)
			} else {
				break
			}
		}
	}
}
