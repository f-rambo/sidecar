package utils

import (
	"bufio"
	"bytes"
	"os/exec"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/pkg/errors"
)

// exec command
func ExecCommand(log *log.Helper, command string, args ...string) (output string, err error) {
	log.Info("exec command: %s %s", command, strings.Join(args, " "))

	cmd := exec.Command(command, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	stdoutStr := stdout.String()
	stderrStr := stderr.String()

	if err != nil {
		return "", errors.Wrapf(err, "command failed: %s\nstdout: %s\nstderr: %s", command, stdoutStr, stderrStr)
	}

	if stderrStr != "" {
		return stdoutStr, errors.WithMessage(errors.New(stderrStr), "command failed")
	}

	log.Info(stdoutStr)

	return stdoutStr, nil
}

func RunCommandWithLogging(log *log.Helper, command string, args ...string) error {
	log.Info("exec command: %s %s", command, strings.Join(args, " "))
	cmd := exec.Command(command, args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "failed to get stdout pipe")
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return errors.Wrap(err, "failed to get stderr pipe")
	}

	if err := cmd.Start(); err != nil {
		return errors.Wrap(err, "failed to start command")
	}

	var stderrBuffer bytes.Buffer

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			log.Info(scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			log.Error(line)
			stderrBuffer.WriteString(line + "\n")
		}
	}()

	if err := cmd.Wait(); err != nil {
		return errors.Wrap(err, "command failed")
	}

	if stderrBuffer.Len() > 0 {
		return errors.Errorf("command wrote to stderr: %s", stderrBuffer.String())
	}

	return nil
}
