package tests

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func Test_Pipeline_Success_Lifecycle(t *testing.T) {
	wd, _ := os.Getwd()
	projectRoot := filepath.Dir(wd)
	fixtureDir := filepath.Join(projectRoot, "tests", "fixtures", "python-app")
	binaryPath := filepath.Join(projectRoot, "bin", "runner")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	
	t.Log("Purging local Docker images and cache folders to guarantee a fresh environment...")
	
	_ = exec.CommandContext(ctx, "docker", "rmi", "-f", "python:3.12-slim").Run()

	preCleanup := exec.CommandContext(ctx, "docker", "run", "--rm",
		"-v", fixtureDir+":/workspace",
		"alpine:latest",
		"rm", "-rf",
		"/workspace/.pip_packages",
		"/workspace/.pytest_cache",
		"/workspace/tests/__pycache__",
	)
	_ = preCleanup.Run()

	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", binaryPath, filepath.Join(projectRoot, "cmd", "runner", "main.go"))
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to compile runner binary: %v", err)
	}

	t.Cleanup(func() {
		os.Remove(binaryPath)
		postCleanup := exec.Command("docker", "run", "--rm",
			"-v", fixtureDir+":/workspace",
			"alpine:latest",
			"rm", "-rf",
			"/workspace/.pip_packages",
			"/workspace/.pytest_cache",
			"/workspace/tests/__pycache__",
		)
		_ = postCleanup.Run()
		
		_ = exec.Command("docker", "rmi", "-f", "python:3.12-slim").Run()
	})

	// Prepare the execution invocation context
	cmd := exec.CommandContext(ctx, binaryPath, "-config", ".runner/workflows/ci.toml")
	cmd.Dir = fixtureDir

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err := cmd.Run()

	t.Logf("STDOUT:\n%s", stdoutBuf.String())
	if err != nil {
		t.Fatalf("Pipeline execution failed!\nSTDERR:\n%s\nError: %v", stderrBuf.String(), err)
	}
}