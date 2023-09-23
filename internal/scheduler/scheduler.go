package scheduler

import (
	"errors"
	"math/rand"
	"sync"

	"github.com/The-Flash/code-racer/internal/manifest"
	"github.com/The-Flash/code-racer/internal/runtime_manager"
	"github.com/docker/docker/api/types"
)

type Scheduler struct {
	Random     *RandomScheduler
	RoundRobin *RoundRobinScheduler
}

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

// Scheduler determines which container to use
// in executions
type Algorithm interface {
	// Get the active containers for the runtime
	GetActiveContainers(rt *manifest.ManifestRuntime) ([]types.Container, error)

	// Get the next container to run the execution
	GetNextExecutor(rt *manifest.ManifestRuntime) (types.Container, error)
}

type RandomScheduler struct {
	r *runtime_manager.RuntimeManager
}

func NewRandomScheduler(r *runtime_manager.RuntimeManager) *RandomScheduler {
	return &RandomScheduler{
		r: r,
	}
}

func (s *RandomScheduler) GetActiveContainers(rt *manifest.ManifestRuntime) ([]types.Container, error) {
	return s.r.GetContainersForRuntime(rt)
}

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

type RoundRobinScheduler struct {
	r                     *runtime_manager.RuntimeManager
	l                     sync.Mutex
	lastUsedExecutorIndex int
}

func NewRoundRobinScheduler(r *runtime_manager.RuntimeManager) *RoundRobinScheduler {
	return &RoundRobinScheduler{
		r: r,
	}
}

func (s *RoundRobinScheduler) GetActiveContainers(rt *manifest.ManifestRuntime) ([]types.Container, error) {
	return s.r.GetContainersForRuntime(rt)
}

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
