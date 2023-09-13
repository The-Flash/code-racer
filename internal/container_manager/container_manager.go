// this package manages containers available for execution
package container_manager

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/The-Flash/code-racer/internal/manifest"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func CheckNumberOfActiveContainersForRuntime(cli *client.Client, r *manifest.ManifestRuntime) int {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		log.Fatal("could not get container list ", err)
	}
	runtimeImageName := r.Image
	numberOfContainers := 0
	for _, container := range containers {
		if container.Image == runtimeImageName {
			numberOfContainers++
		}
	}
	return numberOfContainers
}

func SpinUpContainers(cli *client.Client, r *manifest.ManifestRuntime) {
	ctx := context.Background()
	// spin up containers
	preferredNumberOfInstances := r.Instances
	fmt.Println("Spinning up containers for", r.Language)
	for i := 0; i < preferredNumberOfInstances; i++ {
		// pull image
		reader, err := cli.ImagePull(ctx, r.Image, types.ImagePullOptions{})
		if err != nil {
			panic(err)
		}
		defer reader.Close()
		io.Copy(os.Stdout, reader)
		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Image: r.Image,
			Tty:   true,
		}, nil, nil, nil, "")
		fmt.Println(resp)
		if err != nil {
			panic(err)
		}
		if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			panic(err)
		}
	}
	fmt.Printf("Spinned up %d %s containers", preferredNumberOfInstances, r.Language)
}
