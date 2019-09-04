package jobs

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/freundallein/docker-repo-watcher/logger"
	"github.com/freundallein/docker-repo-watcher/settings"
)

type image struct {
	ID,
	tag string
}

func (i *image) String() string {
	return fmt.Sprintf("%s - %s", i.tag, i.ID)
}

// ImageSelection - clean stale images
func ImageSelection(cli *client.Client, settings *settings.Settings) {
	matchString := fmt.Sprintf("%s:%s/%s*", settings.RegistryIP, settings.RegistryPort, settings.AppPrefix)
	logger.Info("Run image selection")
	images, err := fetchImages(cli, matchString)
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
	}
	imagesByName := separateImages(images)
	err = cleanImages(cli, imagesByName, settings.ImageAmount)
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
	}
}

func fetchImages(cli *client.Client, matchString string) ([]types.ImageSummary, error) {
	args := filters.NewArgs(
		filters.Arg("reference", matchString),
	)
	imageOptions := types.ImageListOptions{Filters: args}
	logger.Debug("Fetch images from local repository")
	images, err := cli.ImageList(context.Background(), imageOptions)
	if err != nil {
		return nil, err
	}
	return images, nil
}

func separateImages(images []types.ImageSummary) map[string][]*image {
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

func cleanImages(cli *client.Client, imagesByName map[string][]*image, amountToStore int) error {
	logger.Debug("Start cleaning images")
	for name, images := range imagesByName {
		if len(images) > amountToStore {
			toDelete := chooseToDelete(images, amountToStore)
			logger.Debug(fmt.Sprintf("Deleting %d %s images", len(toDelete), name))
			logger.Debug("IDs for deletion:")
			for i := range toDelete {
				logger.Debug(fmt.Sprintf("%s", toDelete[i]))
			}

			deleted, err := deleteImages(cli, toDelete)
			if err != nil {
				return err
			}
			logger.Info(fmt.Sprintf("%d images deleted", deleted))
		} else {
			logger.Info(fmt.Sprintf("%s images are OK", name))
		}
	}
	return nil
}

func chooseToDelete(images []*image, amountToStore int) []string {
	toDelete := []string{}
	re := regexp.MustCompile(`20..-..-..*`)
	sort.SliceStable(images, func(i, j int) bool {
		return images[i].tag > images[j].tag
	})
	storedImages := 1
	for _, img := range images {
		if (storedImages > amountToStore) || (!re.MatchString(img.tag) && img.tag != "latest") {
			toDelete = append(toDelete, img.ID)
		} else {
			storedImages++
		}
	}
	return toDelete
}

func deleteImages(cli *client.Client, toDelete []string) (int, error) {
	deleted := 0
	for _, imageID := range toDelete {
		_, err := cli.ImageRemove(context.Background(), imageID, types.ImageRemoveOptions{Force: true})
		if err != nil {
			return deleted, err
		}
		deleted++
	}
	return deleted, nil
}
