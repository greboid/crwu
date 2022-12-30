package containers

import (
	"context"
	"github.com/docker/compose/v2/pkg/api"
	dtypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"os/exec"
)

type ComposeProject struct {
	name string
	file string
	dir  string
}

func UpdateMatchingContainers(name string) (int, error) {
	var composeContainers, standaloneContainers []dtypes.Container
	var err error
	if composeContainers, standaloneContainers, err = GetMatchingContainers(name); err != nil {
		return 0, err
	}
	lo.ForEach(standaloneContainers, func(item dtypes.Container, _ int) {
		log.Info().Strs("Container names", item.Names).Msg("Standalone container needs updating")
	})
	count := 0
	lo.ForEach(GetComposeProjectsFromContainers(composeContainers), func(item ComposeProject, _ int) {
		log.Debug().Interface("project", item).Msg("Starting update")
		err = UpdateComposeProject(item)
		if err != nil {
			log.Error().Err(err).Str("Project", item.name).Msg("Unable to update project")
		} else {
			log.Debug().Interface("project", item).Msg("Finished update")
			count++
		}
	})
	return count, nil
}

func GetMatchingContainers(name string) ([]dtypes.Container, []dtypes.Container, error) {
	var cContainers []dtypes.Container
	var sContainers []dtypes.Container
	dclient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, nil, err
	}
	listOptions := dtypes.ContainerListOptions{}
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

func GetComposeProjectsFromContainers(containers []dtypes.Container) []ComposeProject {
	return lo.Uniq(lo.Map(containers, func(item dtypes.Container, index int) ComposeProject {
		return ComposeProject{
			name: item.Labels[api.ProjectLabel],
			file: item.Labels[api.ConfigFilesLabel],
			dir:  item.Labels[api.WorkingDirLabel],
		}
	}))
}

func UpdateComposeProject(composeProject ComposeProject) error {
	err := exec.Command("/bin/chroot", "/composemount", "/docker-compose", "-f", composeProject.file, "pull").Run()
	if err != nil {
		return err
	}
	return exec.Command("/bin/chroot", "/composemount", "/docker-compose", "-f", composeProject.file, "up", "-d").Run()
}

func isCompose(container dtypes.Container) bool {
	for key := range container.Labels {
		if key == api.ProjectLabel {
			return true
		}
	}
	return false
}
