package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/dodopizza/kubectl-shovel/internal/events"
	"github.com/dodopizza/kubectl-shovel/internal/kubernetes"
	"github.com/dodopizza/kubectl-shovel/internal/utils"
	"github.com/dodopizza/kubectl-shovel/internal/watchdog"
)

func (cb *CommandBuilder) launch() error {
	events.NewEvent(
		events.Status,
		"Looking for and mapping container fs",
	)

	containerInfo := &kubernetes.ContainerInfo{
		Runtime: cb.CommonOptions.containerRuntime,
		ID:      cb.CommonOptions.containerID,
	}
	containerFS, err := containerInfo.GetMountPoint()
	if err != nil {
		return err
	}
	if err := os.RemoveAll("/tmp"); err != nil {
		return err
	}
	if err := os.Symlink(filepath.Join(containerFS, "tmp"), "/tmp"); err != nil {
		return err
	}

	args := []string{"collect"}
	args = append(args, cb.tool.GetArgs()...)
	events.NewEvent(
		events.Status,
		fmt.Sprintf(
			"Running command: %s %s",
			cb.tool.BinaryName(),
			strings.Join(args, " "),
		),
	)

	// if we do not set proper file extension dotnet tools will do it anyway
	// write output file to /tmp, because it's available in target and worker pods
	output := fmt.Sprintf("/tmp/output.%s", cb.tool.ToolName())
	args = append(
		args,
		"--output",
		output,
	)
	if err := utils.ExecCommand(cb.tool.BinaryName(), args...); err != nil {
		return err
	}

	events.NewEvent(
		events.Status,
		"Gathering completed",
	)

	_, err = ioutil.ReadFile(output)
	if err != nil {
		return err
	}

	events.NewEvent(
		events.Completed,
		output,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := watchdog.Watch(ctx); err != nil {
		events.NewEvent(
			events.Error,
			err.Error(),
		)
		return err
	}

	return nil
}
