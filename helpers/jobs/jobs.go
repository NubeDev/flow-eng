package jobs

import (
	"github.com/go-co-op/gocron"
	"time"
)

type Jobs struct {
	cron *gocron.Scheduler
}

func New() *Jobs {
	c := gocron.NewScheduler(time.UTC)
	c.StartAsync()
	return &Jobs{c}
}

func (inst *Jobs) Get() *gocron.Scheduler {
	return inst.cron
}
