package cli

import (
	"sort"
	"time"

	"github.com/urfave/cli"
	"fmt"
)

var (
	Version string
	Build   string
)

func App(c *Config, envVars []string) *cli.App {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("rae version %s (build %s)\n", c.App.Version, Build)
	}

	app := cli.NewApp()

	app.Name = "rae"
	app.Version = Version
	app.Usage = "A docker-compose development environment orchestrator"
	app.EnableBashCompletion = true
	app.Compiled = time.Now()
	app.Authors = []cli.Author{{
		Name:  "J. Iqbal",
		Email: "jade@beamery.com",
	}}

	app.UsageText = "rae [global options] context [context options] verb [verb options] [objects...]"

	app.Commands = append(app.Commands, GenerateContextCommands(c, envVars)...)
	sort.Sort(cli.CommandsByName(app.Commands))

	return app
}
