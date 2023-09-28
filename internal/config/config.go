package config

import "github.com/docker/docker/api/types/container"

type Config struct {
	ManifestPath    string
	FsMount         *FileSystemMount
	RunnersMount    *RunnersDirectoryMount
	NetworkMode     container.NetworkMode
	MemoryLimit     int64
	OutputSizeLimit int
}

type FileSystemMount struct {
	MountSourcePath string
	MountTargetPath string
}

type RunnersDirectoryMount struct {
	MountSourcePath string
	MountTargetPath string
}

func NewConfig(manifestPath string) *Config {
	return &Config{
		ManifestPath:    manifestPath,
		FsMount:         &FileSystemMount{},
		RunnersMount:    &RunnersDirectoryMount{},
		NetworkMode:     container.NetworkMode("none"), // no network access
		MemoryLimit:     1024 * 1024 * 1024,            // ! GiB
		OutputSizeLimit: 2 * 1024 * 1024,               // 2 MiB
	}
}
