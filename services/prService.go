package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/luchojuarez/issue-assigner/dao"
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
	dao                 dao.PrDaoInterface
}
type prSearchResult struct {
	Number int `json:"number"`
}
type prSearchResponse []*prSearchResult

func NewPRService() *PRService {
	return &PRService{
		RestClient:          env.GetEnv().GetResty("PRService"),
		UserServiceInstance: NewUserService0(),
		dao:                 dao.NewLocalPrDao(),
	}
}
func NewPRService0(userDao dao.UserDaoInterface, prDao dao.PrDaoInterface) *PRService {
	return &PRService{
		RestClient:          env.GetEnv().GetResty("PRService"),
		UserServiceInstance: NewUserService(userDao),
		dao:                 prDao,
	}
}

func (this *PRService) GetOpenPRsAsinc(fullRepoName string, topic chan bool, errors chan error) {
	if _, err := this.GetOpenPRs(fullRepoName); err != nil {
		tracerr.Print(err)
		errors <- err
	}
	topic <- true
}

func (this *PRService) GetOpenPRs(fullRepoName string) ([]*models.PR, error) {
	prList, err := this.dao.GetPrByRepo(fullRepoName)
	if err != nil {
		return nil, err
	}
	if prList != nil {
		return prList, nil
	}
	defer this.setEndTime(time.Now())
	req := this.RestClient.R()
	if env.GetEnv().TokenManager.HasToken() {
		req = req.SetHeader("Authorization", "token "+env.GetEnv().TokenManager.Get())
	} else {
		TraceInfo("Not tokent set")
	}

	response, err := req.Get(fmt.Sprintf(getAllPrURL, fullRepoName))
	if err != nil {
		return nil, TraceError0(tracerr.Wrap(err))
	}
	if response.StatusCode() == http.StatusForbidden {
		return nil, tracerr.Errorf("Request over cuota, check github token, resource '%s'", fmt.Sprintf(getAllPrURL, fullRepoName))
	}
	if response.StatusCode() != http.StatusOK {
		return nil, TraceError0(tracerr.Errorf("invalid status code: '%d' for resource '%s'", response.StatusCode(), fmt.Sprintf(getAllPrURL, fullRepoName)))
	}
	searchResult := prSearchResponse{}
	if err := json.Unmarshal([]byte(fmt.Sprintf("%s", response)), &searchResult); err != nil {
		return nil, TraceError0(tracerr.Wrap(err))
	}

	toReturn := []*models.PR{}

	for _, number := range searchResult {
		newPr, err := this.getPrByNumber(fullRepoName, number.Number)
		if err != nil {
			return nil, TraceError0(tracerr.Wrap(err))
		}
		toReturn = append(toReturn, newPr)
		this.dao.SavePr(fullRepoName, newPr)
	}
	return toReturn, nil
}

func (this *PRService) GetPrByNumber(fullRepoName string, number int, dones chan bool, errors chan error) {
	pr, err := this.dao.GetPr(fullRepoName, number)
	if err != nil {
		errors <- err
		return
	}
	if pr != nil {
		log.Printf("lo encontre %s/%d", fullRepoName, number)
		dones <- true
		return
	}
	pr, err = this.getPrByNumber(fullRepoName, number)
	if err != nil {
		errors <- err
		return
	}
	this.dao.SavePr(fullRepoName, pr)
	dones <- true
	return
}

func (this *PRService) getPrByNumber(fullRepoName string, number int) (*models.PR, error) {
	newPr := models.PR{}
	defer newPr.SetEndTime(time.Now())
	defer TraceTime("getPrByNumber", time.Now())

	req := this.RestClient.R()
	if env.GetEnv().TokenManager.HasToken() {
		req = req.SetHeader("Authorization", "token "+env.GetEnv().TokenManager.Get())
	} else {
		TraceInfo("Not tokent set")
	}
	response, err := req.Get(fmt.Sprintf(getPrByNumberURL, fullRepoName, number))

	if err != nil {
		return nil, TraceError0(tracerr.Wrap(err))
	}
	if response.StatusCode() != http.StatusOK {
		return nil, TraceError0(tracerr.Errorf("invalid status code: '%d' for resource '%s'", response.StatusCode(), fmt.Sprintf(getAllPrURL, fullRepoName)))
	}

	if err := json.Unmarshal(([]byte(fmt.Sprintf("%s", response))), &newPr); err != nil {
		return nil, TraceError0(tracerr.Wrap(err))
	}
	newPr.AssigneesSize = len(newPr.Assignees)
	if len(newPr.Body) > 200 {
		newPr.Body = newPr.Body[:200]
	}
	for _, userInterface := range newPr.Assignees {
		//cast interface to map
		userMap := userInterface.(map[string]interface{})
		//cast map value to string
		userName := userMap["login"].(string)
		fetchedUser, err := this.UserServiceInstance.GetUser(userName)
		if err != nil {
			return nil, TraceError0(tracerr.New(fmt.Sprintf("cat get user info '%s' '%v'", userName, err)))
		}
		TraceInfof("OLD assignation found. repo:'%s', PR(%d) '%s' to user '%s'", fullRepoName, newPr.Number, newPr.Body, fetchedUser.NickName)
		fetchedUser.AssingIssue(&newPr)
		newPr.AssignedUsers = append(newPr.AssignedUsers, fetchedUser)
	}

	return &newPr, nil
}

func (this *PRService) setEndTime(initTime time.Time) {
	endMillis := time.Now().UnixNano() / int64(time.Millisecond)
	this.TotalRequestTime = endMillis - (initTime.UnixNano() / int64(time.Millisecond))
}
