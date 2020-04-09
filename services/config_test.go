package services

import (
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	env "github.com/luchojuarez/issue-assigner/environment"
	"github.com/stretchr/testify/assert"
	"github.com/ztrue/tracerr"
)

const (
	patern            = "localhost:%d/%s/foo"
	jsonResourcesPath = "../resources/test/json/"
)

func TestSuccess(t *testing.T) {
	init := time.Now().UnixNano() / int64(time.Millisecond)
	httpmock.Reset()
	env.GetEnv().ClearUserStorage()

	mockConfigSuccessCase()

	config, _ := load("https://api.github.com", jsonResourcesPath+"config_test.json", "sync")

	assert.Equal(t, int(2), len(config.RepoNames))
	// check repo names
	assert.Equal(t, "luchojuarez/crypto", config.RepoNames[0])
	assert.Equal(t, "luchojuarez/issue-assigner", config.RepoNames[1])

	// check info in JSON input and api response
	assert.Equal(t, "luchojuarez", config.UsersNicknames[0])
	assert.Equal(t, "luchojuarez", config.Users[0].NickName)

	assert.Equal(t, "luchojuarez2", config.UsersNicknames[1])
	assert.Equal(t, "luchojuarez2", config.Users[1].NickName)

	// test teken
	assert.Equal(t, "token foo", config.GithubToken)
	log.Printf("total time sync %d", time.Now().UnixNano()/int64(time.Millisecond)-init)
}

func TestSuccessAsync(t *testing.T) {
	init := time.Now().UnixNano() / int64(time.Millisecond)

	httpmock.Reset()
	env.GetEnv().ClearUserStorage()

	mockConfigSuccessCase()

	config, _ := load("https://api.github.com", jsonResourcesPath+"config_test.json", "async")

	assert.Equal(t, int(2), len(config.RepoNames))
	// check repo names
	assert.Equal(t, "luchojuarez/crypto", config.RepoNames[0])
	assert.Equal(t, "luchojuarez/issue-assigner", config.RepoNames[1])

	// check info in JSON input and api response
	assert.Equal(t, "luchojuarez", config.UsersNicknames[0])
	assert.Equal(t, "luchojuarez", config.Users[0].NickName)

	assert.Equal(t, "luchojuarez2", config.UsersNicknames[1])
	assert.Equal(t, "luchojuarez2", config.Users[1].NickName)

	// test teken
	assert.Equal(t, "token foo", config.GithubToken)
	log.Printf("tiempo async %d", time.Now().UnixNano()/int64(time.Millisecond)-init)

}

func TestFileNotFound(t *testing.T) {
	_, err := load("https://api.github.com", "foo.json", "async")
	assert.Equal(t, "open foo.json: no such file or directory", tracerr.Unwrap(err).Error())
}

func TestInvalidJsonFile(t *testing.T) {
	_, err := load("https://api.github.com", jsonResourcesPath+"invalid.json", "async")
	assert.Equal(t, "invalid character 'n' looking for beginning of object key string", tracerr.Unwrap(err).Error())
}

func TestLoadUserError(t *testing.T) {
	_, err := load("https://api.github.com", jsonResourcesPath+"user_no_exist.json", "async")
	if err == nil {
		tracerr.Print(err)
		assert.Fail(t, "Error ist not null")
	}
	assert.Equal(t, "Get https://api.github.com/users/unknow_user: no responder found", tracerr.Unwrap(err).Error())
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

func mockConfigSuccessCase() {
	simpleStringResponderForGithubGetUser("luchojuarez", `{"login": "luchojuarez"}`, 200, 400)
	simpleStringResponderForGithubGetUser("luchojuarez2", `{"login": "luchojuarez2"}`, 200, 400)
	simpleStringResponderForGithubGetUser("user3", `{"login": "user3"}`, 200, 400)
	simpleStringResponderForPrSearch("luchojuarez/crypto", `[{"number": 1},{"number": 2},{"number": 3},{"number": 4}]`, 200, 200)
	simpleStringResponderForGetPR(1, "luchojuarez/crypto", `{"number": 1,"title":"Title 1","body":"description 1","assignees":[{"login":"luchojuarez"}],"user":{"login":"luchojuarez"},"commits": 2,"additions": 353,  "deletions": 2}`, 200, 400)
	simpleStringResponderForGetPR(2, "luchojuarez/crypto", `{"number": 1,"title":"Title 2","body":"description 2","assignees":[{"login":"luchojuarez2"}],"user":{"login":"luchojuarez2"},"commits": 2,"additions": 7,  "deletions": 89}`, 200, 300)

	simpleStringResponderForPrSearch("luchojuarez/issue-assigner", `[{"number": 3},{"number": 8}]`, 200, 200)
	simpleStringResponderForGetPR(3, "luchojuarez/issue-assigner", `{"number": 3,"title":"Title 3","body":"description 3","assignees":[{"login":"luchojuarez"}],"user":{"login":"luchojuarez"},"commits": 1,"additions": 1,  "deletions": 100}`, 200, 250)
	simpleStringResponderForGetPR(8, "luchojuarez/issue-assigner", `{"number": 8,"title":"Title 8","body":"description 8","assignees":null,"user":{"login":"luchojuarez2"},"commits": 2,"additions": 99,  "deletions": 89}`, 200, 300)
}

// oputput for times = 500000 -> 'timeFormatter: 223ms | timeStringConcat: 67ms'
// func TestBenckmarck(t *testing.T) {
// 	times := 500000
// 	timeFormatter := benckFormatter(times, t)
// 	timeStringConcat := benckStringConcat(times, t)
// 	fmt.Printf("timeFormatter: %dms | timeStringConcat: %dms\n", timeFormatter, timeStringConcat)
// }
//
// func benckFormatter(times int, t *testing.T) int64 {
// 	startTime := time.Now().UnixNano() / int64(time.Millisecond)
// 	i := 0
// 	for {
// 		if i == times {
// 			break
// 		}
// 		tempString := fmt.Sprintf(patern, i, "foo")
// 		_ = len(tempString)
// 		i++
// 	}
//
// 	return (time.Now().UnixNano() / int64(time.Millisecond)) - startTime
// }
//
// func benckStringConcat(times int, t *testing.T) int64 {
// 	startTime := time.Now().UnixNano() / int64(time.Millisecond)
// 	i := 0
// 	for {
// 		if i == times {
// 			break
// 		}
// 		tempString := "localhost:" + strconv.Itoa(i) + "/" + "foo" + "/foo"
// 		_ = len(tempString)
// 		i++
// 	}
// 	return (time.Now().UnixNano() / int64(time.Millisecond)) - startTime
// }
