package execution

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	cappedbuffer "github.com/The-Flash/code-racer/internal"
	"github.com/The-Flash/code-racer/internal/config"
	"github.com/The-Flash/code-racer/internal/file_system"
	"github.com/The-Flash/code-racer/internal/manifest"
	"github.com/The-Flash/code-racer/internal/names"
	"github.com/The-Flash/code-racer/internal/runtime_manager"
	"github.com/The-Flash/code-racer/internal/scheduler"
	"github.com/The-Flash/code-racer/pkg/models"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/google/uuid"
	"github.com/sarulabs/di/v2"
)

type ExecutionConfig struct {
	// Task Exeuction Id
	ExecutionId string
	// Entrypoint specified by request
	EntryPoint string
	// Language runtime for task
	Runtime *manifest.ManifestRuntime
}

type Executor struct {
	mfest      *manifest.Manifest
	fp         *file_system.FileProvider
	config     *config.Config
	rm         *runtime_manager.RuntimeManager
	cli        *client.Client
	schedulers map[string]scheduler.Scheduler
}

// Setup the executor struct
func (r *Executor) Setup(ctn di.Container) {
	r.fp = ctn.Get(names.DiFileProvider).(*file_system.FileProvider)
	r.config = ctn.Get(names.DiConfigProvider).(*config.Config)
	r.rm = ctn.Get(names.DiRuntimeManagerProvider).(*runtime_manager.RuntimeManager)
	r.mfest = ctn.Get(names.DiManifestProvider).(*manifest.Manifest)
	r.cli = ctn.Get(names.DiDockerClientProvider).(*client.Client)
	r.schedulers = ctn.Get(names.DiSchedulerProvider).(map[string]scheduler.Scheduler)
}

// Prepare prepares the execution
// It creates a directory with the execution id
// It copies the files to the directory
// It returns the execution id
func (r *Executor) Prepare(containerId string, files []models.ExecutionFile) (executionId string, err error) {
	executionId = uuid.New().String()
	base := filepath.Join(r.config.FsMount.MountTargetPath, executionId)
	for _, file := range files {
		if err := r.fp.CreateFile(containerId, base, file); err != nil {
			return "", err
		}
	}
	return
}

func (r *Executor) exec(container *types.Container, c *ExecutionConfig) (stdout cappedbuffer.CappedBuffer, stderr cappedbuffer.CappedBuffer, err error) {
	defer func(containerId string, executionId string) {
		go r.cleanup(containerId, executionId)
	}(container.ID, c.ExecutionId)
	// create container exec process
	workingDir := filepath.Join(r.config.FsMount.MountTargetPath, c.ExecutionId)

	execCreateResponse, err := r.cli.ContainerExecCreate(context.Background(), container.ID, types.ExecConfig{
		Tty: false,
		Cmd: []string{
			"timeout",
			"-s",
			"SIGKILL",
			fmt.Sprint(r.mfest.TaskTimeoutSeconds),
			"sh",
			c.Runtime.Runner,
			c.EntryPoint,
		},
		WorkingDir:   workingDir,
		AttachStderr: true,
		AttachStdout: true,
		Detach:       true,
	})
	if err != nil {
		log.Println(err)
		return
	}

	execAttachResponse, err := r.cli.ContainerExecAttach(context.Background(), execCreateResponse.ID, types.ExecStartCheck{
		Detach: false,
		Tty:    false,
	})
	if err != nil {
		log.Println(err)
		return
	}

	defer execAttachResponse.Close()

	stdout = *cappedbuffer.New(make([]byte, 0, r.config.OutputSizeLimit), r.config.OutputSizeLimit)
	stderr = *cappedbuffer.New(make([]byte, 0, r.config.OutputSizeLimit), r.config.OutputSizeLimit)

	_, err = stdcopy.StdCopy(&stdout, &stderr, execAttachResponse.Reader)
	if err != nil {
		return
	}
	return
}

// IsExecutorAvailable checks if there is an available executor/container to use
func (r *Executor) IsExecutorAvailable(rt *manifest.ManifestRuntime) bool {
	activeContainers, err := r.rm.GetContainersForRuntime(rt)
	if err != nil {
		return false
	}
	return len(activeContainers) > 0
}

// Execute execute code
func (r *Executor) Execute(files []models.ExecutionFile, c *ExecutionConfig) (*models.ExecutionResponse, error) {
	executionScheduler := r.schedulers[c.Runtime.Language]
	var container types.Container
	var err error
	if c.Runtime.SchedulingAlgorithm == manifest.RoundRobin {
		container, err = executionScheduler.RoundRobin.GetNextExecutor(c.Runtime)
	} else {
		container, err = executionScheduler.Random.GetNextExecutor(c.Runtime)
	}
	if err != nil {
		return nil, err
	}
	prepareTimeStart := time.Now()
	executionId, err := r.Prepare(container.ID, files)
	if err != nil {
		return nil, err
	}
	preparationTime := time.Since(prepareTimeStart)
	c.ExecutionId = executionId
	stdout, stderr, err := r.exec(&container, c)
	if err != nil {
		return nil, err
	}
	return &models.ExecutionResponse{
		Stdout:          stdout.String(),
		Stderr:          stderr.String(),
		PreparationTime: preparationTime.String(),
	}, nil
}

// Cleanup cleanup removes the files created for executionId
func (r *Executor) cleanup(containerId string, executionId string) error {
	base := filepath.Join(r.config.FsMount.MountTargetPath, executionId)
	removeFileResponse, err := r.cli.ContainerExecCreate(context.Background(), containerId, types.ExecConfig{
		Cmd:    []string{"rm", "-rf", base},
		Tty:    false,
		Detach: false,
	})
	if err != nil {
		return err
	}
	if _, err := r.cli.ContainerExecAttach(context.Background(), removeFileResponse.ID, types.ExecStartCheck{}); err != nil {
		return err
	}
	return nil
}
