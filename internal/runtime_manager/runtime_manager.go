// this package manages containers available for execution
package runtime_manager

import (
	"context"
	"fmt"
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
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/sarulabs/di/v2"
)

type RuntimeManager struct {
	config *config.Config
	cli    *client.Client
	mfest  *manifest.Manifest
	mu     sync.Mutex
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
	containers, err := r.GetContainersForRuntime(rt)
	return len(containers), err
}

func (r *RuntimeManager) IsContainerReady(containerId string) bool {
	cmd := []string{"cat", "READY"}
	stdout, stderr, err := exec_utils.ExecCmd(cmd, exec_utils.ExecCmdConfig{
		StdOutSizeLimit: r.config.OutputSizeLimit,
		StdErrSizeLimit: r.config.OutputSizeLimit,
		ExecConfig: &types.ExecConfig{
			AttachStderr: true,
			AttachStdout: true,
			WorkingDir:   "/",
		},
	}, r.cli, containerId)
	if err != nil {
		return false
	}
	if stdout.String() != "" {
		return true
	}
	if stderr.String() != "" {
		return false
	}
	return false
}

func (r *RuntimeManager) GetContainersForRuntime(rt *manifest.ManifestRuntime) ([]types.Container, error) {
	cli := r.cli
	filters := filters.NewArgs()
	labels := rt.Labels
	for k, v := range labels {
		filters.Add("label", fmt.Sprintf("%s=%s", k, v))
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		All:     false,
		Filters: filters,
	})
	if err != nil {
		return nil, err
	}

	c := make([]types.Container, 0, 10)

	for _, container := range containers {
		if r.IsContainerReady(container.ID) {
			r.mu.Lock()
			c = append(c, container)
			r.mu.Unlock()
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
				Labels:       rt.Labels,
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
		// log.Printf("Spinned up %d %s container(s)\n", numberOfContainersToSpinup, rt.Language)
	}
	wg.Wait()
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
	return nil
}

// Run scales up or down containers based on manifest file
func (r *RuntimeManager) Run() {
	for {
		manifestData, err := os.ReadFile(r.config.ManifestPath)
		if err != nil {
			log.Fatal("could not read manifest file", err)
		}
		err = r.mfest.Load(manifestData)
		if err != nil {
			log.Fatal("could not load manifest", err)
		}
		for _, runtime := range r.mfest.Runtimes {
			runningInstances, err := r.checkNumberOfActiveContainersForRuntime(&runtime)
			if err != nil {
				log.Fatal("could not check number of running instances", err)
			}
			if runningInstances < runtime.Instances {
				log.Printf("too few instances of %s running", runtime.Language)
				// spin up containers
				err := r.scaleUpRuntime(&runtime)
				if err != nil {
					log.Println("could not scale up runtime", err)
				}
			} else if runningInstances > runtime.Instances {
				log.Printf("too many instances of %s running\n", runtime.Language)
				err := r.scaleDownRuntime(&runtime)
				if err != nil {
					log.Println("could not scale up runtime", err)
				}
			}
		}
		time.Sleep(time.Minute * time.Duration(r.mfest.PeriodMinutes))
	}
}
