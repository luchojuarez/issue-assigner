package dao

import (
	env "github.com/luchojuarez/issue-assigner/environment"
	"github.com/luchojuarez/issue-assigner/models"

	"github.com/ztrue/tracerr"
)

type LocalUserDao struct {
}

func NewLocalUserDao() LocalUserDao {
	return LocalUserDao{}
}

func (this LocalUserDao) GetUser(nickname string) (*models.User, error) {
	return (*env.GetEnv().GetUserStorage())[nickname], nil
}

func (this LocalUserDao) SaveUser(user *models.User) error {
	if user == nil {
		return tracerr.New("User cant be null")
	}
	nickname := (*user).NickName
	(*env.GetEnv().GetUserStorage())[nickname] = user
	return nil
}

func (this LocalUserDao) GetAllCached() *map[string]*models.User {
	return env.GetEnv().GetUserStorage()
}
