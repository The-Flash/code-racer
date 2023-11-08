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
	Runtimes           []ManifestRuntime `yaml:"runtimes" json:"runtimes"`
	PeriodMinutes      int               `yaml:"periodMinutes" json:"periodMinutes"`
	TaskTimeoutSeconds int               `yaml:"taskTimeoutSeconds" json:"taskTimeoutSeconds"`
}

type ManifestRuntime struct {
	Language            string              `yaml:"language" json:"language"`
	Image               string              `yaml:"image" json:"-"`
	Version             string              `yaml:"version" json:"version"`
	Instances           int                 `yaml:"instances" json:"-"`
	SchedulingAlgorithm SchedulingAlgorithm `yaml:"schedulingAlgorithm" json:"-"`
	Runner              string              `yaml:"runner" json:"-"`
	Setup               string              `yaml:"setup" json:"-"`
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
		}
		runtimes = append(runtimes, r)
	}
	copy(m.Runtimes, runtimes)
}

func (m *Manifest) Load(data []byte) error {
	err := yaml.Unmarshal(data, &m)
	m.setDefaults()
	return err
}
