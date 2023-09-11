package manifest

import (
	"gopkg.in/yaml.v3"
)

type Manifest struct {
	Runtimes []ManifestRuntime `yaml:"runtimes"`
}

type ManifestRuntime struct {
	Language  string   `yaml:"language"`
	Image     string   `yaml:"image"`
	Instances int      `yaml:"instances"`
	Aliases   []string `yaml:"aliases"`
}

func (m *Manifest) Load(f []byte) error {
	err := yaml.Unmarshal(f, &m)
	return err
}
