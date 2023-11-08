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
	m.Load(data)

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
	m.Load(data)

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
	m.Load(data)

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
	m.Load(data)
	runtimes := m.GetRuntimes()
	s.Assert().Len(runtimes, 1)
}

func TestManifest(t *testing.T) {
	suite.Run(t, new(ManifestTestSuite))
}
