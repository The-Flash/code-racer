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
	Runtimes           []ManifestRuntime `yaml:"runtimes" json:"runtimes"`
	PeriodMinutes      int               `yaml:"periodMinutes" json:"periodMinutes"`
	TaskTimeoutSeconds int               `yaml:"taskTimeoutSeconds" json:"taskTimeoutSeconds"`
}

type ManifestRuntime struct {
	Language            string              `yaml:"language" json:"language"`
	Image               string              `yaml:"image" json:"image"`
	Instances           int                 `yaml:"instances" json:"instances"`
	SchedulingAlgorithm SchedulingAlgorithm `yaml:"schedulingAlgorithm" json:"schedulingAlgorithm"`
	Runner              string              `yaml:"runner" json:"runner"`
}

func (m *Manifest) GetRuntimeForLanguage(language string) (*ManifestRuntime, bool) {
	for _, r := range m.Runtimes {
		if r.Language == language {
			return &r, true
		}
	}
	return &ManifestRuntime{}, false
}

func (m *Manifest) GetRuntimes() []ManifestRuntime {
	return m.Runtimes
}

func (m *Manifest) setDefaults() {
	var runtimes []ManifestRuntime = []ManifestRuntime{}
	if m.TaskTimeoutSeconds == 0 {
		m.TaskTimeoutSeconds = 20
	}
	for _, r := range m.Runtimes {
		if r.SchedulingAlgorithm == "" {
			r.SchedulingAlgorithm = Random
			runtimes = append(runtimes, r)
		}
	}
	copy(m.Runtimes, runtimes)
}

func (m *Manifest) Load(manifestPath string) error {
	manifestData, err := os.ReadFile(manifestPath)
	if err != nil {
		log.Fatal("could not read manifest file", err)
	}
	err = yaml.Unmarshal(manifestData, &m)
	m.setDefaults()
	return err
}
