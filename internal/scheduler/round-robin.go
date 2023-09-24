package scheduler

import (
	"errors"
	"sync"

	"github.com/The-Flash/code-racer/internal/manifest"
	"github.com/The-Flash/code-racer/internal/runtime_manager"
	"github.com/docker/docker/api/types"
)

// RoundRobinScheduler round robin scheduling algorithm
type RoundRobinScheduler struct {
	r                     *runtime_manager.RuntimeManager
	l                     sync.Mutex
	lastUsedExecutorIndex int
}

// NewRoundRobinScheduler create a new instance
// of the round robin scheduler
func NewRoundRobinScheduler(r *runtime_manager.RuntimeManager) *RoundRobinScheduler {
	return &RoundRobinScheduler{
		r: r,
	}
}

// GetActiveContainers get the current containers active for the specified
// runtime
func (s *RoundRobinScheduler) GetActiveContainers(rt *manifest.ManifestRuntime) ([]types.Container, error) {
	return s.r.GetContainersForRuntime(rt)
}

// GetNextExecutor uses the round robin algorithm
// To determine which container to use for the next execution
func (s *RoundRobinScheduler) GetNextExecutor(rt *manifest.ManifestRuntime) (c types.Container, err error) {
	activeContainers, err := s.GetActiveContainers(rt)
	if err != nil {
		return
	}
	numberOfActiveContainers := len(activeContainers)
	if numberOfActiveContainers == 0 {
		err = errors.New("no executors available")
		return
	}
	nextExecutorIndex := (s.lastUsedExecutorIndex + 1) % (numberOfActiveContainers)
	c = activeContainers[nextExecutorIndex]
	s.l.Lock()
	s.lastUsedExecutorIndex = nextExecutorIndex
	s.l.Unlock()
	return
}
