package schedly

import (
	"testing"
	"time"
)

func TestStartStop(t *testing.T) {
	sched := NewScheduler(time.Millisecond)
	done := sched.Start()

	sched.Stop()

	<-done
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

	done := sched.Start()

	results := make(map[string]int)
	for len(results) != len(tasks) {
		select {
		case r := <-resultChan:
			results[r]++
		}
	}
	for _, tName := range tasks {
		if _, ok := results[tName]; !ok {
			t.Fatalf("Task %s has not been launched", tName)
		}
	}

	sched.Stop()

	<-done
}
