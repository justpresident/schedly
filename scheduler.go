package schedly

import (
	"sync"
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
	Start()
	Stop()
	WaitForRunningTasks()
	WaitUntilStopped()
}

/*NewScheduler creates a new Scheduler instance with provided precision.
Higher precision means more frequent evaluation of Schedules.
To reduce overhead please avoid using tick time smaller than 1 second unless you really need it.
*/
func NewScheduler(tick time.Duration) Scheduler {
	return &scheduler{
		tick:       tick,
		jobs:       make(map[*schedule][]*job),
		toFinish:   make(chan int),
		globalStop: make(chan bool),
		wg:         new(sync.WaitGroup),
	}
}

type scheduler struct {
	jobs       map[*schedule][]*job
	toFinish   chan int
	globalStop chan bool
	tick       time.Duration
	aligned    bool
	ticker     *time.Ticker
	wg         *sync.WaitGroup
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
	schedule := newSchedule(s.tick).SetEvery(every).SetAligned(s.aligned)
	return schedule
}

func (s *scheduler) Schedule(every time.Duration, name string, jobFunc func()) *job {
	schedule := s.NewSchedule(every)
	return s.AddJob(schedule, name, jobFunc)
}

func (s *scheduler) AddJob(schedule *schedule, name string, jobFunc func()) *job {
	// TODO: Check that job name is unique
	// TODO: Create an interface for retrieving jobs state
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

func (s *scheduler) Start() {

	if s.aligned {
		curTime := time.Now()
		time.Sleep(s.tick - curTime.Sub(curTime.Truncate(s.tick)))
	}

	s.ticker = time.NewTicker(s.tick)

	go func() {
		for {
			select {
			case tick := <-s.ticker.C:
				s.runPending(tick)
			case <-s.toFinish:
				s.ticker.Stop()
				s.globalStop <- true
				return
			}
		}
	}()
}

func (s *scheduler) Stop() {
	s.toFinish <- 1
}

func (s *scheduler) WaitUntilStopped() {
	<-s.globalStop
}

func (s *scheduler) WaitForRunningTasks() {
	s.wg.Wait()
}

func (s *scheduler) runPending(tick time.Time) {
	for schedule, jobs := range s.jobs {
		for _, j := range jobs {
			if schedule.CanRun(tick, j) {
				s.runJob(tick, j)
			}
		}
	}
}

func (s *scheduler) runJob(tick time.Time, j *job) {
	s.wg.Add(1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				// TODO: provide an interface for task state change listener
			} else {
				// TODO: provide an interface for task state change listener
			}
			s.wg.Done()
		}()
		j.Run(tick)
	}()
}
