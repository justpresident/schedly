package schedly

import (
	"time"
)

type Schedule interface {
	CanRun(moment time.Time, lastRun time.Time) bool
}

type ConstrainedSchedule struct {
	tick           time.Duration
	every          time.Duration
	aligned        bool
	constraintFunc func(time.Time) bool
}

func (s *ConstrainedSchedule) Aligned() bool {
	return s.aligned
}

func (s *ConstrainedSchedule) SetAligned(aligned bool) *ConstrainedSchedule {
	s.aligned = aligned
	return s
}

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

func (s *ConstrainedSchedule) ConstraintFunc() func(time.Time) bool {
	return s.constraintFunc
}

func (s *ConstrainedSchedule) SetConstraintFunc(constraintFunc func(time.Time) bool) *ConstrainedSchedule {
	s.constraintFunc = constraintFunc
	return s
}

func (s *ConstrainedSchedule) Every() time.Duration {
	return s.every
}

func (s *ConstrainedSchedule) SetEvery(every time.Duration) *ConstrainedSchedule {
	s.every = every
	return s
}

func NewConstrainedSchedule(tick time.Duration) *ConstrainedSchedule {
	return &ConstrainedSchedule{
		tick:           tick,
		every:          tick,
		constraintFunc: func(moment time.Time) bool { return true },
	}
}
