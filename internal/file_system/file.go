// Purpose of this package is to create the request's files and folders
package file_system

import (
	"os"
	"path/filepath"

	"github.com/The-Flash/code-racer/internal/config"
	"github.com/The-Flash/code-racer/pkg/models"
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
	// cli := fp.cli
	dir := filepath.Dir(file.Name)
	dirPart := filepath.Join(base, dir)
	err := os.MkdirAll(base, os.ModePerm)
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(dirPart, file.Name))
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err = f.WriteString(file.Content); err != nil {
		return err
	}
	f.Sync()
	return nil
}
