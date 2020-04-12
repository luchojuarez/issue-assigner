package services

import (
	"encoding/json"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/luchojuarez/issue-assigner/dao"
	env "github.com/luchojuarez/issue-assigner/environment"
	"github.com/luchojuarez/issue-assigner/models"
	"github.com/ztrue/tracerr"
)

type UserService struct {
	RestClient    *resty.Client
	GithubBaseURL string
	dao           dao.UserDaoInterface
	lock          *sync.Mutex
}

func NewUserService0() *UserService {
	return NewUserService(dao.NewLocalUserDao())
}
func NewUserService(dao dao.UserDaoInterface) *UserService {
	return &UserService{
		RestClient:    env.GetEnv().GetResty("UserService"),
		GithubBaseURL: "https://api.github.com",
		dao:           dao,
		lock:          &sync.Mutex{},
	}
}

func (this *UserService) GetUser(nickname string) (*models.User, error) {
	user, err := this.dao.GetUser(nickname)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}
	user, err = this.getUser(nickname)
	if err != nil {
		return nil, err
	}
	err = this.dao.SaveUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (this *UserService) GetUserAsinc(nickname string, topic chan string, errors chan error) {
	if _, err := this.GetUser(nickname); err != nil {
		tracerr.Print(err)
		errors <- err
	}
	topic <- nickname
}

// this private function call github API to get user info
func (this UserService) getUser(nickname string) (*models.User, error) {
	defer TraceTime("get_user_from_api", time.Now())
	startMillis := time.Now().UnixNano() / int64(time.Millisecond)
	response, err := this.RestClient.
		R().
		SetHeader("Authorization", "token "+env.GetEnv().TokenManager.Get()).
		Get(this.GithubBaseURL + "/users/" + nickname)

	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	if response.StatusCode() != http.StatusOK {
		return nil, tracerr.Errorf("invalid status code: '%d' for resource '%s'", response.StatusCode(), this.GithubBaseURL+"/users/"+nickname)
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

//This func calculate ondeman the sorted list.
// each time this function is called, the result will be calculated
// from users in cache defined in dao.GetAllCached().
func (this *UserService) GetSortedUsersByAssignations(config *JsonConfig) []*models.User {
	userList := []*models.User{}
	for _, u := range *this.dao.GetAllCached() {
		// search if user are in config
		for _, nickname := range config.UsersNicknames {
			if nickname == u.NickName {
				userList = append(userList, u)
				break
			}
		}
	}

	sort.Slice(userList[:], func(i, j int) bool {
		return userList[i].AssignedTaskValue < userList[j].AssignedTaskValue
	})
	return userList
}
