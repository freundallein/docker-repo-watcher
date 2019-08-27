package service

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/docker/docker/api/types/filters"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/freundallein/docker-repo-watcher/logger"
	"github.com/robfig/cron"
)

type image struct {
	ID,
	tag string
}

// Service - watch local docker repository
type Service struct {
	cli *client.Client
	matchString,
	crontab string
	amountToStore,
	period int
}

// NewService - spawn Service struct
func NewService(cli *client.Client, matchString, crontab string, period, amountToStore int) *Service {
	return &Service{
		cli:           cli,
		matchString:   matchString,
		crontab:       crontab,
		period:        period,
		amountToStore: amountToStore,
	}
}

// Run - execute service
func (serv *Service) Run() {
	logger.Debug("Start service")
	period := serv.period
	logger.Debug(fmt.Sprintf("Period - %s", time.Duration(period)*time.Second))
	c := cron.New()
	c.AddFunc(serv.crontab, serv.prune)
	c.Start()
	logger.Debug("Start prune crontab job")
	defer c.Stop()
	for {
		select {
		case <-time.After(time.Duration(period) * time.Second):
			logger.Debug("Run iteration")
			images, err := serv.fetchImages()
			if err != nil {
				logger.Error(fmt.Sprintf("%s", err))
			}
			imagesByName := serv.separateImages(images)
			err = serv.cleanImages(imagesByName)
			if err != nil {
				logger.Error(fmt.Sprintf("%s", err))
			}
		}
	}
}

func (serv *Service) prune() {
	logger.Debug("Start prune")
	serv.cli.NetworksPrune(context.Background(), filters.NewArgs())
	serv.cli.VolumesPrune(context.Background(), filters.NewArgs())
	serv.cli.ImagesPrune(context.Background(), filters.NewArgs())
	serv.cli.ContainersPrune(context.Background(), filters.NewArgs())
	logger.Debug("End prune")
}

func (serv *Service) fetchImages() ([]types.ImageSummary, error) {
	args := filters.NewArgs(
		filters.Arg("reference", serv.matchString),
	)
	imageOptions := types.ImageListOptions{Filters: args}
	logger.Debug("Fetch images from local repository")
	images, err := serv.cli.ImageList(context.Background(), imageOptions)
	if err != nil {
		return nil, err
	}
	return images, nil
}

func (serv *Service) separateImages(images []types.ImageSummary) map[string][]*image {
	logger.Debug("Separate images by name")
	imagesByName := map[string][]*image{}
	for _, img := range images {
		for i := range img.RepoTags {
			imageNameWithTag := strings.Split(img.RepoTags[i], "/")[1]
			splitedImage := strings.Split(imageNameWithTag, ":")
			imageName := splitedImage[0]
			tag := splitedImage[1]
			if _, ok := imagesByName[imageName]; ok {
				imagesByName[imageName] = append(imagesByName[imageName], &image{ID: img.ID, tag: tag})
			} else {
				imagesByName[imageName] = []*image{&image{ID: img.ID, tag: tag}}
			}
		}

	}
	return imagesByName
}
func (serv *Service) chooseToDelete(images []*image) []string {
	toDelete := []string{}
	re := regexp.MustCompile(`20..-..-..*`)
	sort.SliceStable(images, func(i, j int) bool {
		return images[i].tag > images[j].tag
	})
	storedImages := 0
	for _, img := range images {
		if (storedImages > serv.amountToStore) || (!re.MatchString(img.tag) && img.tag != "latest") {
			toDelete = append(toDelete, img.ID)
		}
		storedImages++
	}
	return toDelete
}

func (serv *Service) deleteImages(toDelete []string) (int, error) {
	deleted := 0
	for _, imageID := range toDelete {
		_, err := serv.cli.ImageRemove(context.Background(), imageID, types.ImageRemoveOptions{Force: true})
		deleted++
		return deleted, err
	}
	return deleted, nil
}

func (serv *Service) cleanImages(imagesByName map[string][]*image) error {
	logger.Debug("Start cleaning images")
	for name, images := range imagesByName {
		if len(images) > serv.amountToStore {
			toDelete := serv.chooseToDelete(images)
			logger.Debug("IDs for deletion:")
			for i := range toDelete {
				logger.Debug(fmt.Sprintf("%s", toDelete[i]))
			}
			logger.Debug(fmt.Sprintf("%s images to delete", name))
			logger.Debug(fmt.Sprintf("%d images for deletion", len(toDelete)))
			deleted, err := serv.deleteImages(toDelete)
			if err != nil {
				return err
			}
			logger.Debug(fmt.Sprintf("%d images deleted", deleted))
		} else {
			logger.Debug(fmt.Sprintf("%s images are OK", name))
		}
	}
	return nil
}
