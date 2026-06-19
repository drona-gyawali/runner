package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/drona-gyawali/runner/pkg/executor/utils"
	"github.com/drona-gyawali/runner/pkg/types"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	hostWorkspace := "/tmp/sandbox-test"

	err := os.MkdirAll(hostWorkspace, 0755)
	if err != nil {
		log.Fatalf("Failed to create host workspace directory: %v", err)
	}
	fmt.Printf("Host workspace ready at: %s\n", hostWorkspace)

	sucessSpec := types.ExecReq{
		SandboxId:   "test-sucess-01",
		ProjectPath: hostWorkspace,
		Image:       "ubuntu",
		Shell:       true, 
		Cmd: []string{
			"echo 'Starting build inside sandbox...'",
			"echo 'Creating a mock file in workspace...'",
			"echo 'hello from gvisor sandbox' > /tmp/output.txt",
			"cat /tmp/output.txt",
			"echo 'Step finished completely!'",
		},
	}

	fmt.Println(" Launching sandbox container and streaming real-time logs...")

	err = utils.RunSandboxEnv(ctx, sucessSpec, os.Stdout)
	
	if err != nil {
		fmt.Printf(" Test 1 failed unexpectedly: %s\n", err)
		return
	} 
	
	fmt.Println("Test 1 Result: SUCCESS (Exit Code 0 captured!)")
}