package models

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	name := "a testing task"
	taskTime := time.Second * 2

	executionCount := new(int32)
	testFunc := func() {
		atomic.AddInt32(executionCount, 1)
	}

	task := CreateTask(name, testFunc, taskTime)
	assert.EqualValues(t, 0, atomic.LoadInt32(executionCount))

	time.Sleep(taskTime + time.Second)
	assert.EqualValues(t, 1, atomic.LoadInt32(executionCount))
	assert.Equal(t, name, task.Name)
	assert.Equal(t, taskTime, task.Interval)
	assert.False(t, task.Recurring)
}

func TestCreateRecurringTask(t *testing.T) {
	name := "a testing task"
	taskTime := time.Second * 2

	executionCount := new(int32)
	testFunc := func() {
		atomic.AddInt32(executionCount, 1)
	}

	task := CreateRecurringTask(name, testFunc, taskTime)
	assert.EqualValues(t, 0, atomic.LoadInt32(executionCount))

	time.Sleep(taskTime + time.Second)
	assert.EqualValues(t, 1, atomic.LoadInt32(executionCount))

	time.Sleep(taskTime)
	assert.EqualValues(t, 2, atomic.LoadInt32(executionCount))

	assert.Equal(t, name, task.Name)
	assert.Equal(t, taskTime, task.Interval)
	assert.True(t, task.Recurring)

	task.Cancel()
}

func TestCancelTask(t *testing.T) {
	name := "a testing task"
	taskTime := time.Second * 2

	executionCount := new(int32)
	testFunc := func() {
		atomic.AddInt32(executionCount, 1)
	}

	task := CreateTask(name, testFunc, taskTime)
	assert.EqualValues(t, 0, atomic.LoadInt32(executionCount))
	task.Cancel()

	time.Sleep(taskTime + time.Second)
	assert.EqualValues(t, 0, atomic.LoadInt32(executionCount))
}

func TestCreateRecurringTaskFromNextIntervalTime(t *testing.T) {
	name := "recurring task starting from next interval time"
	taskTime := time.Second * 3

	var executionTime time.Time
	var mu sync.Mutex
	testFunc := func() {
		mu.Lock()
		executionTime = time.Now()
		mu.Unlock()
	}

	task := CreateRecurringTaskFromNextIntervalTime(name, testFunc, taskTime)
	defer task.Cancel()

	time.Sleep(taskTime)
	mu.Lock()
	expectedSeconds := executionTime.Second()
	mu.Unlock()
	// Ideally we would expect 0, but in busy CI environments it can lag
	// by a second. If we see a lag of more than a second, we would need to disable
	// the test entirely.
	assert.LessOrEqual(t, expectedSeconds%3, 1)

	assert.Equal(t, name, task.Name)
	assert.Equal(t, taskTime, task.Interval)
	assert.True(t, task.Recurring)
}
