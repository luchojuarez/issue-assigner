package services

import (
	"fmt"
	"os"
	"time"

	env "github.com/luchojuarez/issue-assigner/environment"
	"github.com/luchojuarez/issue-assigner/models"
	"github.com/ztrue/tracerr"
)

func TraceError(message string) {
	trace(message, models.LevelError, nil)
}

func TraceError0(err error) error {
	trace(err.Error(), models.LevelError, err)
	return err
}

func TraceErrorf(format string, arguments ...interface{}) {
	TraceError(fmt.Sprintf(format, arguments...))
}

func TraceInfo(message string) {
	trace(message, models.LevelInfo, nil)
}

func TraceInfof(format string, arguments ...interface{}) {
	TraceInfo(fmt.Sprintf(format, arguments...))
}
func PrintAndClearWhithBeginTime(logFileName string, startTime time.Time) error {
	TraceInfof("End at (%s) total millis: %d", startTime.Format(time.ANSIC), (time.Now().UnixNano()-startTime.UnixNano())/int64(time.Millisecond))
	return PrintAndClear(logFileName)
}
func PrintAndClear(logFileName string) error {
	defer env.GetEnv().ClearEventTracer()
	return printToFile(logFileName)
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
