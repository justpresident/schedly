# Schedly: A job scheduling library for Golang applications
![Go Test](https://github.com/justpresident/schedly/workflows/Go%20Test/badge.svg?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/justpresident/schedly)](https://goreportcard.com/report/github.com/justpresident/schedly)

## What is it?

Very simple and flexible task scheduling library for Golang applications without extra dependencies.

## Usage

### Install
```bash
go get github.com/justpresident/schedly
```

### Examples
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

### Schedule
Schedule have following configuration options:

**SetConstraintFunc()** to set up an extra constraint. Example:
```go
sched.NewSchedule(time.Minute).
  SetConstraintFunc(func (t time.Time) {
    return t.Weekday() != time.Sunday
  })
```
**SetAligned(true)** to set up aligned feature of the schedule.
It is automatically set to true if Aligned property has already been set for scheduler itself. You can enable/disable this property individually per job.

Schedule has following accessors:
```go
schedule.IntervalMode() // returns interval mode setting (read above)
schedule.Aligned() // returns aligned setting (read above)
schedule.Every() // returns launch interval
```

**SetIntervalMode()** When IntervalMode is set to `false` then next job launch time is calculated as 'last job start time' + configured interval.
When set to `true` the interval for launching the job is measured from the finish time of previous launch.
### Jobs
The job returned by `AddJob` or `Schedule` methods can be further configured:
```go
myJob := sched.Schedule(5*time.Second, "print time", func() {fmt.Print(time.Now())})

/* Avoid parallel runs. If previous run of the job has not finished yet,
 next run will be blocked until it finishes. It simply locks a mutex before starting.
More options of controlling this behaviour will be introduced in next releases.
*/ 
myJob.SetExclusive(true)

// You don't need to assign a job to a variable to configure it. You can do it in one line:
sched.Schedule(...).SetExclusive(true)
```

Also job has following accessors:
```go
myJob.Name() // returns the name of the job
myJob.LastRun() // returns previous start time. Returns 0 if the job has not been launched yet
myJob.LastFinish() // returns finish time of previous launch. Returns 0 if job has not been launched yet or has not finished yet.
myJob.Exclusive() // returns exclusive launch setting (read above)
```

## How to contribute
Please refer to the contribution guidelines here: [/docs/CONTRIBUTING.md](/docs/CONTRIBUTING.md)
