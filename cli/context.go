package cli

import (
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

	switch context.Name {
	case "recipe":
		verbCommands = GenerateRecipeVerbCommands(c, envVars)
	default:
		verbCommands = GenerateVerbCommands(c, envVars)
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
