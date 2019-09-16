package jobs

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/freundallein/docker-repo-watcher/logger"
	"github.com/freundallein/docker-repo-watcher/settings"
)

type repository struct {
	Name,
	Path string
}

func cleanRegistry(cli *client.Client, config *settings.Settings) {
	logger.Info("Registry cleaning")
	repositories, err := discoverRepositories(config.RegistryPath)
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
	}
	var wg sync.WaitGroup
	repositroriesAmount := 0
	for _, repo := range repositories {
		if !strings.HasPrefix(repo.Name, config.AppPrefix) {
			continue
		}
		repositroriesAmount++
	}
	wg.Add(repositroriesAmount)
	for _, repo := range repositories {
		if !strings.HasPrefix(repo.Name, config.AppPrefix) {
			continue
		}
		go cleanRepository(&wg, repo, config.ImageAmount)
	}
	wg.Wait()
	garbageCollect(cli)
}

func discoverRepositories(registryPath string) ([]*repository, error) {
	repoPath := fmt.Sprintf("%s/docker/registry/v2/repositories/", registryPath)
	logger.Debug(fmt.Sprintf("Registry path %s", repoPath))
	files, err := ioutil.ReadDir(repoPath)
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
	}
	repos := []*repository{}
	for _, f := range files {
		if f.IsDir() {
			repos = append(repos, &repository{Name: f.Name(), Path: fmt.Sprintf("%s%s", repoPath, f.Name())})
		}
	}
	return repos, nil
}

func cleanRepository(wg *sync.WaitGroup, repo *repository, imageAMount int) {
	defer wg.Done()
	tagsPath := fmt.Sprintf("%s/_manifests/tags", repo.Path)
	tagDirs, err := ioutil.ReadDir(tagsPath)
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
	}
	tags := []string{}
	revisions := map[string][]string{}
	for _, tag := range tagDirs {
		tags = append(tags, tag.Name())
		revisionns, _ := ioutil.ReadDir(fmt.Sprintf("%s/%s/index/sha256/", tagsPath, tag.Name()))
		for _, rev := range revisionns {
			revisions[tag.Name()] = append(revisions[tag.Name()], rev.Name())
		}

	}
	logger.Debug(fmt.Sprintf("%s: %s", repo.Name, tags))
	tagsToDelete := chooseTagsToDelete(tags, imageAMount)
	revisionsToStay := chooseRevisionsToStay(tags, tagsToDelete, revisions)
	deleteTags(repo, tagsToDelete)
	deleteRevisions(repo, revisionsToStay)
}

func chooseTagsToDelete(tags []string, amountToStore int) map[string]struct{} {
	toDelete := map[string]struct{}{}
	re := regexp.MustCompile(`20..-..-..*`)
	sort.SliceStable(tags, func(i, j int) bool {
		return tags[i] > tags[j]
	})
	storedImages := 1
	for _, tag := range tags {
		if (storedImages > amountToStore) || (!re.MatchString(tag) && tag != "latest") {
			toDelete[tag] = struct{}{}
		} else {
			storedImages++
		}
	}
	return toDelete
}

func chooseRevisionsToStay(tags []string, tagsToDelete map[string]struct{}, revisions map[string][]string) map[string]struct{} {
	revisionsToStay := map[string]struct{}{}
	for _, tag := range tags {
		if _, hasKey := tagsToDelete[tag]; hasKey {
			continue
		}
		for _, rev := range revisions[tag] {
			revisionsToStay[rev] = struct{}{}
		}
	}
	return revisionsToStay
}

func deleteTags(repo *repository, tags map[string]struct{}) {
	for tag := range tags {
		removePath := fmt.Sprintf("%s/_manifests/tags/%s", repo.Path, tag)
		logger.Debug(fmt.Sprintf("DELETE tag %s", removePath))
		err := os.RemoveAll(removePath)
		if err != nil {
			logger.Error(fmt.Sprintf("%s", err))
		}
	}

}

func deleteRevisions(repo *repository, revisionsToStay map[string]struct{}) {
	revisionsPath := fmt.Sprintf("%s/_manifests/revisions/sha256", repo.Path)
	revisionDir, err := ioutil.ReadDir(revisionsPath)
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
	}
	for _, revision := range revisionDir {
		if _, shouldStay := revisionsToStay[revision.Name()]; shouldStay {
			continue
		}
		removePath := fmt.Sprintf("%s/_manifests/revisions/sha256/%s", repo.Path, revision.Name())
		logger.Debug(fmt.Sprintf("DELETE REVISION %s", removePath))
		err = os.RemoveAll(removePath)
		if err != nil {
			logger.Error(fmt.Sprintf("%s", err))
		}
	}
}

func garbageCollect(cli *client.Client) {
	ctx := context.Background()
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		logger.Error(fmt.Sprintf("%s", err))
	}
	for _, container := range containers {
		if container.Names[0][1:] == "registry" {
			logger.Debug(fmt.Sprintf("Registry container %s\n", container.ID))
			gcConfig := types.ExecConfig{
				AttachStdout: true,
				AttachStderr: true,
				Cmd:          []string{"/bin/registry", "garbage-collect", "/etc/docker/registry/config.yml"},
			}
			execID, err := cli.ContainerExecCreate(ctx, container.ID, gcConfig)
			if err != nil {
				logger.Error(fmt.Sprintf("%s", err))
			}
			r, err := cli.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{})
			if err != nil {
				logger.Error(fmt.Sprintf("%s", err))
			}
			content, _, _ := r.Reader.ReadLine()
			logger.Info(string(content))
		}
	}
}
