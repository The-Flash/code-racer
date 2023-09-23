// Purpose of this package is to create the request's files and folders
package file_system

import (
	"os"
	"path/filepath"

	"github.com/The-Flash/code-racer/internal/config"
	"github.com/The-Flash/code-racer/pkg/models"
)

type FileProvider struct {
	config *config.Config
}

func NewFileProvider(config *config.Config) *FileProvider {
	return &FileProvider{
		config: config,
	}
}

// Method to create a single file
func (fp *FileProvider) CreateFile(base string, file models.ExecutionFile) error {
	if err := os.MkdirAll(base, 0755); err != nil {
		return err
	}
	path := filepath.Join(base, file.Name)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = f.WriteString(file.Content)
	if err != nil {
		return err
	}
	return nil
}
