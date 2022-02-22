package globals

import (
	"fmt"
)

const (
	PluginName            = "kubectl-shovel"
	DumperImageName       = "dodopizza/kubectl-shovel-dumper"
	PathTmpFolder         = "/tmp"
	PathHostTmpFolder     = "/host-tmp"
	PathContainerDFS      = "/run/containerd"
	PathContainerDVolumes = "/var/lib/kubelet/pods"
	PathDockerFS          = "/var/lib/docker"
	PathDockerVolumes     = PathDockerFS
)

func GetDumperImage() string {
	return fmt.Sprintf("%s:%s", DumperImageName, GetVersion())
}
