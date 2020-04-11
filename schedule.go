package schedly

import (
	"time"
)

/*schedule has a constrained function, allowing flexible configuration
 */
type schedule struct {
	tick           time.Duration
	every          time.Duration
	aligned        bool
	intervalMode   bool
	constraintFunc func(time.Time) bool
}

/*IntervalMode returns intervalMode setting. When true job schedule receives previous successful finish time as a param in CanRun method.

When false job schedule receives previous start time as a param in CanRun method*/
func (s *schedule) IntervalMode() bool {
	return s.intervalMode
}

/*SetIntervalMode sets intervalMode. When 'true' job schedule receives previous successful finish time as a param in CanRun method.

  When 'false' job schedule receives previous start time as a param in CanRun method*/
func (s *schedule) SetIntervalMode(intervalMode bool) {
	s.intervalMode = intervalMode
}

/*Aligned returns Aligned flag. If set to true tasks are launched at the beginning of configured 'Every' interval
 */
func (s *schedule) Aligned() bool {
	return s.aligned
}

/*SetAligned sets Aligned flag. If set to true tasks are launched at the beginning of configured 'Every' interval
 */
func (s *schedule) SetAligned(aligned bool) *schedule {
	s.aligned = aligned
	return s
}

/*CanRun checks if task can be launched based on current moment and last launch time.
When the job is started in Interval mode, last finish time is supplied as a lastRun parameter
*/
func (s *schedule) CanRun(moment time.Time, lastRun time.Time) bool {
	// safety gap here = 1.1 Making it smaller wouldn't affect execution really. It just needs to be in the interval (1+epsilon,2-epsilon)
	// where epsilon is time.Ticker precision error
	if lastRun.IsZero() && s.aligned && (moment.Sub(moment.Truncate(s.every))) > (11*s.tick)/10 {
		return false
	}
	if moment.UnixNano()-lastRun.Truncate(s.tick).UnixNano() < s.every.Nanoseconds() {
		return false
	}

	if s.constraintFunc != nil {
		return s.constraintFunc(moment)
	}

	return true
}

/*ConstraintFunc returns constraint function used to check if execution can be started.
 */
func (s *schedule) ConstraintFunc() func(time.Time) bool {
	return s.constraintFunc
}

/*SetConstraintFunc sets extra constraint function to check if execution can be started.
 */
func (s *schedule) SetConstraintFunc(constraintFunc func(time.Time) bool) *schedule {
	s.constraintFunc = constraintFunc
	return s
}

/*Every return an interval for running a task
 */
func (s *schedule) Every() time.Duration {
	return s.every
}

/*SetEvery sets an interval for running the task
 */
func (s *schedule) SetEvery(every time.Duration) *schedule {
	s.every = every
	return s
}

/*newSchedule creates new schedule. tick is a tick time from Scheduler. Used for aligning intervals internally
 */
func newSchedule(tick time.Duration) *schedule {
	return &schedule{
		tick:           tick,
		every:          tick,
		constraintFunc: nil,
	}
}
