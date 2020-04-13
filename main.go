package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	env "github.com/luchojuarez/issue-assigner/environment"
	"github.com/luchojuarez/issue-assigner/services"
	"github.com/ztrue/tracerr"
)

const defaultConfigFilePath = "resources/main/json/config.json"

func main() {
	var outputPath, configFile, token string
	if err := bindParams(&outputPath, &configFile, &token); err != nil {
		tracerr.Print(err)
		os.Exit(1)
	}

	defer services.PrintAndClearWhithBeginTime(outputPath, time.Now())
	assignmentService, err := services.NewAssignmentService(configFile)

	if err == nil {
		assignmentService.Run()
	} else {
		tracerr.Print(err)
	}
}

func bindParams(outputPath, configFile, token *string) error {
	argsWithoutProg := os.Args[1:]
	for i, currentArg := range argsWithoutProg {
		// check if param has a valid format
		if (strings.HasPrefix(currentArg, "-")) && i+1 < len(argsWithoutProg) {
			// check for token
			if currentArg == "-t" || currentArg == "--token" {
				*token = argsWithoutProg[i+1]
			}

			// check for output file path
			if currentArg == "-o" || currentArg == "--output" {
				*outputPath = argsWithoutProg[i+1]
			}

			// check for config file
			if currentArg == "-c" || currentArg == "--config-file" {
				*configFile = argsWithoutProg[i+1]
			}
		}
	}
	// if someone is not defined take default values
	bindToken(token)
	bindOutputPath(outputPath)
	bindConfigFile(configFile)
	return nil
}

func bindToken(token *string) {
	if *token == "" {
		if !env.GetEnv().TokenManager.HasToken() {
			log.Printf("Warning: github token not set")
		}
	} else {
		if err := env.GetEnv().TokenManager.Set(*token); err != nil {
			tracerr.Print(err)
			os.Exit(1)
		}
	}
}

func bindConfigFile(configFile *string) {
	if *configFile == "" {
		*configFile = defaultConfigFilePath
	}
}
func bindOutputPath(outputPath *string) {
	if *outputPath == "" {
		*outputPath = fmt.Sprintf("out/run_%s.log", time.Now().Format("2006-01-02_15:04:05"))
	}
}
