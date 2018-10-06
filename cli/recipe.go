package cli

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

func GenerateRecipeVerbCommands(c *Config, envVars []*string) []cli.Command {
	var validRecipeVerbs []*Verb
	var verbCommands []cli.Command
	for _, recipeVerb := range c.GroupVerbs {
		c.Verbs[recipeVerb].Name = recipeVerb
		validRecipeVerbs = append(validRecipeVerbs, c.Verbs[recipeVerb])
	}

	for _, verb := range validRecipeVerbs {
		verbCommands = append(verbCommands, GenerateRecipeVerbCommand(verb, c, envVars))
	}

	return verbCommands
}

func GenerateRecipeVerbCommand(verb *Verb, c *Config, envVars []*string) cli.Command {
	return cli.Command{
		Name:        verb.Name,
		Description: verb.Description,
		Usage:       verb.Usage,
		Category:    verb.Category,
		HideHelp:    true,
		BashComplete: func(ctx *cli.Context) {
			var completions []string

			for name, _ := range c.Groups {
				completions = append(completions, name)
			}

			fmt.Fprintf(ctx.App.Writer, strings.Join(completions, " "))
		},
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() < 1 || ctx.NArg() > 1 {
				cli.ShowCommandHelpAndExit(ctx, ctx.Command.Name, 0)
			}

			recipe := c.Groups[ctx.Args().First()]

			for _, instruction := range recipe.Members {
				for context, objects := range instruction {
					composeFile := fmt.Sprintf("%s.yaml", context)

					for _, command := range verb.Commands {
						if err := ExecuteDockerCommand(c.Home, envVars, []string{composeFile}, command, objects); err != nil {
							return err
						}
					}
				}
			}

			return nil
		},
	}
}
