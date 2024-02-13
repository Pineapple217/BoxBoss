package docker

import (
	"context"
	"io"
	"os"
	"os/exec"

	"github.com/creack/pty"
)

func ComposePull(cf, service string) error {
	cmd := exec.CommandContext(context.TODO(),
		"docker", "compose",
		"--file", cf,
		"pull", service,
	)
	f, err := pty.StartWithSize(cmd, &pty.Winsize{Rows: 10, Cols: 120})
	if err != nil {
		return err
	}
	defer f.Close()
	io.Copy(os.Stdout, f)
	return nil
}

func ComposeStop(cf, service string) error {
	cmd := exec.CommandContext(context.TODO(),
		"docker", "compose",
		"--file", cf,
		"stop", service,
	)
	f, err := pty.StartWithSize(cmd, &pty.Winsize{Rows: 10, Cols: 120})
	if err != nil {
		return err
	}
	defer f.Close()
	io.Copy(os.Stdout, f)
	return nil
}

func ComposeRemove(cf, service string) error {
	cmd := exec.CommandContext(context.TODO(),
		"docker", "compose",
		"--file", cf,
		"rm", "-f", service,
	)
	f, err := pty.StartWithSize(cmd, &pty.Winsize{Rows: 10, Cols: 120})
	if err != nil {
		return err
	}
	defer f.Close()
	io.Copy(os.Stdout, f)
	return nil
}

func ComposeUp(cf, service string) error {
	cmd := exec.CommandContext(context.TODO(),
		"docker", "compose",
		"--file", cf,
		"up", "-d", service,
	)
	f, err := pty.StartWithSize(cmd, &pty.Winsize{Rows: 10, Cols: 120})
	if err != nil {
		return err
	}
	defer f.Close()
	io.Copy(os.Stdout, f)
	return nil
}
