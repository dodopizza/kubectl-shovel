package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

const (
	procPath = "/host/proc"
)

var (
	kubePattern   = regexp.MustCompile(`\d+:.+:/kubepods/[^/]+/pod[^/]+/([0-9a-f]{64})`)
	dockerPattern = regexp.MustCompile(`\d+:.+:/docker/pod[^/]+/([0-9a-f]{64})`)
)

func findPID(containerID, path string) (int, error) {
	if path == "" {
		path = procPath
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return 0, errors.Wrap(err, "Error while reading dir")
	}
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		pid, err := strconv.Atoi(file.Name())
		if err != nil {
			continue
		}
		cgroupFileName := filepath.Join(
			path,
			file.Name(),
			"cgroup",
		)
		f, err := os.Open(cgroupFileName)
		if err != nil {
			if err == os.ErrNotExist {
				continue
			}
			return 0, errors.Wrap(err, "Can't open files")
		}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			parts := dockerPattern.FindStringSubmatch(line)
			if parts != nil {
				if parts[1] == containerID {
					return pid, nil
				}
			}
			parts = kubePattern.FindStringSubmatch(line)
			if parts != nil {
				if parts[1] == containerID {
					return pid, nil
				}
			}
		}
	}

	return 0, errors.New("Container ID was not found")
}
