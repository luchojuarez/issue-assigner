package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/luchojuarez/issue-assigner/models"
	"github.com/luchojuarez/issue-assigner/services"
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
	UserService    *services.UserService
}

func Load() (*JsonConfig, error) {
	return load(GithubBaseURL, ConfigFilePath)
}

func load(githubBaseURL, configFilePath string) (*JsonConfig, error) {
	// read main config file
	file, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}
	newConfig := JsonConfig{UserService: services.NewUserService()}
	// unmarshalling data...
	if err = json.Unmarshal([]byte(file), &newConfig); err != nil {
		return nil, err
	}

	// load user info
	if err := newConfig.loadUsers(); err != nil {
		return nil, err
	}
	// load repository info
	if err := newConfig.loadRepos(); err != nil {
		return nil, err
	}

	return nil, err
}

func (this JsonConfig) loadRepos() error {
	log.Printf("estos repos '%v'", this.RepoNames)

	return nil
}

func (this JsonConfig) loadUsers() error {
	for _, name := range this.UsersNicknames {
		log.Printf("estos users '%s'", name)
		newUser, err := this.UserService.GetUser(name)
		if err != nil {
			return err
		}
		this.Users = append(this.Users, newUser)
	}
	return nil
}
