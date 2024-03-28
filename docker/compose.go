package docker

import (
	"context"
	"io"
	"log/slog"
	"os/exec"
	"strings"

	"github.com/creack/pty"
)

const (
	winRows = 10
	winCols = 120
)

func ComposePull(cf, service string, output io.Writer) error {
	cmd := exec.CommandContext(context.TODO(),
		"docker", "compose",
		"--file", cf,
		"pull", service,
	)
	f, err := pty.StartWithSize(cmd, &pty.Winsize{Rows: winRows, Cols: winCols})
	if err != nil {
		return err
	}
	defer f.Close()
	io.Copy(output, f)
	// io.Copy(os.Stdout, f)
	return nil
}

func ComposeStop(cf, service string, output io.Writer) error {
	cmd := exec.CommandContext(context.TODO(),
		"docker", "compose",
		"--file", cf,
		"stop", service,
	)
	f, err := pty.StartWithSize(cmd, &pty.Winsize{Rows: winRows, Cols: winCols})
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(output, f)
	if err != nil {
		slog.Info("aaa", "err", err)
	}
	// io.Copy(os.Stdout, f)
	return nil
}

func ComposeRemove(cf, service string, output io.Writer) error {
	cmd := exec.CommandContext(context.TODO(),
		"docker", "compose",
		"--file", cf,
		"rm", "-f", service,
	)
	f, err := pty.StartWithSize(cmd, &pty.Winsize{Rows: winRows, Cols: winCols})
	if err != nil {
		return err
	}
	defer f.Close()
	io.Copy(output, f)
	// io.Copy(os.Stdout, f)
	return nil
}

func ComposeUp(cf, service string, output io.Writer) error {
	cmd := exec.CommandContext(context.TODO(),
		"docker", "compose",
		"--file", cf,
		"up", "-d", service,
	)
	f, err := pty.StartWithSize(cmd, &pty.Winsize{Rows: winRows, Cols: winCols})
	if err != nil {
		return err
	}
	defer f.Close()
	io.Copy(output, f)
	// io.Copy(os.Stdout, f)
	return nil
}

func GetServiceContainerId(cf, service string) (string, error) {
	cmd := exec.CommandContext(context.TODO(),
		"docker", "compose",
		"--file", cf,
		"ps", "-q", service,
	)
	f, err := pty.StartWithSize(cmd, &pty.Winsize{Rows: 1, Cols: 200})
	if err != nil {
		return "", err
	}
	defer f.Close()
	// TODO: allways errors but stil works?
	b, _ := io.ReadAll(f)
	id := strings.TrimSpace(string(b))
	return id, nil
}
