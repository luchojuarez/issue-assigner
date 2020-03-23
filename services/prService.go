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

const (
	githubBaseURL = "https://api.github.com"

	getAllPrURL      = githubBaseURL + "/repos/%s/pulls?status=open"
	getPrByNumberURL = githubBaseURL + "/repos/%s/pulls/%d"
)

type PRService struct {
	RestClient       *resty.Client
	GithubBaseURL    string
	TotalRequestTime int64
}
type prSearchResult struct {
	Number int `json:"number"`
}
type prSearchResponse []*prSearchResult

func NewPRService() *PRService {
	return &PRService{
		RestClient: env.GetEnv().GetResty("PRService"),
	}
}

func (this *PRService) GetOpenPRs(fullRepoName string) ([]*models.PR, error) {
	defer this.setEndTime(time.Now())
	response, err := this.RestClient.R().
		Get(fmt.Sprintf(getAllPrURL, fullRepoName))
	if err != nil {
		return nil, err
	}
	if response.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("invalid status code: '%d'", response.StatusCode())
	}
	searchResult := prSearchResponse{}
	if err := json.Unmarshal([]byte(fmt.Sprintf("%s", response)), &searchResult); err != nil {
		return nil, err
	}

	toReturn := []*models.PR{}

	for _, number := range searchResult {
		newPr, err := this.getPrByNumber(fullRepoName, number.Number)
		if err != nil {
			return nil, err
		}
		toReturn = append(toReturn, newPr)
	}

	return toReturn, nil
}

func (this PRService) getPrByNumber(fullRepoName string, number int) (*models.PR, error) {
	newPr := models.PR{}
	defer newPr.SetEndTime(time.Now())

	response, err := this.RestClient.R().
		Get(fmt.Sprintf(getPrByNumberURL, fullRepoName, number))

	if err != nil {
		return nil, err
	}
	if response.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("invalid status code: '%d'", response.StatusCode())
	}

	if err := json.Unmarshal(([]byte(fmt.Sprintf("%s", response))), &newPr); err != nil {
		return nil, err
	}
	newPr.AssigneesSize = len(newPr.Assignees)

	return &newPr, nil
}

func (this PRService) setEndTime(initTime time.Time) {
	endMillis := time.Now().UnixNano() / int64(time.Millisecond)
	this.TotalRequestTime = endMillis - (initTime.UnixNano() / int64(time.Millisecond))
}
