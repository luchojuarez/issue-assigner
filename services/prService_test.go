package services

import (
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"

	env "github.com/luchojuarez/issue-assigner/environment"
	"github.com/luchojuarez/issue-assigner/models"
)

func TestGetAllPRSuccesCase(t *testing.T) {
	httpmock.Reset()
	env.GetEnv().ClearUserStorage()

	// new service instance
	service := NewPRService()

	simpleStringResponderForPrSearch("luchojuarez/issuer", `[{"number": 1},{"number": 2}]`, 200, 0)
	simpleStringResponderForGithubGetUser("luchojuarez", `{"login": "luchojuarez"}`, 200, 0)

	simpleStringResponderForGetPR(1, "luchojuarez/issuer", `{"number": 1,"title":"Title 1","body":"description 1","assignees":null,"user":{"login":"luchojuarez"},"commits": 2,"additions": 353,  "deletions": 2}`, 200, 50)
	simpleStringResponderForGetPR(2, "luchojuarez/issuer", `{"number": 2,"title":"Title 2","body":"description 2","assignees":[{"login":"luchojuarez"}],"user":{"login":"luchojuarez"},"commits": 2,"additions": 353,  "deletions": 2}`, 200, 40)

	repo := models.Repo{FullName: "luchojuarez/issuer"}
	service.GetOpenPRs(&repo, nil)

	// assert.Equal(t, 2, len(prList))
	//
	// assert.Equal(t, int(1), prList[0].Number)
	// assert.Equal(t, "Title 1", prList[0].Title)
	// assert.Equal(t, "description 1", prList[0].Body)
	// assert.Equal(t, int(0), len(prList[0].Assignees))
	// assert.Equal(t, int(353), prList[0].Additions)
	//
	// assert.Equal(t, int(2), prList[1].Number)
	// assert.Equal(t, "Title 2", prList[1].Title)
	// assert.Equal(t, "description 2", prList[1].Body)
	// assert.Equal(t, int(1), len(prList[1].Assignees))
	// assert.Equal(t, int(353), prList[1].Additions)

}

func TestInvalidApiSCResponse2(t *testing.T) {
	httpmock.Reset()
	env.GetEnv().ClearPrStorage()

	// new service instance
	service := NewPRService()

	simpleStringResponderForPrSearch("luchojuarez/issuer", `[{"number": 1},{"number": 2}]`, 500, 0)

	repo := models.Repo{FullName: "luchojuarez/issuer"}
	_, err := service.GetOpenPRs(&repo, nil)

	assert.True(t, strings.Contains(err.Error(), "invalid status code: '500'"))
}

func TestRestErrorListAll(t *testing.T) {

	// clear mocks
	httpmock.Reset()
	// new service instance
	service := NewPRService()
	repo := models.Repo{FullName: "luchojuarez/issuer"}
	_, err := service.GetOpenPRs(&repo, nil)

	assert.Equal(t, "Get https://api.github.com/repos/luchojuarez/issuer/pulls?status=open: no responder found", err.Error())
}

func TestRestErrorGetPR(t *testing.T) {

	// clear mocks
	httpmock.Reset()
	// new service instance
	service := NewPRService()

	simpleStringResponderForPrSearch("luchojuarez/issuer", `[{"number": 1},{"number": 2}]`, 500, 0)

	repo := models.Repo{FullName: "luchojuarez/issuer"}
	_, err := service.GetOpenPRs(&repo, nil)

	assert.True(t, strings.Contains(err.Error(), "invalid status code: '500'"))
}

func TestInvalidJson(t *testing.T) {
	// new service instance
	service := NewPRService()

	httpmock.Reset()
	simpleStringResponderForPrSearch("luchojuarez/issuer", `[{"number": 1,{"number": 2}]`, 500, 0)
	repo := models.Repo{FullName: "luchojuarez/issuer"}
	_, err := service.GetOpenPRs(&repo, nil)
	assert.True(t, strings.Contains(err.Error(), "invalid status code: '500'"))

	httpmock.Reset()
	simpleStringResponderForPrSearch("luchojuarez/issuer", `[{"number": 1},{"number": 2}]`, 500, 0)
	simpleStringResponderForGetPR(1, "luchojuarez/issuer", `{number": 1,"title":"Title 1","body":"description 1","assignees":null,"user":{"login":"luchojuarez"},"commits": 2,"additions": 353,  "deletions": 2}`, 200, 50)
	_, err = service.GetOpenPRs(&repo, nil)
	assert.True(t, strings.Contains(err.Error(), "invalid status code: '500'"))
}
