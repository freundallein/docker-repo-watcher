package service

import (
	"fmt"
	"time"

	"github.com/freundallein/docker-repo-watcher/logger"
	"github.com/robfig/cron"
)

// Job -...
type Job struct {
	Name,
	Crontab string
	Routine func()
}

func (j *Job) String() string {
	return fmt.Sprintf("%s - %s", j.Name, j.Crontab)
}

// Service - watch local docker repository
type Service struct {
	period int
	customJobs,
	crontabJobs []*Job
}

// NewService - spawn Service struct
func NewService(period int, jobs []*Job) *Service {
	crontabJobs := []*Job{}
	customJobs := []*Job{}
	if period < 1 {
		period = 1
	}
	for _, job := range jobs {
		if job.Crontab == "" {
			customJobs = append(customJobs, job)
		} else {
			crontabJobs = append(crontabJobs, job)
		}
	}
	return &Service{
		period:      period,
		customJobs:  customJobs,
		crontabJobs: crontabJobs,
	}
}

// Run - execute service
func (serv *Service) Run() {
	logger.Debug("Start service")
	period := serv.period
	logger.Debug(fmt.Sprintf("Period - %s", time.Duration(period)*time.Second))

	c := cron.New()

	for _, job := range serv.crontabJobs {
		err := c.AddFunc(job.Crontab, job.Routine)
		logger.Debug(fmt.Sprintf("Scheduled %s crontab job on %s", job.Name, job.Crontab))
		if err != nil {
			logger.Error(fmt.Sprintf("%s", err))
		}
	}
	c.Start()

	defer c.Stop()
	for {
		select {
		case <-time.After(time.Duration(period) * time.Second):
			for _, j := range serv.customJobs {
				go j.Routine()
			}
		}
	}
}
