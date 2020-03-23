package utils

import (
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func TestSuccess(t *testing.T) {
	httpmock.Reset()

	simpleStringResponderForGithubGetUser("luchojuarez", `{"login": "luchojuarez"}`, 200, 0)
	simpleStringResponderForGithubGetUser("luchojuarez2", `{"login": "luchojuarez2"}`, 200, 0)

	config, err := load("https://api.github.com", "config_test.json")

	log.Printf("esto retorna '%v' '%v'", config, err)
}

func simpleStringResponderForGithubGetUser(user, responseBody string, statusCode int, responseLag time.Duration) {

	httpmock.RegisterResponder(
		"GET",
		"https://api.github.com/users/"+user,
		func(req *http.Request) (*http.Response, error) {
			time.Sleep(responseLag * time.Millisecond)
			resp := httpmock.NewStringResponse(statusCode, responseBody)

			return resp, nil
		})
}
