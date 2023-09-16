package runtime_manager

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/The-Flash/code-racer/internal/manifest"
	"github.com/docker/docker/client"
)

func RunScaler(cli *client.Client, manifestPath string, mountPointPath string) {
	for {

		manifestFile, err := os.ReadFile(manifestPath)
		if err != nil {
			log.Fatal("could not read manifest file", err)
		}
		m := new(manifest.Manifest)
		err = m.Load(manifestFile)
		if err != nil {
			log.Fatal("could not load manifest", err)
		}
		for _, runtime := range m.Runtimes {
			runningInstances, err := CheckNumberOfActiveContainersForRuntime(cli, &runtime)
			if err != nil {
				log.Fatal("could not check number of running instances", err)
			}
			if runningInstances < runtime.Instances {
				fmt.Println("too few instances running")
				// spin up containers
				err := ScaleUpRuntime(cli, &runtime, &RuntimeConfig{
					MountSource: mountPointPath,
				})
				if err != nil {
					fmt.Println("could not scale up runtime", err)
				}
			} else if runningInstances > runtime.Instances {
				fmt.Println("too many instances running")
				err := ScaleDownRuntime(cli, &runtime)
				if err != nil {
					fmt.Println("could not scale up runtime", err)
				}
			}
		}
		time.Sleep(time.Minute * time.Duration(m.PeriodMinutes))
	}
}
