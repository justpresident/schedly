package schedly

import (
	"sync"
	"time"
)

/*Job represents a configuration of a task scheduled for execution
 */
type Job interface {
	/*IntervalMode returns intervalMode of the job. When true job schedule receives previous successful finish time as a param in CanRun method.

	When false job schedule receives previous start time as a param in CanRun method*/
	IntervalMode() bool

	/*SetIntervalMode sets intervalMode of the job. When 'true' job schedule receives previous successful finish time as a param in CanRun method.

	When 'false' job schedule receives previous start time as a param in CanRun method*/
	SetIntervalMode(intervalMode bool)

	/*Exclusive returns Exclusive setting of the job. When 'true' it prevents from running multiple instances of the job at the same time*/
	Exclusive() bool
	/*SetExclusive configures Exclusive mode. When 'true' it prevents from running multiple instances of the job at the same time*/
	SetExclusive(exclusive bool)

	/*Schedule for running the job.*/
	Schedule() Schedule
	/*SetSchedule sets schedule for running the job*/
	SetSchedule(schedule Schedule)
	shouldRun(tick time.Time) bool
	run(tick time.Time)
}

type job struct {
	jobFunc      func()
	name         string
	schedule     Schedule
	lastRun      time.Time
	lastSuccess  time.Time
	exclusive    bool
	intervalMode bool
	lock         sync.Mutex
}

func (j *job) IntervalMode() bool {
	return j.intervalMode
}

func (j *job) SetIntervalMode(intervalMode bool) {
	j.intervalMode = intervalMode
}

func (j *job) Exclusive() bool {
	return j.exclusive
}

func (j *job) SetExclusive(exclusive bool) {
	j.exclusive = exclusive
}

func (j *job) Schedule() Schedule {
	return j.schedule
}

func (j *job) SetSchedule(schedule Schedule) {
	j.schedule = schedule
}

func (j *job) shouldRun(tick time.Time) bool {
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
func (j *job) run(tick time.Time) {
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
