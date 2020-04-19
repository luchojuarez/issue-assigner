package services

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/ztrue/tracerr"

	"github.com/luchojuarez/issue-assigner/dao"
	"github.com/luchojuarez/issue-assigner/models"

	env "github.com/luchojuarez/issue-assigner/environment"
)

const (
	GithubBaseURL  = "https://api.github.com"
	ConfigFilePath = "config.json"
)

type JsonConfig struct {
	UsersNicknames    []string             `json:"users_niknames"`
	ReviewersPerIssue int                  `json:"reviewers_per_issue"`
	TaskSoruce        []*models.TaskSoruce `json:"task_source"`
	Users             []*models.User
	taskLoaders       []TaskLoader
	IssueList         []models.Issue
	UserService       *UserService
}

func Load(configFilePath string) (*JsonConfig, error) {
	defer TraceTime("initial_data_load", time.Now())
	return load(GithubBaseURL, configFilePath)
}

func load(githubBaseURL, configFilePath string) (*JsonConfig, error) {
	// read main config file
	file, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, TraceError0(tracerr.Wrap(err))
	}
	newConfig := JsonConfig{
		UserService: NewUserService0(),
	}
	// unmarshalling data...
	if err = json.Unmarshal([]byte(file), &newConfig); err != nil {
		return nil, TraceError0(tracerr.Wrap(err))
	}
	// syinchronize tow main task
	mainChan := make(chan bool, 8)

	if err := newConfig.loadTasks(mainChan); err != nil {
		return nil, err
	}

	if err := newConfig.loadUsers(mainChan); err != nil {
		return nil, TraceError0(tracerr.Wrap(err))
	}
	totalTask := 2
	for totalTask > 0 {
		<-mainChan
		totalTask--
	}

	return &newConfig, nil
}

func (this *JsonConfig) loadUsers(mainChan chan bool) error {
	runningRoutings := 0
	topic := make(chan string, len(this.UsersNicknames)+1)
	errors := make(chan error, len(this.UsersNicknames)+1)
	for _, name := range this.UsersNicknames {
		this.UserService.GetUserAsinc(name, topic, errors)

		runningRoutings++
	}

	dao := dao.NewLocalUserDao()
	for runningRoutings > 0 {
		if len(errors) > 0 {
			mainChan <- false
			return <-errors
		}
		fetchedUser := <-topic
		u, err := dao.GetUser(fetchedUser)
		if err != nil {
			mainChan <- false
			return err
		}
		this.Users = append(this.Users, u)
		runningRoutings--
	}
	mainChan <- true
	return nil
}

func (this *JsonConfig) loadTasks(mainChan chan bool) error {
	taskQueue := make(chan bool, 100)
	errorList := make(chan error, 100)
	totalTask := 0

	for _, taskSoruce := range this.TaskSoruce {
		var currentTask TaskLoader
		switch taskSoruce.ResourceType {
		case "all_pr_from_repo":
			currentTask = &AllPrTaskLoader{
				RepoNames: taskSoruce.Resources,
				prService: NewPRService(),
			}
		case "pr_list":
			currentTask = &PrListTaskLoader{
				PrList:    taskSoruce.Resources,
				prService: NewPRService(),
			}
		}
		totalTask += currentTask.GetTotalTask()
		go currentTask.GetAllTask(taskQueue, errorList)
	}
	for totalTask > 0 {
		if len(errorList) > 0 {
			mainChan <- false
			return <-errorList
		}
		<-taskQueue
		totalTask--
	}
	for _, value := range *env.GetEnv().GetPrStorage() {
		for _, pr := range value {
			this.IssueList = append(this.IssueList, pr)
		}
	}
	mainChan <- true
	return nil
}
