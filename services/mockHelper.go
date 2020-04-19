package services

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/ztrue/tracerr"

	env "github.com/luchojuarez/issue-assigner/environment"
)

func mockPRWhit2Reviwers() {
	simpleStringResponderForPrSearch("luchojuarez/crypto", `[{"number": 1},{"number": 2},{"number": 3},{"number": 4}]`, 200, 0)
	simpleStringResponderForGetPR(3, "luchojuarez/crypto", `{"number": 3,"title":"Title 3 no more reviewers","body":"description 1","assignees":[{"login":"luchojuarez"},{"login":"luchojuarez2"},{"login":"user3"}],"user":{"login":"luchojuarez"},"commits": 2,"additions": 353,  "deletions": 18}`, 200, 0)
	simpleStringResponderForGetPR(4, "luchojuarez/crypto", `{"number": 4,"title":"Title 4 no more reviewers","body":"description 1","assignees":[{"login":"luchojuarez2"},{"login":"user3"}],"user":{"login":"luchojuarez"},"commits": 2,"additions": 353,  "deletions": 18}`, 200, 0)
	simpleStringResponderForGetPR(7, "user1/foo", `{"number": 4,"title":"Title 4 no more reviewers","body":"description 1","user":{"login":"luchojuarez"},"commits": 2,"additions": 353,  "deletions": 18}`, 200, 0)
	simpleStringResponderForGetPR(10, "user2/bar", `{"number": 4,"title":"Title 4 no more reviewers","body":"description 1","user":{"login":"luchojuarez"},"commits": 2,"additions": 353,  "deletions": 18}`, 200, 0)
}

func mockConfigSuccessCase() {
	env.GetEnv().TokenManager.Set("40_chars_0101010101010101010101010101010")
	simpleStringResponderForGithubGetUser("luchojuarez", `{"login": "luchojuarez"}`, 200, 400)
	simpleStringResponderForGithubGetUser("luchojuarez2", `{"login": "luchojuarez2"}`, 200, 400)
	simpleStringResponderForGithubGetUser("user3", `{"login": "user3"}`, 200, 400)
	simpleStringResponderForPrSearch("luchojuarez/crypto", `[{"number": 1},{"number": 2},{"number": 3},{"number": 4}]`, 200, 200)
	simpleStringResponderForGetPR(7, "user1/foo", `{"number": 4,"title":"Title 4 no more reviewers","body":"description 1","user":{"login":"luchojuarez"},"commits": 2,"additions": 353,  "deletions": 18}`, 200, 0)
	simpleStringResponderForGetPR(1, "luchojuarez/crypto", `{"number": 1,"title":"Title 1","body":"description 1","assignees":[{"login":"luchojuarez"}],"user":{"login":"luchojuarez"},"commits": 2,"additions": 353,  "deletions": 2, "labels":["needed"]}`, 200, 400)
	simpleStringResponderForGetPR(2, "luchojuarez/crypto", `{"number": 2,"title":"Title 2","body":"description 2","assignees":[{"login":"luchojuarez2"}],"user":{"login":"luchojuarez2"},"commits": 2,"additions": 7,  "deletions": 89, "labels":["needed"]}`, 200, 300)
	simpleStringResponderForGetPR(3, "luchojuarez/crypto", `{"number": 3,"title":"Title 3 no more reviewers","body":"description 1","assignees":[{"login":"luchojuarez"},{"login":"luchojuarez2"},{"login":"user3"}],"user":{"login":"luchojuarez"},"commits": 2,"additions": 353,  "deletions": 18,"labels":["needed"]}`, 200, 0)
	simpleStringResponderForGetPR(4, "luchojuarez/crypto", `{"number": 4,"title":"Excluded","body":"description 2","assignees":[{"login":"luchojuarez2"}],"user":{"login":"luchojuarez2"},"commits": 2,"additions": 7,  "deletions": 89, "labels":["exclude_me_please"]}`, 200, 300)
	simpleStringResponderForGetPR(10, "user2/bar", `{"number": 4,"title":"Title 4 no more reviewers","body":"description 1","user":{"login":"luchojuarez"},"commits": 2,"additions": 353,  "deletions": 18}`, 200, 0)

	simpleStringResponderForPrSearch("luchojuarez/issue-assigner", `[{"number": 3},{"number": 8}]`, 200, 200)
	simpleStringResponderForGetPR(3, "luchojuarez/issue-assigner", `{"number": 3,"title":"Title 3","body":"description 3","assignees":[{"login":"luchojuarez"}],"user":{"login":"luchojuarez"},"commits": 1,"additions": 1,  "deletions": 100,"labels":["needed"]}`, 200, 250)
	simpleStringResponderForGetPR(8, "luchojuarez/issue-assigner", `{"number": 8,"title":"Title 8","body":"description 8","assignees":null,"user":{"login":"luchojuarez2"},"commits": 2,"additions": 99,  "deletions": 89,"labels":[]}`, 200, 300)
}

func simpleStringResponderForGithubGetUser(user, responseBody string, statusCode int, responseLag time.Duration) {

	httpmock.RegisterResponder(
		"GET",
		"https://api.github.com/users/"+user,
		func(req *http.Request) (*http.Response, error) {
			if err := checkToken(req); err != nil {
				return nil, err
			}
			time.Sleep(responseLag * time.Millisecond)
			resp := httpmock.NewStringResponse(statusCode, responseBody)

			return resp, nil
		})
}

func simpleStringResponderForPrSearch(repoFullName, responseBody string, statusCode int, responseLag time.Duration) {
	httpmock.RegisterResponder(
		"GET",
		"https://api.github.com/repos/"+repoFullName+"/pulls?status=open",
		func(req *http.Request) (*http.Response, error) {
			if err := checkToken(req); err != nil {
				return nil, err
			}
			time.Sleep(responseLag * time.Millisecond)
			resp := httpmock.NewStringResponse(statusCode, responseBody)
			return resp, nil
		})
}

func simpleStringResponderForGetPR(number int, repoFullName, responseBody string, statusCode int, responseLag time.Duration) {
	httpmock.RegisterResponder(
		"GET",
		"https://api.github.com/repos/"+repoFullName+"/pulls/"+fmt.Sprintf("%d", number),
		func(req *http.Request) (*http.Response, error) {
			if err := checkToken(req); err != nil {
				return nil, err
			}
			time.Sleep(responseLag * time.Millisecond)
			resp := httpmock.NewStringResponse(statusCode, responseBody)
			return resp, nil
		})
}

func mockUserFromApi(nickname, responseBody string, statusCode int, responseLag time.Duration) {
	httpmock.RegisterResponder(
		"GET",
		"https://api.github.com/users/"+nickname,
		func(req *http.Request) (*http.Response, error) {
			if err := checkToken(req); err != nil {
				return nil, err
			}
			time.Sleep(responseLag * time.Millisecond)
			resp := httpmock.NewStringResponse(statusCode, responseBody)
			return resp, nil
		})
}

func assertNil(t *testing.T, object interface{}, msgAndArgs ...interface{}) {
	if !isNil(object) {
		assert.Fail(t, "Expected value must be nil.", msgAndArgs...)
	}
	return
}

func isNil(object interface{}) bool {
	if object == nil {
		return true
	}

	value := reflect.ValueOf(object)
	kind := value.Kind()
	if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
		return true
	}

	return false
}

func checkToken(req *http.Request) error {
	if req.Header["Authorization"] == nil {
		return tracerr.New("no tiene token")
	}
	header := req.Header["Authorization"][0]
	token := strings.Split(header, " ")[1]
	if err := env.ValidateToken(token); err != nil {
		return TraceError0(err)
	}
	return nil
}
