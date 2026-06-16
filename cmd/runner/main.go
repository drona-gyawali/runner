package main

import (
	"log"

	"github.com/drona-gyawali/runner/pkg/executor/utils"
	"github.com/joho/godotenv"
)


func main () {
	err:= godotenv.Load()
	if err != nil {
		log.Fatalf("Unable to load .env file")
	}
	runConfig := utils.BuildExecLevels()

	for _, levels := range runConfig {
		log.Print(levels)
	}
}