package services

import (
	"encoding/json"
	"io/ioutil"

	"github.com/ztrue/tracerr"

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
	return load(GithubBaseURL, configFilePath)
}

func load(githubBaseURL, configFilePath string) (*JsonConfig, error) {
	// read main config file
	file, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	newConfig := JsonConfig{
		UserService: NewUserService0(),
		prService:   NewPRService(),
	}
	// unmarshalling data...
	if err = json.Unmarshal([]byte(file), &newConfig); err != nil {
		return nil, tracerr.Wrap(err)
	}

	// load user info
	if err := newConfig.loadUsers(); err != nil {
		return nil, tracerr.Wrap(err)
	}
	// load repository info
	if err := newConfig.loadRepos(); err != nil {
		return nil, tracerr.Wrap(err)
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
