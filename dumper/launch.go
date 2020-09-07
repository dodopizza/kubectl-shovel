package main

import (
	"io/ioutil"
	"log"
)

const (
	output = "/output.gcdump"
)

func launch(containerID string) error {
	err := mapContainerTmp(containerID)
	if err != nil {
		return err
	}
	err = makeGcDump(1, output)
	if err != nil {
		return err
	}
	dumpContent, err := ioutil.ReadFile(output)
	if err != nil {
		return err
	}
	log.Println(string(dumpContent))

	return nil
}
