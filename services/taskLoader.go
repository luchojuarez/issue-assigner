package services

import (
	"fmt"
	"log"
	"strconv"
	"strings"
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

		this.prService.GetOpenPRsAsinc(currentRepoName, dones, errors)
	}
}

func (this *PrListTaskLoader) GetAllTask(dones chan bool, errors chan error) {
	for _, currentRepoName := range this.PrList {
		repoName := fmt.Sprintf("%s/%s", strings.Split(currentRepoName, "/")[0], strings.Split(currentRepoName, "/")[1])
		prNumber, _ := strconv.Atoi(strings.Split(currentRepoName, "/")[2])
		log.Printf("largo %s (%s/%d)", currentRepoName, repoName, prNumber)
		this.prService.GetPrByNumber(repoName, prNumber, dones, errors)
	}
}
