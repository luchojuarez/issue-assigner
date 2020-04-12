package localEnvironment

import (
	"os"

	"github.com/ztrue/tracerr"
)

const envTokenName = "GITHUB_USER_TOKEN"

// interface
type TokenManager interface {
	Get() string
	Set(token string)
	HasToken() error
}

// implementation, save token in environment variables
type LocalToken struct{}

// mock implementation, just use por test
type MockedToken struct {
	token string
}

func (this *LocalToken) Get() string {
	return os.Getenv(envTokenName)
}
func (this *LocalToken) Set(token string) {
	os.Setenv(envTokenName, token)
}
func (this *LocalToken) HasToken() error {
	if this.Get() == "" {
		return tracerr.New("can't access to token")
	}
	return nil
}

func (this *MockedToken) Get() string {
	return this.token
}
func (this *MockedToken) Set(token string) {
	this.token = token
}
func (this *MockedToken) HasToken() error {
	if this.Get() == "" {
		return tracerr.New("can't access to token")
	}
	return nil
}
