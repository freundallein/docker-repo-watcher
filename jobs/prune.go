package jobs

import (
	"context"

	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func networksPrune(cli *client.Client) {
	cli.NetworksPrune(context.Background(), filters.NewArgs())
}
func volumesPrune(cli *client.Client) {
	cli.VolumesPrune(context.Background(), filters.NewArgs())
}
func imagesPrune(cli *client.Client) {
	cli.ImagesPrune(context.Background(), filters.NewArgs())
}
func containersPrune(cli *client.Client) {
	cli.ContainersPrune(context.Background(), filters.NewArgs())
}
