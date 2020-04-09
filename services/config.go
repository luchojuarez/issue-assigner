package services

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/ztrue/tracerr"

	"github.com/luchojuarez/issue-assigner/dao"
	"github.com/luchojuarez/issue-assigner/models"
)

const (
	GithubBaseURL  = "https://api.github.com"
	ConfigFilePath = "config.json"
)

type JsonConfig struct {
	RepoNames      []string `json:"repos_full_names"`
	UsersNicknames []string `json:"users_niknames"`
	GithubToken    string   `json:"github_token"`
	ReviewersPerPR int      `json:"reviewers_per_pr"`
	Repos          []*models.Repo
	Users          []*models.User
	UserService    *UserService
	prService      *PRService
}

func Load(configFilePath string) (*JsonConfig, error) {
	defer TraceTime("initial_data_load_sync", time.Now())
	return load(GithubBaseURL, configFilePath, "sync")
}

func LoadAsync(configFilePath string) (*JsonConfig, error) {
	defer TraceTime("initial_data_load_async", time.Now())
	return load(GithubBaseURL, configFilePath, "async")
}

func load(githubBaseURL, configFilePath, strategy string) (*JsonConfig, error) {
	// read main config file
	file, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, TraceError0(tracerr.Wrap(err))
	}
	newConfig := JsonConfig{
		UserService: NewUserService0(),
		prService:   NewPRService(),
	}
	// unmarshalling data...
	if err = json.Unmarshal([]byte(file), &newConfig); err != nil {
		return nil, TraceError0(tracerr.Wrap(err))
	}

	switch strategy {
	case "async":
		// load user info
		if err := newConfig.asyncLoadUsers(); err != nil {
			return nil, TraceError0(tracerr.Wrap(err))
		}
		// load repository info
		if err := newConfig.asyncLoadRepos(); err != nil {
			return nil, TraceError0(tracerr.Wrap(err))
		}
	case "sync":
		// load user info
		if err := newConfig.loadUsers(); err != nil {
			return nil, TraceError0(tracerr.Wrap(err))
		}
		// load repository info
		if err := newConfig.loadRepos(); err != nil {
			return nil, TraceError0(tracerr.Wrap(err))
		}
	default:
		return nil, TraceError0(tracerr.New("invalid load strategy " + strategy))
	}

	return &newConfig, nil
}

func (this *JsonConfig) loadRepos() error {
	for _, repoName := range this.RepoNames {
		newRepo := models.NewRepo(repoName)
		prList, err := this.prService.GetOpenPRs(newRepo.FullName)
		if err != nil {
			return tracerr.Wrap(err)
		}
		newRepo.PullRequests = prList
		this.Repos = append(this.Repos, &newRepo)
	}
	return nil // tracerr.New("not implemented yet!")
}

func (this *JsonConfig) loadUsers() error {
	for _, name := range this.UsersNicknames {
		newUser, err := this.UserService.GetUser(name)
		if err != nil {
			return tracerr.Wrap(err)
		}
		this.Users = append(this.Users, newUser)
	}
	return nil
}

//Asinc methods
func (this *JsonConfig) asyncLoadRepos() error {
	runningRoutings := 0
	topic := make(chan string, len(this.RepoNames)+1)
	errors := make(chan error, len(this.RepoNames)+1)
	// launch rutines
	for _, repoName := range this.RepoNames {
		this.prService.GetOpenPRsAsinc(repoName, topic, errors)
		runningRoutings++
	}

	dao := dao.NewLocalPrDao()

	for runningRoutings > 0 {
		if len(errors) > 0 {
			return <-errors
		}
		repoFetched := <-topic

		newRepo := models.NewRepo(repoFetched)
		newRepo.PullRequests, _ = dao.GetPrByRepo(repoFetched)
		this.Repos = append(this.Repos, &newRepo)

		runningRoutings--
	}

	return nil
}

func (this *JsonConfig) asyncLoadUsers() error {
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
			return <-errors
		}
		fetchedUser := <-topic
		u, err := dao.GetUser(fetchedUser)
		if err != nil {
			return err
		}
		this.Users = append(this.Users, u)
		runningRoutings--
	}

	return nil
}
