package file_system

import (
	"os"
	"path"
	"testing"

	"github.com/The-Flash/code-racer/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestFileProvider_CreateFile(t *testing.T) {
	type Tests struct {
		name           string
		base           string
		executionFile  models.ExecutionFile
		expectError    bool
		expectedOutput interface{}
	}
	fileProvider := NewFileProvider()

	tests := []Tests{
		{
			name: "Create file successfully",
			base: "/tmp/executions",
			executionFile: models.ExecutionFile{
				Name:    "test",
				Content: "",
			},
			expectError:    false,
			expectedOutput: nil,
		},
		{
			name: "Fail when creating file",
			executionFile: models.ExecutionFile{
				Name:    "test",
				Content: "",
			},
			expectError:    true,
			expectedOutput: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := fileProvider.CreateFile(tt.base, tt.executionFile)
			if tt.expectError {
				assert.Error(t, err, "Expected an error but did not get one")
				assert.NoDirExists(t, tt.base, "Expected that directory does existed")
				assert.NoFileExists(t, path.Join(tt.base, tt.executionFile.Name), "Expected that file did not existed")

			} else {
				assert.NoError(t, err, "Did not expect an error but got one")
				assert.DirExists(t, tt.base, "Expected that directory existed")
				assert.FileExists(t, path.Join(tt.base, tt.executionFile.Name), "Expected that file existed")
			}
		})
	}
}

func TestFileProvider_CreateFileWithContent(t *testing.T) {
	fileProvider := NewFileProvider()

	executionFile := models.ExecutionFile{
		Name:    "README.md",
		Content: "package main",
	}

	base := "/tmp/test"

	err := fileProvider.CreateFile(base, executionFile)
	filepath := path.Join(base, executionFile.Name)
	assert.NoError(t, err, "Did not expect an error but got one")
	assert.DirExists(t, base, "Expected that directory existed")
	assert.FileExists(t, filepath, "Expected that file existed")

	c, err := os.ReadFile(filepath)
	assert.NoError(t, err, "Did not expect an error when reading file but got one")
	assert.Equal(t, executionFile.Content, string(c), "Expected that content was the same")
}
