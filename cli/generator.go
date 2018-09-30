package cli

import (
	"fmt"
	"strings"

	"sort"

	"github.com/urfave/cli"
)

func GenerateRecipeCommand(recipe *Recipe, c *Config, envVars []string) cli.Command {
	return cli.Command{
		Name:        recipe.Name,
		Description: recipe.Description,
		Usage:       recipe.Usage,
		Category:    recipe.Category,
		Action: func(ctx *cli.Context) error {
			var err error
			for context, objects := range recipe.Instructions {
				composeFile := fmt.Sprintf("%s.yaml", context)
				if err = ExecuteDockerCommand(c.Home, envVars, []string{composeFile}, []string{"up", "-d"}, objects); err != nil {
					return err
				}
			}

			return nil
		},
	}
}

func GenerateContextCommands(c *Config, envVars []string) []cli.Command {
	var contextCommands []cli.Command

	for name, context := range c.Contexts {
		context.Name = name
		contextCommands = append(contextCommands, GenerateContextCommand(context, c, envVars))
	}

	return contextCommands
}

func GenerateContextCommand(context *Context, c *Config, envVars []string) cli.Command {
	var verbCommands []cli.Command

	if context.Name == "recipe" {
		for name, recipe := range c.Recipes {
			recipe.Name = name
			verbCommands = append(verbCommands, GenerateRecipeCommand(recipe, c, envVars))
		}
	} else {
		for name, verb := range c.Verbs {
			verb.Name = name
			verbCommands = append(verbCommands, GenerateVerbCommand(verb, c, envVars))
		}

		sort.Sort(cli.CommandsByName(verbCommands))
	}

	return cli.Command{
		Name:        context.Name,
		Description: context.Description,
		Usage:       context.Usage,
		Aliases:     context.Aliases,
		Category:    context.Category,
		Subcommands: verbCommands,
	}
}

func GenerateVerbCommand(verb *Verb, c *Config, envVars []string) cli.Command {
	return cli.Command{
		Name:        verb.Name,
		Description: verb.Description,
		Usage:       verb.Usage,
		Category:    verb.Category,
		Action: func(ctx *cli.Context) error {
			var composeFiles []string
			var err error

			context := strings.Split(ctx.Command.FullName(), " ")
			if context[0] == "global" {
				for _, context := range c.Contexts {
					if context.Name != "global" && context.Name != "recipe" {
						composeFiles = append(composeFiles, fmt.Sprintf("%s.yaml", context.Name))
					}
				}
			} else {
				composeFiles = append(composeFiles, fmt.Sprintf("%s.yaml", context[0]))
			}

			for _, command := range verb.Commands {
				if err = ExecuteDockerCommand(c.Home, envVars, composeFiles, command, ctx.Args()); err != nil {
					return err
				}
			}

			return nil
		},
	}
}
