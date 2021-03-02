package commands

import (
	"context"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/jfrog/live-logs/internal"
	"github.com/jfrog/live-logs/internal/constants"
	"strconv"
)

func GetLogsCommand() components.Command {
	return components.Command{
		Name:        "logs",
		Description: "Fetch the log of a desired service",
		Aliases:     []string{"l"},
		Arguments:   getLogsArguments(),
		Flags:       getLogsFlags(),
		Action:      logsCmd,
	}
}

func getLogsArguments() []components.Argument {
	return []components.Argument{
		{Name: "product-id", Description: "JFrog product id, value can be one of the following, \n" +
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
			Name:         "i",
			Description:  "Activate interactive menu",
			DefaultValue: false,
		},
		components.BoolFlag{
			Name:         "f",
			Description:  "Do 'tail -f' on the log",
			DefaultValue: false,
		},
	}
}

func logsCmd(c *components.Context) error {
	isStreaming := c.GetBoolFlagValue("f")
	isInteractive := c.GetBoolFlagValue("i")

	mainCtx, mainCtxCancel := context.WithCancel(context.Background())
	defer mainCtxCancel()

	var liveLogClient livelog.LiveLogs
	liveLogClient = livelog.NewLiveLogs()

	ListenForTermination(mainCtxCancel)

	if !isInteractive {
		if len(c.Arguments) != 4 {
			return fmt.Errorf("wrong number of arguments. Expected: 4, " + "Received: " + strconv.Itoa(len(c.Arguments)))
		}
		productId := c.Arguments[0]
		serverId := c.Arguments[1]
		nodeId := c.Arguments[2]
		logFileName := c.Arguments[3]
		return liveLogClient.LogNonInteractive(mainCtx, productId, serverId, nodeId, logFileName, isStreaming)
	}
	return LogInteractiveMenu(mainCtx, isStreaming, liveLogClient)
}
