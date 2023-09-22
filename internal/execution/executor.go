package execution

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/The-Flash/code-racer/internal/config"
	"github.com/The-Flash/code-racer/internal/file_system"
	"github.com/The-Flash/code-racer/internal/manifest"
	"github.com/The-Flash/code-racer/internal/names"
	"github.com/The-Flash/code-racer/internal/runtime_manager"
	"github.com/The-Flash/code-racer/pkg/models"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/google/uuid"
	"github.com/sarulabs/di/v2"
)

type ExecutionConfig struct {
	ExecutionId string
	EntryPoint  string
	Runtime     *manifest.ManifestRuntime
}

type Executor struct {
	mfest  *manifest.Manifest
	fp     *file_system.FileProvider
	config *config.Config
	rm     *runtime_manager.RuntimeManager
	cli    *client.Client
}

// Setup
func (r *Executor) Setup(ctn di.Container) {
	fp := ctn.Get(names.DiFileProvider).(*file_system.FileProvider)
	config := ctn.Get(names.DiConfigProvider).(*config.Config)
	rm := ctn.Get(names.DiRuntimeManagerProvider).(*runtime_manager.RuntimeManager)
	mfest := ctn.Get(names.DiManifestProvider).(*manifest.Manifest)
	cli := ctn.Get(names.DiDockerClientProvider).(*client.Client)
	r.fp = fp
	r.config = config
	r.rm = rm
	r.mfest = mfest
	r.cli = cli
}

// Prepare prepares the execution
func (r *Executor) Prepare(files []models.ExecutionFile) (executionId string, err error) {
	executionId = uuid.New().String()
	base := filepath.Join(r.config.FsMount.MountSourcePath, executionId)
	for _, file := range files {
		if err := r.fp.CreateFile(base, file); err != nil {
			return "", err
		}
	}
	return
}

func (r *Executor) exec(container types.Container, c *ExecutionConfig) (stdout bytes.Buffer, stderr bytes.Buffer, err error) {
	go r.cleanup(c.ExecutionId)
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

	_, err = stdcopy.StdCopy(&stdout, &stderr, execAttachResponse.Reader)
	if err != nil {
		return
	}
	return
}

// Execute
func (r *Executor) Execute(c *ExecutionConfig) (*models.ExecutionResponse, error) {
	// check the algorithm for this runtime
	// algorithm := c.Runtime.SchedulingAlgorithm
	// find the next runtime to use

	activeContainers, err := r.rm.GetContainersForRuntime(c.Runtime)
	if err != nil {
		return nil, err
	}
	numberOfActiveContainers := len(activeContainers)
	if numberOfActiveContainers == 0 {
		return nil, errors.New("no executors available")
	}
	selectedIndex := rand.Intn(numberOfActiveContainers)
	container := activeContainers[selectedIndex]

	stdout, stderr, err := r.exec(container, c)
	if err != nil {
		return nil, err
	}
	return &models.ExecutionResponse{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}, nil
}

// Cleanup cleanup kills the task inside it's container and deletes its source files
// It does this after sleeping for the timeout seconds specified in the manifest
func (r *Executor) cleanup(executionId string) error {
	time.Sleep(time.Second * time.Duration(r.mfest.TaskTimeoutSeconds))
	base := filepath.Join(r.config.FsMount.MountSourcePath, executionId)
	// TODO: Kill running process
	if err := os.RemoveAll(base); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
