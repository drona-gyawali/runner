package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/drona-gyawali/runner/pkg/config"
	"github.com/drona-gyawali/runner/pkg/executor"
	"github.com/joho/godotenv"
)



func main() {
	_  = godotenv.Load()


	configFlag := flag.String("config", "", "Path to the configuration file")
	flag.Parse()

	configPath := *configFlag
	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	if configPath == "" {
		log.Fatal("Unable to detect configuration file")
	}

	os.Setenv("CONFIG_PATH", configPath)

	workflow := config.MustLoad()
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Unable to detect current directory %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	err = executor.PipelineRunner(ctx, *workflow, currentDir)
	
	if err != nil {
		log.Fatalf("Pipleline Execution terminated %s", err)
	}


	log.Println("Pipleline execution completed with Exit code 0")
}