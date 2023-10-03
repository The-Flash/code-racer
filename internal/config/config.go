package config

import "github.com/docker/docker/api/types/container"

type Config struct {
	ManifestPath    string
	FsMount         *FileSystemMount
	RunnersMount    *RunnersDirectoryMount
	NetworkMode     container.NetworkMode
	LocalBinPath    string // contains binaries such as nosocket
	MemoryLimit     int64
	OutputSizeLimit int
	PullImages      bool
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
		PullImages:      true,
		LocalBinPath:    "/usr/local/bin/nosocket",
	}
}
