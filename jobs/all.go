package jobs

import (
	"github.com/docker/docker/client"
	"github.com/freundallein/docker-repo-watcher/service"
	"github.com/freundallein/docker-repo-watcher/settings"
)

// InitJobs - make slice of all available jobs
func InitJobs(cli *client.Client, settings *settings.Settings) []*service.Job {
	crontab := settings.Crontab
	jobs := []*service.Job{
		&service.Job{
			Name:    "Networks prune",
			Crontab: crontab,
			Routine: func() { networksPrune(cli) },
		},
		&service.Job{
			Name:    "Volumes prune",
			Crontab: crontab,
			Routine: func() { volumesPrune(cli) },
		},
		&service.Job{
			Name:    "Images prune",
			Crontab: crontab,
			Routine: func() { imagesPrune(cli) },
		},
		&service.Job{
			Name:    "Containers prune",
			Crontab: crontab,
			Routine: func() { containersPrune(cli) },
		},
		&service.Job{
			Name:    "Image selection",
			Routine: func() { ImageSelection(cli, settings) },
		},
	}
	if settings.AutoUpdate {
		jobs = append(jobs, &service.Job{
			Name:    "Auto update",
			Routine: func() { AutoUpdate(cli, settings) },
		})
	}
	return jobs
}
