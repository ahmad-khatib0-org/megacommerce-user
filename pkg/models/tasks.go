package models

import (
	"fmt"
	"time"
)

type TaskName string

type TaskFunc func()

const (
	TaskNameEmailBatching   TaskName = "email_batching"
	TaskNameSendVerifyEmail TaskName = "send_verify_email"
)

type TaskSendVerifyEmailPayload struct {
	Ctx   *Context `json:"ctx"`
	Email string   `json:"email"`
	Token string   `json:"token"`
	Hours int      `json:"hours"`
}

type ScheduledTask struct {
	Name                 string        `json:"name"`
	Interval             time.Duration `json:"interval"`
	Recurring            bool          `json:"recuring"`
	function             func()
	cancel               chan struct{}
	cancelled            chan struct{}
	fromNextIntervalTime bool
}

func CreateTask(name string, function TaskFunc, timeToExecution time.Duration) *ScheduledTask {
	return createTask(name, function, timeToExecution, false, false)
}

func CreateRecurringTask(name string, function TaskFunc, interval time.Duration) *ScheduledTask {
	return createTask(name, function, interval, true, false)
}

func CreateRecurringTaskFromNextIntervalTime(name string, function TaskFunc, interval time.Duration) *ScheduledTask {
	return createTask(name, function, interval, true, true)
}

func createTask(name string, function TaskFunc, interval time.Duration, recurring bool, fromNextIntervalTime bool) *ScheduledTask {
	task := &ScheduledTask{
		Name:                 name,
		Interval:             interval,
		Recurring:            recurring,
		function:             function,
		cancel:               make(chan struct{}),
		cancelled:            make(chan struct{}),
		fromNextIntervalTime: fromNextIntervalTime,
	}

	go func() {
		defer close(task.cancelled)

		var firstTick <-chan time.Time
		var ticker *time.Ticker

		if task.fromNextIntervalTime {
			currentTime := time.Now()
			first := currentTime.Truncate(interval)
			if first.Before(currentTime) {
				first = first.Add(interval)
			}
			firstTick = time.After(time.Until(first))
			ticker = &time.Ticker{C: nil}
		} else {
			firstTick = nil
			ticker = time.NewTicker(interval)
		}

		defer func() {
			ticker.Stop()
		}()

		for {
			select {
			case <-firstTick:
				ticker = time.NewTicker(interval)
				function()
			case <-ticker.C:
				function()
			case <-task.cancel:
				return
			}

			if !task.Recurring {
				break
			}
		}
	}()

	return task
}

func (t *ScheduledTask) Cancel() {
	close(t.cancel)
	<-t.cancelled
}

func (task *ScheduledTask) String() string {
	return fmt.Sprintf(
		"%s\nInterval: %s\nRecurring: %t\n",
		task.Name,
		task.Interval.String(),
		task.Recurring,
	)
}
