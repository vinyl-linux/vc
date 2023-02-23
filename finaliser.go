package vc

import (
	"os/exec"

	"github.com/google/shlex"
)

var (
	finaliserNothingCommand  = "echo complete!"
	finaliserRebootCommand   = "vinitctl reboot"
	finaliserShutdownCommand = "vinitctl shutdown"
)

func (c Config) Finalise() error {
	var (
		cmd        string
		binary     string
		binaryArgs []string
	)

	switch c.Finaliser {
	case Finaliser_Nothing:
		cmd = finaliserNothingCommand
	case Finaliser_Reboot:
		cmd = finaliserRebootCommand
	case Finaliser_Shutdown:
		cmd = finaliserShutdownCommand
	}

	args, err := shlex.Split(cmd)
	if err != nil {
		return err
	}

	switch len(args) {
	case 0:
		return nil

	case 1:
		binary = args[0]

	default:
		binary = args[0]
		binaryArgs = args[1:]
	}

	// #nosec: G204
	command := exec.Command(binary, binaryArgs...)

	return command.Run()
}
