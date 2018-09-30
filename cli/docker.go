package cli

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"os/signal"
)

func ExecuteDockerCommand(home string, envVars, composeFiles []string, command []string, objects []string) error {
	cmd := exec.Command("docker-compose")
	cmd.Dir = home

	for _, envVar := range envVars {
		cmd.Env = append(cmd.Env, envVar)
	}

	for _, composeFile := range composeFiles {
		cmd.Args = append(cmd.Args, "-f")
		cmd.Args = append(cmd.Args, composeFile)
	}

	for _, object := range objects {
		if object == "all" {
			objects = []string{}
		}
	}

	cmd.Args = append(cmd.Args, command...)
	cmd.Args = append(cmd.Args, objects...)

	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw
	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		return err
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		if err := cmd.Process.Signal(os.Interrupt); err != nil {
			panic(err)
		}
	}()
	cmd.Wait()

	return nil
}
