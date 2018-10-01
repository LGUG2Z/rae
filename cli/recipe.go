package cli

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

func GenerateRecipeVerbCommand(verb *Verb, c *Config, envVars []*string) cli.Command {
	return cli.Command{
		Name:        verb.Name,
		Description: verb.Description,
		Usage:       verb.Usage,
		Category:    verb.Category,
		HideHelp:    true,
		BashComplete: func(ctx *cli.Context) {
			var completions []string

			for name, _ := range c.Recipes {
				completions = append(completions, name)
			}

			fmt.Fprintf(ctx.App.Writer, strings.Join(completions, " "))
		},
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() < 1 || ctx.NArg() > 1 {
				cli.ShowCommandHelpAndExit(ctx, ctx.Command.Name, 0)
			}

			recipe := c.Recipes[ctx.Args().First()]

			for context, objects := range recipe.Instructions {
				composeFile := fmt.Sprintf("%s.yaml", context)

				for _, command := range verb.Commands {
					if err := ExecuteDockerCommand(c.Home, envVars, []string{composeFile}, command, objects); err != nil {
						return err
					}
				}
			}

			return nil
		},
	}
}
