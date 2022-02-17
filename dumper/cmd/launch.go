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

func launch(
	commonOptions commonOptions,
	executable string,
	args ...string,
) error {
	events.NewEvent(
		events.Status,
		"Looking for and mapping container fs",
	)

	containerInfo := &kubernetes.ContainerInfo{
		Runtime: commonOptions.containerRuntime,
		ID:      commonOptions.containerID,
	}
	containerFS, err := containerInfo.GetMountPoint()
	if err != nil {
		return err
	}
	err = os.RemoveAll("/tmp")
	if err != nil {
		return err
	}
	if err := os.Symlink(filepath.Join(containerFS, "tmp"), "/tmp"); err != nil {
		return err
	}

	events.NewEvent(
		events.Status,
		fmt.Sprintf(
			"Running command: %s %s",
			executable,
			strings.Join(args, " "),
		),
	)

	// if we do not set proper file extension dotnet tools will do it anyway
	// write output file to /tmp, because it's available in target and worker pods
	outputExtension := strings.TrimPrefix(executable, "dotnet-")
	output := fmt.Sprintf("/tmp/output.%s", outputExtension)
	args = append(
		args,
		"--output",
		output,
	)
	if err := utils.ExecCommand(executable, args...); err != nil {
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
