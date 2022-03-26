package flags

import (
	"os"
	"path/filepath"
)

const (
	DotnetFrameworkApp = "Microsoft.NETCore.App"
)

type DotnetToolResolver struct {
	Root string
	Path string
}

type DotnetFramework struct {
	Root    string
	Name    string
	Version string
}

func NewDotnetToolResolver(root string) *DotnetToolResolver {
	return &DotnetToolResolver{
		Root: root,
		Path: "/usr/share/dotnet/shared",
	}
}

func (r *DotnetToolResolver) LocateFrameworks() ([]DotnetFramework, error) {
	entries, err := os.ReadDir(r.FullPath())
	if err != nil {
		return nil, err
	}

	var frameworks []DotnetFramework
	for _, entry := range entries {
		framework := entry.Name()
		frameworkPath := filepath.Join(r.FullPath(), framework)

		entries, err := os.ReadDir(frameworkPath)
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			frameworks = append(frameworks,
				DotnetFramework{
					Name:    framework,
					Root:    r.FullPath(),
					Version: entry.Name(),
				},
			)
		}
	}

	return frameworks, nil
}

func (r *DotnetToolResolver) FullPath() string {
	return filepath.Join(r.Root, r.Path)
}

func (f *DotnetFramework) FullPath() string {
	return filepath.Join(f.Root, f.Name, f.Version)
}

func (f *DotnetFramework) NameVersion() string {
	return filepath.Join(f.Name, f.Version)
}
