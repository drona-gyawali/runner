package executor

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/drona-gyawali/runner/pkg/executor/utils"
	"github.com/drona-gyawali/runner/pkg/types"
	"github.com/google/uuid"
)

func PipelineRunner(Ctx context.Context, Jobs types.Jobs, ProjectPath string) error {
	execLevels := utils.BuildExecLevels(Jobs)

	log.Printf("CI pipeline has been started to execute jobs")

	// we are making developers explicitly configured images. Because we donot want to maintain
	// pre-baked  > approx 4-8GB of file consist of all images. I personally think this is light
	// weight and good for initial use.
	tarImg := Jobs.Image

	for i, level := range execLevels {
		var wg sync.WaitGroup
		errorChan := make(chan error, len(level))
		for _, jobName := range level {
			Jconfig := Jobs.Jobs[jobName]
			sandboxConfig := types.ExecReq{
				SandboxId:   uuid.New().String(),
				ProjectPath: ProjectPath,
				Image:       tarImg,
				Shell:       Jobs.Shell,
				Cmd:         []string{Jconfig.Command},
			}

			wg.Add(1)
			go func(spec types.ExecReq, name string) {

				defer wg.Done()
				log.Printf("[%s] Spawning isolated gVisor container runtime...", name)
				err := utils.RunSandboxEnv(Ctx, spec, os.Stdout)
				if err != nil {
					errorChan <- fmt.Errorf("Workflow %s failed %w", name, err)
				}
				log.Printf("[%s] Container execution completed cleanly.", name)
			}(sandboxConfig, jobName)
		}

		wg.Wait()
		close(errorChan)

		for err := range errorChan {
			if err != nil {
				return fmt.Errorf("Workflow saw error at level %d due to %w", i, err)
			}
		}
	}
	log.Print("Workflow execution has been completed")
	return nil
}
