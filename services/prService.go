package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	env "github.com/luchojuarez/issue-assigner/environment"
	"github.com/luchojuarez/issue-assigner/models"
	"github.com/ztrue/tracerr"
)

const (
	githubBaseURL = "https://api.github.com"

	getAllPrURL      = githubBaseURL + "/repos/%s/pulls?status=open"
	getPrByNumberURL = githubBaseURL + "/repos/%s/pulls/%d"
)

type PRService struct {
	RestClient          *resty.Client
	UserServiceInstance *UserService
	GithubBaseURL       string
	TotalRequestTime    int64
	FetchedPRsByRepo    *map[string][]*models.PR
}
type prSearchResult struct {
	Number int `json:"number"`
}
type prSearchResponse []*prSearchResult

func NewPRService() *PRService {
	return &PRService{
		RestClient:          env.GetEnv().GetResty("PRService"),
		UserServiceInstance: NewUserService(),
		FetchedPRsByRepo:    env.GetEnv().GetPrStorage(),
	}
}

func (this *PRService) GetOpenPRs(fullRepoName string) ([]*models.PR, error) {
	if (*this.FetchedPRsByRepo)[fullRepoName] == nil {
		defer this.setEndTime(time.Now())
		response, err := this.
			RestClient.
			R().
			Get(fmt.Sprintf(getAllPrURL, fullRepoName))
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
		if response.StatusCode() != http.StatusOK {
			return nil, fmt.Errorf("invalid status code: '%d'", response.StatusCode())
		}
		searchResult := prSearchResponse{}
		if err := json.Unmarshal([]byte(fmt.Sprintf("%s", response)), &searchResult); err != nil {
			return nil, tracerr.Wrap(err)
		}

		toReturn := []*models.PR{}

		for _, number := range searchResult {
			newPr, err := this.getPrByNumber(fullRepoName, number.Number)
			if err != nil {
				return nil, tracerr.Wrap(err)
			}
			toReturn = append(toReturn, newPr)
		}
		(*this.FetchedPRsByRepo)[fullRepoName] = toReturn
	}
	return (*this.FetchedPRsByRepo)[fullRepoName], nil
}

func (this *PRService) getPrByNumber(fullRepoName string, number int) (*models.PR, error) {
	newPr := models.PR{}
	defer newPr.SetEndTime(time.Now())

	response, err := this.RestClient.R().
		Get(fmt.Sprintf(getPrByNumberURL, fullRepoName, number))

	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	if response.StatusCode() != http.StatusOK {
		return nil, tracerr.Errorf("invalid status code: '%d'", response.StatusCode())
	}

	if err := json.Unmarshal(([]byte(fmt.Sprintf("%s", response))), &newPr); err != nil {
		return nil, tracerr.Wrap(err)
	}
	newPr.AssigneesSize = len(newPr.Assignees)
	for _, userInterface := range newPr.Assignees {
		//cast interface to map
		userMap := userInterface.(map[string]interface{})
		//cast map value to string
		userName := userMap["login"].(string)
		fetchedUser, err := this.UserServiceInstance.GetUser(userName)
		if err != nil {
			log.Panicf("cat get user info '%s' '%v'", userName, err)
		}
		// TODO is OK this calc?
		fetchedUser.AssignedPRLines += newPr.Deletions + newPr.Additions
		newPr.AssignedUsers = append(newPr.AssignedUsers, fetchedUser)
	}

	return &newPr, nil
}

func (this *PRService) setEndTime(initTime time.Time) {
	endMillis := time.Now().UnixNano() / int64(time.Millisecond)
	this.TotalRequestTime = endMillis - (initTime.UnixNano() / int64(time.Millisecond))
}
