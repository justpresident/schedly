package schedly

import (
	"time"
)

type Scheduler struct {
	jobs     []*Job
	finished chan int
	tick     time.Duration
	aligned  bool
	ticker   *time.Ticker
}

func (s *Scheduler) Tick() time.Duration {
	return s.tick
}

func (s *Scheduler) Aligned() bool {
	return s.aligned
}

func (s *Scheduler) SetAligned(aligned bool) *Scheduler {
	s.aligned = aligned
	return s
}

func NewScheduler(tick time.Duration) *Scheduler {
	return &Scheduler{
		tick:     tick,
		jobs:     []*Job{},
		finished: make(chan int),
	}
}

func (s *Scheduler) AddJob(name string, jobFunc func()) *Job {
	job := &Job{
		jobFunc:  jobFunc,
		name:     name,
		schedule: NewConstrainedSchedule(s.tick).SetAligned(s.aligned),
	}
	s.jobs = append(s.jobs, job)

	return job
}

func (s *Scheduler) Start() chan int {

	if s.aligned {
		curTime := time.Now()
		time.Sleep(s.tick - curTime.Sub(curTime.Truncate(s.tick)))
	}

	s.ticker = time.NewTicker(1 * time.Second)

	stopped := make(chan int)
	go func() {
		for {
			select {
			case tick := <-s.ticker.C:
				s.runPending(tick)
			case <-s.finished:
				s.ticker.Stop()
				stopped <- 1
				return
			}
		}
	}()

	return stopped
}

func (s *Scheduler) Stop() {

	s.finished <- 1
}

func (s *Scheduler) WaitAndStop() {
	panic("not implemented")

	// TODO: use waitgroup to wait for running tasks
	s.finished <- 1
}

func (s *Scheduler) runPending(tick time.Time) {
	for _, job := range s.jobs {
		if job.shouldRun(tick) {
			job.run(tick)
		}
	}
}
