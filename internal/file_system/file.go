// Purpose of this package is to create the request's files and folders
package file_system

import (
	"os"
	"path/filepath"

	"github.com/The-Flash/code-racer/pkg/models"
)

type FileProvider struct{}

func NewFileProvider() *FileProvider {
	return &FileProvider{}
}

// Method to create a single file
func (fp *FileProvider) CreateFile(base string, file models.ExecutionFile) error {
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
