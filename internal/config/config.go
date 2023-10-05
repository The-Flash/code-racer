package config

type Config struct {
	ManifestPath      string
	FsMount           *FileSystemMount
	RunnersMount      *RunnersDirectoryMount
	NosocketFileMount *NosocketMount
	MemoryLimit       int64
	OutputSizeLimit   int
	PullImages        bool
	DisableNetworking bool
	PrLimits          *PrLimitsConfig
}

type PrLimitsConfig struct {
	MaxProcesses int
	MaxOpenFiles int
	MaxFileSize  int
}

type FileSystemMount struct {
	MountSourcePath string
	MountTargetPath string
}

type RunnersDirectoryMount struct {
	MountSourcePath string
	MountTargetPath string
}

type NosocketMount struct {
	MountSourcePath string
	MountTargetPath string
}

func NewConfig() *Config {
	return &Config{
		ManifestPath:      "./manifest.yml",
		FsMount:           &FileSystemMount{},
		RunnersMount:      &RunnersDirectoryMount{},
		MemoryLimit:       1024 * 1024 * 1024, // ! GiB
		OutputSizeLimit:   2 * 1024 * 1024,    // 2 MiB
		PullImages:        true,
		NosocketFileMount: &NosocketMount{},
		DisableNetworking: false,
		PrLimits: &PrLimitsConfig{
			MaxProcesses: 256,
			MaxOpenFiles: 2048,
			MaxFileSize:  10 * 1024 * 1024, // 10 MiB
		},
	}
}
