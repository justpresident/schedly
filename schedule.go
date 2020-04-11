package schedly

import (
	"time"
)

/*Schedule is used by Scheduler to understand when tasks need to be launched
 */
type Schedule interface {
	CanRun(moment time.Time, lastRun time.Time) bool
}

/*ConstrainedSchedule has a constrained function, allowing flexible configuration
 */
type ConstrainedSchedule struct {
	tick           time.Duration
	every          time.Duration
	aligned        bool
	constraintFunc func(time.Time) bool
}

/*Aligned returns Aligned flag. If set to true tasks are launched at the beginning of configured 'Every' interval
 */
func (s *ConstrainedSchedule) Aligned() bool {
	return s.aligned
}

/*SetAligned sets Aligned flag. If set to true tasks are launched at the beginning of configured 'Every' interval
 */
func (s *ConstrainedSchedule) SetAligned(aligned bool) *ConstrainedSchedule {
	s.aligned = aligned
	return s
}

/*CanRun checks if task can be launched based on current moment and last launch time.
	When the job is started in Interval mode, last finish time is supplied as a lastRun parameter
*/
func (s *ConstrainedSchedule) CanRun(moment time.Time, lastRun time.Time) bool {
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
func (s *ConstrainedSchedule) ConstraintFunc() func(time.Time) bool {
	return s.constraintFunc
}

/*SetConstraintFunc sets extra constraint function to check if execution can be started.
*/
func (s *ConstrainedSchedule) SetConstraintFunc(constraintFunc func(time.Time) bool) *ConstrainedSchedule {
	s.constraintFunc = constraintFunc
	return s
}

/*Every return an interval for running a task
*/
func (s *ConstrainedSchedule) Every() time.Duration {
	return s.every
}

/*SetEvery sets an interval for running the task
*/
func (s *ConstrainedSchedule) SetEvery(every time.Duration) *ConstrainedSchedule {
	s.every = every
	return s
}

/*NewConstrainedSchedule creates new ConstrainedSchedule. tick is a tick time from Scheduler. Used for aligning intervals internally
 */
func NewConstrainedSchedule(tick time.Duration) *ConstrainedSchedule {
	return &ConstrainedSchedule{
		tick:           tick,
		every:          tick,
		constraintFunc: func(moment time.Time) bool { return true },
	}
}
