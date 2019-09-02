package main

import (
	"fmt"
	"os"

	"github.com/docker/docker/client"
	"github.com/freundallein/docker-repo-watcher/jobs"
	"github.com/freundallein/docker-repo-watcher/logger"
	"github.com/freundallein/docker-repo-watcher/service"
	"github.com/freundallein/docker-repo-watcher/settings"
)

var dockerAPIVersion = os.Getenv("DOCKER_API_VERSION")

func main() {
	settings := settings.NewSettings()
	logger.InitLogger("D-REPO-WATCHER", settings.LogLevel)
	logger.Debug("Create docker client")
	cli, err := client.NewClientWithOpts(client.WithVersion(dockerAPIVersion))
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
		panic(err)
	}
	logger.Debug("Create service")
	jobs := jobs.InitJobs(cli, settings)
	service := service.NewService(
		settings.Period,
		jobs,
	)
	logger.Info("Start service")
	service.Run()
}
