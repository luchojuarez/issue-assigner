package services

import (
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	env "github.com/luchojuarez/issue-assigner/environment"
	"github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
)

func TestSuccesCase(t *testing.T) {
	httpmock.Reset()

	// new service instance
	userService := NewUserService()

	simpleStringResponderForGithubGetUser("luchojuarez", `{"login": "luchojuarez"}`, 200, 500)

	user, err := userService.GetUser("luchojuarez")
	if err != nil {
		log.Printf("esto trae el error %v", err)
	}
	assert.Equal(t, "luchojuarez", user.NickName)
}

func TestInvalidJsonResponse(t *testing.T) {
	httpmock.Reset()
	// new service instance
	userService := NewUserService()

	simpleStringResponderForGithubGetUser("luchojuarez", `{"login": luchojuarez"}`, 200, 0)

	_, err := userService.GetUser("luchojuarez")
	assert.Equal(t, "invalid character 'l' looking for beginning of value", err.Error())
}

func TestInvalidApiSCResponse(t *testing.T) {
	httpmock.Reset()
	// new service instance
	userService := NewUserService()

	simpleStringResponderForGithubGetUser("luchojuarez", `{"login": "luchojuarez"}`, 404, 0)

	_, err := userService.GetUser("luchojuarez")
	assert.Equal(t, "invalid status code: '404'", err.Error())
}

func TestRestError(t *testing.T) {

	env.GetEnv().IsTestEnv = false
	// clear mocks
	httpmock.Reset()
	// new service instance
	userService := NewUserService()

	_, err := userService.GetUser("luchojuarez")
	assert.Equal(t, "Get https://api.github.com/users/luchojuarez: no responder found", err.Error())
}

func simpleStringResponderForGithubGetUser(user, responseBody string, statusCode int, responseLag time.Duration) {
	httpmock.Reset()

	httpmock.RegisterResponder(
		"GET",
		"https://api.github.com/users/"+user,
		func(req *http.Request) (*http.Response, error) {
			time.Sleep(responseLag * time.Millisecond)
			resp := httpmock.NewStringResponse(statusCode, responseBody)

			return resp, nil
		})
}

var _ = ginkgo.AfterEach(func() {
	log.Printf("After each")

	httpmock.DeactivateAndReset()
})
