package config

type Config struct {
	ManifestPath    string
	MountSourcePath string
	MountTargetPath string
}

func NewConfig(manifestPath string, mountSourcePath string, mountTargetPath string) *Config {
	return &Config{
		ManifestPath:    manifestPath,
		MountSourcePath: mountSourcePath,
		MountTargetPath: mountTargetPath,
	}
}
