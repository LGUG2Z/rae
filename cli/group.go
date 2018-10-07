package cli

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

func GenerateGroupVerbCommands(c *Config, envVars []*string) []cli.Command {
	var validGroupVerbs []*Verb
	var verbCommands []cli.Command
	for _, groupVerb := range c.GroupVerbs {
		c.Verbs[groupVerb].Name = groupVerb
		validGroupVerbs = append(validGroupVerbs, c.Verbs[groupVerb])
	}

	for _, verb := range validGroupVerbs {
		verbCommands = append(verbCommands, GenerateGroupVerbCommand(verb, c, envVars))
	}

	return verbCommands
}

func GenerateGroupVerbCommand(verb *Verb, c *Config, envVars []*string) cli.Command {
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
			return GroupVerbCommandAction(ctx, c, ctx.Args().First(), verb, envVars)
		},
	}
}

func GroupVerbCommandAction(ctx *cli.Context, c *Config, g string, verb *Verb, envVars []*string) error {
	context := strings.Split(ctx.Command.FullName(), " ")

	if context[0] != "recipe" {
		if ctx.NArg() < 1 || ctx.NArg() > 1 {
			cli.ShowCommandHelpAndExit(ctx, ctx.Command.Name, 0)
		}
	}

	group := c.Groups[g]

	for _, instruction := range group.Members {
		for context, objects := range instruction {
			for key, value := range c.Contexts[context].Env {
				envVar := fmt.Sprintf("%s=%s", key, *value)
				envVars = append(envVars, &envVar)
			}

			for _, command := range verb.Commands {
				if err := ExecuteComposeCommand(c.Home, envVars, command, objects); err != nil {
					return err
				}
			}
		}
	}

	return nil

}
