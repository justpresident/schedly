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

"fmt"
"github.com/justpresident/schedly"
"time"

)

func main() {
  // Create a scheduler with a tick = second
  sched := schedly.NewScheduler(time.Second).
    SetAligned(true) // start launching tasks at the beginning of their period - e.g.
                     // minutely tasks at the beginning of every minute.
  
   // Print current time every 5 seconds
  sched.Schedule(5*time.Second, "print time", func() {fmt.Print(time.Now())})

  sched.Start()

  sched.WaitUntilStopped() // to wait indefinitely until program exits or someone calls Stop()
}
```
schedly doesn't take any assumptions and allows you to implement any schedule easily.

Let's define a schedule for our job - every minute of NASDAQ trading time.
```go
  nasdaqTimeSchedule := sched.NewSchedule(time.Minute). // execute tasks minutely
	SetConstraintFunc(isNasdaqTime) // where isNasdaqTime := func(moment time.Time) bool {...} 
  
  sched.AddJob(
    nasdaqTimeSchedule,
    "download TSLA stock price",
    func() { updateQuote("TSLA") },
  )
```

To stop scheduler from running new tasks:
```go
  ... // do whatever your application needs to do and then call
  sched.Stop() // to stop scheduling new tasks
  // and optionally
  sched.WaitForRunningTasks() // to gracefully wait for all the running tasks to finish
```
