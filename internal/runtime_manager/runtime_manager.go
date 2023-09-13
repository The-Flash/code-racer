// this package manages containers available for execution
package runtime_manager

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/The-Flash/code-racer/internal/manifest"
	"github.com/docker/docker/api/types"
	containerTypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

var (
	ctx = context.Background()
)

func CheckNumberOfActiveContainersForRuntime(cli *client.Client, r *manifest.ManifestRuntime) (int, error) {
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return 0, err
	}
	runtimeImageName := r.Image
	numberOfContainers := 0
	for _, container := range containers {
		if container.Image == runtimeImageName {
			numberOfContainers++
		}
	}
	return numberOfContainers, nil
}

func getContainersForRuntime(cli *client.Client, r *manifest.ManifestRuntime) ([]types.Container, error) {
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{
		All: false,
	})
	if err != nil {
		return nil, err
	}

	c := make([]types.Container, 0)

	for _, container := range containers {
		if container.Image == r.Image {
			c = append(c, container)
		}
	}
	return c, nil
}

func ScaleUpRuntime(cli *client.Client, r *manifest.ManifestRuntime) error {
	// spin up containers
	preferredNumberOfInstances := r.Instances
	numberOfActiveContainers, err := CheckNumberOfActiveContainersForRuntime(cli, r)
	if err != nil {
		return err
	}
	// if there are enough containers, return nil
	if numberOfActiveContainers > preferredNumberOfInstances {
		return nil
	}
	numberOfContainersToSpinup := preferredNumberOfInstances - numberOfActiveContainers
	fmt.Println("Spinning up containers for", r.Language)
	for i := 0; i < numberOfContainersToSpinup; i++ {
		// pull image
		reader, err := cli.ImagePull(ctx, r.Image, types.ImagePullOptions{})
		if err != nil {
			return err
		}
		defer reader.Close()
		io.Copy(os.Stdout, reader)
		resp, err := cli.ContainerCreate(ctx, &containerTypes.Config{
			Image: r.Image,
			Tty:   true,
		}, nil, nil, nil, "")
		if err != nil {
			return err
		}
		if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			return err
		}
	}
	fmt.Printf("Spinned up %d %s containers\n", numberOfContainersToSpinup, r.Language)
	return nil
}

func ScaleDownRuntime(cli *client.Client, r *manifest.ManifestRuntime) error {
	numberOfActiveContainers, err := CheckNumberOfActiveContainersForRuntime(cli, r)
	if err != nil {
		return err
	}
	excessContainers := numberOfActiveContainers - r.Instances
	if excessContainers == 0 {
		return nil
	}
	if excessContainers < 0 {
		return nil
	}
	runningContainers, err := getContainersForRuntime(cli, r)
	if err != nil {
		return err
	}
	containersToRemove := runningContainers[0:excessContainers]
	fmt.Println("Removing containers for", r.Language)
	noWaitTimeout := 0
	for _, container := range containersToRemove {
		err := cli.ContainerStop(ctx, container.ID, containerTypes.StopOptions{
			Timeout: &noWaitTimeout,
		})
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Removing container", container.ID)
	}
	fmt.Printf("Removed %d %s containers\n", excessContainers, r.Language)
	return nil
}
