package cli

import (
	"sort"

	"fmt"

	"strings"

	"github.com/urfave/cli"
)

func GenerateContextCommands(c *Config, envVars []*string) []cli.Command {
	var contextCommands []cli.Command

	for name, context := range c.Contexts {
		context.Name = name
		contextCommands = append(contextCommands, GenerateContextCommand(context, c, envVars))
	}

	return contextCommands
}

func GenerateContextCommand(context *Context, c *Config, envVars []*string) cli.Command {
	var verbCommands []cli.Command

	if context.Name == "recipe" {
		var validRecipeVerbs []*Verb
		for _, recipeVerb := range c.RecipeVerbs {
			c.Verbs[recipeVerb].Name = recipeVerb
			validRecipeVerbs = append(validRecipeVerbs, c.Verbs[recipeVerb])
		}

		for _, verb := range validRecipeVerbs {
			verbCommands = append(verbCommands, GenerateRecipeVerbCommand(verb, c, envVars))
		}
	} else {

		for name, verb := range c.Verbs {
			verb.Name = name
			verbCommands = append(verbCommands, GenerateVerbCommand(verb, c, envVars))
		}

		sort.Sort(cli.CommandsByName(verbCommands))
	}

	var flags []cli.Flag
	for flag, envMap := range context.EnvFlags {
		var envMapUsage []string
		for key, value := range *envMap {
			envMapUsage = append(envMapUsage, fmt.Sprintf("%s=%s", key, *value))
		}

		flags = append(
			flags,
			cli.BoolFlag{
				Name:  flag,
				Usage: fmt.Sprintf("Sets %s", strings.Join(envMapUsage, ", ")),
			},
		)
	}

	return cli.Command{
		Aliases:     context.Aliases,
		Category:    context.Category,
		Description: context.Description,
		Flags:       flags,
		HideHelp:    true,
		Name:        context.Name,
		Subcommands: verbCommands,
		Usage:       context.Usage,
	}
}
