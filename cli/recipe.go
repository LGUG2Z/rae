package cli

import (
	"sort"

	"github.com/urfave/cli"
)

func GenerateRecipeCommands(c *Config, envVars []*string) []cli.Command {
	var recipeCommands []cli.Command

	for name, recipe := range c.Recipes {
		recipe.Name = name
		recipeCommands = append(recipeCommands, GenerateRecipeCommand(recipe, c, envVars))
	}

	sort.Sort(cli.CommandsByName(recipeCommands))

	return recipeCommands
}

func GenerateRecipeCommand(recipe *Recipe, c *Config, envVars []*string) cli.Command {
	return cli.Command{
		Name:        recipe.Name,
		Description: recipe.Description,
		Usage:       recipe.Usage,
		HideHelp:    true,
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() != 0 {
				cli.ShowCommandHelpAndExit(ctx, ctx.Command.Name, 0)
			}

			for _, instructions := range recipe.Instructions {
				if instructions.Healthcheck != nil {

				}

				if err := GroupVerbCommandAction(ctx, c, instructions.Group, c.Verbs[instructions.Verb], envVars); err != nil {
					return err
				}
			}

			return nil
		},
	}
}
