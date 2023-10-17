// Purpose of this package is to create the request's files and folders
package file_system

import (
	"fmt"
	"path/filepath"

	"github.com/The-Flash/code-racer/internal/config"
	"github.com/The-Flash/code-racer/internal/exec_utils"
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
	_, _, err := exec_utils.ExecCmd([]string{"mkdir", "-p", dirPart}, exec_utils.ExecCmdConfig{
		ExecConfig: &types.ExecConfig{
			Tty:    false,
			Detach: true,
		},
	}, cli, containerId)
	if err != nil {
		return err
	}
	name := filepath.Base(file.Name)
	_, _, err = exec_utils.ExecCmd([]string{"sh", "-c", fmt.Sprintf("echo '%s' > %s", file.Content, name)}, exec_utils.ExecCmdConfig{
		ExecConfig: &types.ExecConfig{
			Tty:        false,
			Detach:     true,
			WorkingDir: dirPart,
		},
	}, cli, containerId)
	if err != nil {
		return err
	}
	return nil
}
