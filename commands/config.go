package commands

import (
	"context"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/jfrog/live-logs/internal"
	"github.com/jfrog/live-logs/internal/constants"
	"strconv"
)

func GetConfigCommand() components.Command {
	return components.Command{
		Name:        "config",
		Description: "Display the list of nodes and log file names",
		Aliases:     []string{"c"},
		Arguments:   getConfigArguments(),
		Flags:       getConfigFlags(),
		EnvVars:     getConfigEnvVar(),
		Action:      configCmd,
	}
}

func getConfigArguments() []components.Argument {
	return []components.Argument{
		{Name: "product-id", Description: "JFrog product id; the value can be one of the following, \n" +
			"\t\t\t" + constants.ArtifactoryId + " - Artifactory\n" +
			"\t\t\t" + constants.XrayId + " - Xray\n" +
			"\t\t\t" + constants.McId + " - Mission Control\n" +
			"\t\t\t" + constants.DistributionId + " - Distribution\n" +
			"\t\t\t" + constants.PipelinesId + " - Pipelines"},
		{Name: "server-id", Description: "JFrog CLI Artifactory server id"},
	}
}

func getConfigFlags() []components.Flag {
	return []components.Flag{
		components.BoolFlag{
			Name:         constants.InteractiveFlag,
			Description:  "Activate the interactive menu",
			DefaultValue: false,
		},
	}
}

func getConfigEnvVar() []components.EnvVar {
	return []components.EnvVar{
		{
			Name:        constants.VersionCheckEnv,
			Default:     "true",
			Description: "Set this to \"false\" to disable validation on minimum supported version of the product.",
		},
	}
}

func configCmd(c *components.Context) error {
	isInteractive := c.GetBoolFlagValue(constants.InteractiveFlag)

	mainCtx, mainCtxCancel := context.WithCancel(context.Background())
	defer mainCtxCancel()

	var liveLogClient livelog.LiveLogs
	liveLogClient = livelog.NewLiveLogs()

	if !isInteractive {
		if len(c.Arguments) != 2 {
			return fmt.Errorf("incorrect number of arguments were passed: expected: 2," + " received: " + strconv.Itoa(len(c.Arguments)))
		}
		productId := c.Arguments[0]
		serverId := c.Arguments[1]
		return liveLogClient.ConfigNonInteractive(mainCtx, productId, serverId)
	}
	return ConfigInteractive(mainCtx, liveLogClient)
}
