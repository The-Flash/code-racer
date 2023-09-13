package main

import (
	"flag"
	"log"
	"os"

	"github.com/The-Flash/code-racer/internal/container_manager"
	"github.com/The-Flash/code-racer/internal/manifest"
	"github.com/docker/docker/client"
)

var (
	manifestPtr = flag.String("f", "", "Path to manifest file")
)

func main() {
	flag.Parse()
	manifestFile, err := os.ReadFile(*manifestPtr)
	if err != nil {
		log.Fatal("could not read manifest file", err)
	}

	m := new(manifest.Manifest)
	m.Load(manifestFile)
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal("could not connect to docker ", err)
	}
	defer cli.Close()
	for _, runtime := range m.Runtimes {
		runningInstances := container_manager.CheckNumberOfActiveContainersForRuntime(cli, &runtime)
		if runningInstances < runtime.Instances {
			// spin up containers
			container_manager.SpinUpContainers(cli, &runtime)
		}

	}
}
