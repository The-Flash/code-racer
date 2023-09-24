package scheduler

import (
	"errors"
	"math/rand"

	"github.com/The-Flash/code-racer/internal/manifest"
	"github.com/The-Flash/code-racer/internal/runtime_manager"
	"github.com/docker/docker/api/types"
)

// RandomScheduler random scheduling algorithm
type RandomScheduler struct {
	r *runtime_manager.RuntimeManager
}

// NewRandomScheduler returns an instance of RandomScheduler
func NewRandomScheduler(r *runtime_manager.RuntimeManager) *RandomScheduler {
	return &RandomScheduler{
		r: r,
	}
}

// GetActiveContainers get the containers available for execution
func (s *RandomScheduler) GetActiveContainers(rt *manifest.ManifestRuntime) ([]types.Container, error) {
	return s.r.GetContainersForRuntime(rt)
}

// GetNextExecutor randomly selects a container from
// the available containers
func (s *RandomScheduler) GetNextExecutor(rt *manifest.ManifestRuntime) (c types.Container, err error) {
	activeContainers, err := s.GetActiveContainers(rt)
	if err != nil {
		return
	}
	numberOfActiveContainers := len(activeContainers)
	if numberOfActiveContainers == 0 {
		err = errors.New("no executors available")
		return
	}
	selectedIndex := rand.Intn(numberOfActiveContainers)
	c = activeContainers[selectedIndex]
	return
}
