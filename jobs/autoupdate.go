package jobs

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/freundallein/docker-repo-watcher/logger"
	"github.com/freundallein/docker-repo-watcher/settings"
)

type contextKey string

// AutoUpdate - pull latest docker image and restart container
func AutoUpdate(cli *client.Client, config *settings.Settings) {
	reference := fmt.Sprintf("%s:latest", settings.AppImageName)
	logger.Info("Auto update check")
	ctx := context.WithValue(context.Background(), contextKey("cli"), cli)
	ctx = context.WithValue(ctx, contextKey("reference"), reference)
	ctx = context.WithValue(ctx, contextKey("settings"), config)
	containerID, err := getContainerID(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
	}
	ctx = context.WithValue(ctx, contextKey("contID"), containerID)
	pullOptions := types.ImagePullOptions{}
	oldImageID, err := getCurrentImageID(ctx)
	logger.Debug(fmt.Sprintf("Current image %s", oldImageID))
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
		return
	}
	_, err = cli.ImagePull(ctx, reference, pullOptions)
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
		return
	}
	newImageID, err := checkRepoImageID(ctx)
	logger.Debug(fmt.Sprintf("Repository latest image %s", newImageID))
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
		return
	}
	logger.Debug(fmt.Sprintf("Image diff %s %s", oldImageID, newImageID))
	if oldImageID == newImageID {
		logger.Info("Latest image used")
		return
	}
	logger.Info("Re-deploy application.")
	err = redeploy(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
	}
}

func getCurrentImageID(ctx context.Context) (string, error) {
	cli := ctx.Value(contextKey("cli")).(*client.Client)
	containerID := ctx.Value(contextKey("contID")).(string)
	containerJSON, err := cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return "", err
	}
	return containerJSON.Image, nil
}

func checkRepoImageID(ctx context.Context) (string, error) {
	cli := ctx.Value(contextKey("cli")).(*client.Client)
	reference := ctx.Value(contextKey("reference")).(string)
	filters := filters.NewArgs(
		filters.Arg("reference", reference),
	)
	imageOptions := types.ImageListOptions{Filters: filters}
	list, err := cli.ImageList(ctx, imageOptions)
	if err != nil {
		return "", err
	}
	if len(list) > 0 {
		return list[0].ID, nil
	}
	return "", nil
}

func getContainerID(ctx context.Context) (string, error) {
	cmd := "cat /proc/self/cgroup | head -1"
	output, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return "", nil
	}
	containerID := strings.Split(string(output), "/")[2][:12]
	logger.Debug(fmt.Sprintf("Container ID %s", containerID))
	return containerID, nil
}

func redeploy(ctx context.Context) error {
	cli := ctx.Value(contextKey("cli")).(*client.Client)
	containerID := ctx.Value(contextKey("contID")).(string)
	image := ctx.Value(contextKey("reference")).(string)
	settings := ctx.Value(contextKey("settings")).(*settings.Settings)
	newContainer, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: image,
			Env:   settings.ToEnvString(),
		},
		&container.HostConfig{
			Binds: []string{"/var/run/docker.sock:/var/run/docker.sock"},
		}, nil, "")
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
	}
	cli.ContainerStart(ctx, newContainer.ID, types.ContainerStartOptions{})
	err = cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true})
	logger.Debug(fmt.Sprintf("Container %s restart", containerID))
	return err
}
