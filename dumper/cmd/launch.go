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
	events.NewStatusEvent("Looking for and mapping container fs")

	containerInfo := &kubernetes.ContainerInfo{
		Runtime: cb.ContainerOptions.Runtime,
		ID:      cb.ContainerOptions.ID,
	}
	containerFS, err := containerInfo.GetMountPoint()
	if err != nil {
		events.NewErrorEvent(err, "unable to find mount point for container")
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
	events.NewStatusEvent(
		fmt.Sprintf("Running command: %s %s",
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
		events.NewErrorEvent(err, "failed to execute tool command")
		return err
	}
	events.NewStatusEvent("Gathering completed")

	_, err = ioutil.ReadFile(output)
	if err != nil {
		events.NewErrorEvent(err, "failed to locate output result")
		return err
	}

	events.NewCompletedEvent(output)

	// wait until output file to be copied from
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := watchdog.Watch(ctx); err != nil {
		events.NewErrorEvent(err, "failed to watch copy progress")
		return err
	}

	return nil
}
