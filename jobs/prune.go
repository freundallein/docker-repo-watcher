package jobs

import (
	"context"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/freundallein/docker-repo-watcher/logger"
)

func networksPrune(cli *client.Client) {
	cli.NetworksPrune(context.Background(), filters.NewArgs())
	logger.Info("Networks pruned")
}

func volumesPrune(cli *client.Client) {
	cli.VolumesPrune(context.Background(), filters.NewArgs())
	logger.Info("Volumes pruned")
}

func imagesPrune(cli *client.Client) {
	cli.ImagesPrune(context.Background(), filters.NewArgs())
	logger.Info("Images pruned")
}

func containersPrune(cli *client.Client) {
	cli.ContainersPrune(context.Background(), filters.NewArgs())
	logger.Info("Containers pruned")
}
