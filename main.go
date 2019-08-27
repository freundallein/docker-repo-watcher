package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/docker/docker/client"
	"github.com/freundallein/docker-repo-watcher/logger"
	"github.com/freundallein/docker-repo-watcher/service"
)

var dockerAPIVersion = os.Getenv("DOCKER_API_VERSION")

func main() {
	registryAddress := flag.String("raddr", "192.168.20.126", "Docker registry name")
	registryPort := flag.Int("rport", 5000, "Docker registry port")
	appPrefix := flag.String("appPrefix", "", "Image name prefix")
	crontab := flag.String("crontab", "0 0 * * *", "Crontab for docker prune")
	period := flag.Int("period", 5, "Request period in seconds")
	logLevel := flag.String("logLevel", "ERROR", "LogLevel - DEBUG/ERROR")
	imageAmount := flag.Int("amount", 5, "Amount of single app images to store")
	flag.Parse()

	logger.InitLogger("D-REPO-WATCHER", *logLevel)
	logger.Debug("Create docker client")
	cli, err := client.NewClientWithOpts(client.WithVersion(dockerAPIVersion))
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
		panic(err)
	}
	logger.Debug("Create service")
	service := service.NewService(
		cli,
		fmt.Sprintf("%s:%d/%s*", *registryAddress, *registryPort, *appPrefix),
		*crontab,
		*period,
		*imageAmount,
	)
	service.Run()
}
