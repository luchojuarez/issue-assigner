package services

import (
	"fmt"
	"os"
	"strings"
	"time"

	env "github.com/luchojuarez/issue-assigner/environment"
	"github.com/luchojuarez/issue-assigner/models"
	"github.com/ztrue/tracerr"
)

const (
	logPath = "../out"
)

func TraceError(message string) {
	trace(message, models.LevelError, nil)
}

func TraceError0(err error) {
	trace(err.Error(), models.LevelError, err)
}

func TraceErrorf(format string, arguments ...interface{}) {
	TraceError(fmt.Sprintf(format, arguments...))
}

func TraceInfo(message string) {
	trace(message, models.LevelInfo, nil)
}

func PrintAndClear(logFileName string) error {
	defer env.GetEnv().ClearEventTracer()
	return printToFile(fmt.Sprintf("%s/%s.log", logPath, strings.Split(logFileName, ".")[0]))
}

func printToFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return tracerr.Wrap(err)
	}
	totalBytes := int(0)

	for _, event := range *(env.GetEnv().GetAllEvents()) {
		event.Print()
		l, err := f.WriteString(fmt.Sprintf("%s\n", event.GetString()))
		if err != nil {
			f.Close()
			return tracerr.Wrap(err)
		}
		totalBytes += l
	}

	return tracerr.Wrap(f.Close())
}

func trace(message, t string, e error) {
	newEvent := models.Event{
		Time:    time.Now(),
		Type:    t,
		Content: message,
		Err:     e,
	}
	env.GetEnv().AddEventSafe(&newEvent)
}
