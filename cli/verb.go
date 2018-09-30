package cli

import (
	"fmt"
	"strings"

	"path"

	"io/ioutil"

	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

func GenerateVerbCommand(verb *Verb, c *Config, envVars []string) cli.Command {
	return cli.Command{
		Name:        verb.Name,
		Description: verb.Description,
		Usage:       verb.Usage,
		Category:    verb.Category,
		BashComplete: func(ctx *cli.Context) {
			var completions []string
			var err error

			context := strings.Split(ctx.Command.FullName(), " ")
			if context[0] == "global" {
				var files []string
				if context[0] == "global" {
					for _, context := range c.Contexts {
						if context.Name != "global" && context.Name != "recipe" {
							files = append(files, path.Join(c.Home, fmt.Sprintf("%s.yaml", context.Name)))
						}
					}
				}

				for _, file := range files {
					completions, err = appendCompletions(file, completions)
					if err != nil {
						panic(err)
					}
				}
			} else {
				file := path.Join(c.Home, fmt.Sprintf("%s.yaml", context[0]))
				completions, err = appendCompletions(file, completions)
			}

			fmt.Fprintf(ctx.App.Writer, strings.Join(completions, " "))
		},
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

func appendCompletions(file string, completions []string) ([]string, error) {
	composeFile := map[interface{}]interface{}{}
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(bytes, composeFile); err != nil {
		return nil, err
	}

	for service, _ := range composeFile["services"].(map[interface{}]interface{}) {
		completions = append(completions, service.(string))
	}

	return completions, nil
}
