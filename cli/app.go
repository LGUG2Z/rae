package cli

import (
	"sort"
	"time"

	"fmt"

	"github.com/urfave/cli"
)

var (
	Version string
	Commit  string
)

func App(c *Config, envVars []*string) *cli.App {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("rae version %s (commit %s)\n", c.App.Version, Commit)
	}

	app := cli.NewApp()
	app.Name = "rae"
	app.Usage = "A docker-compose development environment orchestrator"
	app.UsageText = "rae [global options] context [context options] verb [verb options] [objects...]"

	app.Version = Version
	app.Compiled = time.Now()
	app.Authors = []cli.Author{{
		Name:  "J. Iqbal",
		Email: "jade@beamery.com",
	}}

	app.EnableBashCompletion = true
	app.HideHelp = true

	app.Commands = append(app.Commands, GenerateContextCommands(c, envVars)...)
	sort.Sort(cli.CommandsByName(app.Commands))

	return app
}
