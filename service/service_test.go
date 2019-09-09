package service

import (
	"testing"
)

func TestJobString(t *testing.T) {
	job := Job{Name: "TEST", Crontab: "* * * * *"}
	if job.String() != "TEST - * * * * *" {
		t.Error("Expected 'TEST - * * * * * got ", job.String())
	}
}

func TestNewService(t *testing.T) {
	expectedPeriod := 100
	expectedJobs := []*Job{&Job{Name: "TEST_1"}, &Job{Name: "TEST_2", Crontab: "* * * * *"}}
	s := NewService(expectedPeriod, expectedJobs)
	if s.period != expectedPeriod {
		t.Error("Expected", expectedPeriod, "got ", s.period)
	}

	if s.customJobs[0] != expectedJobs[0] {
		t.Error("Expected", expectedJobs[0], "got ", s.customJobs[0])
	}
	if s.crontabJobs[0] != expectedJobs[1] {
		t.Error("Expected", expectedJobs[1], "got ", s.crontabJobs[0])
	}
}
