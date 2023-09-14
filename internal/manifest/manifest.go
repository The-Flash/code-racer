package manifest

import (
	"gopkg.in/yaml.v3"
)

type SchedulingAlgorithm = string

const (
	Random     SchedulingAlgorithm = "random"
	RoundRobin SchedulingAlgorithm = "round-robin"
)

type Manifest struct {
	Runtimes      []ManifestRuntime `yaml:"runtimes"`
	PeriodMinutes int               `yaml:"periodMinutes"`
}

type ManifestRuntime struct {
	Language            string              `yaml:"language"`
	Image               string              `yaml:"image"`
	Instances           int                 `yaml:"instances"`
	Aliases             []string            `yaml:"aliases"`
	SchedulingAlgorithm SchedulingAlgorithm `yaml:"schedulingAlgorithm" default:"random"`
}

func (m *Manifest) Load(f []byte) error {
	err := yaml.Unmarshal(f, &m)
	return err
}
