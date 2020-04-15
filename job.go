package schedly

import (
	"sync"
	"time"
)

/*job represents a configuration of a task scheduled for execution
 */
type job struct {
	jobFunc    func()
	name       string
	lastRun    time.Time
	lastFinish time.Time
	exclusive  bool
	lock       sync.Mutex
}

/*Name returns the job name*/
func (j *job) Name() string {
	return j.name
}

/*LastFinish return the last time the job finished execution*/
func (j *job) LastFinish() time.Time {
	return j.lastFinish
}

/*LastRun returns the last time the job has been launched*/
func (j *job) LastRun() time.Time {
	return j.lastRun
}

/*Exclusive returns Exclusive setting of the job. When 'true' it prevents from running multiple instances of the job at the same time*/
func (j *job) Exclusive() bool {
	return j.exclusive
}

/*SetExclusive configures Exclusive mode. When 'true' it prevents from running multiple instances of the job at the same time*/
func (j *job) SetExclusive(exclusive bool) *job {
	j.exclusive = exclusive
	return j
}

/* Run the job.
Parameter `tick` is time from internal ticker. Last job run time is set to this value before launch*/
func (j *job) Run(tick time.Time) {
	j.lastRun = tick

	if j.exclusive {
		j.lock.Lock()
		defer j.lock.Unlock()
	}
	j.jobFunc()
	j.lastFinish = time.Now()
}
