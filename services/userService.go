package services

import (
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/go-resty/resty/v2"
	env "github.com/luchojuarez/issue-assigner/environment"
	"github.com/luchojuarez/issue-assigner/models"
	"github.com/ztrue/tracerr"
)

type UserService struct {
	RestClient    *resty.Client
	userStorage   *map[string]*models.User
	GithubBaseURL string
}

func NewUserService() *UserService {
	return NewUserServiceCapacity(0)
}
func NewUserServiceCapacity(cap int) *UserService {
	return &UserService{
		RestClient:    env.GetEnv().GetResty("UserService"),
		userStorage:   env.GetEnv().GetUserStorage(),
		GithubBaseURL: "https://api.github.com",
	}
}

func (this *UserService) GetUser(nickname string) (*models.User, error) {
	// chek if user already fetched
	if (*this.userStorage)[nickname] == nil {
		// fetch and save locally
		newUser, err := this.getUser(nickname)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
		(*this.userStorage)[nickname] = newUser
	}

	return (*this.userStorage)[nickname], nil
}

// this private function call github API to get user info
func (this UserService) getUser(nickname string) (*models.User, error) {
	startMillis := time.Now().UnixNano() / int64(time.Millisecond)
	response, err := this.RestClient.
		R().
		Get(this.GithubBaseURL + "/users/" + nickname)

	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	if response.StatusCode() != http.StatusOK {
		return nil, tracerr.Errorf("invalid status code: '%d'", response.StatusCode())
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

func (this *UserService) GetSortedUsersByAssignations() []*models.User {
	userList := []*models.User{}
	for _, u := range *this.userStorage {
		userList = append(userList, u)
	}

	sort.Slice(userList[:], func(i, j int) bool {
		return userList[i].AssignedPRLines < userList[j].AssignedPRLines
	})
	return userList
}
