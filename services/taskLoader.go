package services

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/luchojuarez/issue-assigner/models"
)

type TaskLoader interface {
	// returns number of task will be loaded
	GetAllTask(dones chan bool, errors chan error)

	//
	GetTotalTask() int
}

type AllPrTaskLoader struct {
	RepoNames []string
	prService *PRService
	source    *models.TaskSource
}
type PrListTaskLoader struct {
	PrList    []string
	prService *PRService
}

func (this *AllPrTaskLoader) GetTotalTask() int {
	return len(this.RepoNames)
}

func (this *PrListTaskLoader) GetTotalTask() int {
	return len(this.PrList)
}

func (this *AllPrTaskLoader) GetAllTask(dones chan bool, errors chan error) {
	for _, currentRepoName := range this.RepoNames {
		log.Printf("largo %s", currentRepoName)
		repo := models.Repo{FullName: currentRepoName}

		this.prService.GetOpenPRsAsinc(&repo, this.source, dones, errors)
	}
}

func (this *PrListTaskLoader) GetAllTask(dones chan bool, errors chan error) {
	for _, currentRepoName := range this.PrList {
		repoName := fmt.Sprintf("%s/%s", strings.Split(currentRepoName, "/")[0], strings.Split(currentRepoName, "/")[1])
		prNumber, _ := strconv.Atoi(strings.Split(currentRepoName, "/")[2])
		log.Printf("largo %s (%s/%d)", currentRepoName, repoName, prNumber)
		repo := models.Repo{FullName: repoName}
		this.prService.GetPrByNumber(&repo, prNumber, dones, errors)
	}
}
