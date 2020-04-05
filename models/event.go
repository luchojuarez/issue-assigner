package models

import (
	"fmt"
	"time"

	"github.com/ztrue/tracerr"
)

type Event struct {
	Time    time.Time
	Type    string
	Content string
	Err     error
}

const (
	LevelError = "error"
	LevelInfo  = "info"
)

func (this *Event) GetString() string {
	return fmt.Sprintf("[%v][level:%s] %s", this.Time.Format(time.RFC3339), this.Type, this.Content)
}

func (this *Event) Print() {
	switch t := this.Type; t {
	case LevelError:
		fmt.Println(this.GetString())
		if this.Err != nil {
			tracerr.Print(this.Err)
		}
	case LevelInfo:
		fmt.Println(this.GetString())
	}
}
