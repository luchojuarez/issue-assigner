package localEnvironment

import (
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	"github.com/ztrue/tracerr"
)

const envTokenName = "GITHUB_USER_TOKEN"

// interface
type TokenManager interface {
	Get() string
	Set(token string) error
	HasToken() bool
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
func (this *LocalToken) Set(token string) error {
	if err := ValidateToken(token); err != nil {
		return err
	}
	os.Setenv(envTokenName, token)
	return nil
}
func (this *LocalToken) HasToken() bool {
	return this.Get() != ""
}

func (this *MockedToken) Get() string {
	return this.token
}
func (this *MockedToken) Set(token string) error {
	if err := ValidateToken(token); err != nil {
		return err
	}
	this.token = token
	return nil
}
func (this *MockedToken) HasToken() bool {
	return this.Get() != ""
}

func ValidateToken(token string) error {
	if leng := len(token); leng != 40 {
		return tracerr.New(fmt.Sprintf("invalid token len: %d, invalid token '%s'", leng, token))
	}
	// ToDo: call https://api.github.com/applications/grants to know access
	return nil
}

func (this *LocalToken) LoadTokenFromClipboard() error {
	tokenCandidate, err := clipboard.ReadAll()
	if err != nil {
		return err
	}
	return this.Set(tokenCandidate)
}
