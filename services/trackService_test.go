package services

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ztrue/tracerr"

	env "github.com/luchojuarez/issue-assigner/environment"
)

func TestSimpleTrace(t *testing.T) {
	env.GetEnv().CleanAll()
	echoFile("../out/out.log")
	TraceInfo("init...")
	TraceInfo("other")
	TraceErrorf("%s,%s %d", "format", "error", 1)
	//generate error "invalid character 'n' looking for beginning of object key string"
	load("https://api.github.com", jsonResourcesPath+"invalid.json")

	if err := PrintAndClear("out.log"); err != nil {
		tracerr.Print(err)
	}
	assert.Equal(t, 0, len(*(env.GetEnv().GetAllEvents())))

	file, _ := ioutil.ReadFile("../out/out.log")
	logAsArray := strings.Split(string(file), "\n")
	assert.Equal(t, 5, len(logAsArray))
	assert.True(t, strings.Contains(logAsArray[0], "init..."))
	assert.True(t, strings.Contains(logAsArray[1], "other"))
	assert.True(t, strings.Contains(logAsArray[2], "format,error 1"))
	assert.True(t, strings.Contains(logAsArray[3], "invalid character 'n' looking for beginning of object key string"))
	assert.Equal(t, logAsArray[4], "")
}

func echoFile(file string) {
	f, _ := os.Create(file)
	f.WriteString("")
	f.Close()
}
