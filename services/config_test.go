package services

import (
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
	httpmock.Reset()
	env.GetEnv().ClearUserStorage()

	simpleStringResponderForGithubGetUser("luchojuarez", `{"login": "luchojuarez"}`, 200, 0)
	simpleStringResponderForGithubGetUser("luchojuarez2", `{"login": "luchojuarez2"}`, 200, 0)

	config, _ := load("https://api.github.com", jsonResourcesPath+"config_test.json")

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

}

func TestFileNotFound(t *testing.T) {
	_, err := load("https://api.github.com", "foo.json")
	assert.Equal(t, "open foo.json: no such file or directory", tracerr.Unwrap(err).Error())
}

func TestInvalidJsonFile(t *testing.T) {
	_, err := load("https://api.github.com", jsonResourcesPath+"invalid.json")
	assert.Equal(t, "invalid character 'n' looking for beginning of object key string", tracerr.Unwrap(err).Error())
}

func TestLoadUserError(t *testing.T) {
	_, err := load("https://api.github.com", jsonResourcesPath+"user_no_exist.json")
	if err == nil {
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
