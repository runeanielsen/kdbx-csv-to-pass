package parser

import (
	"log"
	"os/exec"
)

type CmdExec interface {
	ExecuteCmd(cmd string)
}

type CmdExecBash struct{}

func (CmdExecBash) ExecuteCmd(cmd string) {
	err := exec.Command("bash", "-c", cmd).Run()

	if err != nil {
		log.Fatalf("Failed to execute command: %s - %s\n", cmd, err)
	}
}
