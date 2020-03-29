package dao

import (
	"github.com/luchojuarez/issue-assigner/models"
)

type UserDaoInterface interface {
	GetUser(nickname string) (*models.User, error)
	SaveUser(*models.User) error
	GetAllCached() *map[string]*models.User
}

type PrDaoInterface interface {
	GetPr(repoName string, prNumber int) (*models.PR, error)
	GetPrByRepo(repoName string) ([]*models.PR, error)
	SavePr(repoName string, pr *models.PR) error
}
