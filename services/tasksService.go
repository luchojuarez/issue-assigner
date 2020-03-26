package services

import (
	env "github.com/luchojuarez/issue-assigner/environment"

	"github.com/go-resty/resty/v2"
)

const ()

type TaskService struct {
	RestClient *resty.Client
}

func NewTaskService() *TaskService {
	return &TaskService{
		RestClient: env.GetEnv().GetResty("UserService"),
	}
}

// func (this TaskService) AssignTasks() error {
// 	config, err := Load()
// 	if err != nil {
// 		return err
// 	}
// 	config.LoadRepos()
// }
