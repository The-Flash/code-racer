package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/The-Flash/code-racer/internal/manifest"
	"github.com/The-Flash/code-racer/internal/runtime_manager"
	"github.com/docker/docker/client"
)

var (
	manifestPtr   = flag.String("f", "", "Path to manifest file")
	mountPointPtr = flag.String("m", "", "Path to mount point")
)

func main() {
	flag.Parse()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal("could not connect to docker ", err)
	}
	defer cli.Close()

	for {

		manifestFile, err := os.ReadFile(*manifestPtr)
		if err != nil {
			log.Fatal("could not read manifest file", err)
		}
		m := new(manifest.Manifest)
		err = m.Load(manifestFile)
		if err != nil {
			log.Fatal("could not load manifest", err)
		}
		for _, runtime := range m.Runtimes {
			runningInstances, err := runtime_manager.CheckNumberOfActiveContainersForRuntime(cli, &runtime)
			if err != nil {
				log.Fatal("could not check number of running instances", err)
			}
			if runningInstances < runtime.Instances {
				fmt.Println("too few instances running")
				// spin up containers
				err := runtime_manager.ScaleUpRuntime(cli, &runtime, &runtime_manager.RuntimeConfig{
					MountSource: *mountPointPtr,
				})
				if err != nil {
					fmt.Println("could not scale up runtime", err)
				}
			} else if runningInstances > runtime.Instances {
				fmt.Println("too many instances running")
				err := runtime_manager.ScaleDownRuntime(cli, &runtime)
				if err != nil {
					fmt.Println("could not scale up runtime", err)
				}
			}
		}
		time.Sleep(time.Minute * time.Duration(m.PeriodMinutes))
	}
}
