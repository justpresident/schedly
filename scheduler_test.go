package schedly

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestStartStop(t *testing.T) {
	sched := NewScheduler(time.Millisecond)
	sched.Start()

	sched.Stop()

	sched.WaitUntilStopped()
}

func TestAligned(t *testing.T) {
	sched := NewScheduler(time.Millisecond).SetAligned(true)

	var startedAt time.Time
	sched.Schedule(5*time.Millisecond, "x", func() {
		startedAt = time.Now()
		sched.Stop()
	})

	sched.Start()

	sched.WaitUntilStopped()
	assert.Less(t, startedAt.Sub(startedAt.Truncate(5*time.Millisecond)).Nanoseconds(), time.Millisecond.Nanoseconds())
}

func TestAddJob(t *testing.T) {
	sched := NewScheduler(time.Millisecond)
	tasks := []string{"x", "y"}

	resultChan := make(chan string)
	for tNum := 0; tNum < len(tasks); tNum++ {
		tName := tasks[tNum]

		sched.Schedule(time.Millisecond, tName, func() {
			resultChan <- tName
		})
	}

	sched.Start()

	results := make(map[string]int)
	for len(results) != len(tasks) {
		select {
		case r := <-resultChan:
			results[r]++
		}
	}
	sched.Stop()

	for _, tName := range tasks {
		if _, ok := results[tName]; !ok {
			t.Fatalf("Task %s has not been launched", tName)
		}
	}


	sched.WaitUntilStopped()
}
