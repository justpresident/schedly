package schedly

import (
	"time"
)

/*Scheduler is a simple task scheduling tool.
 */
type Scheduler interface {
	Tick() time.Duration
	Aligned() bool
	SetAligned(aligned bool) Scheduler
	NewSchedule(every time.Duration) *schedule
	Schedule(every time.Duration, name string, jobFunc func()) *job
	AddJob(schedule *schedule, name string, jobFunc func()) *job
	Start() chan int
	Stop()
	WaitAndStop()
}

/*NewScheduler creates a new Scheduler instance with provided precision.
Higher precision means more frequent evaluation of Schedules.
To reduce overhead please avoid using tick time smaller than 1 second unless you really need it.
*/
func NewScheduler(tick time.Duration) Scheduler {
	return &scheduler{
		tick:     tick,
		jobs:     make(map[*schedule][]*job),
		finished: make(chan int),
	}
}

type scheduler struct {
	jobs     map[*schedule][]*job
	finished chan int
	tick     time.Duration
	aligned  bool
	ticker   *time.Ticker
}

func (s *scheduler) Tick() time.Duration {
	return s.tick
}

func (s *scheduler) Aligned() bool {
	return s.aligned
}

func (s *scheduler) SetAligned(aligned bool) Scheduler {
	s.aligned = aligned
	return s
}

func (s *scheduler) NewSchedule(every time.Duration) *schedule {
	schedule := newSchedule(s.tick).SetEvery(every)
	return schedule
}

func (s *scheduler) Schedule(every time.Duration, name string, jobFunc func()) *job {
	schedule := s.NewSchedule(every)
	return s.AddJob(schedule, name, jobFunc)
}

func (s *scheduler) AddJob(schedule *schedule, name string, jobFunc func()) *job {
	j := &job{
		jobFunc: jobFunc,
		name:    name,
	}
	if jList, ok := s.jobs[schedule]; ok {
		s.jobs[schedule] = append(jList, j)
	} else {
		s.jobs[schedule] = []*job{j}
	}

	return j
}

func (s *scheduler) Start() chan int {

	if s.aligned {
		curTime := time.Now()
		time.Sleep(s.tick - curTime.Sub(curTime.Truncate(s.tick)))
	}

	s.ticker = time.NewTicker(s.tick)

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

func (s *scheduler) Stop() {

	s.finished <- 1
}

func (s *scheduler) WaitAndStop() {
	panic("not implemented")

	// TODO: use waitgroup to wait for running tasks
	//s.finished <- 1
}

func (s *scheduler) runPending(tick time.Time) {
	for schedule, jobs := range s.jobs {
		for _, job := range jobs {
			lastTime := job.LastRun()
			if schedule.IntervalMode() {
				lastTime = job.LastFinish()
			}
			if schedule.CanRun(tick, lastTime) {
				job.Run(tick)
			}
		}
	}
}
