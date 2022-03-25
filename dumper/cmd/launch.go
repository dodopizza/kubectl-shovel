package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dodopizza/kubectl-shovel/internal/events"
	"github.com/dodopizza/kubectl-shovel/internal/flags"
	"github.com/dodopizza/kubectl-shovel/internal/globals"
	"github.com/dodopizza/kubectl-shovel/internal/kubernetes"
	"github.com/dodopizza/kubectl-shovel/internal/utils"
	"github.com/dodopizza/kubectl-shovel/internal/watchdog"
)

func (cb *CommandBuilder) launch() error {
	events.NewStatusEvent("Looking for and mapping container fs")

	container, err := kubernetes.NewContainerConfigInfo(cb.CommonOptions.ContainerRuntime, cb.CommonOptions.ContainerID)
	if err != nil {
		return err
	}
	// remove /tmp directory,
	// because it will be mounted from container /tmp directory
	if err := os.RemoveAll(globals.PathTmpFolder); err != nil {
		return err
	}

	// for dotnet tools, in /tmp folder must exists sockets to running dotnet apps
	// https://github.com/dotnet/diagnostics/blob/main/documentation/design-docs/ipc-protocol.md#diagnostic-ipc-protocol
	if err := os.Symlink(container.GetTmpSource(), globals.PathTmpFolder); err != nil {
		events.NewErrorEvent(err, "unable to mount tmp folder for container")
		return err
	}

	// write output file to /tmp, because it's available in target and worker pods
	output := fmt.Sprintf("%s/output.%s", globals.PathTmpFolder, cb.tool.ToolName())
	cb.tool.SetOutput(output)

	if cb.tool.IsPrivileged() {
		// host process required for privileged commands
		cb.tool.SetProcessID(container.HostProcessID)

		// exact same versions required
		resolver := flags.NewDotnetToolResolver(container.RootFS)
		frameworks, err := resolver.LocateFrameworks()
		if err != nil {
			events.NewErrorEvent(err, "unable to locate runtime folders for target container")
			return err
		}

		for _, framework := range frameworks {
			if framework.Name != flags.DotnetFrameworkApp {
				continue
			}

			err := os.Symlink(framework.FullPath(), filepath.Join(resolver.Path, framework.NameVersion()))
			if err != nil {
				events.NewErrorEvent(err, "unable to mount runtime folders for target container")
				return err
			}
		}
	}

	args := flags.NewArgs()
	cb.tool.FormatArgs(args, flags.FormatArgsTypeBinary)

	events.NewStatusEvent(
		fmt.Sprintf("Running command: %s %s",
			cb.tool.BinaryName(),
			strings.Join(args.Get(), " "),
		),
	)

	if err := utils.ExecCommand(cb.tool.BinaryName(), args.Get()...); err != nil {
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
		outputHost := fmt.Sprintf("%s/%s.%s.%s.%s.%s",
			globals.PathHostOutputFolder,
			cb.CommonOptions.PodNamespace,
			cb.CommonOptions.PodName,
			cb.CommonOptions.ContainerName,
			filepath.Base(output),
			time.Now().UTC().Format("2006-01-02-15-04-05"),
		)

		if err := utils.MoveFile(output, outputHost); err != nil {
			events.NewErrorEvent(err, "failed to copy output on host")
			return err
		}

		events.NewCompletedEvent(filepath.Base(outputHost))
		return nil
	}

	events.NewCompletedEvent(output)

	// wait until output file to be copied from
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	watcher := watchdog.NewWatcher()
	if err := watcher.Run(ctx); err != nil {
		events.NewErrorEvent(err, "failed to watch copy progress")
		return err
	}

	return nil
}
