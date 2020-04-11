# schedly: A job scheduling library for Golang applications
![Go Test](https://github.com/justpresident/schedly/workflows/Go%20Test/badge.svg?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/justpresident/schedly)](https://goreportcard.com/report/github.com/justpresident/schedly)

Very simple and flexible task scheduling library for Golang applications without extra dependencies.

**This is a WIP project. The interface will most probably change in the nearest future. Please stay tuned for a final release**

Install:
```bash
go get github.com/justpresident/schedly
```

Examples:
```go
package main

import (
  "time"
  "github.com/justpresident/schedly"
)

func main() {
  // Create a scheduler with a tick = second
  sched := schedly.NewScheduler(time.Second).
    SetAligned(true) // start launching tasks at the beginning of their period - e.g.
                     // minutely tasks at the beginning of every minute.
  
  // Let's define a schedule for our job - every minute of NASDAQ trading time
  nasdaqTimeSchedule := schedly.NewConstrainedSchedule(sched.Tick()).
		SetEvery(time.Minute). // execute minutely
		SetAligned(true).      // at the beginning of the minute
		SetConstraintFunc(isNasdaqTime) // at isNasdaqTime := func(moment time.Time) bool {...} 
  
  sched.AddJob(
		"download TSLA stock price",
		func() { updateQuote("TSLA") },
	).SetSchedule(nasdaqTimeSchedule)
  
  stopped := sched.Start()
  
  ... // do whatever your application needs to do and then call
  sched.Stop() // to stop the scheduler
  // or
  sched.WaitAndStop() // to gracefully wait for all the running tasks to finish and then stop
  
  <-stopped // wait until scheduler stops
}
```
