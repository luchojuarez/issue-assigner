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
	env.GetEnv().ClearUserStorage()

	usersMap := env.GetEnv().GetUserStorage()
	userReference := (*usersMap)["luchojuarez"]
	if userReference != nil {
		assert.Fail(t, "user found", userReference)
	}
	// new service instance
	userService := NewUserService()

	simpleStringResponderForGithubGetUser("luchojuarez", `{"login": "luchojuarez"}`, 200, 0)

	user, err := userService.GetUser("luchojuarez")
	if err != nil {
		log.Printf("esto trae el error %v", err)
	}
	assert.Equal(t, "luchojuarez", user.NickName)

	userReference = (*usersMap)["luchojuarez"]
	assert.NotNil(t, userReference, nil)
}

func TestInvalidJsonResponse(t *testing.T) {
	httpmock.Reset()
	env.GetEnv().ClearUserStorage()
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

var _ = ginkgo.AfterEach(func() {
	log.Printf("After each")

	httpmock.DeactivateAndReset()
})

func mockUserFromApi(nickname, responseBody string, statusCode int, responseLag time.Duration) {
	httpmock.RegisterResponder(
		"GET",
		"https://api.github.com/users/"+nickname,
		func(req *http.Request) (*http.Response, error) {
			time.Sleep(responseLag * time.Millisecond)
			resp := httpmock.NewStringResponse(statusCode, responseBody)
			return resp, nil
		})
}
