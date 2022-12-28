package containers

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func GetMatchingContainers(name string) ([]types.Container, []types.Container, error) {
	var cContainers []types.Container
	var sContainers []types.Container
	dclient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, nil, err
	}
	listOptions := types.ContainerListOptions{}
	listOptions.Filters = filters.NewArgs(filters.Arg("ancestor", name))
	containers, err := dclient.ContainerList(context.Background(), listOptions)
	if err != nil {
		return nil, nil, err
	}
	for index := range containers {
		if isCompose(containers[index]) {
			cContainers = append(cContainers, containers[index])
		} else {
			sContainers = append(sContainers, containers[index])
		}
	}
	return cContainers, sContainers, nil
}

func isCompose(container types.Container) bool {
	for key := range container.Labels {
		if key == "com.docker.compose.project" {
			return true
		}
	}
	return false
}
