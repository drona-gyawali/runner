package utils

import (
	"context"
	"log"

	"github.com/drona-gyawali/runner/pkg/config"
	"github.com/drona-gyawali/runner/pkg/types"
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

func  RunSandboxEnv (Ctx context.Context, CfgInitialization types.ExecReq) error {
}