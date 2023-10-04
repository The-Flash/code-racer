// this package manages containers available for execution
package runtime_manager

import (
	"context"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/The-Flash/code-racer/internal/config"
	"github.com/The-Flash/code-racer/internal/exec_utils"
	"github.com/The-Flash/code-racer/internal/manifest"
	"github.com/The-Flash/code-racer/internal/names"
	"github.com/docker/docker/api/types"
	containerTypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/sarulabs/di/v2"
)

type RuntimeManager struct {
	config *config.Config
	cli    *client.Client
	mfest  *manifest.Manifest
}

type RuntimeConfig struct {
	MountSource string
}

func (r *RuntimeManager) NewRuntimeManager(ctn di.Container) error {
	cli := ctn.Get(names.DiDockerClientProvider).(*client.Client)
	r.cli = cli
	r.mfest = ctn.Get(names.DiManifestProvider).(*manifest.Manifest)
	r.config = ctn.Get(names.DiConfigProvider).(*config.Config)
	return nil
}

func (r *RuntimeManager) checkNumberOfActiveContainersForRuntime(rt *manifest.ManifestRuntime) (int, error) {
	cli := r.cli
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return 0, err
	}
	runtimeImageName := rt.Image
	numberOfContainers := 0
	for _, container := range containers {
		if container.Image == runtimeImageName {
			numberOfContainers++
		}
	}
	return numberOfContainers, nil
}

func (r *RuntimeManager) GetContainersForRuntime(rt *manifest.ManifestRuntime) ([]types.Container, error) {
	cli := r.cli
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		All: false,
	})
	if err != nil {
		return nil, err
	}

	c := make([]types.Container, 0)

	for _, container := range containers {
		if container.Image == rt.Image {
			c = append(c, container)
		}
	}
	return c, nil
}

func (r *RuntimeManager) scaleUpRuntime(rt *manifest.ManifestRuntime) error {

	cli := r.cli
	// spin up containers
	preferredNumberOfInstances := rt.Instances
	numberOfActiveContainers, err := r.checkNumberOfActiveContainersForRuntime(rt)
	if err != nil {
		return err
	}
	// if there are enough containers, return nil
	if numberOfActiveContainers > preferredNumberOfInstances {
		return nil
	}
	numberOfContainersToSpinup := preferredNumberOfInstances - numberOfActiveContainers
	log.Println("Spinning up containers for", rt.Language)
	var wg sync.WaitGroup
	for i := 0; i < numberOfContainersToSpinup; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if r.config.PullImages {
				// pull image
				reader, err := cli.ImagePull(context.Background(), rt.Image, types.ImagePullOptions{})
				if err != nil {
					log.Println(err)
					return
				}
				defer reader.Close()
				io.Copy(os.Stdout, reader)
			}
			mounts := []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: r.config.FsMount.MountSourcePath,
					Target: r.config.FsMount.MountTargetPath,
				},
				{
					Type:   mount.TypeBind,
					Source: r.config.RunnersMount.MountSourcePath,
					Target: r.config.RunnersMount.MountTargetPath,
				},
			}
			if !r.config.DisableNetworking {
				mounts = append(mounts, mount.Mount{
					Type:   mount.TypeBind,
					Source: r.config.NosocketFileMount.MountSourcePath,
					Target: r.config.NosocketFileMount.MountTargetPath,
				},
				)
			}

			resp, err := cli.ContainerCreate(context.Background(), &containerTypes.Config{
				Image:        rt.Image,
				Tty:          true,
				AttachStdout: true,
				AttachStderr: true,
				AttachStdin:  true,
			}, &containerTypes.HostConfig{
				// NetworkMode: r.config.NetworkMode,
				Resources: containerTypes.Resources{
					Memory: r.config.MemoryLimit,
				},
				Mounts: mounts,
			}, nil, nil, "")
			if err != nil {
				log.Println(err)
				return
			}
			if err := cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
				return
			}
			err = r.setupContainer(resp.ID, rt)
			if err != nil {
				log.Println("Error setting up container")
				return
			}
			log.Println("Finished setting up container")

		}()
	}
	wg.Wait()
	log.Printf("Spinned up %d %s container(s)\n", numberOfContainersToSpinup, rt.Language)
	return nil
}

func (r *RuntimeManager) setupContainer(containerID string, runtime *manifest.ManifestRuntime) error {
	cmd := []string{
		"sh",
		runtime.Setup,
	}
	_, _, err := exec_utils.ExecCmd(cmd, exec_utils.ExecCmdConfig{
		ExecConfig: &types.ExecConfig{
			AttachStderr: false,
			AttachStdout: false,
			Tty:          false,
			Detach:       true,
		},
	}, r.cli, containerID)
	if err != nil {
		return err
	}
	return nil
}

func (r *RuntimeManager) scaleDownRuntime(rt *manifest.ManifestRuntime) error {
	cli := r.cli
	numberOfActiveContainers, err := r.checkNumberOfActiveContainersForRuntime(rt)
	if err != nil {
		return err
	}
	excessContainers := numberOfActiveContainers - rt.Instances
	if excessContainers == 0 {
		return nil
	}
	if excessContainers < 0 {
		return nil
	}
	runningContainers, err := r.GetContainersForRuntime(rt)
	if err != nil {
		return err
	}
	containersToRemove := runningContainers[0:excessContainers]
	log.Println("Removing container(s) for", rt.Language)
	noWaitTimeout := 0
	for _, container := range containersToRemove {
		err := cli.ContainerStop(context.Background(), container.ID, containerTypes.StopOptions{
			Timeout: &noWaitTimeout,
		})
		if err != nil {
			log.Println(err)
			return err
		}

		err = cli.ContainerRemove(context.Background(), container.ID, types.ContainerRemoveOptions{})
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println("Removing container", container.ID)

	}
	log.Printf("Removed %d %s container(s)\n", excessContainers, rt.Language)
	return nil
}

// Run scales up or down containers based on manifest file
func (r *RuntimeManager) Run() {
	for {
		err := r.mfest.Load(r.config.ManifestPath)
		if err != nil {
			log.Fatal("could not load manifest", err)
		}
		for _, runtime := range r.mfest.Runtimes {
			runningInstances, err := r.checkNumberOfActiveContainersForRuntime(&runtime)
			if err != nil {
				log.Fatal("could not check number of running instances", err)
			}
			if runningInstances < runtime.Instances {
				log.Println("too few instances running")
				// spin up containers
				err := r.scaleUpRuntime(&runtime)
				if err != nil {
					log.Println("could not scale up runtime", err)
				}
			} else if runningInstances > runtime.Instances {
				log.Println("too many instances running")
				err := r.scaleDownRuntime(&runtime)
				if err != nil {
					log.Println("could not scale up runtime", err)
				}
			}
		}
		time.Sleep(time.Minute * time.Duration(r.mfest.PeriodMinutes))
	}
}
