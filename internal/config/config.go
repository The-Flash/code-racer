package config

type Config struct {
	ManifestPath string
	FsMount      *FileSystemMount
	RunnersMount *RunnersDirectoryMount
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
		ManifestPath: manifestPath,
		FsMount:      &FileSystemMount{},
		RunnersMount: &RunnersDirectoryMount{},
	}
}
