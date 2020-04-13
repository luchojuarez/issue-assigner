package localEnvironment

import (
	"log"
	"os"
	"strings"
	"testing"

	"github.com/atotto/clipboard"
	"github.com/stretchr/testify/assert"
	"github.com/ztrue/tracerr"
)

const validToken = "40_chars_0101010101010101010101010101010"
const invalidToken = "40_char"

func TestSuccessLoadFromClipBoard(t *testing.T) {

	//success case
	mockClipboardValidToken()

	tokenLocal := &LocalToken{}
	asd, _ := clipboard.ReadAll()
	log.Printf("esto tiene el clipboard %s", asd)
	if err := tokenLocal.LoadTokenFromClipboard(); err != nil {
		tracerr.Print(err)
		assert.Fail(t, err.Error())
	}

	assert.True(t, tokenLocal.HasToken())

	tokenLocal = &LocalToken{}
	os.Setenv(envTokenName, "")
	// invalid case
	mockClipboardInvalidToken()
	if err := tokenLocal.LoadTokenFromClipboard(); err == nil {
		assert.Fail(t, "error are expected")
	}

	assert.False(t, tokenLocal.HasToken())
}

func TestFooo(t *testing.T) {
	os.Setenv(envTokenName, "")

	tokenLocal := &LocalToken{}
	// invalid case
	mockClipboardInvalidToken()
	if err := tokenLocal.LoadTokenFromClipboard(); err == nil {
		assert.Fail(t, "error are expected")
	} else {
		assert.True(t, strings.Contains(err.Error(), "invalid token"))
	}

	assert.False(t, tokenLocal.HasToken())
}

func TestMockedTokens(t *testing.T) {

	to := MockedToken{token: invalidToken}
	assert.True(t, to.HasToken())

	to.Set(validToken)
	assert.True(t, to.HasToken())

	to.token = ""

	if err := to.Set(validToken); err != nil {
		assert.Fail(t, err.Error())
	}

	if err := to.Set(invalidToken); err == nil {
		assert.Fail(t, "error expected")
	}
}

func mockClipboardValidToken() {
	clipboard.WriteAll(validToken)
}

func mockClipboardInvalidToken() {
	clipboard.WriteAll(invalidToken)
}
