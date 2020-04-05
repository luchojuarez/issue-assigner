package dao

import (
	env "github.com/luchojuarez/issue-assigner/environment"
	"github.com/luchojuarez/issue-assigner/models"

	"github.com/ztrue/tracerr"
)

type LocalPrDao struct {
}

func NewLocalPrDao() LocalPrDao {
	return LocalPrDao{}
}

func (this LocalPrDao) GetPr(repoName string, prNumber int) (*models.PR, error) {
	for _, pr := range (*env.GetEnv().GetPrStorage())[repoName] {
		if pr.Number == prNumber {
			return pr, nil
		}
	}
	return nil, nil
}

func (this LocalPrDao) GetPrByRepo(repoName string) ([]*models.PR, error) {
	return (*env.GetEnv().GetPrStorage())[repoName], nil
}

func (this LocalPrDao) SavePr(repoName string, pr *models.PR) error {
	if pr == nil {
		return tracerr.New("Pr cant be null")
	}
	prStorage := *env.GetEnv().GetPrStorage()

	prStorage[repoName] = append(prStorage[repoName], pr)

	return nil
}
