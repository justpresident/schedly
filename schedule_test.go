package schedly

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConstrainedSchedule_Aligned(t *testing.T) {
	assert.True(t, (&schedule{aligned: true}).Aligned())
	assert.False(t, (&schedule{aligned: false}).Aligned())
}

func mkTime(v string) time.Time {
	tm, err := time.Parse(time.RFC3339, v)
	if err != nil {
		panic(fmt.Sprintf("Wrong time format: %s", v))
	}
	return tm
}

func TestConstrainedSchedule_CanRun(t *testing.T) {
	type fields struct {
		tick           time.Duration
		every          time.Duration
		aligned        bool
		constraintFunc func(time.Time) bool
	}
	type args struct {
		moment  time.Time
		j *job
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Wait for alignment",
			fields: fields{
				tick:           time.Second,
				every:          time.Minute,
				aligned:        true,
				constraintFunc: nil,
			},
			args: args{
				moment:  mkTime("2006-01-02T15:04:05Z"),
				j : &job{},
			},
			want: false,
		},
		{
			name: "AlignmentPass",
			fields: fields{
				tick:           time.Second,
				every:          time.Minute,
				aligned:        true,
				constraintFunc: nil,
			},
			args: args{
				moment:  mkTime("2006-01-02T15:04:00Z"),
				j : &job{},
			},
			want: true,
		},
		{
			name: "ConstraintFuncEvenMinute",
			fields: fields{
				tick:           time.Second,
				every:          time.Minute,
				aligned:        false,
				constraintFunc: func(tm time.Time) bool { return tm.Minute()%2 == 0 },
			},
			args: args{
				moment:  mkTime("2006-01-02T15:04:00Z"),
				j : &job{},
			},
			want: true,
		},
		{
			name: "ConstraintFuncOddMinute",
			fields: fields{
				tick:           time.Second,
				every:          time.Minute,
				aligned:        false,
				constraintFunc: func(tm time.Time) bool { return tm.Minute()%2 == 0 },
			},
			args: args{
				moment:  mkTime("2006-01-02T15:05:00Z"),
				j : &job{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &schedule{
				tick:           tt.fields.tick,
				every:          tt.fields.every,
				aligned:        tt.fields.aligned,
				constraintFunc: tt.fields.constraintFunc,
			}
			if got := s.CanRun(tt.args.moment, tt.args.j); got != tt.want {
				t.Errorf("CanRun() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConstrainedSchedule_ConstraintFunc(t *testing.T) {

	constraintFunc := func(tm time.Time) bool { return tm.IsZero() }

	gotFunc := (&schedule{constraintFunc: constraintFunc}).ConstraintFunc()
	assert.Equal(t, gotFunc(time.Time{}), true)
	assert.Equal(t, gotFunc(time.Now()), false)

}

func TestConstrainedSchedule_Every(t *testing.T) {
	assert.Equal(t, (&schedule{}).SetEvery(time.Minute).Every(), time.Minute)
	assert.Equal(t, (&schedule{every: time.Millisecond}).SetEvery(time.Second).Every(), time.Second)
}
