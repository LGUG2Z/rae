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
		for key, value := range context.Env {
			envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
		}

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
