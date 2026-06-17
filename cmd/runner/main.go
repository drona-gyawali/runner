package main

import (
	"log"
	"os"

	"context"
	"time"

	"fmt"

	"github.com/drona-gyawali/runner/pkg/executor/utils"
	"github.com/drona-gyawali/runner/pkg/types"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Unable to load .env file")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	hostWorkspace := "/tmp/sandbox-test"

	sucessSpec := types.ExecReq{
		SandboxId:   "test-sucess-01",
		ProjectPath: hostWorkspace,
		Image:       "alpine:latest",
		Cmd: []string{
			"echo 'Starting build inside sandbox...'",
			"echo 'Creating a mock file in workspace...'",
			"echo 'hello from gvisor sandbox' > /workspace/output.txt",
			"cat /workspace/output.txt",
			"echo 'Step finished completely!'",
		},
	}

	err = utils.RunSandboxEnv(ctx, sucessSpec, os.Stdout)
	if err != nil {
		fmt.Printf(" Test 1 failed unexpectedly: %s\n", err)
	} else {
		fmt.Println("Test 1 Result: SUCCESS (Exit Code 0 captured!)")
	}

}
