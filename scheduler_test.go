package schedly

import (
	"github.com/stretchr/testify/assert"
	"sync/atomic"
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

func TestScheduler_WaitForRunningTasks(t *testing.T) {
	sched := NewScheduler(time.Millisecond)

	var result int32 = 1
	sched.Schedule(time.Millisecond, "x", func() {
		time.Sleep(5 * time.Millisecond)
		atomic.CompareAndSwapInt32(&result, 1, 2)
	})

	sched.Start()

	time.Sleep(time.Millisecond)

	sched.Stop()
	waitStarted := time.Now()
	sched.WaitForRunningTasks()
	waited := time.Now().Sub(waitStarted).Milliseconds()
	assert.Less(t, int64(4), waited, "Shoul've waited for the task to finish")
	t.Logf("Waited for tasks to finish for %d ms", waited)

	assert.Equal(t, int32(2), atomic.LoadInt32(&result))
}
