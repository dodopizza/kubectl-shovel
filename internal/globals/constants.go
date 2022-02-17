package globals

import "fmt"

const (
	PluginName      = "kubectl-shovel"
	DumperImageName = "dodopizza/kubectl-shovel-dumper"
)

func GetDumperImage() string {
	return fmt.Sprintf("%s:%s", DumperImageName, GetVersion())
}
