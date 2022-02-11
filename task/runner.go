package task

import (
  "context"
  "github.com/xpwu/timer/scheduler"
)

type taskRunner interface {
  Run(ctx context.Context, runtime scheduler.UnixTimeSecond)
}

type runner struct {
}

func (r *runner) Run(ctx context.Context, schedulerTime scheduler.UnixTimeSecond, tasks []scheduler.Task) {
  // todo 1、crash; 2、async
  for _,task := range tasks {
    newTaskRunner(task).Run(ctx, schedulerTime)
  }
}

func init() {
  scheduler.SetRunner(&runner{})
}


