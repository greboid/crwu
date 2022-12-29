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
	"github.com/samber/lo"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func UpdateComposeProject(name string, composeFile string, workingDir string) error {
	env, err := composeEnv(workingDir)
	if err != nil {
		return err
	}
	composeContents, err := os.ReadFile(composeFile)
	project, err := loader.Load(types.ConfigDetails{
		WorkingDir: workingDir,
		ConfigFiles: []types.ConfigFile{{
			Filename: composeFile,
			Content:  composeContents,
		}},
		Environment: env,
	}, func(options *loader.Options) {
		options.SetProjectName(name, true)
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
