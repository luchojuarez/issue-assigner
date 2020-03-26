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
	RepoNames      []string `json:"repos"`
	UsersNicknames []string `json:"users_niknames"`
	GithubToken    string   `json:"github_token"`
	Repos          []*models.Repo
	Users          []*models.User
	UserService    *UserService
}

func Load() (*JsonConfig, error) {
	return load(GithubBaseURL, ConfigFilePath)
}

func load(githubBaseURL, configFilePath string) (*JsonConfig, error) {
	// read main config file
	file, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	newConfig := JsonConfig{UserService: NewUserService()}
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
