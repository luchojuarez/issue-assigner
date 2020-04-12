package services

import (
	"log"
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
	mockPRWhit2Reviwers()

	config, err := load("https://api.github.com", jsonResourcesPath+"config_test.json")
	log.Printf("vamo con el error")
	tracerr.Print(err)

	// check info in JSON input and api response
	assert.Equal(t, "luchojuarez", config.UsersNicknames[0])
	assert.Equal(t, "luchojuarez", config.Users[0].NickName)

	assert.Equal(t, "luchojuarez2", config.UsersNicknames[1])
	assert.Equal(t, "luchojuarez2", config.Users[1].NickName)

	// test teken
	assert.Equal(t, "token foo", config.GithubToken)
	log.Printf("total time sync %d", time.Now().UnixNano()/int64(time.Millisecond)-init)
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
		tracerr.Print(err)
		assert.Fail(t, "Error ist not null")
	}
	assert.Equal(t, "Get https://api.github.com/users/unknow_user: no responder found", tracerr.Unwrap(err).Error())
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
