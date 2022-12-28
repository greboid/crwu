package containers

import (
	"context"
	"github.com/compose-spec/compose-go/types"
	"github.com/docker/cli/cli/command"
	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/compose"
	"github.com/rs/zerolog/log"
	"time"
)

func doComposeThings() {
	dcli, err := command.NewDockerCli()
	if err != nil {
		log.Fatal().Err(err).Msg("Can't create Docker CLI")
	}
	s := compose.NewComposeService(dcli)
	project := &types.Project{
		Name:             "",
		WorkingDir:       "",
		Services:         nil,
		Networks:         nil,
		Volumes:          nil,
		Secrets:          nil,
		Configs:          nil,
		Extensions:       nil,
		ComposeFiles:     nil,
		Environment:      nil,
		DisabledServices: nil,
	}
	err = s.Pull(context.Background(), project, api.PullOptions{Quiet: true})
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to pull")
	}
	timeout := 10 * time.Second
	err = s.Up(context.Background(), project, api.UpOptions{
		Create: api.CreateOptions{
			Services:             nil,
			Recreate:             api.RecreateDiverged,
			RecreateDependencies: api.RecreateDiverged,
			Inherit:              true,
			Timeout:              &timeout,
			QuietPull:            true,
		},
		Start: api.StartOptions{
			Project:  nil,
			Services: nil,
		},
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to up")
	}
}
