package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"path"

	"github.com/LGUG2Z/rae/cli"
	"github.com/joho/godotenv"
)

func main() {
	if err := preFlightChecks(); err != nil {
		log.Fatal(err)
	}

	home := os.Getenv("RAE_HOME")

	c := &cli.Config{}
	if err := c.Load(home); err != nil {
		log.Fatal(err)
	}

	o := &cli.Config{}
	if err := o.LoadOverride(home); err != nil {
		log.Fatal(err)
	}

	if err := c.MergeOverride(o); err != nil {
		log.Fatal(err)
	}

	var envVars []string
	var err error

	if len(c.EnvFiles) > 0 {
		envVars, err = collectEnvVars(c)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := cli.App(c, envVars).Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func preFlightChecks() error {
	_, ok := os.LookupEnv("RAE_HOME")
	if !ok {
		return fmt.Errorf("RAE_HOME is not set")
	}

	applications := []string{"docker", "docker-compose"}

	for _, application := range applications {
		if !isInstalled(application) {
			return fmt.Errorf("%s is not installed", application)
		}
	}

	return nil
}

func collectEnvVars(c *cli.Config) ([]string, error) {
	var envVars []string

	for index, value := range c.EnvFiles {
		c.EnvFiles[index] = path.Join(c.Home, value)
	}

	envMap, err := godotenv.Read(c.EnvFiles...)
	if err != nil {
		return nil, err
	}

	for key, value := range envMap {
		envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
	}

	for key, value := range c.Env {
		envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
	}

	return envVars, nil
}

func isInstalled(application string) bool {
	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("command -v %s", application))
	if err := cmd.Run(); err != nil {
		return false
	}

	return true
}
