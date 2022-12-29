package containers

import (
	"bufio"
	"context"
	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"
	dtypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ComposeProject struct {
	name string
	file string
	dir  string
}

func UpdateMatchingContainers(name string) error {
	var composeContainers, standaloneContainers []dtypes.Container
	var err error
	if composeContainers, standaloneContainers, err = GetMatchingContainers(name); err != nil {
		return err
	}
	lo.ForEach(standaloneContainers, func(item dtypes.Container, _ int) {
		log.Info().Strs("Container names", item.Names).Msg("Standalone container needs updating")
	})
	lo.ForEach(GetComposeProjectsFromContainers(composeContainers), func(item ComposeProject, _ int) {
		err = UpdateComposeProject(item)
		if err != nil {
			log.Error().Err(err).Str("Project", item.name).Msg("Unable to update project")
		}
	})
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
	env, err := composeEnv(composeProject.dir)
	if err != nil {
		return err
	}
	composeContents, err := os.ReadFile(composeProject.file)
	project, err := loader.Load(types.ConfigDetails{
		WorkingDir: composeProject.dir,
		ConfigFiles: []types.ConfigFile{{
			Filename: composeProject.file,
			Content:  composeContents,
		}},
		Environment: env,
	}, func(options *loader.Options) {
		options.SetProjectName(composeProject.name, true)
	})
	if err != nil {
		return err
	}
	dcli, err := command.NewDockerCli()
	if err != nil {
		return err
	}
	err = dcli.Initialize(&flags.ClientOptions{Common: flags.NewCommonOptions()})
	if err != nil {
		return err
	}
	s := compose.NewComposeService(dcli)
	err = s.Pull(context.Background(), project, api.PullOptions{Quiet: true})
	if err != nil {
		return err
	}
	timeout := 10 * time.Second
	err = s.Up(context.Background(), project, api.UpOptions{
		Create: api.CreateOptions{
			Recreate:             api.RecreateDiverged,
			RecreateDependencies: api.RecreateDiverged,
			Inherit:              true,
			Timeout:              &timeout,
			QuietPull:            true,
		},
		Start: api.StartOptions{
			Project: project,
		},
	})
	if err != nil {
		return err
	}
	return nil
}

func composeEnv(wd string) (map[string]string, error) {
	file, err := filepath.Abs(filepath.Join(wd, ".env"))
	if err != nil {
		return nil, err
	}
	var dotenv []string
	bytes, _ := os.Open(file)
	scanner := bufio.NewScanner(bytes)
	for scanner.Scan() {
		dotenv = append(dotenv, scanner.Text())
	}
	return lo.Assign(envSliceToMap(os.Environ()), envSliceToMap(dotenv)), nil
}

func envSliceToMap(input []string) map[string]string {
	return lo.Associate[string, string, string](input, func(item string) (string, string) {
		split := strings.SplitN(item, "=", 2)
		if len(split) == 1 {
			return split[0], ""
		}
		return split[0], split[1]
	})
}

func isCompose(container dtypes.Container) bool {
	for key := range container.Labels {
		if key == api.ProjectLabel {
			return true
		}
	}
	return false
}
