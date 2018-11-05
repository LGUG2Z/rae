package cli

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"

	"fmt"

	"github.com/fatih/color"
)

func ExecuteHealthCheck(object string) error {
	ps := exec.Command("docker")

	dockerHost, ok := os.LookupEnv("DOCKER_HOST")
	if ok {
		ps.Env = append(ps.Env, fmt.Sprintf("DOCKER_HOST=%s", dockerHost))
	}

	ps.Args = append(ps.Args, "ps", "-a", "-q", "--filter", fmt.Sprintf("name=_%s_", object))
	b, err := ps.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error getting container id: %s", err)
	}

	containerId := strings.TrimSpace(string(b))

	inspect := exec.Command("docker")
	inspect.Args = append(inspect.Args, "inspect", containerId)
	b, err = inspect.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error inspecting container: %s (%s)", err, string(b))
	}

	c := Containers{}
	if err := json.Unmarshal(b, &c); err != nil {
		return err
	}

	fmt.Printf("Waiting for %s to pass healthcheck ... ", strings.TrimPrefix(c[0].Name, "/"))
	healthy := c[0].State.Health.Status == "healthy"

	for !healthy {
		c = nil

		inspect := exec.Command("docker")
		inspect.Args = append(inspect.Args, "inspect", containerId)
		b, err := inspect.CombinedOutput()
		if err != nil {
			return fmt.Errorf("error inspecting container: %s (%s)", err, string(b))
		}

		c = Containers{}
		if err := json.Unmarshal(b, &c); err != nil {
			return err
		}

		if c[0].State.Health.Status == "healthy" {
			healthy = true
		} else {
			time.Sleep(5 * time.Second)
		}
	}
	color.Green("done")
	return nil
}

func ExecuteComposeCommand(home string, envVars []*string, command []string, objects []string) error {
	cmd := exec.Command("docker-compose")
	cmd.Dir = home

	p, _ := os.LookupEnv("PATH")
	cmd.Env = append(cmd.Env, fmt.Sprintf("PATH=%s", p))

	dockerHost, ok := os.LookupEnv("DOCKER_HOST")
	if ok {
		cmd.Env = append(cmd.Env, fmt.Sprintf("DOCKER_HOST=%s", dockerHost))
	}

	for _, envVar := range envVars {
		cmd.Env = append(cmd.Env, *envVar)
	}

	var cleanedObjects []string
	for _, object := range objects {
		if object == "all" {
			break
		}

		cleanedObjects = append(cleanedObjects, strings.Split(object, " ")...)
	}

	cmd.Args = append(cmd.Args, command...)
	cmd.Args = append(cmd.Args, cleanedObjects...)

	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw
	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		return err
	}

	signalChan := make(chan os.Signal, 1)
	doneChan := make(chan struct{}, 1)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		for {
			select {
			case sig := <-signalChan:
				if err := cmd.Process.Signal(sig); err != nil {
					panic(err)
				}
				return
			case <-doneChan:
				close(doneChan)
				close(signalChan)
				signal.Reset(os.Interrupt)
				return
			}
		}
	}()

	cmd.Wait()
	doneChan <- struct{}{}

	return nil
}
