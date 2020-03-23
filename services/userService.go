package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	env "github.com/luchojuarez/issue-assigner/environment"
	"github.com/luchojuarez/issue-assigner/models"
)

type UserService struct {
	RestClient    *resty.Client
	fetchedUsers  map[string]*models.User
	GithubBaseURL string
}

func NewUserService() *UserService {
	return NewUserServiceCapacity(0)
}
func NewUserServiceCapacity(cap int) *UserService {
	return &UserService{
		RestClient:    env.GetEnv().GetResty("UserService"),
		fetchedUsers:  make(map[string]*models.User, cap),
		GithubBaseURL: "https://api.github.com",
	}
}

func (this UserService) GetUser(nickname string) (*models.User, error) {
	// chek if user already fetched
	if this.fetchedUsers[nickname] == nil {
		// fetch and save locally
		newUser, err := this.getUser(nickname)
		if err != nil {
			return nil, err
		}
		this.fetchedUsers[nickname] = newUser
	}

	return this.fetchedUsers[nickname], nil
}

// this private function call github API to get user info
func (this UserService) getUser(nickname string) (*models.User, error) {
	startMillis := time.Now().UnixNano() / int64(time.Millisecond)
	response, err := this.RestClient.
		R().
		Get(this.GithubBaseURL + "/users/" + nickname)

	if err != nil {
		return nil, err
	}

	if response.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("invalid status code: '%d'", response.StatusCode())
	}

	var newUser models.User

	if err = json.Unmarshal([]byte(response.Body()), &newUser); err != nil {
		return nil, err
	}
	endMillis := time.Now().UnixNano() / int64(time.Millisecond)
	newUser.FetchedAt = time.Now()
	newUser.RequestTime = endMillis - startMillis

	return &newUser, nil
}
