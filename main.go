package main

import (
	"github.com/jfrog/jfrog-cli-core/plugins"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/jfrog/live-logs/commands"
)

func main() {
	plugins.PluginMain(getApp())
}

func getApp() components.App {
	app := components.App{}
	app.Name = "live-logs"
	app.Description = "Print logs from a remote JFrog product."
	app.Version = "v1.0.2"
	app.Commands = getCommands()
	return app
}


func getCommands() []components.Command {
	return []components.Command{
		commands.GetLogsCommand(),
		commands.GetConfigCommand(),
	}
}
