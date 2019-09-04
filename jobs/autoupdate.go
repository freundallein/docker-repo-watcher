package jobs

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/freundallein/docker-repo-watcher/logger"
	"github.com/freundallein/docker-repo-watcher/settings"
)

type contextKey string

// AutoUpdate - pull latest docker image and restart container
func AutoUpdate(cli *client.Client) {
	reference := fmt.Sprintf("%s:latest", settings.AppImageName)
	logger.Info("Auto update check")
	ctx := context.WithValue(context.Background(), contextKey("cli"), cli)
	ctx = context.WithValue(ctx, contextKey("reference"), reference)
	pullOptions := types.ImagePullOptions{}
	oldImageID, err := checkImageID(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
		return
	}
	_, err = cli.ImagePull(ctx, reference, pullOptions)
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
		return
	}
	newImageID, err := checkImageID(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
		return
	}
	fmt.Println(oldImageID, newImageID)
	if oldImageID == newImageID {
		return
	}
	runningContainers, err := getContainerIDs(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
	}
	for _, containerID := range runningContainers {
		containerContext := context.WithValue(ctx, contextKey("contID"), containerID)
		err = modifyRestartPolicy(containerContext)
		if err != nil {
			logger.Error(fmt.Sprintf("%s", err))
		}
		err = restart(containerContext)
		if err != nil {
			logger.Error(fmt.Sprintf("%s", err))
		}
	}
}

func checkImageID(ctx context.Context) (string, error) {
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

func getContainerIDs(ctx context.Context) ([]string, error) {
	cli := ctx.Value(contextKey("cli")).(*client.Client)
	filters := filters.NewArgs(
		filters.Arg("ancestor", settings.AppImageName),
	)
	containerOptions := types.ContainerListOptions{Filters: filters}
	runningContainers, err := cli.ContainerList(context.Background(), containerOptions)
	if err != nil {
		return nil, err
	}
	containerIDs := []string{}
	for _, container := range runningContainers {
		fmt.Println(container.ImageID, container.Image)
		if container.Image == "freundallein/drwatcher" {
			containerIDs = append(containerIDs, container.ID)
		}
	}
	fmt.Println(containerIDs)
	return containerIDs, nil
}

func modifyRestartPolicy(ctx context.Context) error {
	cli := ctx.Value(contextKey("cli")).(*client.Client)
	containerID := ctx.Value(contextKey("contID")).(string)
	containerJSON, err := cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return err
	}
	restartPolicy := containerJSON.HostConfig.RestartPolicy.Name
	if restartPolicy != "always" {
		updateConfig := container.UpdateConfig{RestartPolicy: container.RestartPolicy{Name: "always"}}
		_, err := cli.ContainerUpdate(ctx, containerID, updateConfig)
		if err != nil {
			return err
		}
	}
	return nil
}

func restart(ctx context.Context) error {
	cli := ctx.Value(contextKey("cli")).(*client.Client)
	containerID := ctx.Value(contextKey("contID")).(string)
	err := cli.ContainerRestart(ctx, containerID, nil)
	return err
}
