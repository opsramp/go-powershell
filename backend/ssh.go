// Copyright (c) 2017 Gorillalabs. All rights reserved.

package backend

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"errors"
)

// sshSession exists so we don't create a hard dependency on crypto/ssh.
type sshSession interface {
	Waiter

	StdinPipe() (io.WriteCloser, error)
	StdoutPipe() (io.Reader, error)
	StderrPipe() (io.Reader, error)
	Start(string) error
}

type SSH struct {
	Session sshSession
}

func (b *SSH) StartProcess(cmd string, args ...string) (Waiter, io.Writer, io.Reader, io.Reader, error) {
	stdin, err := b.Session.StdinPipe()
	if err != nil {
		return nil, nil, nil, nil, errors.New("Could not get hold of the SSH session's stdin stream, error info:" + err.Error())
	}

	stdout, err := b.Session.StdoutPipe()
	if err != nil {
		return nil, nil, nil, nil, errors.New("Could not get hold of the SSH session's stdout stream, error info:" + err.Error())
	}

	stderr, err := b.Session.StderrPipe()
	if err != nil {
		return nil, nil, nil, nil, errors.New("Could not get hold of the SSH session's stderr stream, error info:" + err.Error())
	}

	err = b.Session.Start(b.createCmd(cmd, args))
	if err != nil {
		return nil, nil, nil, nil, errors.New("Could not spawn process via SSH, error info:" + err.Error())
	}

	return b.Session, stdin, stdout, stderr, nil
}

func (b *SSH) createCmd(cmd string, args []string) string {
	parts := []string{cmd}
	simple := regexp.MustCompile(`^[a-z0-9_/.~+-]+$`)

	for _, arg := range args {
		if !simple.MatchString(arg) {
			arg = b.quote(arg)
		}

		parts = append(parts, arg)
	}

	return strings.Join(parts, " ")
}

func (b *SSH) quote(s string) string {
	return fmt.Sprintf(`"%s"`, s)
}
