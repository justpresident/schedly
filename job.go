package schedly

import (
	"sync"
	"time"
)

type Job struct {
	jobFunc      func()
	name         string
	schedule     Schedule
	lastRun      time.Time
	lastSuccess  time.Time
	exclusive    bool
	intervalMode bool
	lock         sync.Mutex
}

/*
	When true job schedule receives previous successful finish time as a param in CanRun method.

	When false job schedule receives previous start time as a param in CanRun method
*/
func (j *Job) IntervalMode() bool {
	return j.intervalMode
}

func (j *Job) SetIntervalMode(intervalMode bool) {
	j.intervalMode = intervalMode
}

/*
	Exclusive mode prevents from running multiple instances of the job at the same time
*/
func (j *Job) Exclusive() bool {
	return j.exclusive
}

func (j *Job) SetExclusive(exclusive bool) {
	j.exclusive = exclusive
}

func (j *Job) Schedule() Schedule {
	return j.schedule
}

func (j *Job) SetSchedule(schedule Schedule) {
	j.schedule = schedule
}

func (j *Job) shouldRun(tick time.Time) bool {
	lastTime := j.lastRun
	if j.intervalMode {
		lastTime = j.lastSuccess
	}

	shouldRun := j.schedule.CanRun(tick, lastTime)

	return shouldRun
}

/*
	Run the job.

	Parameter `tick` is time from internal ticker
*/
func (j *Job) run(tick time.Time) {
	j.lastRun = tick
	go func() {
		if j.exclusive {
			j.lock.Lock()
			defer j.lock.Unlock()
		}
		j.jobFunc()
		j.lastSuccess = time.Now()
	}()
}
