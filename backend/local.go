// Copyright (c) 2017 Gorillalabs. All rights reserved.

package backend

import (
	"io"
	"os/exec"

	"errors"
)

type Local struct{}

func (b *Local) StartProcess(cmd string, args ...string) (Waiter, io.Writer, io.Reader, io.Reader, error) {
	command := exec.Command(cmd, args...)

	stdin, err := command.StdinPipe()
	if err != nil {
		return nil, nil, nil, nil, errors.New("Could not get hold of the PowerShell's stdin stream, error info:" + err.Error())
	}

	stdout, err := command.StdoutPipe()
	if err != nil {
		return nil, nil, nil, nil, errors.New("Could not get hold of the PowerShell's stdout stream, error info:" + err.Error())
	}

	stderr, err := command.StderrPipe()
	if err != nil {
		return nil, nil, nil, nil, errors.New("Could not get hold of the PowerShell's stderr stream, error info:" + err.Error())
	}

	err = command.Start()
	if err != nil {
		return nil, nil, nil, nil, errors.New("Could not spawn PowerShell process, error info:" + err.Error())
	}

	return command, stdin, stdout, stderr, nil
}
