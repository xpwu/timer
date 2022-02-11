package task

import (
  "context"
  "github.com/xpwu/timer/scheduler"
  "github.com/xpwu/timer/task/delay"
  "github.com/xpwu/timer/task/fixed"
  "github.com/xpwu/timer/task/flag"
)


type noopTask struct{}

func (n *noopTask) Run(ctx context.Context, runtime scheduler.UnixTimeSecond) {
}

func newTaskRunner(task scheduler.Task) taskRunner {
  switch task[0] {
  case flag.Fixed:
    return fixed.FromBytes(task[1:])
  case flag.Delay:
    return delay.FromBytes(task[1:])
  default:
    return &noopTask{}
  }
}
