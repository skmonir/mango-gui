package scheduler

import (
	"context"
	"fmt"
	"github.com/procyon-projects/chrono"
	"github.com/skmonir/mango-gui/backend/judge-framework/logger"
	"sync"
	"time"
)

var once sync.Once
var scheduleTasks map[string]chrono.ScheduledTask

func initScheduler() {
	if scheduleTasks == nil {
		once.Do(func() {
			scheduleTasks = map[string]chrono.ScheduledTask{}
		})
	}
}

func ScheduleOneTimeTask(taskId string, action func(), timeToSchedule time.Time) error {
	initScheduler()

	taskScheduler := chrono.NewDefaultTaskScheduler()
	task, err := taskScheduler.Schedule(func(ctx context.Context) {
		action()
	}, chrono.WithTime(timeToSchedule))
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	scheduleTasks[taskId] = task
	fmt.Println("Task scheduled successfully")
	return nil
}

func RemoveScheduledTask(taskId string) {
	initScheduler()

	if task, found := scheduleTasks[taskId]; found {
		delete(scheduleTasks, taskId)
		if !task.IsCancelled() {
			task.Cancel()
		}
	}
	fmt.Println("Scheduled task canceled successfully")
}
