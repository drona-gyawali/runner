package utils

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/drona-gyawali/runner/pkg/config"
	"github.com/drona-gyawali/runner/pkg/types"
)

const (
	SandboxMemory   = 2 * 1024 * 1024 * 1024
	SandboxNanoCPUs = 2000000000
)

func BuildExecLevels() [][]string {
	cfg := config.MustLoad()

	forward_graphs := make(map[string][]string)
	in_degree := make(map[string]int)

	for JobName := range cfg.Jobs {
		in_degree[JobName] = 0
	}

	for JobName, JobData := range cfg.Jobs {
		for _, dependency := range JobData.Needs {

			forward_graphs[string(dependency)] = append(forward_graphs[string(dependency)], JobName)

			in_degree[JobName]++
		}
	}

	ready_state := []string{}

	for jobName, DependencyCount := range in_degree {
		if DependencyCount == 0 {
			ready_state = append(ready_state, jobName)
		}
	}

	execLevels := [][]string{}

	processedCount := 0

	for len(ready_state) > 0 {
		running_state := ready_state
		execLevels = append(execLevels, running_state)

		next_state := []string{}
		for _, isReady := range running_state {
			processedCount++
			for _, dependencyState := range forward_graphs[isReady] {
				in_degree[dependencyState]--
				if in_degree[dependencyState] == 0 {
					next_state = append(next_state, dependencyState)
				}
			}
		}

		ready_state = next_state
	}

	if len(in_degree) != processedCount {
		log.Fatal("Configuration has a Circular dependency state")
	}

	return execLevels

}

func RunSandboxEnv(Ctx context.Context, CfgInitialization types.ExecReq, OutputLogStream io.Writer) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Sandbox initilization failed %s", err)
	}

	defer cli.Close()

	_, err = cli.ImageInspect(Ctx, CfgInitialization.Image)
	
	if err != nil {
		log.Printf("Image not found locally")

		reader, err := cli.ImagePull(Ctx, CfgInitialization.Image, image.PullOptions{})
		if err != nil {
			return fmt.Errorf("Failed to auto-pull image for runner %s", err)
		}

		_, _ = io.Copy(io.Discard, reader)
		reader.Close()
		log.Printf("Image Pulled Successfully")
	}

	var joinedCmd string
	for i, c := range CfgInitialization.Cmd {
		if i > 0 {
			joinedCmd += " && "
		}
		joinedCmd += c
	}
	dynamicCommand := []string{"sh", "-c", joinedCmd}
	sandboxConfig := &container.Config{
		Image:        CfgInitialization.Image,
		Cmd:          dynamicCommand,
		WorkingDir:   "/workspace",
		Tty:          false,
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
	}

	bindRepo := []string{fmt.Sprintf("%s:/workspace", CfgInitialization.ProjectPath)}
	hostConfig := &container.HostConfig{
		Runtime: "runsc",
		Binds:   bindRepo,
		Resources: container.Resources{
			Memory:   SandboxMemory,
			NanoCPUs: SandboxNanoCPUs,
		},
		CapDrop: []string{"ALL"},
	}

	containerName := fmt.Sprintf("isolated-runner-%s", CfgInitialization.SandboxId)

	_ = cli.ContainerRemove(Ctx, containerName, container.RemoveOptions{
		Force:         true,
		RemoveVolumes: true,
	})

	resp, err := cli.ContainerCreate(Ctx, sandboxConfig, hostConfig, nil, nil, containerName)
	if err != nil {
		return fmt.Errorf("Failed to provision sandbox environment %w", err)
	}

	defer func() {
		removeOpts := container.RemoveOptions{
			RemoveVolumes: true,
			RemoveLinks:   false,
			Force:         true,
		}
		err = cli.ContainerRemove(context.Background(), resp.ID, removeOpts)
		if err != nil {
			log.Printf("Failed to delete resources %s", err)
		}
		_, err = cli.ImageRemove(Ctx, CfgInitialization.Image , image.RemoveOptions{
			Force: true,
			PruneChildren: true,
		})
		if err != nil {
			log.Printf("Unable to remove images from system")
		}
	}()

	err = cli.ContainerStart(Ctx, resp.ID, container.StartOptions{})
	if err != nil {
		return fmt.Errorf("Failed to boot sandbox %w", err)
	}

	logOpts := container.LogsOptions{ShowStdout: true, ShowStderr: true, Follow: true, Timestamps: true}
	logReader, err := cli.ContainerLogs(Ctx, resp.ID, logOpts)
	if err == nil {
		_, _ = io.Copy(OutputLogStream, logReader)
		logReader.Close()
	}

	statusCh, errorCh := cli.ContainerWait(Ctx, resp.ID, container.WaitConditionNotRunning)

	select {
	case err := <- errorCh:
		if err != nil {
			return fmt.Errorf("Unhandled termination error occured in sanbox execution %w", err)
		}
	case success := <- statusCh:
		if success.StatusCode != 0{
			return fmt.Errorf("Sandbox return non-zero termination error %d", success.StatusCode)
		}
	}


	return nil
}
