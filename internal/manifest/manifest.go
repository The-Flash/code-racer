package manifest

import (
	"log"
	"os"

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

func (m *Manifest) Load(manifestPath string) error {
	manifestData, err := os.ReadFile(manifestPath)
	if err != nil {
		log.Fatal("could not read manifest file", err)
	}
	err = yaml.Unmarshal(manifestData, &m)
	return err
}
