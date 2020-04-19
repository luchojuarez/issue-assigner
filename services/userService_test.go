package services

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/luchojuarez/issue-assigner/dao"
	env "github.com/luchojuarez/issue-assigner/environment"
	"github.com/luchojuarez/issue-assigner/utils"
	"github.com/stretchr/testify/assert"
	"github.com/ztrue/tracerr"
)

func TestSuccesCase(t *testing.T) {
	httpmock.Reset()
	dao := dao.NewLocalUserDao()

	// new service instance
	userService := NewUserService(dao)

	simpleStringResponderForGithubGetUser("luchojuarez", `{"login": "luchojuarez"}`, 200, 0)

	user, err := userService.GetUser("luchojuarez")
	if err != nil {
		log.Printf("esto trae el error %v", err)
	}
	assert.Equal(t, "luchojuarez", user.NickName)

	userReference, err := dao.GetUser("luchojuarez")
	tracerr.Print(err)
	assert.NotNil(t, userReference, nil)
}

func TestInvalidJsonResponse(t *testing.T) {
	httpmock.Reset()
	env.GetEnv().ClearUserStorage()
	// new service instance
	dao := dao.NewLocalUserDao()
	userService := NewUserService(dao)

	simpleStringResponderForGithubGetUser("luchojuarez", `{"login": luchojuarez"}`, 200, 0)

	_, err := userService.GetUser("luchojuarez")
	assert.Equal(t, "invalid character 'l' looking for beginning of value", err.Error())
}

func TestInvalidApiSCResponse(t *testing.T) {
	httpmock.Reset()
	// new service instance
	dao := dao.NewLocalUserDao()
	userService := NewUserService(dao)

	simpleStringResponderForGithubGetUser("luchojuarez", `{"login": "luchojuarez"}`, 404, 0)

	_, err := userService.GetUser("luchojuarez")
	assert.True(t, strings.Contains(err.Error(), "invalid status code: '404'"))
}

func TestRestError(t *testing.T) {

	env.GetEnv().IsTestEnv = false
	// clear mocks
	httpmock.Reset()
	// new service instance
	dao := dao.NewLocalUserDao()
	userService := NewUserService(dao)

	_, err := userService.GetUser("luchojuarez")
	assert.True(t, strings.Contains(err.Error(), "Get https://api.github.com/users/luchojuarez: no responder found"))
}

func TestCOncurrencyGetUser(t *testing.T) {
	simpleStringResponderForGithubGetUser("slow", `{"login": "user1"}`, 200, 500)
	simpleStringResponderForGithubGetUser("slow1", `{"login": "user1"}`, 200, 500)
	simpleStringResponderForGithubGetUser("slow2", `{"login": "user1"}`, 200, 500)
	service := NewUserService0()
	go service.GetUser("slow")
	go service.GetUser("slow1")
	go service.GetUser("slow2")
	go service.GetUser("slow")
	go service.GetUser("slow")

}

func TestSortUsers(t *testing.T) {
	// clear mocks
	httpmock.Reset()
	env.GetEnv().ClearUserStorage()
	// MOCK USERS
	simpleStringResponderForGithubGetUser("user1", `{"login": "user1"}`, 200, 0)
	simpleStringResponderForGithubGetUser("user2", `{"login": "user2"}`, 200, 0)
	simpleStringResponderForGithubGetUser("user3", `{"login": "user3"}`, 200, 0)
	simpleStringResponderForGithubGetUser("user4", `{"login": "user4"}`, 200, 0)

	// MOCK SEARCH PR
	crypto1, _ := utils.GetJsonPrFromSearch(1, "pr 1 from crypto", `{"login": "user1"}`, "description")
	crypto2, _ := utils.GetJsonPrFromSearch(2, "pr 2 from crypto", `{"login": "user2"}`, "description")
	crypto3, _ := utils.GetJsonPrFromSearch(3, "pr 3 from crypto", `{"login": "user1"}`, "description")

	issue1, _ := utils.GetJsonPrFromSearch(1, "pr 1 from issue", `{"login": "user1"}`, "description")
	issue2, _ := utils.GetJsonPrFromSearch(2, "pr 2 from issue", `{"login": "user3"}`, "description")
	issue3, _ := utils.GetJsonPrFromSearch(3, "pr 3 from issue", `{"login": "user3"}`, "description")

	simpleStringResponderForPrSearch("luchojuarez/crypto", fmt.Sprintf("[%s,%s,%s]", crypto1, crypto2, crypto3), 200, 0)
	simpleStringResponderForPrSearch("luchojuarez/issue", fmt.Sprintf("[%s,%s,%s]", issue1, issue2, issue3), 200, 0)

	// MOCK GET PR BY NUMBER
	adds := 5
	dels := 40
	crypto1, _ = utils.GetJsonPrFromGETWithAssignees(1, "pr 1 from crypto", `{"login": "user1"}`, "description", adds, dels, `[{"login": "user1"}]`)
	crypto2, _ = utils.GetJsonPrFromGET(2, "pr 2 from crypto", `{"login": "user2"}`, "description", 3, 5)
	crypto3, _ = utils.GetJsonPrFromGET(3, "pr 3 from crypto", `{"login": "user1"}`, "description", 3, 5)

	adds2 := 21
	dels2 := 3
	issue1, _ = utils.GetJsonPrFromGET(1, "pr 1 from issue", `{"login": "user1"}`, "description", 3, 5)
	issue2, _ = utils.GetJsonPrFromGETWithAssignees(2, "pr 2 from issue", `{"login": "user3"}`, "description", adds2, dels2, `[{"login": "user2"}]`)
	issue3, _ = utils.GetJsonPrFromGET(3, "pr 3 from issue", `{"login": "user3"}`, "description", 3, 5)
	simpleStringResponderForGetPR(1, "luchojuarez/crypto", crypto1, 200, 0)
	simpleStringResponderForGetPR(2, "luchojuarez/crypto", crypto2, 200, 0)
	simpleStringResponderForGetPR(3, "luchojuarez/crypto", crypto3, 200, 0)
	simpleStringResponderForGetPR(1, "luchojuarez/issue", issue1, 200, 0)
	simpleStringResponderForGetPR(2, "luchojuarez/issue", issue2, 200, 0)
	simpleStringResponderForGetPR(3, "luchojuarez/issue", issue3, 200, 0)

	service, _ := NewAssignmentService(jsonResourcesPath + "a_lot_of_users.json")
	service.Run()
	config := service.config
	dao := dao.NewLocalUserDao()
	userService := NewUserService(dao)

	list := userService.GetSortedUsersByAssignations(&config)
	assert.Equal(t, 4, len(list))
	for i, user := range list {
		if i == len(list)-1 {
			break
		}
		if user.AssignedTaskValue > list[i+1].AssignedTaskValue {
			assert.Fail(t, "list not sorted")
		}
	}
}
