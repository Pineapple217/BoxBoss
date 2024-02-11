package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/Pineapple217/harbor-hawk/database"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
)

var (
	cli *client.Client
)

func GetClient() *client.Client {
	return cli
}

func Init() {
	var err error
	cli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
}

func Ps() []types.Container {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		panic(err)
	}

	// for _, container := range containers {
	// 	fmt.Sprintf("%s %s\n", container.ID[:10], container.Image)
	// }
	return containers
}

// , removeOldImage bool
func UpdateContainer(containerID string) error {
	if cli == nil {
		panic("docker client is nil")
	}

	// Get container details
	fmt.Println("container details")
	containerInfo, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return err
	}

	// Extract image name from container details
	oldImageName := containerInfo.Config.Image

	// Pull the new image
	out, err := cli.ImagePull(context.Background(), oldImageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()

	// Stop the container
	err = cli.ContainerStop(context.Background(), containerID, container.StopOptions{})
	if err != nil {
		return err
	}

	// Remove the old image if specified
	// if removeOldImage {
	// 	err = cli.ImageRemove(context.Background(), oldImageName, types.ImageRemoveOptions{})
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// Start the container with the updated image
	err = cli.ContainerStart(context.Background(), containerID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	fmt.Println("DONE")
	return nil
}

type BuildSettings struct {
	Repo *database.Repository
}

func BuildAndUploadImage(buildSettings BuildSettings, ch chan<- string) error {
	repo := buildSettings.Repo
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Step 1: Start the build container
	containerID, err := startBuildContainer(ctx, cli, ch)
	if err != nil {
		return err
	}
	ch <- "start build container"

	// Step 2: Pull the GitHub repo inside the container
	err = pullRepo(ctx, cli, containerID, repo.Url, ch)
	if err != nil {
		return err
	}
	ch <- "pull repo"

	// Step 3: Build the Dockerfile in the root of the given repo
	err = buildDockerfile(ctx, cli, containerID, repo.ContainerRepo.String, repo.ContainerTag.String, ch)
	if err != nil {
		return err
	}
	ch <- "build docker file"

	err = cli.ContainerStop(ctx, containerID, container.StopOptions{})
	// err = cli.ContainerStop(ctx, containerID, container.StopOptions{}, ch)
	if err != nil {
		return err
	}
	ch <- "stop container"

	err = cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{})
	// err = cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{}, ch)
	if err != nil {
		return err
	}
	ch <- "remove container"
	return nil
}

type chanWriter struct {
	ch chan<- string
}

func newChanWriter(ch chan<- string) *chanWriter {
	return &chanWriter{ch: ch}
}

func (w *chanWriter) Write(p []byte) (int, error) {
	n := len(p)
	j, _ := json.Marshal(p)
	w.ch <- string(j)
	return n, nil
}

func startBuildContainer(ctx context.Context, cli *client.Client, ch chan<- string) (string, error) {
	const buildContainerImg = "moby/buildkit:master"
	{
		fmt.Println("pull image")
		out, err := cli.ImagePull(ctx, buildContainerImg, types.ImagePullOptions{})
		if err != nil {
			return "", err
		}
		defer out.Close()
		writer := newChanWriter(ch)
		jsonmessage.DisplayJSONMessagesStream(out, writer, 1, true, nil)
	}
	fmt.Println("creating container")
	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: buildContainerImg,
			Tty:   true,
		},
		&container.HostConfig{
			Privileged: true,
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: "/home/vagrant/.docker/config.json",
					Target: "/root/.docker/config.json",
				},
			},
		},
		// TODO: volumes stapelen op maar ik dacht dat ik er geen aanmaakt???
		nil,
		nil,
		"HH-build",
	)
	if err != nil {
		return "", err
	}

	fmt.Println("starting container")
	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func pullRepo(ctx context.Context, cli *client.Client, containerID, githubRepoURL string, ch chan<- string) error {
	// Use 'git clone' to pull the GitHub repo inside the container
	fmt.Println("pulling repo")
	cmd := []string{"git", "clone", githubRepoURL, "/app"}
	consoleSize := [2]uint{800, 600}
	execResp, err := cli.ContainerExecCreate(ctx, containerID, types.ExecConfig{
		Cmd:          cmd,
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
		ConsoleSize:  &consoleSize,
	})
	if err != nil {
		return err
	}

	startResp, err := cli.ContainerExecAttach(context.Background(), execResp.ID, types.ExecStartCheck{})
	if err != nil {
		panic(err)
	}
	defer startResp.Close()
	io.Copy(newChanWriter(ch), startResp.Reader)
	return nil
}

// func Test() {
// 	l, _ := cli.PluginList(context.Background(), filters.Args{})
// 	// r, _ := cli.PluginInstall(context.Background(), "compose", types.PluginInstallOptions{})
// 	// buf := new(strings.Builder)
// 	// io.Copy(buf, r)
// 	// // check errors
// 	// fmt.Println(buf.String())
// 	fmt.Printf("%v", l)
// }

func buildDockerfile(ctx context.Context, cli *client.Client, containerID string, repo string, tag string, ch chan<- string) error {
	fmt.Println("Building container...")
	repoUrl := fmt.Sprintf("docker.io/%s:%s", repo, tag)
	cmd := []string{
		"buildctl-daemonless.sh",
		"build",
		"--frontend", "dockerfile.v0",
		"--local", "context=/app",
		"--local", "dockerfile=/app",
		"--progress", "tty",
		"--output", fmt.Sprintf("type=image,name=%s,push=true", repoUrl),
	}
	consoleSize := [2]uint{100, 80}
	execResp, err := cli.ContainerExecCreate(ctx, containerID, types.ExecConfig{
		Cmd:          cmd,
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
		ConsoleSize:  &consoleSize,
	})
	if err != nil {
		return err
	}

	execID := execResp.ID
	startResp, err := cli.ContainerExecAttach(context.Background(), execID, types.ExecStartCheck{
		ConsoleSize: &consoleSize,
		// removing Tty: true makes the output slightly corrupted, DO NOT REMOVE
		Tty: true,
	})
	if err != nil {
		panic(err)
	}
	defer startResp.Close()
	io.Copy(newChanWriter(ch), startResp.Reader)
	fmt.Println("Build done!")

	return nil
}
