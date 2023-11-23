package manifest

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ManifestTestSuite struct {
	suite.Suite
}

func (s *ManifestTestSuite) TestLoadSuccess() {
	data := []byte(`
---
periodMinutes: 1 # periodic check to scale up or down containers
taskTimeoutSeconds: 10 # maximum time for task execution
runtimes:
  - language: python3
    image: python:3-alpine3.17
    version: 3.12.0
    instances: 2
    runner: "/runners/python3/run.sh"
    setup: "/runners/python3/setup.sh"
    schedulingAlgorithm: round-robin
`)
	m := Manifest{}
	err := m.Load(data)
	if err != nil {
		s.T().Error(err)
	}

	s.Assert().Equal(m.PeriodMinutes, 1)
	s.Assert().Equal(m.TaskTimeoutSeconds, 10)
}

func (s *ManifestTestSuite) TestGetRuntimeForLanguage() {
	data := []byte(`
---
periodMinutes: 1 # periodic check to scale up or down containers
taskTimeoutSeconds: 10 # maximum time for task execution
runtimes:
  - language: python3
    image: python:3-alpine3.17
    version: 3.12.0
    instances: 2
    runner: "/runners/python3/run.sh"
    setup: "/runners/python3/setup.sh"
    schedulingAlgorithm: round-robin
`)
	m := Manifest{}
	err := m.Load(data)
	if err != nil {
		s.T().Error(err)
	}

	runtime, _ := m.GetRuntimeForLanguage("python3")
	s.Assert().NotEqual(runtime, nil)
	s.Assert().Equal(runtime.Language, "python3")
	s.Assert().Equal(runtime.Image, "python:3-alpine3.17")
	s.Assert().Equal(runtime.Instances, 2)
	s.Assert().Equal(runtime.Version, "3.12.0")
	s.Assert().Equal(runtime.Runner, "/runners/python3/run.sh")
	s.Assert().Equal(runtime.Setup, "/runners/python3/setup.sh")
	s.Assert().Equal(runtime.SchedulingAlgorithm, RoundRobin)
}

func (s *ManifestTestSuite) TestManifestDefaults() {
	data := []byte(`
---
periodMinutes: 1 # periodic check to scale up or down containers
runtimes:
  - language: python3
    image: python:3-alpine3.17
    version: 3.12.0
    instances: 2
    runner: "/runners/python3/run.sh"
    setup: "/runners/python3/setup.sh"
`)
	m := Manifest{}
	err := m.Load(data)
	if err != nil {
		s.T().Error(err)
	}
	runtime, _ := m.GetRuntimeForLanguage("python3")
	s.Assert().NotEqual(runtime, nil)
	s.Assert().Equal(m.TaskTimeoutSeconds, 20)
	s.Assert().Equal(runtime.SchedulingAlgorithm, Random)
}

func (s *ManifestTestSuite) TestGetRuntimes() {
	data := []byte(`
---
periodMinutes: 1 # periodic check to scale up or down containers
runtimes:
  - language: python3
    image: python:3-alpine3.17
    version: 3.12.0
    instances: 2
    runner: "/runners/python3/run.sh"
    setup: "/runners/python3/setup.sh"
`)
	m := Manifest{}
	err := m.Load(data)
	if err != nil {
		s.T().Error(err)
	}
	runtimes := m.GetRuntimes()
	s.Assert().Len(runtimes, 1)
}

func (s *ManifestTestSuite) TestGetLabels() {
	data := []byte(`
---
periodMinutes: 1 # periodic check to scale up or down containers
runtimes:
  - language: python3
    image: python:3-alpine3.17
    version: 3.12.0
    instances: 2
    runner: "/runners/python3/run.sh"
    setup: "/runners/python3/setup.sh"
    labels:
      runtime: python3
      version: 3.12
`)
	m := Manifest{}
	err := m.Load(data)
	if err != nil {
		s.T().Error(err)
	}
	labels := m.GetLabels("python3")
	s.Assert().Equal("python3", labels["runtime"])
	s.Assert().Equal("3.12", labels["version"])
}

func TestManifest(t *testing.T) {
	suite.Run(t, new(ManifestTestSuite))
}
