package scheduler

import (
	"github.com/The-Flash/code-racer/internal/manifest"
	"github.com/The-Flash/code-racer/internal/runtime_manager"
	"github.com/docker/docker/api/types"
)

// Scheduler determines which container to use
// in executions
type Scheduler struct {
	Random     *RandomScheduler
	RoundRobin *RoundRobinScheduler
}

// Load loads scheduling algorithms for each runtime specified
// in the manifest
func Load(m []manifest.ManifestRuntime, rm *runtime_manager.RuntimeManager) map[string]Scheduler {
	r := make(map[string]Scheduler)
	for _, runtime := range m {
		r[runtime.Language] = Scheduler{
			Random:     NewRandomScheduler(rm),
			RoundRobin: NewRoundRobinScheduler(rm),
		}
	}
	return r
}

// Algorithm interface for scheduling algorithms
type Algorithm interface {
	// Get the active containers for the runtime
	GetActiveContainers(rt *manifest.ManifestRuntime) ([]types.Container, error)

	// Get the next container to run the execution
	GetNextExecutor(rt *manifest.ManifestRuntime) (types.Container, error)
}
