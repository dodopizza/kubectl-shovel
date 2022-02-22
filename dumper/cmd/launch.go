package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/dodopizza/kubectl-shovel/internal/events"
	"github.com/dodopizza/kubectl-shovel/internal/globals"
	"github.com/dodopizza/kubectl-shovel/internal/kubernetes"
	"github.com/dodopizza/kubectl-shovel/internal/utils"
	"github.com/dodopizza/kubectl-shovel/internal/watchdog"
)

func (cb *CommandBuilder) launch() error {
	events.NewStatusEvent("Looking for and mapping container fs")

	container := kubernetes.NewContainerInfoRaw(cb.CommonOptions.ContainerRuntime, cb.CommonOptions.ContainerID)
	// remove /tmp directory,
	// because will be mounted either from rootfs or container mounts
	if err := os.RemoveAll(globals.PathTmpFolder); err != nil {
		return err
	}

	tmpSource, err := container.GetTmpSource()
	if err != nil {
		events.NewErrorEvent(err, "unable to find mount point for container")
		return err
	}

	// for dotnet tools, in /tmp folder must exists sockets to running dotnet apps
	// https://github.com/dotnet/diagnostics/blob/main/documentation/design-docs/ipc-protocol.md#naming-and-location-conventions
	if err := os.Symlink(tmpSource, globals.PathTmpFolder); err != nil {
		events.NewErrorEvent(err, "unable to mount tmp folder for container")
		return err
	}

	// if we do not set proper file extension dotnet tools will do it anyway
	// write output file to /tmp, because it's available in target and worker pods
	output := fmt.Sprintf("%s/output.%s", globals.PathTmpFolder, cb.tool.ToolName())
	cb.tool.
		SetAction("collect").
		SetOutput(output)

	events.NewStatusEvent(
		fmt.Sprintf("Running command: %s %s",
			cb.tool.BinaryName(),
			strings.Join(cb.tool.FormatArgs(), " "),
		),
	)

	args := cb.tool.FormatArgs()
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

	if cb.CommonOptions.StoreOutputOnHost {
		outputHost := fmt.Sprintf("%s/%s.%s.%s.%s",
			globals.PathHostTmpFolder,
			cb.CommonOptions.PodNamespace,
			cb.CommonOptions.PodName,
			cb.CommonOptions.ContainerName,
			output,
			time.Now().UTC().Format("2006-04-02-15-04-05"),
		)

		if err := utils.MoveFile(output, outputHost); err != nil {
			events.NewErrorEvent(err, "failed to copy output on host")
			return err
		}

		events.NewCompletedEvent(outputHost)
		return nil
	}

	events.NewCompletedEvent(output)

	// otherwise, wait until output file to be copied from
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	watcher := watchdog.NewWatcher()
	if err := watcher.Run(ctx); err != nil {
		events.NewErrorEvent(err, "failed to watch copy progress")
		return err
	}

	return nil
}
