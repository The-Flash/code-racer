// Purpose of this package is to create the request's files and folders
package file_system

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/The-Flash/code-racer/internal/config"
	"github.com/The-Flash/code-racer/pkg/models"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type FileProvider struct {
	config *config.Config
	cli    *client.Client
}

func NewFileProvider(config *config.Config, cli *client.Client) *FileProvider {
	return &FileProvider{
		config: config,
		cli:    cli,
	}
}

// Method to create a single file
func (fp *FileProvider) CreateFile(containerId string, base string, file models.ExecutionFile) error {
	cli := fp.cli
	dir := filepath.Dir(file.Name)
	dirPart := filepath.Join(base, dir)
	mkdirExecResponse, err := cli.ContainerExecCreate(context.Background(), containerId, types.ExecConfig{
		Cmd:    []string{"mkdir", "-p", dirPart},
		Tty:    false,
		Detach: false,
	})
	if err != nil {
		return err
	}
	if _, err := cli.ContainerExecAttach(context.Background(), mkdirExecResponse.ID, types.ExecStartCheck{}); err != nil {
		return err
	}
	name := filepath.Base(file.Name)
	createFileResponse, err := cli.ContainerExecCreate(context.Background(), containerId, types.ExecConfig{
		Cmd:        []string{"sh", "-c", fmt.Sprintf("echo \"%s\" > %s", file.Content, name)},
		Tty:        false,
		Detach:     false,
		WorkingDir: dirPart,
	})
	if err != nil {
		return err
	}
	if _, err := cli.ContainerExecAttach(context.Background(), createFileResponse.ID, types.ExecStartCheck{}); err != nil {
		return err
	}

	return nil
}
