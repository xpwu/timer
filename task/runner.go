package task

import (
  "context"
  "fmt"
  "github.com/xpwu/go-log/log"
  "github.com/xpwu/timer/scheduler"
  "time"
)

type taskRunner interface {
  Run(ctx context.Context, runtime scheduler.UnixTimeSecond)
}

type runner struct {
}

func (r *runner) Run(ctx context.Context, schedulerTime scheduler.UnixTimeSecond, tasks []scheduler.Task) {
  ctx,logger := log.WithCtx(ctx)
  logger.PushPrefix(fmt.Sprintf("run %d", schedulerTime))

  doneC := 0
  done := make(chan scheduler.Task, len(tasks))

  for _,task := range tasks {
    go func(task scheduler.Task) {
      defer func() {
        if r := recover(); r != nil {
          logger.Fatal(r)
          done <- task
          return
        }

        done <- nil
      }()
      newTaskRunner(task).Run(ctx, schedulerTime)
    }(task)

    retry := make([]scheduler.Task, 0, len(tasks))
    for t := range done {
      if t != nil {
        retry = append(retry, t)
      }
      doneC++
      if doneC >= len(tasks) {
        break
      }
    }

    if len(retry) == 0 {
      return
    }

    // 延后5s再试
    next := scheduler.UnixTimeSecond(time.Now().Unix() + 5)
    scheduler.AddTask(next, retry)
  }
}

func init() {
  scheduler.SetRunner(&runner{})
}


