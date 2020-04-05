package main

import (
	"log"
	"time"

	"github.com/luchojuarez/issue-assigner/dao"
	"github.com/luchojuarez/issue-assigner/services"
	"github.com/ztrue/tracerr"
)

func main() {
	defer services.PrintAndClearWhithBeginTime("new_run", time.Now())
	assignmentService, err := services.NewAssignmentService("resources/main/json/config.json")

	if err == nil {
		assignmentService.Run()
	} else {
		tracerr.Print(err)
	}
	for nickname, user := range *(dao.NewLocalUserDao().GetAllCached()) {
		log.Printf("el user:(%s)  -> %v", nickname, user)
	}
}
