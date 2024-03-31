package docker

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/Pineapple217/BoxBoss/pkg/database"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/volume"
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
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		panic(err)
	}

	// for _, container := range containers {
	// 	fmt.Sprintf("%s %s\n", container.ID[:10], container.Image)
	// }
	return containers
}

type BuildSettings struct {
	Repo *database.Repository
}

func BuildAndUploadImage(buildSettings BuildSettings, ch chan<- string) error {
	repo := buildSettings.Repo
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	containerID, err := startBuildContainer(ctx, ch)
	if err != nil {
		return err
	}
	ch <- "start build container\\r\\n"

	err = pullRepo(ctx, containerID, repo.Url, ch)
	if err != nil {
		return err
	}
	ch <- "pull repo\\r\\n"

	err = buildDockerfile(ctx, containerID, repo.ContainerRepo.String, repo.ContainerTag.String, ch)
	if err != nil {
		return err
	}
	ch <- "build docker file\\r\\n"

	err = cli.ContainerStop(ctx, containerID, container.StopOptions{})
	if err != nil {
		return err
	}
	ch <- "stop container\\r\\n"

	err = cli.ContainerRemove(ctx, containerID, container.RemoveOptions{})
	if err != nil {
		return err
	}
	ch <- "remove container\\r\\n"
	return nil
}

func startBuildContainer(ctx context.Context, ch chan<- string) (string, error) {
	const buildContainerImg = "moby/buildkit:latest"
	{
		out, err := cli.ImagePull(ctx, buildContainerImg, types.ImagePullOptions{})
		if err != nil {
			return "", err
		}
		defer out.Close()
		writer := NewFixLinebreakMiddleware(NewChanWriter(ch))
		jsonmessage.DisplayJSONMessagesStream(out, writer, 1, true, nil)
	}
	_, err := cli.VolumeCreate(context.Background(), volume.CreateOptions{
		Name: "bb-buildcache",
	})
	if err != nil {
		fmt.Println(err)
	}

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
				{
					Type:   mount.TypeVolume,
					Source: "bb-buildcache",
					Target: "/var/lib/buildkit",
				},
			},
		},
		// TODO: volumes stapelen op maar ik dacht dat ik er geen aanmaakt???
		nil,
		nil,
		"bb-build",
	)
	if err != nil {
		return "", err
	}

	fmt.Println("starting container")
	err = cli.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func pullRepo(ctx context.Context, containerID, githubRepoURL string, ch chan<- string) error {
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
	io.Copy(NewChanWriter(ch), startResp.Reader)
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

func buildDockerfile(ctx context.Context, containerID string, repo string, tag string, ch chan<- string) error {
	fmt.Println("Building container...")
	repoUrl := fmt.Sprintf("%s:%s", repo, tag)
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
	io.Copy(NewChanWriter(ch), startResp.Reader)
	fmt.Println("Build done!")

	return nil
}
