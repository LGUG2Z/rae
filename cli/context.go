package cli

import (
	"sort"

	"fmt"

	"github.com/urfave/cli"
)

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
		var validRecipeVerbs []*Verb
		for _, recipeVerb := range c.RecipeVerbs {
			c.Verbs[recipeVerb].Name = recipeVerb
			validRecipeVerbs = append(validRecipeVerbs, c.Verbs[recipeVerb])
		}

		for _, verb := range validRecipeVerbs {
			verbCommands = append(verbCommands, GenerateRecipeVerbCommand(verb, c, envVars))
		}
	} else {
		for key, value := range context.Env {
			envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
		}

		for name, verb := range c.Verbs {
			verb.Name = name
			verbCommands = append(verbCommands, GenerateVerbCommand(verb, c, envVars))
		}

		sort.Sort(cli.CommandsByName(verbCommands))
	}

	var flags []cli.Flag
	for flag, _ := range context.EnvFlags {
		flags = append(flags, cli.BoolFlag{Name: flag})
	}

	return cli.Command{
		Aliases:     context.Aliases,
		Category:    context.Category,
		Description: context.Description,
		Flags:       flags,
		Name:        context.Name,
		Subcommands: verbCommands,
		Usage:       context.Usage,
	}
}
