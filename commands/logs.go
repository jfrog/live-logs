package commands

import (
	"context"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"github.com/jfrog/live-logs/internal"
	"github.com/jfrog/live-logs/internal/constants"
	"strconv"
)

func GetLogsCommand() components.Command {
	return components.Command{
		Name:        "logs",
		Description: "Fetch the log of a desired service" +
		"\n\nNote:" +
		"\n\t- Xray, Mission Control, Pipelines, and Distribution only support admin access token authentication, while, Artifactory supports all types of authentication." +
		"\n\t- The scope of the generated access token is limited to the corresponding product." +
		"\n\t- For every product, a new dedicated entry will need to be added. " +
		"For example, if you want to stream logs from 3 products, a separate entry will need to be configured for each product in the JFrog CLI (so is 3 entries).",
		Aliases:     []string{"l"},
		Arguments:   getLogsArguments(),
		EnvVars:     getLogsEnvVar(),
		Flags:       getLogsFlags(),
		Action:      logsCmd,
	}
}

func getLogsArguments() []components.Argument {
	return []components.Argument{
		{Name: "product-id", Description: "JFrog product id; the value can be one of the following, \n" +
			"\t\t\t" + constants.ArtifactoryId + " - Artifactory\n" +
			"\t\t\t" + constants.XrayId + " - Xray\n" +
			"\t\t\t" + constants.McId + " - Mission Control\n" +
			"\t\t\t" + constants.DistributionId + " - Distribution\n" +
			"\t\t\t" + constants.PipelinesId + " - Pipelines"},
		{Name: "server-id", Description: "JFrog CLI Artifactory server id"},
		{Name: "node-id", Description: "Selected node id"},
		{Name: "log-name", Description: "Selected log name"},
	}
}

func getLogsFlags() []components.Flag {
	return []components.Flag{
		components.BoolFlag{
			Name:         constants.InteractiveFlag,
			Description:  "Activate the interactive menu",
			DefaultValue: false,
		},
		components.BoolFlag{
			Name:         constants.TailFlag,
			Description:  "Perform a 'tail " + constants.TailFlag + "' on the log",
			DefaultValue: false,
		},
	}
}

func getLogsEnvVar() []components.EnvVar {
	return []components.EnvVar{
		{
			Name:        constants.VersionCheckEnv,
			Default:     "true",
			Description: "Set this to \"false\" to disable validation on the minimum supported version of the product.",
		},
	}
}

func logsCmd(c *components.Context) error {
	isStreaming := c.GetBoolFlagValue(constants.TailFlag)
	isInteractive := c.GetBoolFlagValue(constants.InteractiveFlag)

	mainCtx, mainCtxCancel := context.WithCancel(context.Background())
	defer mainCtxCancel()

	var liveLogClient livelog.LiveLogs
	liveLogClient = livelog.NewLiveLogs()

	ListenForTermination(mainCtxCancel)

	if !isInteractive {
		if len(c.Arguments) != 4 {
			return fmt.Errorf("incorrect number of arguments were passed: expected: 4," + " received: " + strconv.Itoa(len(c.Arguments)))
		}
		productId := c.Arguments[0]
		serverId := c.Arguments[1]
		nodeId := c.Arguments[2]
		logFileName := c.Arguments[3]
		return liveLogClient.LogNonInteractive(mainCtx, productId, serverId, nodeId, logFileName, isStreaming)
	}
	return LogInteractiveMenu(mainCtx, isStreaming, liveLogClient)
}
